package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupResourceToPortBodyProperties(t *testing.T) {
	t.Parallel()

	state := &ScorecardGroupModel{
		Identifier:          types.StringValue("my-group"),
		Title:               types.StringValue("My Group"),
		GroupProperties:     types.StringValue(`{"category":"production"}`),
		ScorecardProperties: types.StringValue(`{"owner":"platform-team"}`),
		Blueprints:          []types.String{types.StringValue("svc")},
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

	group, err := scorecardGroupResourceToPortBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.GroupProperties["category"] != "production" {
		t.Fatalf("unexpected groupProperties: %v", group.GroupProperties)
	}
	if group.ScorecardProperties["owner"] != "platform-team" {
		t.Fatalf("unexpected scorecardProperties: %v", group.ScorecardProperties)
	}
}
