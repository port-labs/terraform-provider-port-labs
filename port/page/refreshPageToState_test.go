package page

import (
	"reflect"
	"testing"
)

func TestStripEmptyArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "strips empty array at top level",
			input: map[string]interface{}{
				"type":        "widget",
				"hiddenQuery": []interface{}{},
			},
			expected: map[string]interface{}{
				"type": "widget",
			},
		},
		{
			name: "preserves non-empty array",
			input: map[string]interface{}{
				"type":  "widget",
				"items": []interface{}{"a", "b"},
			},
			expected: map[string]interface{}{
				"type":  "widget",
				"items": []interface{}{"a", "b"},
			},
		},
		{
			name: "strips nested empty arrays",
			input: map[string]interface{}{
				"groups": []interface{}{
					map[string]interface{}{
						"widgets": []interface{}{
							map[string]interface{}{
								"type":        "table",
								"hiddenQuery": []interface{}{},
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"groups": []interface{}{
					map[string]interface{}{
						"widgets": []interface{}{
							map[string]interface{}{
								"type": "table",
							},
						},
					},
				},
			},
		},
		{
			name: "handles deeply nested structures",
			input: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"emptyArray": []interface{}{},
						"value":      "keep",
					},
				},
			},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"value": "keep",
					},
				},
			},
		},
		{
			name:     "handles nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "handles primitive types",
			input:    "string",
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripEmptyArrays(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("stripEmptyArrays() = %v, want %v", result, tt.expected)
			}
		})
	}
}
