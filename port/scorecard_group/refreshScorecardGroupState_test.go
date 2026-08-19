package scorecard_group

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestShouldRefreshGroupLevels(t *testing.T) {
	t.Parallel()

	customLevels := []scorecard.Level{
		{Title: types.StringValue("Ready"), Color: types.StringValue("green")},
	}
	customCliLevels := []cli.Level{{Title: "Ready", Color: "green"}}

	tests := []struct {
		name        string
		stateLevels []scorecard.Level
		cliLevels   []cli.Level
		want        bool
	}{
		{
			name:        "omit default levels when state has none",
			stateLevels: nil,
			cliLevels:   scorecard.DefaultCliLevels(),
			want:        false,
		},
		{
			name:        "refresh custom levels from api",
			stateLevels: customLevels,
			cliLevels:   customCliLevels,
			want:        true,
		},
		{
			name:        "refresh non-default api levels when state has none",
			stateLevels: nil,
			cliLevels:   customCliLevels,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := shouldRefreshGroupLevels(tt.stateLevels, tt.cliLevels); got != tt.want {
				t.Fatalf("shouldRefreshGroupLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshScorecardGroupStateNilMetadata(t *testing.T) {
	t.Parallel()

	resource := &ScorecardGroupResource{
		portClient: &cli.PortClient{},
	}
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-1"),
		Title:      types.StringValue("Group 1"),
	}
	group := &cli.ScorecardGroup{
		Identifier: "group-1",
		Title:      "Group 1",
		Blueprints: []string{"bp-1"},
		Rules: []cli.Rule{
			{
				Identifier: "rule-1",
				Title:      "Rule 1",
				Level:      "Gold",
				Query: cli.Query{
					Combinator: "and",
					Conditions: []interface{}{
						map[string]interface{}{
							"property": "$team",
							"operator": "isNotEmpty",
						},
					},
				},
			},
		},
	}

	resource.refreshScorecardGroupState(context.Background(), state, group)

	if !state.CreatedAt.IsNull() || !state.UpdatedAt.IsNull() {
		t.Fatalf("expected null timestamps, got created_at=%v updated_at=%v", state.CreatedAt, state.UpdatedAt)
	}
	if !state.CreatedBy.IsNull() || !state.UpdatedBy.IsNull() {
		t.Fatalf("expected null metadata users, got created_by=%v updated_by=%v", state.CreatedBy, state.UpdatedBy)
	}
}

func TestRefreshScorecardGroupStateSharedRulesFilters(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	resource := &ScorecardGroupResource{
		portClient: &cli.PortClient{},
	}
	state := &ScorecardGroupModel{}
	group := &cli.ScorecardGroup{
		Meta: cli.Meta{
			CreatedAt: &createdAt,
			CreatedBy: "user-1",
		},
		Identifier: "group-1",
		Title:      "Group 1",
		Blueprints: []string{"bp-1"},
		Rules: []cli.Rule{
			{
				Identifier: "rule-1",
				Title:      "Rule 1",
				Level:      "Gold",
				Query: cli.Query{
					Combinator: "and",
					Conditions: []interface{}{
						map[string]interface{}{
							"property": "$team",
							"operator": "isNotEmpty",
						},
					},
				},
			},
		},
		Filters: map[string]*cli.Query{
			"bp-1": {
				Combinator: "and",
				Conditions: []interface{}{
					map[string]interface{}{
						"property": "author",
						"operator": "isNotEmpty",
					},
				},
			},
		},
	}

	resource.refreshScorecardGroupState(context.Background(), state, group)

	if len(state.Filters) != 1 {
		t.Fatalf("expected one filter in state, got %d", len(state.Filters))
	}
	if state.Filters["bp-1"] == nil || state.Filters["bp-1"].Combinator.ValueString() != "and" {
		t.Fatalf("unexpected filters in state: %v", state.Filters)
	}
}

func TestRulesFromStateAndCLIPreservesOrder(t *testing.T) {
	t.Parallel()

	stateRules := []scorecard.Rule{
		{Identifier: types.StringValue("zebra-rule")},
		{Identifier: types.StringValue("alpha-rule")},
		{Identifier: types.StringValue("beta-rule")},
	}
	apiRules := []cli.Rule{
		{Identifier: "alpha-rule", Title: "Alpha Rule", Level: "Silver", Query: cli.Query{Combinator: "and"}},
		{Identifier: "beta-rule", Title: "Beta Rule", Level: "Bronze", Query: cli.Query{Combinator: "and"}},
		{Identifier: "zebra-rule", Title: "Zebra Rule", Level: "Gold", Query: cli.Query{Combinator: "and"}},
	}

	got := rulesFromStateAndCLI(stateRules, apiRules, false)

	if len(got) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(got))
	}
	if got[0].Identifier.ValueString() != "zebra-rule" {
		t.Fatalf("expected zebra-rule first, got %s", got[0].Identifier.ValueString())
	}
	if got[1].Identifier.ValueString() != "alpha-rule" {
		t.Fatalf("expected alpha-rule second, got %s", got[1].Identifier.ValueString())
	}
	if got[2].Identifier.ValueString() != "beta-rule" {
		t.Fatalf("expected beta-rule third, got %s", got[2].Identifier.ValueString())
	}
}

func TestRefreshScorecardGroupStatePreservesPerBlueprintRuleOrder(t *testing.T) {
	t.Parallel()

	resource := &ScorecardGroupResource{
		portClient: &cli.PortClient{},
	}
	state := &ScorecardGroupModel{
		Scorecards: map[string]MemberSpecModel{
			"bp-1": {
				Rules: []scorecard.Rule{
					{Identifier: types.StringValue("zebra-rule")},
					{Identifier: types.StringValue("alpha-rule")},
				},
			},
		},
	}
	group := &cli.ScorecardGroup{
		Identifier: "group-1",
		Title:      "Group 1",
		Scorecards: map[string]cli.ScorecardGroupMemberSpec{
			"bp-1": {
				Rules: []cli.Rule{
					{Identifier: "alpha-rule", Title: "Alpha Rule", Level: "Silver", Query: cli.Query{Combinator: "and"}},
					{Identifier: "zebra-rule", Title: "Zebra Rule", Level: "Gold", Query: cli.Query{Combinator: "and"}},
				},
			},
		},
	}

	resource.refreshScorecardGroupState(context.Background(), state, group)

	rules := state.Scorecards["bp-1"].Rules
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Identifier.ValueString() != "zebra-rule" {
		t.Fatalf("expected zebra-rule first, got %s", rules[0].Identifier.ValueString())
	}
	if rules[1].Identifier.ValueString() != "alpha-rule" {
		t.Fatalf("expected alpha-rule second, got %s", rules[1].Identifier.ValueString())
	}
}
