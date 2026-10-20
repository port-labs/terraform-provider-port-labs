package entity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func newRefreshTestEntity(props map[string]interface{}, relations map[string]interface{}) *cli.Entity {
	return &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "test-entity",
		Title:      "test-title",
		Blueprint:  "test-blueprint",
		Properties: props,
		Relations:  relations,
	}
}

func refreshTestBlueprint() *cli.Blueprint {
	return &cli.Blueprint{
		Identifier: "test-blueprint",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"configured_prop": {Type: "string"},
				"drifted_prop":    {Type: "string"},
				"null_string":     {Type: "string"},
				"null_number":     {Type: "number"},
				"null_boolean":    {Type: "boolean"},
				"null_object":     {Type: "object"},
				"null_array":      {Type: "array"},
				"obj":             {Type: "object"},
				"num":             {Type: "number"},
				"tags":            {Type: "array"},
			},
		},
	}
}

func newRefreshTestResource() *EntityResource {
	return &EntityResource{portClient: &cli.PortClient{}}
}

// Null values the API returns for properties/relations that are not in state are dropped, so unset properties do not
// produce a perpetual diff.
func TestRefreshEntityState_DropsNullPropertiesNotInState(t *testing.T) {
	state := &EntityModel{
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{"configured_prop": types.StringValue("value")},
		},
	}
	apiEntity := newRefreshTestEntity(map[string]interface{}{
		"configured_prop": "value",
		"null_string":     nil,
		"null_number":     nil,
		"null_boolean":    nil,
		"null_object":     nil,
		"null_array":      nil,
	}, map[string]interface{}{"null_rel": nil})

	require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

	require.NotNil(t, state.Properties)
	assert.Equal(t, map[string]types.String{"configured_prop": types.StringValue("value")}, state.Properties.StringProps)
	assert.Nil(t, state.Properties.NumberProps)
	assert.Nil(t, state.Properties.BooleanProps)
	assert.Nil(t, state.Properties.ObjectProps)
	assert.Nil(t, state.Properties.ArrayProps)
	assert.Nil(t, state.Relations, "a relations block that only had null relations must stay null")
}

// When state has no properties block and the API only returns nulls, properties must stay null (not an empty object).
func TestRefreshEntityState_OnlyNullPropertiesKeepsPropertiesNull(t *testing.T) {
	state := &EntityModel{}
	apiEntity := newRefreshTestEntity(map[string]interface{}{"null_string": nil, "null_array": nil}, nil)

	require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

	assert.Nil(t, state.Properties)
}

// Non-null values the API returns for properties that are not in state are kept, so out-of-band changes show as drift.
func TestRefreshEntityState_KeepsNonNullPropertiesNotInState(t *testing.T) {
	target := "target"
	state := &EntityModel{
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{"configured_prop": types.StringValue("value")},
		},
		Relations: &RelationModel{SingleRelation: map[string]*string{"configured_rel": &target}},
	}
	apiEntity := newRefreshTestEntity(map[string]interface{}{
		"configured_prop": "value",
		"drifted_prop":    "set-elsewhere",
		"num":             float64(3),
	}, map[string]interface{}{"configured_rel": "target", "drifted_rel": "other"})

	require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

	require.NotNil(t, state.Properties)
	assert.Equal(t, types.StringValue("set-elsewhere"), state.Properties.StringProps["drifted_prop"])
	assert.Equal(t, types.Float64Value(3), state.Properties.NumberProps["num"])
	require.NotNil(t, state.Relations)
	require.Contains(t, state.Relations.SingleRelation, "drifted_rel")
	assert.Equal(t, "other", *state.Relations.SingleRelation["drifted_rel"])
}

// Import starts from an empty state, so every non-null property and relation from the API must be kept.
func TestRefreshEntityState_ImportKeepsNonNullProperties(t *testing.T) {
	state := &EntityModel{
		Identifier: types.StringValue("test-entity"),
		Blueprint:  types.StringValue("test-blueprint"),
	}
	apiEntity := newRefreshTestEntity(map[string]interface{}{
		"configured_prop": "My String Value",
		"obj":             map[string]interface{}{"a": float64(1)},
		"tags":            []interface{}{"x"},
		"null_string":     nil,
	}, map[string]interface{}{"many_rel": []interface{}{"a", "b"}, "null_rel": nil})

	require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

	require.NotNil(t, state.Properties)
	assert.Equal(t, map[string]types.String{"configured_prop": types.StringValue("My String Value")}, state.Properties.StringProps)
	assert.Equal(t, types.StringValue(`{"a":1}`), state.Properties.ObjectProps["obj"])
	require.NotNil(t, state.Properties.ArrayProps)
	assert.Contains(t, state.Properties.ArrayProps.StringItems.Elements(), "tags")
	require.NotNil(t, state.Relations)
	assert.Equal(t, []string{"a", "b"}, state.Relations.ManyRelations["many_rel"])
	assert.NotContains(t, state.Relations.SingleRelation, "null_rel")
}

// A property the user explicitly set to null stays in state as null, matching the config.
func TestRefreshEntityState_KeepsExplicitNullsInState(t *testing.T) {
	state := &EntityModel{
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{"null_string": types.StringNull()},
			ObjectProps: map[string]types.String{"null_object": types.StringNull()},
		},
	}
	apiEntity := newRefreshTestEntity(map[string]interface{}{"null_string": nil, "null_object": nil}, nil)

	require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

	require.NotNil(t, state.Properties)
	assert.Equal(t, map[string]types.String{"null_string": types.StringNull()}, state.Properties.StringProps)
	assert.Equal(t, map[string]types.String{"null_object": types.StringNull()}, state.Properties.ObjectProps)
}

// An object property the API returns as {} is equivalent to null unless state has {}; and the reverse.
func TestRefreshEntityState_EmptyObjectVsNull(t *testing.T) {
	tests := []struct {
		name      string
		oldObject map[string]types.String
		apiValue  interface{}
		wantProps map[string]types.String
	}{
		{"api {} and property unset", nil, map[string]interface{}{}, nil},
		{"api {} and config null", map[string]types.String{"obj": types.StringNull()}, map[string]interface{}{}, map[string]types.String{"obj": types.StringNull()}},
		{"api {} and config {}", map[string]types.String{"obj": types.StringValue("{}")}, map[string]interface{}{}, map[string]types.String{"obj": types.StringValue("{}")}},
		{"api null and config {}", map[string]types.String{"obj": types.StringValue("{ }")}, nil, map[string]types.String{"obj": types.StringValue("{ }")}},
		{"api null and config non-empty object", map[string]types.String{"obj": types.StringValue(`{"a":1}`)}, nil, map[string]types.String{"obj": types.StringNull()}},
		{"api non-empty object and config null", map[string]types.String{"obj": types.StringNull()}, map[string]interface{}{"a": float64(1)}, map[string]types.String{"obj": types.StringValue(`{"a":1}`)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &EntityModel{
				Properties: &EntityPropertiesModel{
					StringProps: map[string]types.String{"configured_prop": types.StringValue("value")},
					ObjectProps: tt.oldObject,
				},
			}
			apiEntity := newRefreshTestEntity(map[string]interface{}{"configured_prop": "value", "obj": tt.apiValue}, nil)

			require.NoError(t, newRefreshTestResource().refreshEntityState(context.Background(), state, apiEntity, refreshTestBlueprint()))

			require.NotNil(t, state.Properties)
			assert.Equal(t, tt.wantProps, state.Properties.ObjectProps)
		})
	}
}
