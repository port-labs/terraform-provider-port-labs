package scorecard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithModifyPlan = &ScorecardResource{}

func (r *ScorecardResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ScorecardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !isRulesOrderOnlyChange(plan.Rules, state.Rules) {
		return
	}

	plan.Rules = state.Rules
	plan.ID = state.ID
	plan.UpdatedAt = state.UpdatedAt
	plan.UpdatedBy = state.UpdatedBy
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func isRulesOrderOnlyChange(plan, state []Rule) bool {
	if len(plan) != len(state) {
		return false
	}
	byID := make(map[string]Rule, len(state))
	isSameOrder := true
	for i, rule := range state {
		byID[rule.Identifier.ValueString()] = rule
		if plan[i].Identifier.ValueString() != rule.Identifier.ValueString() {
			isSameOrder = false
		}
	}
	if isSameOrder {
		return false
	}
	for _, rule := range plan {
		other, ok := byID[rule.Identifier.ValueString()]
		if !ok || !areRulesEqual(rule, other) {
			return false
		}
	}
	return true
}

func areRulesEqual(a, b Rule) bool {
	if !a.Identifier.Equal(b.Identifier) || !a.Title.Equal(b.Title) || !a.Level.Equal(b.Level) ||
		!areStringsEqual(a.Description, b.Description) || !areStringsEqual(a.Properties, b.Properties) || !areStringsEqual(a.Relations, b.Relations) {
		return false
	}
	if a.Query == nil || b.Query == nil {
		return a.Query == b.Query
	}
	if !a.Query.Combinator.Equal(b.Query.Combinator) || len(a.Query.Conditions) != len(b.Query.Conditions) {
		return false
	}
	for i := range a.Query.Conditions {
		if !a.Query.Conditions[i].Equal(b.Query.Conditions[i]) {
			return false
		}
	}
	return true
}

func areStringsEqual(a, b types.String) bool {
	if (a.IsNull() || a.ValueString() == "") && (b.IsNull() || b.ValueString() == "") {
		return true
	}
	return a.Equal(b)
}
