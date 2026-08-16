package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadScorecardGroup(ctx context.Context, identifier string) (*ScorecardGroup, int, error) {
	pb := &PortBody{}
	url := "v1/scorecard-groups/{scorecard_group_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("scorecard_group_identifier", identifier).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read scorecard group, got: %s", resp.Body())
	}
	return &pb.ScorecardGroup, resp.StatusCode(), nil
}

func (c *PortClient) CreateScorecardGroup(ctx context.Context, scorecardGroup *ScorecardGroup) (*ScorecardGroup, error) {
	url := "v1/scorecard-groups"
	resp, err := c.Client.R().
		SetBody(scorecardGroup).
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
		return nil, fmt.Errorf("failed to create scorecard group, got: %s", resp.Body())
	}
	return &pb.ScorecardGroup, nil
}
