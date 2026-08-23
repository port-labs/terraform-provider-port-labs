package entity

import (
	"context"
	"encoding/json"
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

func ptrTime(t time.Time) *time.Time {
	return &t
}
