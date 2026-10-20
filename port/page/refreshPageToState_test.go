package page

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
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
			result := stripEmptyArrays(tt.input, nil)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("stripEmptyArrays() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func refreshWidgetsForTest(t *testing.T, oldWidgets []string, apiWidgets []map[string]any) []string {
	t.Helper()
	r := &PageResource{portClient: &cli.PortClient{}}
	pm := &PageModel{Widgets: types.ListNull(types.StringType)}
	if oldWidgets != nil {
		vals := make([]attr.Value, len(oldWidgets))
		for i, w := range oldWidgets {
			vals[i] = types.StringValue(w)
		}
		pm.Widgets = types.ListValueMust(types.StringType, vals)
	}
	if err := r.refreshPageToState(pm, &cli.Page{Identifier: "p", Widgets: &apiWidgets}); err != nil {
		t.Fatalf("refreshPageToState: %v", err)
	}
	out := make([]string, 0, len(pm.Widgets.Elements()))
	for _, v := range pm.Widgets.Elements() {
		out = append(out, v.(types.String).ValueString())
	}
	return out
}

func assertJSONEqual(t *testing.T, want, got string) {
	t.Helper()
	var w, g any
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("bad want JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(got), &g); err != nil {
		t.Fatalf("bad got JSON: %v", err)
	}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("widget JSON mismatch\nwant: %s\ngot:  %s", want, got)
	}
}

// I5: an empty array the API injects for a field the user never wrote must not cause a diff.
func TestRefreshPageToState_StripsInjectedEmptyArrays(t *testing.T) {
	old := `{"id":"w1","type":"markdown","markdown":"# Hello"}`
	got := refreshWidgetsForTest(t, []string{old}, []map[string]any{
		{"id": "w1", "type": "markdown", "markdown": "# Hello", "hiddenQuery": []any{}},
	})
	assertJSONEqual(t, old, got[0])
}

// A user-written [] must be kept, at the top level and nested inside groups/datasets, while injected ones next to it
// are still stripped.
func TestRefreshPageToState_KeepsUserWrittenEmptyArrays(t *testing.T) {
	old := `{"id":"g","type":"grouper","hiddenQuery":[],"dataset":{"combinator":"and","rules":[]},` +
		`"groups":[{"title":"Graph","widgets":[{"id":"graph","hiddenBlueprints":[]}]}]}`
	api := map[string]any{
		"id": "g", "type": "grouper", "hiddenQuery": []any{},
		"dataset": map[string]any{"combinator": "and", "rules": []any{}},
		"groups": []any{map[string]any{
			"title":   "Graph",
			"widgets": []any{map[string]any{"id": "graph", "hiddenBlueprints": []any{}, "hiddenQuery": []any{}}},
		}},
	}
	got := refreshWidgetsForTest(t, []string{old}, []map[string]any{api})
	assertJSONEqual(t, old, got[0])
}

// When the user wrote a non-empty array and the API now returns [], the change must stay visible as drift.
func TestRefreshPageToState_EmptiedArrayShowsDrift(t *testing.T) {
	old := `{"id":"w1","hiddenBlueprints":["service"]}`
	got := refreshWidgetsForTest(t, []string{old}, []map[string]any{{"id": "w1", "hiddenBlueprints": []any{}}})
	assertJSONEqual(t, `{"id":"w1","hiddenBlueprints":[]}`, got[0])
}

// With no prior state (import, or a widget added out of band) there is no user intent to compare with, so injected
// empty arrays are stripped as before.
func TestRefreshPageToState_NoPriorStateStripsEmptyArrays(t *testing.T) {
	got := refreshWidgetsForTest(t, nil, []map[string]any{{"id": "w1", "hiddenQuery": []any{}}})
	assertJSONEqual(t, `{"id":"w1"}`, got[0])
}
