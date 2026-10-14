package blueprint

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelationsResourceToBodyIncludesUnion(t *testing.T) {
	relations := RelationsResourceToBody(map[string]RelationModel{
		"members": {
			Target: types.StringValue("user"),
			Many:   types.BoolValue(true),
			Union:  types.BoolValue(true),
		},
	})

	relation, ok := relations["members"]
	require.True(t, ok)
	require.NotNil(t, relation.Union)
	assert.True(t, *relation.Union)
	assert.Equal(t, "user", *relation.Target)
	assert.True(t, *relation.Many)
}

func TestRelationsResourceToBodyOmitsUnionWhenUnset(t *testing.T) {
	relations := RelationsResourceToBody(map[string]RelationModel{
		"owner": {
			Target: types.StringValue("team"),
			Many:   types.BoolValue(false),
			Union:  types.BoolNull(),
		},
	})

	relation, ok := relations["owner"]
	require.True(t, ok)
	assert.Nil(t, relation.Union)
}

func TestRelationsResourceToBodyOmitsUnionWhenFalse(t *testing.T) {
	relations := RelationsResourceToBody(map[string]RelationModel{
		"members": {
			Target: types.StringValue("user"),
			Many:   types.BoolValue(true),
			Union:  types.BoolValue(false),
		},
	})

	relation, ok := relations["members"]
	require.True(t, ok)
	assert.Nil(t, relation.Union)
}
