package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadTeam(ctx context.Context, teamName string) (*Team, int, error) {
	url := "v1/teams/{name}?fields=name&fields=provider&fields=description&fields=createdAt&fields=updatedAt&fields=users.firstName&fields=users.status&fields=users.email"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("name", teamName).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}

	var pt PortTeamBody
	err = json.Unmarshal(resp.Body(), &pt)
	if err != nil {
		return nil, resp.StatusCode(), err
	}

	if !pt.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read team, got: %s", resp.Body())
	}
	team := &Team{
		Name:        pt.Team.Name,
		Description: pt.Team.Description,
		CreatedAt:   pt.Team.CreatedAt,
		UpdatedAt:   pt.Team.UpdatedAt,
		Provider:    pt.Team.Provider,
	}

	team.Users = make([]string, len(pt.Team.Users))

	for i, u := range pt.Team.Users {
		team.Users[i] = u.Email
	}

	return team, resp.StatusCode(), nil
}

func (c *PortClient) CreateTeam(ctx context.Context, team *Team) (*Team, error) {
	url := "v1/teams"
	resp, err := c.Client.R().
		SetBody(team).
		SetContext(ctx).
		Post(url)

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
	return &pb.Team, nil
}

func (c *PortClient) UpdateTeam(ctx context.Context, teamName string, team *Team) (*Team, error) {
	url := "v1/teams/{name}"
	resp, err := c.Client.R().
		SetBody(team).
		SetContext(ctx).
		SetPathParam("name", teamName).
		Patch(url)

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
	return &pb.Team, nil
}

func (c *PortClient) DeleteTeam(ctx context.Context, teamName string) error {
	url := "v1/teams/{name}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("name", teamName).
		Delete(url)

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
