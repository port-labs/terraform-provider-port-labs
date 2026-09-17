package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupResourceToPortBodyExtendedFields(t *testing.T) {
	ctx := context.Background()
	state := &ScorecardGroupModel{
		Identifier:          types.StringValue("group-1"),
		Title:               types.StringValue("My Group"),
		GroupProperties:     types.StringValue(`{"category":"production"}`),
		ScorecardProperties: types.StringValue(`{"owner":"platform"}`),
		GroupRelations:      types.StringValue(`{"owner":"team-abc"}`),
		ScorecardRelations:  types.StringValue(`{"sponsors":["entity-1","entity-2"]}`),
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {
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
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(ctx, state)
	if err != nil {
		t.Fatalf("scorecardGroupResourceToPortBody: %v", err)
	}

	if group.GroupProperties["category"] != "production" {
		t.Fatalf("groupProperties: got %#v", group.GroupProperties)
	}
	if group.ScorecardProperties["owner"] != "platform" {
		t.Fatalf("scorecardProperties: got %#v", group.ScorecardProperties)
	}
	if group.GroupRelations["owner"] != "team-abc" {
		t.Fatalf("groupRelations: got %#v", group.GroupRelations)
	}
	sponsors, ok := group.ScorecardRelations["sponsors"].([]interface{})
	if !ok || len(sponsors) != 2 || sponsors[0] != "entity-1" {
		t.Fatalf("scorecardRelations: got %#v", group.ScorecardRelations)
	}
}
