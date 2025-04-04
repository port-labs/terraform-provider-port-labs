package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoObjectToTerraformString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "handles greater than operator",
			input: map[string]interface{}{
				"property": "sum",
				"operator": ">",
				"value":    2,
			},
			expected: `{"operator":">","property":"sum","value":2}`,
		},
		{
			name: "handles less than operator",
			input: map[string]interface{}{
				"property": "sum",
				"operator": "<",
				"value":    2,
			},
			expected: `{"operator":"<","property":"sum","value":2}`,
		},
		{
			name: "handles multiple operators in same object",
			input: map[string]interface{}{
				"conditions": []map[string]interface{}{
					{
						"property": "sum",
						"operator": ">",
						"value":    2,
					},
					{
						"property": "sum",
						"operator": "<",
						"value":    2,
					},
				},
			},
			expected: `{"conditions":[{"operator":">","property":"sum","value":2},{"operator":"<","property":"sum","value":2}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GoObjectToTerraformString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.ValueString())
		})
	}
}
