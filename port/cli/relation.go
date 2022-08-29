package cli

import (
	"context"
	"fmt"
)

func (c *PortClient) CreateRelation(ctx context.Context, bpID string, r *Relation) (string, error) {
	url := "v1/blueprints/{identifier}/relations"
	result := map[string]interface{}{}
	resp, err := c.Client.R().
		SetBody(r).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("identifier", bpID).
		Post(url)
	if err != nil {
		return "", err
	}
	if !result["ok"].(bool) {
		return "", fmt.Errorf("failed to create relation, got: %s", resp.Body())
	}
	return result["identifier"].(string), nil
}

func (c *PortClient) ReadRelations(ctx context.Context, blueprintID string) ([]*Relation, error) {
	url := "v1/relations"
	result := map[string]interface{}{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(url)
	if err != nil {
		return nil, err
	}
	if !result["ok"].(bool) {
		return nil, fmt.Errorf("failed to create relation, got: %s", resp.Body())
	}
	allRelations := result["relations"].([]interface{})
	bpRelations := make([]*Relation, 0)
	for _, relation := range allRelations {
		r := relation.(map[string]interface{})
		if r["source"] != blueprintID {
			continue
		}
		bpRelations = append(bpRelations, &Relation{
			Target:     r["target"].(string),
			Required:   r["required"].(bool),
			Many:       r["many"].(bool),
			Title:      r["title"].(string),
			Identifier: r["identifier"].(string),
		})
	}
	return bpRelations, nil
}
