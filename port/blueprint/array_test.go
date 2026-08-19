package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestArrayPropResourceToBodyWithUnion(t *testing.T) {
	ctx := context.Background()
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"tags": {
				Title: types.StringValue("Tags"),
				Union: types.BoolValue(true),
				StringItems: &StringItems{
					Format: types.StringNull(),
					Enum:   types.ListNull(types.StringType),
				},
			},
		},
	}

	props := map[string]cli.BlueprintProperty{}
	required := []string{}

	err := arrayPropResourceToBody(ctx, state, props, &required)
	require.NoError(t, err)

	prop, ok := props["tags"]
	require.True(t, ok)
	require.NotNil(t, prop.Union)
	require.True(t, *prop.Union)
	require.Nil(t, prop.IncludeDuplicates)
}

func TestArrayPropResourceToBodyWithUnionAndIncludeDuplicates(t *testing.T) {
	ctx := context.Background()
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"tags": {
				Title:             types.StringValue("Tags"),
				Union:             types.BoolValue(true),
				IncludeDuplicates: types.BoolValue(true),
				StringItems: &StringItems{
					Format: types.StringNull(),
					Enum:   types.ListNull(types.StringType),
				},
			},
		},
	}

	props := map[string]cli.BlueprintProperty{}
	required := []string{}

	err := arrayPropResourceToBody(ctx, state, props, &required)
	require.NoError(t, err)

	prop, ok := props["tags"]
	require.True(t, ok)
	require.NotNil(t, prop.Union)
	require.True(t, *prop.Union)
	require.NotNil(t, prop.IncludeDuplicates)
	require.True(t, *prop.IncludeDuplicates)
}

func TestArrayPropResourceToBodyWithoutUnion(t *testing.T) {
	ctx := context.Background()
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"tags": {
				Title: types.StringValue("Tags"),
				Union: types.BoolNull(),
				StringItems: &StringItems{
					Format: types.StringNull(),
					Enum:   types.ListNull(types.StringType),
				},
			},
		},
	}

	props := map[string]cli.BlueprintProperty{}
	required := []string{}

	err := arrayPropResourceToBody(ctx, state, props, &required)
	require.NoError(t, err)

	prop, ok := props["tags"]
	require.True(t, ok)
	require.Nil(t, prop.Union)
	require.Nil(t, prop.IncludeDuplicates)
}

func TestAddArrayPropertiesToStateWithUnion(t *testing.T) {
	ctx := context.Background()
	union := true
	includeDuplicates := false

	arrayProp := AddArrayPropertiesToState(ctx, &cli.BlueprintProperty{
		Type: "array",
		Items: map[string]any{
			"type": "string",
		},
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
	}, false)

	require.True(t, arrayProp.Union.ValueBool())
	require.False(t, arrayProp.IncludeDuplicates.ValueBool())
}
