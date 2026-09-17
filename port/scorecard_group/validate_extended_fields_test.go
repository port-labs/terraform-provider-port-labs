package scorecard_group

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestValidateScorecardGroupConfigurationRejectsReservedScorecardRelation(t *testing.T) {
	state := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"x","operator":"isNotEmpty"}`),
					},
				},
			},
		},
		ScorecardRelations: types.StringValue(`{"group":"some-group-id"}`),
	}

	err := validateScorecardGroupConfiguration(state)
	if err == nil {
		t.Fatal("expected error for reserved scorecard_relations key")
	}
}
