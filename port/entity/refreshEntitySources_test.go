package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestRefreshRelationSourcesEntityState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	state := &EntityModel{}
	entity := &cli.Entity{
		RelationSources: map[string]any{
			"dependencies": map[string]any{
				"terraform": []any{"dep-a"},
			},
		},
	}

	err := refreshRelationSourcesEntityState(ctx, state, entity)
	require.NoError(t, err)

	encoded := state.RelationSources.Elements()["dependencies"].(types.String).ValueString()
	require.JSONEq(t, `{"terraform":["dep-a"]}`, encoded)
}
