package scorecard_group

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestValidateScorecardRelationsRejectsGroupKey(t *testing.T) {
	t.Parallel()

	state := &ScorecardGroupModel{
		ScorecardRelations: types.StringValue(`{"group":"some-group"}`),
		Blueprints:         []types.String{types.StringValue("svc")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{types.StringValue(`{"property":"x","operator":"isNotEmpty"}`)},
				},
			},
		},
	}

	if err := validateScorecardGroupConfiguration(state); err == nil {
		t.Fatal("expected error for reserved group relation key")
	}
}
