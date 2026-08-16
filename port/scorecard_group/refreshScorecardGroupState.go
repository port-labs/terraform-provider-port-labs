package scorecard_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func queryFromCLI(q *cli.Query, jsonEscapeHTML bool) *scorecard.Query {
	if q == nil {
		return nil
	}
	stateQuery := &scorecard.Query{
		Combinator: types.StringValue(q.Combinator),
	}
	stateQuery.Conditions = make([]types.String, len(q.Conditions))
	for i, condition := range q.Conditions {
		cond, _ := utils.GoObjectToTerraformString(condition, jsonEscapeHTML)
		stateQuery.Conditions[i] = cond
	}
	return stateQuery
}

func rulesFromCLI(rules []cli.Rule, jsonEscapeHTML bool) []scorecard.Rule {
	stateRules := make([]scorecard.Rule, 0, len(rules))
	for _, rule := range rules {
		stateRule := scorecard.Rule{
			Title:      types.StringValue(rule.Title),
			Level:      types.StringValue(rule.Level),
			Identifier: types.StringValue(rule.Identifier),
		}
		if rule.Description != "" {
			stateRule.Description = types.StringValue(rule.Description)
		} else {
			stateRule.Description = types.StringNull()
		}
		stateRule.Query = queryFromCLI(&rule.Query, jsonEscapeHTML)
		stateRules = append(stateRules, stateRule)
	}
	return stateRules
}

func levelsFromCLI(levels []cli.Level) []scorecard.Level {
	stateLevels := make([]scorecard.Level, 0, len(levels))
	for _, level := range levels {
		stateLevels = append(stateLevels, scorecard.Level{
			Color: types.StringValue(level.Color),
			Title: types.StringValue(level.Title),
		})
	}
	return stateLevels
}

func memberSpecFromCLI(spec cli.ScorecardGroupMemberSpec, jsonEscapeHTML bool) MemberSpecModel {
	return MemberSpecModel{
		Filter: queryFromCLI(spec.Filter, jsonEscapeHTML),
		Rules:  rulesFromCLI(spec.Rules, jsonEscapeHTML),
	}
}

func (r *ScorecardGroupResource) refreshScorecardGroupState(ctx context.Context, state *ScorecardGroupModel, group *cli.ScorecardGroup) {
	state.ID = types.StringValue(group.Identifier)
	state.Identifier = types.StringValue(group.Identifier)
	state.Title = types.StringValue(group.Title)
	state.CreatedAt = types.StringValue(group.CreatedAt.String())
	state.CreatedBy = types.StringValue(group.CreatedBy)
	state.UpdatedAt = types.StringValue(group.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(group.UpdatedBy)

	if len(group.Levels) > 0 {
		state.Levels = levelsFromCLI(group.Levels)
	} else {
		state.Levels = nil
	}

	if len(group.Scorecards) > 0 {
		state.Scorecards = make(map[string]MemberSpecModel, len(group.Scorecards))
		for blueprintID, memberSpec := range group.Scorecards {
			state.Scorecards[blueprintID] = memberSpecFromCLI(memberSpec, r.portClient.JSONEscapeHTML)
		}
		state.Blueprints = nil
		state.Rules = nil
		state.Filter = nil
		return
	}

	state.Scorecards = nil
	if len(group.Blueprints) > 0 {
		state.Blueprints = make([]types.String, len(group.Blueprints))
		for i, blueprint := range group.Blueprints {
			state.Blueprints[i] = types.StringValue(blueprint)
		}
	} else {
		state.Blueprints = nil
	}

	state.Rules = rulesFromCLI(group.Rules, r.portClient.JSONEscapeHTML)
	state.Filter = queryFromCLI(group.Filter, r.portClient.JSONEscapeHTML)
}

func blueprintIdentifiersFromState(state *ScorecardGroupModel) []string {
	if len(state.Scorecards) > 0 {
		blueprints := make([]string, 0, len(state.Scorecards))
		for blueprint := range state.Scorecards {
			blueprints = append(blueprints, blueprint)
		}
		return blueprints
	}

	blueprints := make([]string, 0, len(state.Blueprints))
	for _, blueprint := range state.Blueprints {
		blueprints = append(blueprints, blueprint.ValueString())
	}
	return blueprints
}

func (r *ScorecardGroupResource) deleteMemberScorecards(ctx context.Context, state *ScorecardGroupModel) error {
	identifier := state.Identifier.ValueString()
	for _, blueprintIdentifier := range blueprintIdentifiersFromState(state) {
		if err := r.portClient.DeleteScorecard(ctx, blueprintIdentifier, identifier); err != nil {
			return fmt.Errorf("failed to delete scorecard %q on blueprint %q: %w", identifier, blueprintIdentifier, err)
		}
	}
	return nil
}
