package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const FeatureFlagUsersAndTeamsV2 = "USERS_AND_TEAMS_OWNERSHIP_V2"

func (c *PortClient) ReadTeam(ctx context.Context, teamName string) (*Team, int, error) {
	var pt PortTeamBody
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("name", teamName).
		SetQueryParamsFromValues(url.Values{
			"fields": []string{"name", "provider", "description", "createdAt", "updatedAt", "users.firstName",
				"users.status", "users.email"},
		}).
		SetResult(&pt).
		Get("v1/teams/{name}")
	if err != nil {
		return nil, 0, err
	} else if resp.IsError() || !pt.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read team, got: %s", resp.Body())
	}

	portTeam := &PortTeam{
		Name:        pt.Team.Name,
		Description: pt.Team.Description,
		CreatedAt:   pt.Team.CreatedAt,
		UpdatedAt:   pt.Team.UpdatedAt,
		Provider:    pt.Team.Provider,
	}

	team, err := c.enrichTeamFromTeamEntityWithRetry(ctx, portTeam)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to enrich team from entity: %w", err)
	}

	team.Users = make([]string, len(pt.Team.Users))

	for i, u := range pt.Team.Users {
		team.Users[i] = u.Email
	}

	return team, resp.StatusCode(), nil
}

const teamsBaseUrl = "v1/teams"
const teamSpecificUrl = teamsBaseUrl + "/{name}"

func (c *PortClient) CreateTeam(ctx context.Context, team *PortTeam) (*Team, error) {
	resp, err := c.Client.R().
		SetBody(team).
		SetContext(ctx).
		Post(teamsBaseUrl)

	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create team, got: %s", resp.Body())
	}

	return c.enrichTeamFromTeamEntityWithRetry(ctx, &pb.Team)
}

func (c *PortClient) UpdateTeam(ctx context.Context, teamName string, team *PortTeam) (*Team, error) {
	resp, err := c.Client.R().
		SetBody(team).
		SetContext(ctx).
		SetPathParam("name", teamName).
		Put(teamSpecificUrl)

	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update team, got: %s", resp.Body())
	}

	return c.enrichTeamFromTeamEntityWithRetry(ctx, &pb.Team)
}

func (c *PortClient) DeleteTeam(ctx context.Context, teamName string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("name", teamName).
		Delete(teamSpecificUrl)

	if err != nil {
		return err
	}
	var pb PortBodyDelete
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return err
	}

	if !(pb.Ok) {
		return fmt.Errorf("failed to delete team. got:\n%s", string(resp.Body()))
	}
	return nil
}

const MissingTeamEntityError = utils.StringErr("team entity is missing")

const teamEntityRetryAttempts = 10

// Retry timing for a lagging _team entity. Exponential backoff from teamEntityRetryDelay, capped per wait at
// teamEntityRetryMaxDelay so the worst case stays well under a minute. Vars so tests can shrink them.
var (
	teamEntityRetryDelay     = time.Second
	teamEntityRetryMaxJitter = time.Second
	teamEntityRetryMaxDelay  = 5 * time.Second
)

func (c *PortClient) enrichTeamFromTeamEntityWithRetry(ctx context.Context, portTeam *PortTeam) (*Team, error) {
	return retry.DoWithData(
		func() (*Team, error) { return c.enrichTeamFromTeamEntity(ctx, portTeam) },
		retry.Context(ctx),
		retry.LastErrorOnly(true),
		// Only a missing _team entity is retried: on ownership-v2 orgs it is created asynchronously right after the
		// team, so it can lag behind a create/update. Every other error fails immediately.
		retry.Attempts(teamEntityRetryAttempts),
		retry.RetryIf(func(err error) bool { return errors.Is(err, MissingTeamEntityError) }),
		retry.Delay(teamEntityRetryDelay),
		retry.MaxJitter(teamEntityRetryMaxJitter),
		retry.MaxDelay(teamEntityRetryMaxDelay),
	)
}

func (c *PortClient) enrichTeamFromTeamEntity(ctx context.Context, portTeam *PortTeam) (*Team, error) {
	team := &Team{PortTeam: *portTeam}

	isUsersAndTeamsV2, err := c.HasFeatureFlags(ctx, FeatureFlagUsersAndTeamsV2)
	if err != nil {
		return nil, fmt.Errorf("failed to read feature flags: %w", err)
	}

	// Look up the team's _team entity to expose its identifier, which is what entity teams[] arrays and permission
	// configurations reference.
	searchResults, searchErr := c.Search(ctx, &SearchRequestQuery{
		Query: &map[string]any{
			"combinator": "and",
			"rules": []map[string]any{
				{"property": "$blueprint", "operator": "=", "value": "_team"},
				{"property": "$title", "operator": "=", "value": portTeam.Name},
			},
		},
		Include: []string{"identifier"},
	})

	if !isUsersAndTeamsV2 {
		// Orgs without ownership v2 are not guaranteed to have a _team blueprint or _team entities, so a missing
		// blueprint (404) or no match leaves the identifier null instead of failing the operation. Any other search
		// failure (auth, 5xx, timeouts) is returned so a transient error can't drop a known identifier from state.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if searchErr != nil {
			var se *SearchError
			if errors.As(searchErr, &se) && se.StatusCode == http.StatusNotFound {
				return team, nil
			}
			return nil, searchErr
		}
		switch len(searchResults.Entities) {
		case 0:
		case 1:
			team.Identifier = &searchResults.Entities[0].Identifier
		default:
			tflog.Warn(ctx, "found multiple _team entities for team, leaving its identifier null", map[string]any{
				"team":    portTeam.Name,
				"matches": len(searchResults.Entities),
			})
		}
		return team, nil
	}

	if searchErr != nil {
		return nil, searchErr
	}
	if len(searchResults.Entities) == 0 {
		return nil, MissingTeamEntityError
	}
	if len(searchResults.Entities) != 1 {
		return nil, fmt.Errorf("failed to read team results, got %d results instead of 1", len(searchResults.Entities))
	}
	team.Identifier = &searchResults.Entities[0].Identifier

	return team, nil
}
