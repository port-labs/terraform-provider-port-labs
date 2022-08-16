package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) CreatePermissions(ctx context.Context, clientID string, scopes ...string) error {
	url := "v0.1/apps/{app_id}/permissions"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetPathParam("app_id", clientID).
		SetBody(map[string]interface{}{
			"permissions": scopes,
		}).
		Post(url)
	if err != nil {
		return err
	}
	b := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &b)
	if err != nil {
		return err
	}
	if !b["ok"].(bool) {
		return fmt.Errorf("failed to create permissions: %s", resp.Body())
	}
	return nil
}
