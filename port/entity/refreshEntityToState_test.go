package entity

import (
	"context"
	"testing"
	"time"

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
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// TestRefreshEntityStateFiltersUnconfiguredProperties verifies that properties
// returned by the API but not present in the user's config are excluded from state,
// preventing perpetual diffs when the API returns null values or defaults.
func TestRefreshEntityStateFiltersUnconfiguredProperties(t *testing.T) {
	state := &EntityModel{
		Identifier: types.StringValue("test-entity"),
		Blueprint:  types.StringValue("test-blueprint"),
		Title:      types.StringValue("test-title"),
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{
				"configured_prop": types.StringValue("value"),
			},
		},
	}

	apiEntity := &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "test-entity",
		Title:      "test-title",
		Blueprint:  "test-blueprint",
		Properties: map[string]interface{}{
			"configured_prop":   "new-value",
			"unconfigured_prop": "should-be-ignored",
			"unconfigured_null": nil,
		},
	}

	blueprint := &cli.Blueprint{
		Identifier: "test-blueprint",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"configured_prop":   {Type: "string"},
				"unconfigured_prop": {Type: "string"},
				"unconfigured_null": {Type: "string"},
			},
		},
	}

	resource := &EntityResource{}
	err := resource.refreshEntityState(context.Background(), state, apiEntity, blueprint)
	if err != nil {
		t.Fatalf("refreshEntityState returned error: %v", err)
	}

	if state.Properties == nil {
		t.Fatal("expected properties to not be nil")
	}

	if state.Properties.StringProps == nil {
		t.Fatal("expected string_props to not be nil")
	}

	if len(state.Properties.StringProps) != 1 {
		t.Fatalf("expected exactly 1 property in state, got %d: %#v", len(state.Properties.StringProps), state.Properties.StringProps)
	}

	if v, ok := state.Properties.StringProps["configured_prop"]; !ok || v.ValueString() != "new-value" {
		t.Fatalf("expected configured_prop to be 'new-value', got: %v", v)
	}

	if _, ok := state.Properties.StringProps["unconfigured_prop"]; ok {
		t.Fatal("unconfigured_prop should not be in state")
	}

	if _, ok := state.Properties.StringProps["unconfigured_null"]; ok {
		t.Fatal("unconfigured_null should not be in state")
	}
}

// TestRefreshEntityStateFiltersUnconfiguredRelations verifies that relations
// returned by the API but not present in the user's config are excluded from state.
func TestRefreshEntityStateFiltersUnconfiguredRelations(t *testing.T) {
	configuredRelation := "configured-target"
	state := &EntityModel{
		Identifier: types.StringValue("test-entity"),
		Blueprint:  types.StringValue("test-blueprint"),
		Title:      types.StringValue("test-title"),
		Relations: &RelationModel{
			SingleRelation: map[string]*string{
				"configured_rel": &configuredRelation,
			},
		},
	}

	apiEntity := &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "test-entity",
		Title:      "test-title",
		Blueprint:  "test-blueprint",
		Relations: map[string]interface{}{
			"configured_rel":   "new-target",
			"unconfigured_rel": nil,
		},
	}

	blueprint := &cli.Blueprint{
		Identifier: "test-blueprint",
	}

	resource := &EntityResource{}
	err := resource.refreshEntityState(context.Background(), state, apiEntity, blueprint)
	if err != nil {
		t.Fatalf("refreshEntityState returned error: %v", err)
	}

	if state.Relations == nil {
		t.Fatal("expected relations to not be nil")
	}

	if state.Relations.SingleRelation == nil {
		t.Fatal("expected single_relations to not be nil")
	}

	if len(state.Relations.SingleRelation) != 1 {
		t.Fatalf("expected exactly 1 relation in state, got %d: %#v", len(state.Relations.SingleRelation), state.Relations.SingleRelation)
	}

	if v, ok := state.Relations.SingleRelation["configured_rel"]; !ok || v == nil || *v != "new-target" {
		t.Fatalf("expected configured_rel to be 'new-target', got: %v", v)
	}

	if _, ok := state.Relations.SingleRelation["unconfigured_rel"]; ok {
		t.Fatal("unconfigured_rel should not be in state")
	}
}
