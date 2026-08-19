package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupResourceToPortBodySharedRules(t *testing.T) {
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-1"),
		Title:      types.StringValue("Group 1"),
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"$team","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(context.TODO(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Identifier != "group-1" {
		t.Fatalf("expected identifier group-1, got %s", group.Identifier)
	}
	if len(group.Blueprints) != 1 || group.Blueprints[0] != "bp-1" {
		t.Fatalf("unexpected blueprints: %v", group.Blueprints)
	}
	if len(group.Rules) != 1 || group.Rules[0].Identifier != "rule-1" {
		t.Fatalf("unexpected rules: %v", group.Rules)
	}
	if group.Scorecards != nil {
		t.Fatalf("expected scorecards to be unset in shared rules mode")
	}
}

func TestScorecardGroupResourceToPortBodySharedRulesWithFilters(t *testing.T) {
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-1"),
		Title:      types.StringValue("Group 1"),
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"$team","operator":"isNotEmpty"}`),
					},
				},
			},
		},
		Filters: map[string]*scorecard.Query{
			"bp-1": {
				Combinator: types.StringValue("and"),
				Conditions: []types.String{
					types.StringValue(`{"property":"$team","operator":"isNotEmpty"}`),
				},
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(context.TODO(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(group.Filters) != 1 {
		t.Fatalf("expected one filter, got %d", len(group.Filters))
	}
	if group.Filters["bp-1"] == nil || group.Filters["bp-1"].Combinator != "and" {
		t.Fatalf("unexpected filters: %v", group.Filters)
	}
}

func TestScorecardGroupResourceToPortBodyPerBlueprint(t *testing.T) {
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-2"),
		Title:      types.StringValue("Group 2"),
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {
				Rules: []scorecard.Rule{
					{
						Identifier: types.StringValue("rule-1"),
						Title:      types.StringValue("Rule 1"),
						Level:      types.StringValue("Gold"),
						Query: &scorecard.Query{
							Combinator: types.StringValue("and"),
							Conditions: []types.String{
								types.StringValue(`{"property":"author","operator":"isNotEmpty"}`),
							},
						},
					},
				},
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(context.TODO(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(group.Scorecards) != 1 {
		t.Fatalf("expected one scorecard member, got %d", len(group.Scorecards))
	}
	if _, ok := group.Scorecards["bp-1"]; !ok {
		t.Fatalf("expected scorecard for bp-1")
	}
	if len(group.Blueprints) > 0 || len(group.Rules) > 0 {
		t.Fatalf("expected shared rules fields to be unset in per-blueprint mode")
	}
}
