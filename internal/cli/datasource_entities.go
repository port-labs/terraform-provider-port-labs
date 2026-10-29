package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

type DatasourceEntitiesRequest struct {
	DatasourcePrefix string  `json:"datasource_prefix"`
	DatasourceSuffix string  `json:"datasource_suffix"`
	From             *string `json:"from,omitempty"`
	Limit            *int    `json:"limit,omitempty"`
	Before           *string `json:"before,omitempty"`
}

type DatasourceEntityIdentifier struct {
	Identifier string `json:"identifier"`
	Blueprint  string `json:"blueprint"`
}

type DatasourceEntitiesResponse struct {
	OK       bool                         `json:"ok"`
	Entities []DatasourceEntityIdentifier `json:"entities"`
	Next     *string                      `json:"next,omitempty"`
}

func (c *PortClient) GetDatasourceEntities(ctx context.Context, request *DatasourceEntitiesRequest) ([]DatasourceEntityIdentifier, error) {
	url := "v1/blueprints/entities/datasource-entities"

	var aggregatedEntities []DatasourceEntityIdentifier
	nextFrom := request.From

	for {
		body := DatasourceEntitiesRequest{
			DatasourcePrefix: request.DatasourcePrefix,
			DatasourceSuffix: request.DatasourceSuffix,
			Limit:            request.Limit,
			Before:           request.Before,
			From:             nextFrom,
		}

		resp, err := c.Client.R().
			SetContext(ctx).
			SetHeader("Accept", "application/json").
			SetBody(body).
			Post(url)
		if err != nil {
			return nil, err
		}

		var response DatasourceEntitiesResponse
		if err := json.Unmarshal(resp.Body(), &response); err != nil {
			return nil, err
		}
		if !response.OK {
			return nil, fmt.Errorf("failed to get datasource entities, got: %s", resp.Body())
		}

		aggregatedEntities = append(aggregatedEntities, response.Entities...)

		if response.Next == nil || *response.Next == "" {
			break
		}
		nextFrom = response.Next
	}

	return aggregatedEntities, nil
}
