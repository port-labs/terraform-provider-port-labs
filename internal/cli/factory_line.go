package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadFactoryLine(ctx context.Context, identifier string) (*FactoryLineResponse, int, error) {
	pb := &FactoryLineBody{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("line_identifier", identifier).
		Get("v1/factory/lines/{line_identifier}")
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if resp.StatusCode() == 404 {
		return nil, resp.StatusCode(), fmt.Errorf("factory line not found")
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read factory line, got: %s", resp.Body())
	}
	return &pb.Line, resp.StatusCode(), nil
}

func (c *PortClient) CreateFactoryLine(ctx context.Context, line *FactoryLine) (*FactoryLineResponse, error) {
	resp, err := c.Client.R().
		SetBody(line).
		SetContext(ctx).
		Post("v1/factory/lines")
	if err != nil {
		return nil, err
	}
	var pb FactoryLineBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create factory line, got (HTTP %d): %s", resp.StatusCode(), resp.Body())
	}
	return &pb.Line, nil
}

func (c *PortClient) UpdateFactoryLine(ctx context.Context, identifier string, change *FactoryLineChange) (*FactoryLineResponse, error) {
	resp, err := c.Client.R().
		SetBody(change).
		SetContext(ctx).
		SetPathParam("line_identifier", identifier).
		Put("v1/factory/lines/{line_identifier}")
	if err != nil {
		return nil, err
	}
	var pb FactoryLineBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update factory line, got (HTTP %d): %s", resp.StatusCode(), resp.Body())
	}
	return &pb.Line, nil
}

func (c *PortClient) DeleteFactoryLine(ctx context.Context, identifier string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("line_identifier", identifier).
		Delete("v1/factory/lines/{line_identifier}")
	if err != nil {
		return err
	}
	var pb PortBodyDelete
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return err
	}
	if !pb.Ok {
		return fmt.Errorf("failed to delete factory line, got: %s", resp.Body())
	}
	return nil
}
