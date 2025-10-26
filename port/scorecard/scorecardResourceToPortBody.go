package scorecard

import (
	"context"
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func fromTerraformLevelsToCliLevels(tfLevels []Level) []cli.Level {
	var levels []cli.Level
	for _, stateLevel := range tfLevels {
		level := &cli.Level{
			Color: stateLevel.Color.ValueString(),
			Title: stateLevel.Title.ValueString(),
		}
		levels = append(levels, *level)
	}
	return levels
}

func scorecardResourceToPortBody(ctx context.Context, state *ScorecardModel) (*cli.Scorecard, error) {
	s := &cli.Scorecard{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	if state.Filter != nil {
		filter := &cli.Query{
			Combinator: state.Filter.Combinator.ValueString(),
		}
		var conditions []interface{}
		for _, stateCondition := range state.Filter.Conditions {
			if !stateCondition.IsNull() {
				stringCond := stateCondition.ValueString()
				cond := map[string]interface{}{}
				err := json.Unmarshal([]byte(stringCond), &cond)
				if err != nil {
					return nil, err
				}
				conditions = append(conditions, cond)
			}
		}
		filter.Conditions = conditions
		s.Filter = filter
	}

	// Sort rules by identifier to ensure consistent ordering
	sortedStateRules := make([]Rule, len(state.Rules))
	copy(sortedStateRules, state.Rules)

	for i := 0; i < len(sortedStateRules)-1; i++ {
		for j := i + 1; j < len(sortedStateRules); j++ {
			if sortedStateRules[i].Identifier.ValueString() > sortedStateRules[j].Identifier.ValueString() {
				sortedStateRules[i], sortedStateRules[j] = sortedStateRules[j], sortedStateRules[i]
			}
		}
	}

	var rules []cli.Rule

	for _, stateRule := range sortedStateRules {
		rule := &cli.Rule{
			Level:      stateRule.Level.ValueString(),
			Identifier: stateRule.Identifier.ValueString(),
			Title:      stateRule.Title.ValueString(),
		}

		if !stateRule.Description.IsNull() {
			rule.Description = stateRule.Description.ValueString()
		}

		query := &cli.Query{
			Combinator: stateRule.Query.Combinator.ValueString(),
		}
		var conditions []interface{}
		for _, stateCondition := range stateRule.Query.Conditions {
			if !stateCondition.IsNull() {
				stringCond := stateCondition.ValueString()
				cond := map[string]interface{}{}
				err := json.Unmarshal([]byte(stringCond), &cond)
				if err != nil {
					return nil, err
				}
				conditions = append(conditions, cond)
			}
		}
		query.Conditions = conditions
		rule.Query = *query

		rules = append(rules, *rule)
	}

	s.Rules = rules

	if len(state.Levels) > 0 {
		s.Levels = fromTerraformLevelsToCliLevels(state.Levels)
	}

	return s, nil
}
