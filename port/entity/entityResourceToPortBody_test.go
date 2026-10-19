package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestWriteUnionManyRelationsToBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	items, diags := types.ListValueFrom(ctx, types.StringType, []string{"dep-a", "dep-b"})
	require.False(t, diags.HasError())

	state := &EntityModel{
		Relations: &RelationModel{
			UnionManyRelations: map[string]UnionManyRelationSliceModel{
				"dependencies": {
					SourceKey: types.StringValue("terraform"),
					Items:     items,
				},
			},
		},
	}

	body, err := writeRelationsToBody(ctx, state.Relations)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"dependencies": map[string]interface{}{
			"terraform": []interface{}{"dep-a", "dep-b"},
		},
	}, body)
}

func TestWriteTeamToBodyUnionSlice(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	teams, diags := types.SetValueFrom(ctx, types.StringType, []string{"team-a"})
	require.False(t, diags.HasError())

	state := &EntityModel{
		UnionTeamSlice: &UnionTeamSliceModel{
			SourceKey: types.StringValue("terraform"),
			Teams:     teams,
		},
	}

	team, err := writeTeamToBody(ctx, state, &cli.Blueprint{})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"terraform": []string{"team-a"},
	}, team)
}

func TestRefreshRelationsEntityStateSkipsUnionRelations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	state := &EntityModel{}
	blueprint := &cli.Blueprint{
		Relations: map[string]cli.Relation{
			"dependencies": {
				Union: boolPointer(true),
			},
		},
	}

	entity := &cli.Entity{
		Relations: map[string]any{
			"dependencies": []any{"dep-a"},
			"owner":        "team-a",
		},
	}

	refreshRelationsEntityState(ctx, state, entity, blueprint)
	require.Nil(t, state.Relations.ManyRelations)
	require.Equal(t, stringPointer("team-a"), state.Relations.SingleRelation["owner"])
}

func boolPointer(value bool) *bool {
	return &value
}

func stringPointer(value string) *string {
	return &value
}
