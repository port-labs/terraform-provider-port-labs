package entity

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshEntityStateClearsStaleCollectionsOnEmptyAPIResponse(t *testing.T) {
	staleRelation := "old-relation"
	state := &EntityModel{
		Identifier: types.StringValue("old"),
		Blueprint:  types.StringValue("old-blueprint"),
		Title:      types.StringValue("old-title"),
		Teams:      []types.String{types.StringValue("old-team")},
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{
				"name": types.StringValue("old-value"),
			},
		},
		Relations: &RelationModel{
			SingleRelation: map[string]*string{
				"owner": &staleRelation,
			},
		},
	}

	apiEntity := &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "new-identifier",
		Title:      "new-title",
		Blueprint:  "new-blueprint",
		Team:       nil,
		Properties: nil,
		Relations:  nil,
	}

	blueprint := &cli.Blueprint{
		Identifier: "new-blueprint",
	}

	resource := &EntityResource{}
	err := resource.refreshEntityState(context.Background(), state, apiEntity, blueprint)
	if err != nil {
		t.Fatalf("refreshEntityState returned error: %v", err)
	}

	if state.Teams != nil {
		t.Fatalf("expected teams to be nil when API returns empty, got: %#v", state.Teams)
	}

	if state.Properties != nil {
		t.Fatalf("expected properties to be nil when API returns empty, got: %#v", state.Properties)
	}

	if state.Relations != nil {
		t.Fatalf("expected relations to be nil when API returns empty, got: %#v", state.Relations)
	}

	if !state.PropertySources.IsNull() {
		t.Fatalf("expected property_sources to be null when API returns empty, got: %#v", state.PropertySources)
	}
}

func TestRefreshPropertySourcesEntityState(t *testing.T) {
	ctx := context.Background()
	state := &EntityModel{}

	apiEntity := &cli.Entity{
		PropertySources: map[string]any{
			"tags": map[string]any{
				"integration-a": []any{"tag-a", "tag-b"},
				"manual":        []any{"tag-c"},
			},
		},
	}

	err := refreshPropertySourcesEntityState(ctx, state, apiEntity)
	if err != nil {
		t.Fatalf("refreshPropertySourcesEntityState returned error: %v", err)
	}

	if state.PropertySources.IsNull() {
		t.Fatal("expected property_sources to be populated")
	}

	tagsValue, ok := state.PropertySources.Elements()["tags"]
	if !ok {
		t.Fatal("expected tags property source in state")
	}

	var tagsSources map[string][]string
	err = json.Unmarshal([]byte(tagsValue.(types.String).ValueString()), &tagsSources)
	if err != nil {
		t.Fatalf("failed to decode tags property sources: %v", err)
	}

	if len(tagsSources["integration-a"]) != 2 || tagsSources["integration-a"][0] != "tag-a" {
		t.Fatalf("unexpected integration-a sources: %#v", tagsSources["integration-a"])
	}
}

func TestRefreshEntityStatePreservesUnionArraySlices(t *testing.T) {
	union := true
	state := &EntityModel{
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				UnionStringSlices: map[string]UnionStringSliceModel{
					"vulnerabilities": {
						SourceKey: types.StringValue("terraform"),
						Items:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("CVE-1")}),
					},
				},
			},
		},
	}

	apiEntity := &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "service-1",
		Title:      "Service",
		Blueprint:  "service",
		Properties: map[string]any{
			"vulnerabilities": []any{"CVE-1"},
			"tags":            []any{"alpha"},
		},
		PropertySources: map[string]any{
			"vulnerabilities": map[string]any{
				"terraform": []any{"CVE-1"},
			},
		},
	}

	blueprint := &cli.Blueprint{
		Identifier: "service",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"vulnerabilities": {
					Type:  "array",
					Union: &union,
					Items: map[string]any{"type": "string"},
				},
				"tags": {
					Type:  "array",
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	resource := &EntityResource{}
	err := resource.refreshEntityState(context.Background(), state, apiEntity, blueprint)
	if err != nil {
		t.Fatalf("refreshEntityState returned error: %v", err)
	}

	if state.Properties == nil || state.Properties.ArrayProps == nil {
		t.Fatal("expected array props to be preserved")
	}

	slice, ok := state.Properties.ArrayProps.UnionStringSlices["vulnerabilities"]
	if !ok {
		t.Fatalf("expected union_string_slices to be preserved, got %#v", state.Properties.ArrayProps.UnionStringSlices)
	}

	if slice.SourceKey.ValueString() != "terraform" {
		t.Fatalf("expected terraform source key, got %q", slice.SourceKey.ValueString())
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
