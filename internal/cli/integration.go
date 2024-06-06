package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

type PortBodyForIntegration struct {
	OK          bool        `json:"ok"`
	Integration Integration `json:"integration"`
}

func (c *PortClient) GetIntegration(ctx context.Context, id string) (*Integration, error) {
	pb := &PortBodyForIntegration{}
	url := "v1/integration/{identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("identifier", id).
		SetQueryParam("byField", "installationId").
		Get(url)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to read migration, got: %s", resp.Body())
	}
	return &pb.Integration, nil
}

func (c *PortClient) UpdateIntegration(ctx context.Context, id string, integration *Integration) (*Integration, error) {
	url := "v1/integration/{identifier}"

	resp, err := c.Client.R().
		SetBody(integration).
		SetContext(ctx).
		SetPathParam("identifier", id).
		Patch(url)
	if err != nil {
		return nil, err
	}
	var pppb PortBodyForIntegration
	err = json.Unmarshal(resp.Body(), &pppb)
	if err != nil {
		return nil, err
	}
	if !pppb.OK {
		return nil, fmt.Errorf("failed to update integration, got: %s", resp.Body())
	}
	return &pppb.Integration, nil
}

func (c *PortClient) CreateIntegration(ctx context.Context, integration *Integration) (*Integration, error) {
	url := "v1/integration"

	resp, err := c.Client.R().
		SetBody(integration).
		SetContext(ctx).
		Post(url)
	if err != nil {
		return nil, err
	}

	var pppb PortBodyForIntegration
	err = json.Unmarshal(resp.Body(), &pppb)
	if err != nil {
		return nil, err
	}
	if !pppb.OK {
		return nil, fmt.Errorf("failed to create integration, got: %s", resp.Body())
	}

	return &pppb.Integration, nil
}

func (c *PortClient) DeleteIntegration(ctx context.Context, id string) (int, error) {
	url := "v1/integration/{identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetPathParam("identifier", id).
		Delete(url)
	if err != nil {
		return resp.StatusCode(), err
	}
	var pb PortBodyForIntegration
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return resp.StatusCode(), err
	}
	if !pb.OK {
		return resp.StatusCode(), fmt.Errorf("failed to delete integration, got: %s", resp.Body())
	}
	return resp.StatusCode(), nil
}
