package scorecard_group

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupMembershipChangedSharedRulesBlueprints(t *testing.T) {
	previous := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules:      []scorecard.Rule{{Identifier: types.StringValue("rule-1")}},
	}
	current := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1"), types.StringValue("bp-2")},
		Rules:      []scorecard.Rule{{Identifier: types.StringValue("rule-1")}},
	}

	if !scorecardGroupMembershipChanged(previous, current) {
		t.Fatal("expected membership change when blueprints are added")
	}
}

func TestScorecardGroupMembershipChangedPerBlueprintScorecards(t *testing.T) {
	previous := &ScorecardGroupModel{
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {},
		},
	}
	current := &ScorecardGroupModel{
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {},
			"bp-2": {},
		},
	}

	if !scorecardGroupMembershipChanged(previous, current) {
		t.Fatal("expected membership change when scorecards keys are added")
	}
}

func TestScorecardGroupMembershipChangedModeSwitch(t *testing.T) {
	previous := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules:      []scorecard.Rule{{Identifier: types.StringValue("rule-1")}},
	}
	current := &ScorecardGroupModel{
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {},
		},
	}

	if !scorecardGroupMembershipChanged(previous, current) {
		t.Fatal("expected membership change when switching between shared-rules and per-blueprint modes")
	}
}

func TestScorecardGroupMembershipNotChangedForRuleUpdate(t *testing.T) {
	previous := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules:      []scorecard.Rule{{Identifier: types.StringValue("rule-1"), Title: types.StringValue("Old")}},
	}
	current := &ScorecardGroupModel{
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules:      []scorecard.Rule{{Identifier: types.StringValue("rule-1"), Title: types.StringValue("New")}},
	}

	if scorecardGroupMembershipChanged(previous, current) {
		t.Fatal("did not expect membership change when only rules content changes")
	}
}
