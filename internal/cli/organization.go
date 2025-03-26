package cli

import (
	"context"
	"fmt"
)

const orgUrl = "/v1/organization"

func (c *PortClient) ReadOrganization(ctx context.Context) (*Organization, int, error) {
	pb := &PortBody{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		Get(orgUrl)
	if err != nil {
		return nil, 0, err
	} else if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read organization, got: %s", resp.Body())
	}
	if pb.Organization == nil {
		return nil, 0, fmt.Errorf("port-api returned an invalid response: Organization is nil")
	}
	return pb.Organization, resp.StatusCode(), nil
}
