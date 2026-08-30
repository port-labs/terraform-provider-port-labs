package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *PortClient) ReadBlueprint(ctx context.Context, id string) (*Blueprint, int, error) {
	pb := &PortBody{}
	const url = "v1/blueprints/{identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("exclude_calculated_properties", "true").
		SetResult(pb).
		SetPathParam("identifier", id).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read blueprint, got: %s", resp.Body())
	}
	return &pb.Blueprint, resp.StatusCode(), nil
}

func (c *PortClient) ReadSystemBlueprintStructure(ctx context.Context, id string) (*Blueprint, int, error) {
	pb := &PortBody{}
	const url = "v1/blueprints/system/{identifier}/structure"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("identifier", id).
		Get(url)
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read system blueprint structure, got: %s", resp.Body())
	}
	return &pb.Blueprint, resp.StatusCode(), nil
}

func (c *PortClient) CreateBlueprint(ctx context.Context, b *Blueprint, createCatalogPage *bool) (*Blueprint, error) {
	const url = "v1/blueprints"
	request := c.Client.R().
		SetBody(b).
		SetContext(ctx)
	if createCatalogPage != nil {
		request.SetQueryParam("create_catalog_page", fmt.Sprintf("%t", *createCatalogPage))
	}
	resp, err := request.Post(url)
	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create blueprint, got: %s", resp.Body())
	}
	return &pb.Blueprint, nil
}

func (c *PortClient) UpdateBlueprint(ctx context.Context, b *Blueprint, id string) (*Blueprint, error) {
	defer c.LockBlueprint(id)()
	return c.updateBlueprint(ctx, b, id)
}

// updateBlueprint performs the write without taking the blueprint lock. Callers that already
// hold it must use this variant, as the lock is not reentrant.
func (c *PortClient) updateBlueprint(ctx context.Context, b *Blueprint, id string) (*Blueprint, error) {
	const url = "v1/blueprints/{identifier}"
	resp, err := c.Client.R().
		SetBody(b).
		SetContext(ctx).
		SetPathParam("identifier", id).
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
		return nil, fmt.Errorf("failed to update blueprint, got: %s", resp.Body())
	}
	return &pb.Blueprint, nil
}

// PatchBlueprintRelation creates or updates a single relation on a blueprint.
//
// Port's PATCH deep-merges the request body into the existing blueprint, so sending only
// the one relation leaves every other relation - and the rest of the blueprint - untouched.
// This is a single atomic call, so no blueprint lock is needed.
func (c *PortClient) PatchBlueprintRelation(ctx context.Context, blueprintID string, relationID string, relation *Relation) (*Blueprint, error) {
	const url = "v1/blueprints/{identifier}"
	body := map[string]any{
		"relations": map[string]*Relation{relationID: relation},
	}
	resp, err := c.Client.R().
		SetBody(body).
		SetContext(ctx).
		SetPathParam("identifier", blueprintID).
		Patch(url)
	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to write relation %q on blueprint %q, got: %s", relationID, blueprintID, resp.Body())
	}
	return &pb.Blueprint, nil
}

// PutBlueprintRelation replaces a single relation on a blueprint, dropping any field that is
// no longer set.
//
// PATCH cannot be used to clear an optional field: it deep-merges, and the API rejects null
// for `title` and `description`. Removing such a field therefore requires rewriting the whole
// blueprint, which makes this a read-modify-write and requires the lock.
func (c *PortClient) PutBlueprintRelation(ctx context.Context, blueprintID string, relationID string, relation *Relation) (*Blueprint, error) {
	defer c.LockBlueprint(blueprintID)()

	b, _, err := c.ReadBlueprint(ctx, blueprintID)
	if err != nil {
		return nil, err
	}

	if b.Relations == nil {
		b.Relations = map[string]Relation{}
	}
	b.Relations[relationID] = *relation

	return c.updateBlueprint(ctx, b, blueprintID)
}

// DeleteBlueprintRelation removes a single relation from a blueprint.
//
// Port exposes no endpoint for deleting an individual relation, and PATCH is additive, so the
// only way to remove one is to write the blueprint back without it. Read-modify-write, so the
// lock is required.
//
// Returns statusCode 404 if the blueprint itself no longer exists.
func (c *PortClient) DeleteBlueprintRelation(ctx context.Context, blueprintID string, relationID string) (int, error) {
	defer c.LockBlueprint(blueprintID)()

	b, statusCode, err := c.ReadBlueprint(ctx, blueprintID)
	if err != nil {
		return statusCode, err
	}

	if _, ok := b.Relations[relationID]; !ok {
		// already gone, nothing to write
		return statusCode, nil
	}

	delete(b.Relations, relationID)

	_, err = c.updateBlueprint(ctx, b, blueprintID)
	return statusCode, err
}

func (c *PortClient) DeleteBlueprint(ctx context.Context, id string) error {
	const url = "v1/blueprints/{identifier}"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("identifier", id).
		Delete(url)
	if err != nil {
		return err
	}
	responseBody := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &responseBody)
	if err != nil {
		return err
	}
	if !(responseBody["ok"].(bool)) {
		return fmt.Errorf("failed to delete blueprint. got:\n%s", string(resp.Body()))
	}
	return nil
}

func (c *PortClient) DeleteBlueprintWithAllEntities(ctx context.Context, id string) (*string, error) {
	const url = "v1/blueprints/{identifier}/all-entities?delete_blueprint=true"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("identifier", id).
		Delete(url)
	if err != nil {
		return nil, err
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to trigger blueprint deletion with all entities, got: %s", resp.Body())
	}

	return &pb.MigrationId, nil
}
