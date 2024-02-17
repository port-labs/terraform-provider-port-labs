package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) GetPage(ctx context.Context, pageId string) (*Page, int, error) {
	pb := &PortBody{}
	url := "v1/pages/{page_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("page_identifier", pageId).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to get page, got: %s", resp.Body())
	}
	return &pb.Page, resp.StatusCode(), nil

}

func (c *PortClient) CreatePage(ctx context.Context, page *Page) (*Page, error) {
	url := "v1/pages"
	resp, err := c.Client.R().
		SetBody(page).
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
		if resp.IsSuccess() {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to create page, got: %s", resp.Body())
	}
	return &pb.Page, nil
}

func (c *PortClient) UpdatePage(ctx context.Context, pageId string, page *Page) (*Page, error) {
	url := "v1/pages/{page_identifier}"
	resp, err := c.Client.R().
		SetBody(page).
		SetContext(ctx).
		SetPathParam("page_identifier", pageId).
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
		return nil, fmt.Errorf("failed to update page, got: %s", resp.Body())
	}
	return &pb.Page, nil
}

func (c *PortClient) DeletePage(ctx context.Context, pageId string) (int, error) {
	url := "v1/pages/{page_identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetPathParam("page_identifier", pageId).
		Delete(url)
	if err != nil {
		return resp.StatusCode(), err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return resp.StatusCode(), err
	}
	if !pb.OK {
		return resp.StatusCode(), fmt.Errorf("failed to delete page, got: %s", resp.Body())
	}
	return resp.StatusCode(), nil
}
