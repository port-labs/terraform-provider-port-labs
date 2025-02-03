package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

var sidebarRoute = "v1/sidebars"

func (c *PortClient) GetFolder(ctx context.Context, sidebarId string, folderId string) (*Folder, int, error) {
	pb := &SidebarDTO{}

	url := fmt.Sprintf("%s/%s", sidebarRoute, sidebarId)

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		Get(url)

	if err != nil {
		return nil, resp.StatusCode(), err
	}

	if resp.StatusCode() != 200 {
		return nil, resp.StatusCode(), fmt.Errorf("failed to get sidebar, got: %s", resp.Body())
	}

	for _, item := range pb.Items {
		if item.SidebarType == "FOLDER" && item.Identifier == folderId {
			folder := &Folder{
				FolderIdentifier:  item.Identifier,
				SidebarIdentifier: sidebarId,
				Title:             item.Title,
				After:             item.After,
				Parent:            item.Parent,
			}
			return folder, resp.StatusCode(), nil
		}
	}

	return nil, resp.StatusCode(), fmt.Errorf("folder with identifier %s not found", folderId)
}

func (c *PortClient) CreateFolder(ctx context.Context, folder *Folder) (*Folder, error) {
	encodedSidebarId := url.QueryEscape(folder.SidebarIdentifier)
	url := fmt.Sprintf("%s/%s/folders", sidebarRoute, encodedSidebarId)

	resp, err := c.Client.R().
		SetBody(folder).
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
		return nil, fmt.Errorf("failed to create folder, got: %s", resp.Body())
	}
	// For forward compatibility, handle cases where the response body is empty when a folder is created.
	// The current API response body is { "ok": true, "identifier": "folder_identifier" },
	// but it is expected to be the folder object in the future to align with other API endpoints.
	if pb.Folder.FolderIdentifier != "" {
		return &pb.Folder, nil
	}
	return nil, nil
}

func (c *PortClient) UpdateFolder(ctx context.Context, folder *Folder) (*Folder, error) {
	encodedSidebarId := url.QueryEscape(folder.SidebarIdentifier)
	encodedFolderId := url.QueryEscape(folder.FolderIdentifier)

	url := fmt.Sprintf("%s/%s/folders/%s", sidebarRoute, encodedSidebarId, encodedFolderId)

	resp, err := c.Client.R().
		SetBody(folder).
		SetContext(ctx).
		Patch(url)
	if err != nil {
		return nil, err
	}

	var pb PortBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}

	if !pb.OK {
		return nil, fmt.Errorf("failed to update folder, got: %s", resp.Body())
	}

	return &pb.Folder, nil
}

func (c *PortClient) DeleteFolder(ctx context.Context, sidebarId string, folderId string) (int, error) {
	encodedSidebarId := url.QueryEscape(sidebarId)
	encodedFolderId := url.QueryEscape(folderId)

	url := fmt.Sprintf("%s/%s/folders/%s", sidebarRoute, encodedSidebarId, encodedFolderId)

	resp, err := c.Client.R().
		SetContext(ctx).
		Delete(url)
	if err != nil {
		return resp.StatusCode(), err
	}

	var pb PortBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return resp.StatusCode(), err
	}

	if !pb.OK {
		return resp.StatusCode(), fmt.Errorf("failed to delete folder, got: %s", resp.Body())
	}

	return resp.StatusCode(), nil
}
