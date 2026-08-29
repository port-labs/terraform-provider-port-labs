package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestArrayPropResourceToBodyUnionFields(t *testing.T) {
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"vulnerabilities": {
				Title: types.StringValue("Vulnerabilities"),
				Union: types.BoolValue(true),
				StringItems: &StringItems{
					Enum: types.ListNull(types.StringType),
				},
			},
		},
	}

	props := make(map[string]cli.BlueprintProperty)
	var required []string
	err := arrayPropResourceToBody(context.Background(), state, props, &required)
	require.NoError(t, err)

	prop, ok := props["vulnerabilities"]
	require.True(t, ok)
	require.NotNil(t, prop.Union)
	require.True(t, *prop.Union)
}

func TestAddArrayPropertiesToStateUnionFields(t *testing.T) {
	union := true
	prop := AddArrayPropertiesToState(context.Background(), &cli.BlueprintProperty{
		Type:  "array",
		Union: &union,
		Items: map[string]any{
			"type": "string",
		},
	}, false, nil)

	require.True(t, prop.Union.ValueBool())
}

func TestArrayPropResourceToBodyUnionRequiresStringOrNumberItems(t *testing.T) {
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"configs": {
				Title: types.StringValue("Configs"),
				Union: types.BoolValue(true),
				ObjectItems: &ObjectItems{
					Format: types.StringNull(),
				},
			},
		},
	}

	props := make(map[string]cli.BlueprintProperty)
	var required []string
	err := arrayPropResourceToBody(context.Background(), state, props, &required)
	require.Error(t, err)
	require.Contains(t, err.Error(), "union is only supported for string or number items")
}
