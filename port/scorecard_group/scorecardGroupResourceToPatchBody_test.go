package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupResourceToPatchBodySharedRules(t *testing.T) {
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-1"),
		Title:      types.StringValue("Updated Group"),
		Blueprints: []types.String{types.StringValue("bp-1"), types.StringValue("bp-2")},
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
		Filters: map[string]*scorecard.Query{
			"bp-1": {
				Combinator: types.StringValue("and"),
				Conditions: []types.String{
					types.StringValue(`{"property":"environment","operator":"=","value":"production"}`),
				},
			},
		},
	}

	patch, err := scorecardGroupResourceToPatchBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if patch.Title != "Updated Group" {
		t.Fatalf("expected title Updated Group, got %q", patch.Title)
	}
	if len(patch.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(patch.Rules))
	}
	if len(patch.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(patch.Filters))
	}
	if len(patch.Scorecards) != 0 {
		t.Fatalf("expected no scorecards in shared-rules patch body")
	}
}

func TestScorecardGroupResourceToPatchBodyPerBlueprint(t *testing.T) {
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-2"),
		Title:      types.StringValue("Per Blueprint Group"),
		Properties: types.StringValue(`{"owner":"platform-team"}`),
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

	patch, err := scorecardGroupResourceToPatchBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if patch.Properties["owner"] != "platform-team" {
		t.Fatalf("expected properties owner=platform-team, got %v", patch.Properties)
	}
	if len(patch.Scorecards) != 1 {
		t.Fatalf("expected 1 scorecard member, got %d", len(patch.Scorecards))
	}
	if len(patch.Rules) != 0 {
		t.Fatalf("expected no shared rules in per-blueprint patch body")
	}
	if len(patch.Filters) != 0 {
		t.Fatalf("expected no filters in per-blueprint patch body")
	}
}

func TestScorecardGroupResourceToPatchBodyOmitsBlueprints(t *testing.T) {
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
						types.StringValue(`{"property":"author","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	patch, err := scorecardGroupResourceToPatchBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(patch.Rules) != 1 {
		t.Fatalf("expected rules in patch body, got %d", len(patch.Rules))
	}
}
