package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestEntityResourceToBodyWrapsUnionArrayProperties(t *testing.T) {
	ctx := context.Background()
	union := true
	alpha := "alpha"
	beta := "beta"
	stringItems, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, map[string][]*string{
		"tags": {&alpha, &beta},
	})
	require.False(t, diags.HasError())

	state := &EntityModel{
		Title:               types.StringValue("service"),
		Blueprint:           types.StringValue("service"),
		UnionArraySourceKey: types.StringValue("terraform"),
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				StringItems: stringItems,
			},
		},
	}
	bp := &cli.Blueprint{
		Identifier: "service",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"tags": {
					Type:  "array",
					Union: &union,
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	entity, err := entityResourceToBody(ctx, state, bp)
	require.NoError(t, err)

	value, ok := entity.Properties["tags"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, []interface{}{"alpha", "beta"}, value["terraform"])
}

func TestEntityResourceToBodyLeavesPlainArrayPropertiesUnwrapped(t *testing.T) {
	ctx := context.Background()
	alpha := "alpha"
	stringItems, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, map[string][]*string{
		"tags": {&alpha},
	})
	require.False(t, diags.HasError())

	state := &EntityModel{
		Title:     types.StringValue("service"),
		Blueprint: types.StringValue("service"),
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				StringItems: stringItems,
			},
		},
	}
	bp := &cli.Blueprint{
		Identifier: "service",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"tags": {
					Type:  "array",
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	entity, err := entityResourceToBody(ctx, state, bp)
	require.NoError(t, err)
	require.Equal(t, []interface{}{"alpha"}, entity.Properties["tags"])
}
