package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

const oauthAppsUrl = "/v1/organization/oauth-apps"

func (c *PortClient) ReadOAuthApp(ctx context.Context, oauthAppID string) (*OAuthApp, int, error) {
	pb := &OAuthAppBody{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("oauth_app_id", oauthAppID).
		Get(oauthAppsUrl + "/{oauth_app_id}")
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read oauth app, got: %s", resp.Body())
	}
	return &pb.OAuthApp, resp.StatusCode(), nil
}

func (c *PortClient) CreateOAuthApp(ctx context.Context, oauthApp *OAuthAppCreate) (*OAuthApp, error) {
	resp, err := c.Client.R().
		SetBody(oauthApp).
		SetContext(ctx).
		Post(oauthAppsUrl)
	if err != nil {
		return nil, err
	}
	var pb OAuthAppBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create oauth app, got: %s", resp.Body())
	}
	return &pb.OAuthApp, nil
}

func (c *PortClient) UpdateOAuthApp(ctx context.Context, oauthAppID string, oauthApp *OAuthAppUpdate) (*OAuthApp, error) {
	resp, err := c.Client.R().
		SetBody(oauthApp).
		SetContext(ctx).
		SetPathParam("oauth_app_id", oauthAppID).
		Patch(oauthAppsUrl + "/{oauth_app_id}")
	if err != nil {
		return nil, err
	}
	var pb OAuthAppBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update oauth app, got: %s", resp.Body())
	}
	return &pb.OAuthApp, nil
}

func (c *PortClient) DeleteOAuthApp(ctx context.Context, oauthAppID string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("oauth_app_id", oauthAppID).
		Delete(oauthAppsUrl + "/{oauth_app_id}")
	if err != nil {
		return err
	}
	var pb PortBodyDelete
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return err
	}
	if !pb.Ok {
		return fmt.Errorf("failed to delete oauth app, got: %s", string(resp.Body()))
	}
	return nil
}
