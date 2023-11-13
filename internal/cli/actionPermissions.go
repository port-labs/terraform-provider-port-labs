package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadActionPermissions(ctx context.Context, actionID string, blueprintID string) (*ActionPermissions, int, error) {
	pb := &PortActionPermissionsBody{}
	url := "/v1/blueprints/{blueprint_identifier}/actions/{action_identifier}/permissions"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("action_identifier", actionID).
		SetPathParam("blueprint_identifier", blueprintID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read action permissions, got: %s", resp.Body())
	}
	return &pb.Permissions, resp.StatusCode(), nil
}

func (c *PortClient) UpdateActionPermissions(ctx context.Context, actionID string, blueprintID string, ap *ActionPermissions) (*ActionPermissions, error) {
	url := "/v1/blueprints/{blueprint_identifier}/actions/{action_identifier}/permissions"
	resp, err := c.Client.R().
		SetBody(ap).
		SetContext(ctx).
		SetPathParam("blueprint_identifier", blueprintID).
		SetPathParam("action_identifier", actionID).
		Patch(url)
	if err != nil {
		return nil, err
	}
	var pb PortActionPermissionsBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update action permissions, got: %s", resp.Body())
	}
	return &pb.Permissions, nil
}
