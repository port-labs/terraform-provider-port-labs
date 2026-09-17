package page

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageFiltersToPortBodyWithLast3DaysPreset(t *testing.T) {
	pageFilter := `{
		"identifier": "recent-entities-filter",
		"title": "Created in the last 3 days",
		"query": {
			"combinator": "and",
			"rules": [
				{
					"property": "$createdAt",
					"operator": "between",
					"value": {
						"preset": "last3Days"
					}
				}
			],
			"blueprint": "microservice"
		}
	}`

	pageFilters, diags := types.ListValue(types.StringType, []attr.Value{
		types.StringValue(pageFilter),
	})
	require.False(t, diags.HasError())

	pm := &PageModel{
		Identifier:  types.StringValue("recent_entities_page"),
		Type:        types.StringValue("blueprint-entities"),
		PageFilters: pageFilters,
	}

	body, err := PageToPortBody(pm)
	require.NoError(t, err)
	require.NotNil(t, body.PageFilters)
	require.Len(t, *body.PageFilters, 1)

	filter := (*body.PageFilters)[0]
	query, ok := filter["query"].(map[string]any)
	require.True(t, ok)

	rules, ok := query["rules"].([]any)
	require.True(t, ok)
	require.Len(t, rules, 1)

	rule, ok := rules[0].(map[string]any)
	require.True(t, ok)

	value, ok := rule["value"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "last3Days", value["preset"])
}
