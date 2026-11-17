package action

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestDatasetRuleRoundTrip(t *testing.T) {
	model := &DatasetModel{
		Combinator: types.StringValue("and"),
		Rules: []Rule{
			{
				PropertyJson: types.StringValue(`{"context":"user","property":"port_role"}`),
				Operator:     types.StringValue("="),
				Value:        &Value{JqQuery: types.StringValue(`"Admin"`)},
			},
			{
				PropertyJson: types.StringValue(`{"path":["team_memberships","team"],"fromBlueprint":"person"}`),
				Operator:     types.StringValue("matchAny"),
				Value:        &Value{JqQuery: types.StringValue(".user.id")},
			},
			{
				Property: types.StringValue("language"),
				Operator: types.StringValue("="),
				Value:    &Value{JqQuery: types.StringValue(".form.language")},
			},
		},
	}

	body, err := actionDataSetToPortBody(model)
	require.NoError(t, err)

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"combinator":"and","rules":[
		{"property":{"context":"user","property":"port_role"},"operator":"=","value":{"jqQuery":"\"Admin\""}},
		{"property":{"path":["team_memberships","team"],"fromBlueprint":"person"},"operator":"matchAny","value":{"jqQuery":".user.id"}},
		{"property":"language","operator":"=","value":{"jqQuery":".form.language"}}
	]}`, string(encoded))

	// Simulate the API response and convert back to state.
	var fromApi cli.Dataset
	require.NoError(t, json.Unmarshal(encoded, &fromApi))
	state := writeDatasetToResource(&fromApi)
	require.Len(t, state.Rules, 3)
	assert.JSONEq(t, `{"context":"user","property":"port_role"}`, state.Rules[0].PropertyJson.ValueString())
	assert.Equal(t, `"Admin"`, state.Rules[0].Value.JqQuery.ValueString())
	assert.True(t, state.Rules[0].Property.IsNull())
	assert.Equal(t, "language", state.Rules[2].Property.ValueString())
	assert.Equal(t, ".form.language", state.Rules[2].Value.JqQuery.ValueString())
}
