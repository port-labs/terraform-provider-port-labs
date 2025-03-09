package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadWebhook(ctx context.Context, webhookID string) (*Webhook, int, error) {
	pb := &PortBody{}
	url := "v1/webhooks/{webhook_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("webhook_identifier", webhookID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read webhook, got: %s", resp.Body())
	}
	return &pb.Webhook, resp.StatusCode(), nil
}

func (c *PortClient) CreateWebhook(ctx context.Context, webhook *Webhook) (*Webhook, error) {
	url := "v1/webhooks"
	resp, err := c.Client.R().
		SetBody(webhook).
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
		return nil, fmt.Errorf("failed to create webhook, got: %s", resp.Body())
	}

	return &pb.Webhook, nil
}

func (c *PortClient) UpdateWebhook(ctx context.Context, webhookID string, webhook *Webhook) (*Webhook, error) {
	url := "v1/webhooks/{webhook_identifier}"
	resp, err := c.Client.R().
		SetBody(webhook).
		SetContext(ctx).
		SetPathParam("webhook_identifier", webhookID).
		Put(url)
	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update webhook, got: %s", resp.Body())
	}

	return &pb.Webhook, nil
}

func (c *PortClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	url := "v1/webhooks/{webhook_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("webhook_identifier", webhookID).
		Delete(url)
	if err != nil {
		return err
	}
	responseBody := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &responseBody)
	if err != nil {
		return err
	}
	if !(responseBody["ok"].(bool)) {
		return fmt.Errorf("failed to delete webhook. got:\n%s", string(resp.Body()))
	}
	return nil
}
