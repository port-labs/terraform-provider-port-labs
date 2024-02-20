package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) GetPagePermissions(ctx context.Context, pageID string) (*PagePermissions, int, error) {
	pppb := &PortPagePermissionsBody{}
	url := "v1/pages/{page_identifier}/permissions"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pppb).
		SetPathParam("page_identifier", pageID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pppb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to get page permissions, got: %s", resp.Body())
	}
	return &pppb.PagePermissions, resp.StatusCode(), nil

}

func (c *PortClient) UpdatePagePermissions(ctx context.Context, pageID string, permissions *PagePermissions) (*PagePermissions, error) {
	url := "v1/pages/{page_identifier}/permissions"

	resp, err := c.Client.R().
		SetBody(permissions).
		SetContext(ctx).
		SetPathParam("page_identifier", pageID).
		Patch(url)
	if err != nil {
		return nil, err
	}
	var pppb PortPagePermissionsBody
	err = json.Unmarshal(resp.Body(), &pppb)
	if err != nil {
		return nil, err
	}
	if !pppb.OK {
		return nil, fmt.Errorf("failed to update page permissions, got: %s", resp.Body())
	}
	return &pppb.PagePermissions, nil
}
