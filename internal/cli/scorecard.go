package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadScorecard(ctx context.Context, blueprintID string, scorecardID string) (*Scorecard, int, error) {
	pb := &PortBody{}
	url := "v1/blueprints/{blueprint_identifier}/scorecards/{scorecard_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("blueprint_identifier", blueprintID).
		SetPathParam("scorecard_identifier", scorecardID).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read scorecard, got: %s", resp.Body())
	}
	return &pb.Scorecard, resp.StatusCode(), nil
}

func (c *PortClient) CreateScorecard(ctx context.Context, blueprintID string, scorecard *Scorecard) (*Scorecard, error) {
	url := "v1/blueprints/{blueprint_identifier}/scorecards"
	resp, err := c.Client.R().
		SetBody(scorecard).
		SetContext(ctx).
		SetPathParam("blueprint_identifier", blueprintID).
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
		return nil, fmt.Errorf("failed to create scorecard, got: %s", resp.Body())
	}
	return &pb.Scorecard, nil
}

func (c *PortClient) UpdateScorecard(ctx context.Context, blueprintID string, scorecardId string, scorecard *Scorecard) (*Scorecard, error) {
	url := "v1/blueprints/{blueprint_identifier}/scorecards/{scorecard_identifier}"
	resp, err := c.Client.R().
		SetBody(scorecard).
		SetContext(ctx).
		SetPathParam("blueprint_identifier", blueprintID).
		SetPathParam("scorecard_identifier", scorecardId).
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
		return nil, fmt.Errorf("failed to update scorecard, got: %s", resp.Body())
	}
	return &pb.Scorecard, nil
}

func (c *PortClient) DeleteScorecard(ctx context.Context, blueprintID string, scorecardID string) error {
	url := "v1/blueprints/{blueprint_identifier}/scorecards/{scorecard_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("blueprint_identifier", blueprintID).
		SetPathParam("scorecard_identifier", scorecardID).
		Delete(url)
	if err != nil {
		return err
	}
	var pb PortBodyDelete
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return err
	}

	if !(pb.Ok) {
		return fmt.Errorf("failed to delete scorecard. got:\n%s", string(resp.Body()))
	}
	return nil
}
