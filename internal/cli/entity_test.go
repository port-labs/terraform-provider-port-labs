package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntityUnmarshalJSONPropertySources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		payload  string
		expected map[string]any
	}{
		{
			name: "propertySources",
			payload: `{
				"identifier": "entity-1",
				"title": "Entity",
				"blueprint": "service",
				"properties": {},
				"propertySources": {
					"tags": {
						"integration-a": ["tag-a"],
						"manual": ["tag-b"]
					}
				}
			}`,
			expected: map[string]any{
				"tags": map[string]any{
					"integration-a": []any{"tag-a"},
					"manual":        []any{"tag-b"},
				},
			},
		},
		{
			name: "_propertySources",
			payload: `{
				"identifier": "entity-1",
				"title": "Entity",
				"blueprint": "service",
				"properties": {},
				"_propertySources": {
					"owners": {
						"integration-b": [1, 2]
					}
				}
			}`,
			expected: map[string]any{
				"owners": map[string]any{
					"integration-b": []any{float64(1), float64(2)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var entity Entity
			err := json.Unmarshal([]byte(tt.payload), &entity)
			require.NoError(t, err)
			require.Equal(t, tt.expected, entity.PropertySources)
		})
	}
}
