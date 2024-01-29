package cli

import (
	"context"
	"fmt"
)

func (c *PortClient) GetMigration(ctx context.Context, id string) (*Migration, error) {
	pb := &PortBody{}
	url := "v1/migrations/{identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("identifier", id).
		Get(url)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to read migration, got: %s", resp.Body())
	}
	return &pb.Migration, nil
}
