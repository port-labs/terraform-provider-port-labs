package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadWorkflow(ctx context.Context, id string) (*Workflow, int, error) {
	pb := &PortBody{}
	url := "v1/workflows/{workflow_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("workflow_identifier", id).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read workflow, got: %s", resp.Body())
	}
	return &pb.Workflow, resp.StatusCode(), nil
}

func (c *PortClient) CreateWorkflow(ctx context.Context, workflow *Workflow) (*Workflow, error) {
	url := "v1/workflows"
	resp, err := c.Client.R().
		SetBody(workflow).
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
		return nil, fmt.Errorf("failed to create workflow, got (HTTP %d): %s", resp.StatusCode(), resp.Body())
	}
	return &pb.Workflow, nil
}

func (c *PortClient) UpdateWorkflow(ctx context.Context, workflowID string, workflow *Workflow) (*Workflow, error) {
	url := "v1/workflows/{workflow_identifier}"
	resp, err := c.Client.R().
		SetBody(workflow).
		SetContext(ctx).
		SetPathParam("workflow_identifier", workflowID).
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
		return nil, fmt.Errorf("failed to update workflow, got (HTTP %d): %s", resp.StatusCode(), resp.Body())
	}
	return &pb.Workflow, nil
}

func (c *PortClient) DeleteWorkflow(ctx context.Context, workflowID string) error {
	url := "v1/workflows/{workflow_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("workflow_identifier", workflowID).
		Delete(url)
	if err != nil {
		return err
	}
	responseBody := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &responseBody)
	if err != nil {
		return err
	}
	if ok, _ := responseBody["ok"].(bool); !ok {
		return fmt.Errorf("failed to delete workflow. got:\n%s", string(resp.Body()))
	}
	return nil
}
