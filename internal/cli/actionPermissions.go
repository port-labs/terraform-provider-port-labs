package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) GetActionPermissions(ctx context.Context, actionID string) (*ActionPermissions, int, error) {
	pb := &PortBody{}
	url := "v1/actions/{action_identifier}/permissions"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("action_identifier", actionID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to get action permissions, got: %s", resp.Body())
	}
	return &pb.ActionPermissions, resp.StatusCode(), nil

}

func (c *PortClient) UpdateActionPermissions(ctx context.Context, actionID string, permissions *ActionPermissions) (*ActionPermissions, error) {
	url := "v1/actions/{action_identifier}/permissions"

	resp, err := c.Client.R().
		SetBody(permissions).
		SetContext(ctx).
		SetPathParam("action_identifier", actionID).
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
		return nil, fmt.Errorf("failed to update action permissions, got: %s", resp.Body())
	}
	return &pb.ActionPermissions, nil
}
