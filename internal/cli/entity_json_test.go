package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntityUnmarshalJSONRelationSources(t *testing.T) {
	t.Parallel()

	payload := `{
		"identifier": "entity-1",
		"title": "Entity",
		"blueprint": "service",
		"properties": {},
		"relationSources": {
			"dependencies": {
				"integration-a": ["dep-a"],
				"terraform": ["dep-b"]
			}
		}
	}`

	var entity Entity
	err := json.Unmarshal([]byte(payload), &entity)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"dependencies": map[string]any{
			"integration-a": []any{"dep-a"},
			"terraform":     []any{"dep-b"},
		},
	}, entity.RelationSources)
}

func TestEntityTeamIdentifiers(t *testing.T) {
	t.Parallel()

	require.Nil(t, EntityTeamIdentifiers(nil))
	require.Equal(t, []string{"a", "b"}, EntityTeamIdentifiers([]any{"a", "b"}))
	require.True(t, EntityTeamIsUnionSlice(map[string]any{"terraform": []any{"team-a"}}))
}
