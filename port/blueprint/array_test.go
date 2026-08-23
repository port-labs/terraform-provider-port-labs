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
				Title:             types.StringValue("Vulnerabilities"),
				Union:             types.BoolValue(true),
				IncludeDuplicates: types.BoolValue(false),
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
	require.NotNil(t, prop.IncludeDuplicates)
	require.False(t, *prop.IncludeDuplicates)
}

func TestAddArrayPropertiesToStateUnionFields(t *testing.T) {
	union := true
	includeDuplicates := true
	prop := AddArrayPropertiesToState(context.Background(), &cli.BlueprintProperty{
		Type:              "array",
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
		Items: map[string]any{
			"type": "string",
		},
	}, false, nil)

	require.True(t, prop.Union.ValueBool())
	require.True(t, prop.IncludeDuplicates.ValueBool())
}
