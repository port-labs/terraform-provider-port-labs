package scorecard_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
	"github.com/stretchr/testify/require"
)

func TestApplyRuleCustomFieldsToCLI(t *testing.T) {
	rule := &cli.Rule{
		Identifier: "has-owner",
		Title:      "Has Owner",
		Level:      "Gold",
	}
	stateRule := scorecard.Rule{
		Properties: types.StringValue(`{"sla":2,"severity":"high"}`),
		Relations:  types.StringValue(`{"owner":"team-a"}`),
	}

	err := scorecard.ApplyRuleCustomFieldsToCLI(stateRule, rule)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"sla": float64(2), "severity": "high"}, rule.Properties)
	require.Equal(t, map[string]any{"owner": "team-a"}, rule.Relations)
}

func TestSyncRuleCustomFieldsFromAPI(t *testing.T) {
	stateRule := scorecard.Rule{
		Properties: types.StringValue(`{"sla":2}`),
		Relations:  types.StringValue(`{"owner":"team-a"}`),
	}
	apiRule := cli.Rule{
		Properties: map[string]any{"sla": float64(2)},
		Relations:  map[string]any{"owner": "team-a"},
	}

	err := scorecard.SyncRuleCustomFieldsFromAPI(&stateRule, apiRule, false, false)
	require.NoError(t, err)
	require.JSONEq(t, `{"sla":2}`, stateRule.Properties.ValueString())
	require.JSONEq(t, `{"owner":"team-a"}`, stateRule.Relations.ValueString())
}

func TestSyncRuleCustomFieldsFromAPIRejectsMissingValues(t *testing.T) {
	stateRule := scorecard.Rule{
		Properties: types.StringValue(`{"sla":2}`),
	}
	apiRule := cli.Rule{
		Properties: map[string]any{},
	}

	err := scorecard.SyncRuleCustomFieldsFromAPI(&stateRule, apiRule, false, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rule properties were not applied")
}
