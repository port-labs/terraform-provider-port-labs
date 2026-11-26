package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// SearchError is returned by Search when Port answers with ok=false. It keeps the HTTP status so callers can tell a
// missing blueprint (404) apart from auth or server failures.
type SearchError struct {
	StatusCode int
	Body       []byte
}

func (e *SearchError) Error() string {
	return fmt.Sprintf("failed to search, got: %s", e.Body)
}

func (c *PortClient) Search(ctx context.Context, searchRequest *SearchRequestQuery) (*SearchResult, error) {
	url := "v1/entities/search"

	req := c.Client.R().
		SetContext(ctx).
		SetBody(*searchRequest.Query).
		SetHeader("Accept", "application/json")

	if searchRequest.ExcludeCalculatedProperties != nil {
		req.SetQueryParam("exclude_calculated_properties", fmt.Sprintf("%v", &searchRequest.ExcludeCalculatedProperties))
	}

	if len(searchRequest.Include) > 0 {
		req.SetQueryParam("include", strings.Join(searchRequest.Include, ","))
	}

	if len(searchRequest.Exclude) > 0 {
		req.SetQueryParam("exclude", strings.Join(searchRequest.Exclude, ","))
	}

	if searchRequest.AttachTitleToRelation != nil {
		req.SetQueryParam("attach_title_to_relation", fmt.Sprintf("%v", &searchRequest.AttachTitleToRelation))
	}

	resp, err := req.Post(url)

	if err != nil {
		return nil, err
	}
	var searchResult SearchResult
	err = json.Unmarshal(resp.Body(), &searchResult)
	if err != nil {
		return nil, err
	}
	if !searchResult.OK {
		return nil, &SearchError{StatusCode: resp.StatusCode(), Body: resp.Body()}
	}
	return &searchResult, nil
}
