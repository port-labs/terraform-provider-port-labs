package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"net/url"
	"time"
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

func (c *PortClient) enrichTeamFromTeamEntityWithRetry(ctx context.Context, portTeam *PortTeam) (*Team, error) {
	return retry.DoWithData(
		func() (*Team, error) { return c.enrichTeamFromTeamEntity(ctx, portTeam) },
		retry.LastErrorOnly(true),
		retry.Attempts(1),
		retry.AttemptsForError(10, MissingTeamEntityError),
		retry.Delay(time.Second),
		retry.MaxJitter(time.Second),
	)
}

func (c *PortClient) enrichTeamFromTeamEntity(ctx context.Context, portTeam *PortTeam) (*Team, error) {
	team := &Team{PortTeam: *portTeam}

	isUsersAndTeamsV2, err := c.HasFeatureFlags(ctx, FeatureFlagUsersAndTeamsV2)
	if err != nil {
		return nil, fmt.Errorf("failed to read feature flags: %w", err)
	}

	if isUsersAndTeamsV2 {
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
	}

	return team, nil
}
