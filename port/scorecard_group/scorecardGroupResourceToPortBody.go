package scorecard_group

import (
	"context"
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func queryToCLI(q *scorecard.Query) (*cli.Query, error) {
	if q == nil {
		return nil, nil
	}

	query := &cli.Query{
		Combinator: q.Combinator.ValueString(),
	}
	var conditions []interface{}
	for _, stateCondition := range q.Conditions {
		if !stateCondition.IsNull() {
			cond := map[string]interface{}{}
			err := json.Unmarshal([]byte(stateCondition.ValueString()), &cond)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, cond)
		}
	}
	query.Conditions = conditions
	return query, nil
}

func rulesToCLI(rules []scorecard.Rule) ([]cli.Rule, error) {
	var cliRules []cli.Rule
	for _, stateRule := range rules {
		rule := cli.Rule{
			Level:      stateRule.Level.ValueString(),
			Identifier: stateRule.Identifier.ValueString(),
			Title:      stateRule.Title.ValueString(),
		}
		if !stateRule.Description.IsNull() {
			rule.Description = stateRule.Description.ValueString()
		}
		query, err := queryToCLI(stateRule.Query)
		if err != nil {
			return nil, err
		}
		if query != nil {
			rule.Query = *query
		}
		cliRules = append(cliRules, rule)
	}
	return cliRules, nil
}

func levelsToCLI(levels []scorecard.Level) []cli.Level {
	var cliLevels []cli.Level
	for _, level := range levels {
		cliLevels = append(cliLevels, cli.Level{
			Color: level.Color.ValueString(),
			Title: level.Title.ValueString(),
		})
	}
	return cliLevels
}

func memberSpecToCLI(spec MemberSpecModel) (cli.ScorecardGroupMemberSpec, error) {
	filter, err := queryToCLI(spec.Filter)
	if err != nil {
		return cli.ScorecardGroupMemberSpec{}, err
	}
	rules, err := rulesToCLI(spec.Rules)
	if err != nil {
		return cli.ScorecardGroupMemberSpec{}, err
	}
	return cli.ScorecardGroupMemberSpec{
		Filter: filter,
		Rules:  rules,
	}, nil
}

func memberSpecToPatchCLI(spec MemberSpecModel) (cli.PatchScorecardGroupMemberSpec, error) {
	filter, err := queryToCLI(spec.Filter)
	if err != nil {
		return cli.PatchScorecardGroupMemberSpec{}, err
	}
	rules, err := rulesToCLI(spec.Rules)
	if err != nil {
		return cli.PatchScorecardGroupMemberSpec{}, err
	}
	return cli.PatchScorecardGroupMemberSpec{
		Filter: filter,
		Rules:  rules,
	}, nil
}

func scorecardGroupResourceToPortBody(ctx context.Context, state *ScorecardGroupModel) (*cli.ScorecardGroup, error) {
	group := &cli.ScorecardGroup{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	if len(state.Levels) > 0 {
		group.Levels = levelsToCLI(state.Levels)
	}

	if !state.Properties.IsNull() && !state.Properties.IsUnknown() {
		properties, err := utils.TerraformJsonStringToGoObject(state.Properties.ValueStringPointer())
		if err != nil {
			return nil, err
		}
		if properties != nil {
			group.Properties = *properties
		}
	}

	if len(state.Scorecards) > 0 {
		scorecards := make(map[string]cli.ScorecardGroupMemberSpec, len(state.Scorecards))
		for blueprintID, memberSpec := range state.Scorecards {
			spec, err := memberSpecToCLI(memberSpec)
			if err != nil {
				return nil, err
			}
			scorecards[blueprintID] = spec
		}
		group.Scorecards = scorecards
		return group, nil
	}

	if len(state.Blueprints) > 0 {
		blueprints := make([]string, len(state.Blueprints))
		for i, bp := range state.Blueprints {
			blueprints[i] = bp.ValueString()
		}
		group.Blueprints = blueprints
	}

	rules, err := rulesToCLI(state.Rules)
	if err != nil {
		return nil, err
	}
	group.Rules = rules

	if len(state.Filters) > 0 {
		filters := make(map[string]*cli.Query, len(state.Filters))
		for blueprintID, filter := range state.Filters {
			query, err := queryToCLI(filter)
			if err != nil {
				return nil, err
			}
			filters[blueprintID] = query
		}
		group.Filters = filters
	}

	return group, nil
}

func scorecardGroupResourceToPatchBody(ctx context.Context, state *ScorecardGroupModel) (*cli.PatchScorecardGroup, error) {
	patch := &cli.PatchScorecardGroup{
		Title: state.Title.ValueString(),
	}

	if len(state.Levels) > 0 {
		patch.Levels = levelsToCLI(state.Levels)
	}

	if !state.Properties.IsNull() && !state.Properties.IsUnknown() {
		properties, err := utils.TerraformJsonStringToGoObject(state.Properties.ValueStringPointer())
		if err != nil {
			return nil, err
		}
		if properties != nil {
			patch.Properties = *properties
		}
	}

	if len(state.Scorecards) > 0 {
		scorecards := make(map[string]cli.PatchScorecardGroupMemberSpec, len(state.Scorecards))
		for blueprintID, memberSpec := range state.Scorecards {
			spec, err := memberSpecToPatchCLI(memberSpec)
			if err != nil {
				return nil, err
			}
			scorecards[blueprintID] = spec
		}
		patch.Scorecards = scorecards
		return patch, nil
	}

	rules, err := rulesToCLI(state.Rules)
	if err != nil {
		return nil, err
	}
	patch.Rules = rules

	if len(state.Filters) > 0 {
		filters := make(map[string]*cli.Query, len(state.Filters))
		for blueprintID, filter := range state.Filters {
			query, err := queryToCLI(filter)
			if err != nil {
				return nil, err
			}
			filters[blueprintID] = query
		}
		patch.Filters = filters
	}

	return patch, nil
}
