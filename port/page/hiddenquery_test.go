package page

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestNormalizeHiddenQueryArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "sorts hiddenQuery at top level",
			input: map[string]interface{}{
				"type":        "widget",
				"hiddenQuery": []interface{}{"zebra", "alpha", "beta"},
			},
			expected: map[string]interface{}{
				"type":        "widget",
				"hiddenQuery": []interface{}{"alpha", "beta", "zebra"},
			},
		},
		{
			name: "same hiddenQuery in different order becomes equal",
			input: map[string]interface{}{
				"id":          "w1",
				"hiddenQuery": []interface{}{"b", "a", "c"},
			},
			expected: map[string]interface{}{
				"id":          "w1",
				"hiddenQuery": []interface{}{"a", "b", "c"},
			},
		},
		{
			name: "preserves order of non-hiddenQuery arrays",
			input: map[string]interface{}{
				"type":    "widget",
				"widgets": []interface{}{"first", "second", "third"},
				"rules":   []interface{}{1, 2, 3},
			},
			expected: map[string]interface{}{
				"type":    "widget",
				"widgets": []interface{}{"first", "second", "third"},
				"rules":   []interface{}{1, 2, 3},
			},
		},
		{
			name: "normalizes nested hiddenQuery in widgets",
			input: map[string]interface{}{
				"groups": []interface{}{
					map[string]interface{}{
						"widgets": []interface{}{
							map[string]interface{}{
								"type":        "table",
								"hiddenQuery": []interface{}{"z", "a"},
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
								"type":        "table",
								"hiddenQuery": []interface{}{"a", "z"},
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
						"hiddenQuery": []interface{}{"c", "b", "a"},
						"otherArray":  []interface{}{3, 2, 1},
					},
				},
			},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"hiddenQuery": []interface{}{"a", "b", "c"},
						"otherArray":  []interface{}{3, 2, 1},
					},
				},
			},
		},
		{
			name: "handles empty hiddenQuery",
			input: map[string]interface{}{
				"hiddenQuery": []interface{}{},
			},
			expected: map[string]interface{}{
				"hiddenQuery": []interface{}{},
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
		{
			name: "handles object elements in hiddenQuery",
			input: map[string]interface{}{
				"hiddenQuery": []interface{}{
					map[string]interface{}{"property": "z"},
					map[string]interface{}{"property": "a"},
				},
			},
			expected: map[string]interface{}{
				"hiddenQuery": []interface{}{
					map[string]interface{}{"property": "a"},
					map[string]interface{}{"property": "z"},
				},
			},
		},
		{
			name: "handles numeric elements in hiddenQuery",
			input: map[string]interface{}{
				"hiddenQuery": []interface{}{3.0, 1.0, 2.0},
			},
			expected: map[string]interface{}{
				"hiddenQuery": []interface{}{1.0, 2.0, 3.0},
			},
		},
		{
			name: "non-array hiddenQuery is preserved",
			input: map[string]interface{}{
				"hiddenQuery": "not-an-array",
			},
			expected: map[string]interface{}{
				"hiddenQuery": "not-an-array",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeHiddenQueryArrays(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeHiddenQueryArrays() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortStringArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "sorts strings",
			input:    []interface{}{"c", "a", "b"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "sorts numbers",
			input:    []interface{}{3.0, 1.0, 2.0},
			expected: []interface{}{1.0, 2.0, 3.0},
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "single element",
			input:    []interface{}{"only"},
			expected: []interface{}{"only"},
		},
		{
			name: "sorts objects by JSON representation",
			input: []interface{}{
				map[string]interface{}{"z": 1},
				map[string]interface{}{"a": 2},
			},
			expected: []interface{}{
				map[string]interface{}{"a": 2},
				map[string]interface{}{"z": 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortStringArray(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortStringArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormalizeWidgetForComparison(t *testing.T) {
	t.Run("normalizes widget map", func(t *testing.T) {
		widget := map[string]any{
			"id":          "w1",
			"type":        "table",
			"hiddenQuery": []any{"z", "a"},
		}
		result := NormalizeWidgetForComparison(widget)
		expected := map[string]any{
			"id":          "w1",
			"type":        "table",
			"hiddenQuery": []interface{}{"a", "z"},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("NormalizeWidgetForComparison() = %v, want %v", result, expected)
		}
	})
}

func TestHiddenQueryOrderInsensitiveComparison(t *testing.T) {
	t.Run("same hiddenQuery different order yields equal JSON", func(t *testing.T) {
		config := map[string]any{
			"id":          "w1",
			"type":        "table",
			"hiddenQuery": []any{"alpha", "beta", "gamma"},
		}
		apiResponse := map[string]any{
			"id":          "w1",
			"type":        "table",
			"hiddenQuery": []any{"gamma", "alpha", "beta"},
		}

		normalizedConfig := NormalizeWidgetForComparison(config)
		normalizedAPI := NormalizeWidgetForComparison(apiResponse)

		configJSON, _ := json.Marshal(normalizedConfig)
		apiJSON, _ := json.Marshal(normalizedAPI)

		if string(configJSON) != string(apiJSON) {
			t.Errorf("normalized JSON should be equal:\nconfig: %s\napi: %s", configJSON, apiJSON)
		}
	})

	t.Run("different hiddenQuery contents are not equal", func(t *testing.T) {
		config := map[string]any{
			"id":          "w1",
			"hiddenQuery": []any{"alpha", "beta"},
		}
		apiResponse := map[string]any{
			"id":          "w1",
			"hiddenQuery": []any{"alpha", "gamma"},
		}

		normalizedConfig := NormalizeWidgetForComparison(config)
		normalizedAPI := NormalizeWidgetForComparison(apiResponse)

		configJSON, _ := json.Marshal(normalizedConfig)
		apiJSON, _ := json.Marshal(normalizedAPI)

		if string(configJSON) == string(apiJSON) {
			t.Errorf("normalized JSON should NOT be equal when contents differ")
		}
	})

	t.Run("nested widgets with reordered hiddenQuery", func(t *testing.T) {
		config := map[string]any{
			"id":   "grouper",
			"type": "grouper",
			"groups": []any{
				map[string]any{
					"title": "Overview",
					"widgets": []any{
						map[string]any{
							"id":          "dashboard",
							"hiddenQuery": []any{"prop1", "prop2", "prop3"},
						},
					},
				},
			},
		}
		apiResponse := map[string]any{
			"id":   "grouper",
			"type": "grouper",
			"groups": []any{
				map[string]any{
					"title": "Overview",
					"widgets": []any{
						map[string]any{
							"id":          "dashboard",
							"hiddenQuery": []any{"prop3", "prop1", "prop2"},
						},
					},
				},
			},
		}

		normalizedConfig := NormalizeWidgetForComparison(config)
		normalizedAPI := NormalizeWidgetForComparison(apiResponse)

		configJSON, _ := json.Marshal(normalizedConfig)
		apiJSON, _ := json.Marshal(normalizedAPI)

		if string(configJSON) != string(apiJSON) {
			t.Errorf("nested widgets with reordered hiddenQuery should be equal:\nconfig: %s\napi: %s", configJSON, apiJSON)
		}
	})
}

// TestRefreshPageToState_HiddenQueryOrderInsensitive tests the full integration of hiddenQuery
// order-insensitive comparison through refreshPageToState. This simulates the scenario where
// the Port backend reorders hiddenQuery arrays asynchronously per organization.
func TestRefreshPageToState_HiddenQueryOrderInsensitive(t *testing.T) {
	t.Run("reordered hiddenQuery does not cause diff", func(t *testing.T) {
		r := &PageResource{portClient: &cli.PortClient{}}
		oldWidgetJSON := `{"id":"w1","type":"table","hiddenQuery":["prop1","prop2","prop3"]}`
		pm := &PageModel{
			Widgets: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(oldWidgetJSON),
			}),
		}

		apiWidget := map[string]any{
			"id":          "w1",
			"type":        "table",
			"hiddenQuery": []any{"prop3", "prop1", "prop2"},
		}

		err := r.refreshPageToState(pm, &cli.Page{
			Identifier: "test-page",
			Widgets:    &[]map[string]any{apiWidget},
		})
		if err != nil {
			t.Fatalf("refreshPageToState: %v", err)
		}

		gotJSON := pm.Widgets.Elements()[0].(types.String).ValueString()
		if gotJSON != oldWidgetJSON {
			t.Errorf("widget JSON changed when only hiddenQuery order differed:\nold: %s\ngot: %s", oldWidgetJSON, gotJSON)
		}
	})

	t.Run("different hiddenQuery contents still shows change", func(t *testing.T) {
		r := &PageResource{portClient: &cli.PortClient{}}
		oldWidgetJSON := `{"id":"w1","hiddenQuery":["prop1","prop2"]}`
		pm := &PageModel{
			Widgets: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(oldWidgetJSON),
			}),
		}

		apiWidget := map[string]any{
			"id":          "w1",
			"hiddenQuery": []any{"prop1", "prop3"},
		}

		err := r.refreshPageToState(pm, &cli.Page{
			Identifier: "test-page",
			Widgets:    &[]map[string]any{apiWidget},
		})
		if err != nil {
			t.Fatalf("refreshPageToState: %v", err)
		}

		gotJSON := pm.Widgets.Elements()[0].(types.String).ValueString()
		if gotJSON == oldWidgetJSON {
			t.Errorf("widget JSON should change when hiddenQuery contents differ")
		}
	})

	t.Run("nested widgets with reordered hiddenQuery do not cause diff", func(t *testing.T) {
		r := &PageResource{portClient: &cli.PortClient{}}
		oldWidgetJSON := `{"groups":[{"title":"Overview","widgets":[{"hiddenQuery":["a","b","c"],"id":"nested"}]}],"id":"grouper","type":"grouper"}`
		pm := &PageModel{
			Widgets: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(oldWidgetJSON),
			}),
		}

		apiWidget := map[string]any{
			"id":   "grouper",
			"type": "grouper",
			"groups": []any{
				map[string]any{
					"title": "Overview",
					"widgets": []any{
						map[string]any{
							"id":          "nested",
							"hiddenQuery": []any{"c", "a", "b"},
						},
					},
				},
			},
		}

		err := r.refreshPageToState(pm, &cli.Page{
			Identifier: "test-page",
			Widgets:    &[]map[string]any{apiWidget},
		})
		if err != nil {
			t.Fatalf("refreshPageToState: %v", err)
		}

		gotJSON := pm.Widgets.Elements()[0].(types.String).ValueString()
		if gotJSON != oldWidgetJSON {
			t.Errorf("nested widget JSON changed when only hiddenQuery order differed:\nold: %s\ngot: %s", oldWidgetJSON, gotJSON)
		}
	})

	t.Run("other array order changes are still detected", func(t *testing.T) {
		r := &PageResource{portClient: &cli.PortClient{}}
		oldWidgetJSON := `{"id":"w1","rules":["rule1","rule2"]}`
		pm := &PageModel{
			Widgets: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(oldWidgetJSON),
			}),
		}

		apiWidget := map[string]any{
			"id":    "w1",
			"rules": []any{"rule2", "rule1"},
		}

		err := r.refreshPageToState(pm, &cli.Page{
			Identifier: "test-page",
			Widgets:    &[]map[string]any{apiWidget},
		})
		if err != nil {
			t.Fatalf("refreshPageToState: %v", err)
		}

		gotJSON := pm.Widgets.Elements()[0].(types.String).ValueString()
		if gotJSON == oldWidgetJSON {
			t.Errorf("non-hiddenQuery array order change should be detected:\nold: %s\ngot: %s", oldWidgetJSON, gotJSON)
		}
	})
}
