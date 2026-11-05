package page

import (
	"encoding/json"
	"fmt"
	"slices"
)

// normalizeHiddenQueryArrays recursively traverses a JSON structure and sorts any array found at
// a key named "hiddenQuery". This ensures order-insensitive comparison for hiddenQuery fields,
// which the Port backend may reorder asynchronously per organization.
//
// Other arrays (e.g., widgets, rules, layout) are preserved in their original order since their
// order is semantically meaningful.
func normalizeHiddenQueryArrays(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, value := range v {
			normalized := normalizeHiddenQueryArrays(value)
			if key == "hiddenQuery" {
				if arr, ok := normalized.([]interface{}); ok {
					result[key] = sortStringArray(arr)
				} else {
					result[key] = normalized
				}
			} else {
				result[key] = normalized
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeHiddenQueryArrays(item)
		}
		return result
	default:
		return data
	}
}

// sortStringArray sorts an array of interface{} values by their JSON string representation.
// This provides a deterministic order for hiddenQuery arrays regardless of the original order.
func sortStringArray(arr []interface{}) []interface{} {
	if len(arr) == 0 {
		return arr
	}

	type sortable struct {
		value   interface{}
		sortKey string
	}

	items := make([]sortable, len(arr))
	for i, v := range arr {
		key, err := json.Marshal(v)
		if err != nil {
			key = []byte(fmt.Sprintf("%v", v))
		}
		items[i] = sortable{value: v, sortKey: string(key)}
	}

	slices.SortStableFunc(items, func(a, b sortable) int {
		if a.sortKey < b.sortKey {
			return -1
		}
		if a.sortKey > b.sortKey {
			return 1
		}
		return 0
	})

	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item.value
	}
	return result
}

// NormalizeWidgetForComparison normalizes a widget map by sorting hiddenQuery arrays.
// This should be called before comparing widget JSON to ensure order-insensitive comparison
// for hiddenQuery fields while preserving order sensitivity for other arrays.
func NormalizeWidgetForComparison(widget map[string]any) map[string]any {
	normalized := normalizeHiddenQueryArrays(widget)
	if m, ok := normalized.(map[string]interface{}); ok {
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result
	}
	return widget
}
