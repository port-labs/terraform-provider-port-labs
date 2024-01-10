package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadAction(ctx context.Context, id string) (*Action, int, error) {
	pb := &PortBody{}
	url := "v1/actions/{action_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("action_identifier", id).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read action, got: %s", resp.Body())
	}
	return &pb.Action, resp.StatusCode(), nil
}

func (c *PortClient) CreateAction(ctx context.Context, action *Action) (*Action, error) {
	url := "v1/actions"
	resp, err := c.Client.R().
		SetBody(action).
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
		return nil, fmt.Errorf("failed to create action, got: %s", resp.Body())
	}
	return &pb.Action, nil
}

func (c *PortClient) UpdateAction(ctx context.Context, actionID string, action *Action) (*Action, error) {
	url := "v1/actions/{action_identifier}"
	resp, err := c.Client.R().
		SetBody(action).
		SetContext(ctx).
		SetPathParam("action_identifier", actionID).
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
		return nil, fmt.Errorf("failed to create action, got: %s", resp.Body())
	}
	return &pb.Action, nil
}

func (c *PortClient) DeleteAction(ctx context.Context, actionID string) error {
	url := "v1/actions/{action_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("action_identifier", actionID).
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
		return fmt.Errorf("failed to delete action. got:\n%s", string(resp.Body()))
	}
	return nil
}
