package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadEntity(ctx context.Context, id string, blueprint string) (*Entity, int, error) {
	url := "v1/blueprints/{blueprint}/entities/{identifier}"
	resp, err := c.Client.R().
		SetHeader("Accept", "application/json").
		// we don't want to include those properties as they are calculated by the backend
		// and not part of the state, pulling them would cause a diff
		SetQueryParam("exclude_calculated_properties", "true").
		SetPathParam(("blueprint"), blueprint).
		SetPathParam("identifier", id).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read entity, got: %s", resp.Body())
	}
	return &pb.Entity, resp.StatusCode(), nil
}

func (c *PortClient) CreateEntity(ctx context.Context, e *Entity, runID string) (*Entity, error) {
	url := "v1/blueprints/{blueprint}/entities"
	pb := &PortBody{}
	resp, err := c.Client.R().
		SetBody(e).
		SetPathParam(("blueprint"), e.Blueprint).
		SetQueryParam("upsert", "true").
		SetQueryParam("run_id", runID).
		SetResult(&pb).
		Post(url)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create entity, got: %s", resp.Body())
	}
	return &pb.Entity, nil
}

func (c *PortClient) UpdateEntity(ctx context.Context, id string, blueprint string, e *Entity, runID string) (*Entity, error) {
	url := "v1/blueprints/{blueprint}/entities/{identifier}"
	pb := &PortBody{}
	resp, err := c.Client.R().
		SetBody(e).
		SetPathParam(("blueprint"), e.Blueprint).
		SetPathParam("identifier", id).
		SetQueryParam("run_id", runID).
		SetResult(&pb).
		Put(url)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update entity, got: %s", resp.Body())
	}
	return &pb.Entity, nil
}

func (c *PortClient) DeleteEntity(ctx context.Context, id string, blueprint string) error {
	url := "v1/blueprints/{blueprint}/entities/{identifier}"
	pb := &PortBody{}
	resp, err := c.Client.R().
		SetHeader("Accept", "application/json").
		SetPathParam("blueprint", blueprint).
		SetPathParam("identifier", id).
		SetResult(pb).
		Delete(url)
	if err != nil {
		return err
	}
	if !pb.OK {
		return fmt.Errorf("failed to delete entity, got: %s", resp.Body())
	}
	return nil
}
