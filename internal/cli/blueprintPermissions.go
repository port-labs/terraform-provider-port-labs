package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) GetBlueprintPermissions(ctx context.Context, blueprintID string) (*BlueprintPermissions, int, error) {
	pppb := &PortBlueprintPermissionsBody{}
	url := "v1/blueprints/{blueprint_identifier}/permissions"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pppb).
		SetPathParam("blueprint_identifier", blueprintID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pppb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to get blueprint permissions, got: %s", resp.Body())
	}
	return &pppb.BlueprintPermissions, resp.StatusCode(), nil

}

func (c *PortClient) UpdateBlueprintPermissions(ctx context.Context, blueprintID string, permissions *BlueprintPermissions) (*BlueprintPermissions, error) {
	url := "v1/blueprints/{blueprint_identifier}/permissions"

	resp, err := c.Client.R().
		SetBody(permissions).
		SetContext(ctx).
		SetPathParam("blueprint_identifier", blueprintID).
		Patch(url)
	if err != nil {
		return nil, err
	}
	var pppb PortBlueprintPermissionsBody
	err = json.Unmarshal(resp.Body(), &pppb)
	if err != nil {
		return nil, err
	}
	if !pppb.OK {
		return nil, fmt.Errorf("failed to update blueprint permissions, got: %s", resp.Body())
	}
	return &pppb.BlueprintPermissions, nil
}
