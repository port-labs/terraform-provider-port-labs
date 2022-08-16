package cli

import (
	"context"
	"fmt"
)

func (c *PortClient) CreatePermissions(ctx context.Context, clientID string, scopes ...string) error {
	url := "v0.1/apps/{app_id}/permissions"
	responseBody := make(map[string]interface{})
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetPathParam("app_id", clientID).
		SetBody(map[string]interface{}{
			"permissions": scopes,
		}).
		SetResult(&responseBody).
		Post(url)
	if err != nil {
		return err
	}
	if statusOK, ok := responseBody["ok"].(bool); ok && !statusOK {
		return fmt.Errorf("failed to create permissions: %s", resp.Body())
	} else if !ok {
		return fmt.Errorf("failed to create permissionz: %s", resp.Body())
	}
	return nil
}
