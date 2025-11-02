package scorecard

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func shouldRefreshLevels(stateLevels []Level, cliLevels []cli.Level) bool {
	// When you create a scorecard in Port, the scorecard gets created with default levels.
	// If your scorecard doesn't have the "levels" attribute, it means the scorecard is created with default levels behind the scenes.
	//
	// If the TF state has no levels and the Port existing levels are the default levels, This means both are considered
	// to have default levels, And so we don't need to update them.
	if len(stateLevels) == 0 && reflect.DeepEqual(cliLevels, DefaultCliLevels()) {
		return false
	}
	// If the TF state has defined levels, we have to make sure that Port's existing levels are the same as the TF state levels.
	// also,
	// If TF state doesn't have levels and the Port existing levels are not the default ones,
	// this means we have to make sure that Port's defined levels are the default levels,
	// as the state without levels is considered to have default levels.
	if len(stateLevels) > 0 || (len(stateLevels) == 0 && !reflect.DeepEqual(cliLevels, DefaultCliLevels())) {
		return true
	}

	return false
}

func fromCliLevelsToTerraformLevels(cliLevels []cli.Level) []Level {
	terraformLevels := []Level{}
	for _, cliLevel := range cliLevels {
		level := &Level{
			Color: types.StringValue(cliLevel.Color),
			Title: types.StringValue(cliLevel.Title),
		}
		terraformLevels = append(terraformLevels, *level)
	}
	return terraformLevels
}

func DefaultCliLevels() []cli.Level {
	return []cli.Level{
		{
			Color: "paleBlue",
			Title: "Basic",
		},
		{
			Color: "bronze",
			Title: "Bronze",
		},
		{
			Color: "silver",
			Title: "Silver",
		},
		{
			Color: "gold",
			Title: "Gold",
		},
	}
}

func (r *ScorecardResource) refreshScorecardState(ctx context.Context, state *ScorecardModel, s *cli.Scorecard, blueprintIdentifier string) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintIdentifier, s.Identifier))
	state.Identifier = types.StringValue(s.Identifier)
	state.Blueprint = types.StringValue(blueprintIdentifier)
	state.Title = types.StringValue(s.Title)
	state.CreatedAt = types.StringValue(s.CreatedAt.String())
	state.CreatedBy = types.StringValue(s.CreatedBy)
	state.UpdatedAt = types.StringValue(s.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(s.UpdatedBy)

	if s.Filter != nil {
		stateFilter := &Query{
			Combinator: types.StringValue(s.Filter.Combinator),
		}
		stateFilter.Conditions = make([]types.String, len(s.Filter.Conditions))
		for i, u := range s.Filter.Conditions {
			cond, _ := utils.GoObjectToTerraformString(u, r.portClient.JSONEscapeHTML)
			stateFilter.Conditions[i] = cond
		}
		state.Filter = stateFilter
	}

	stateRules := []Rule{}
	for _, rule := range s.Rules {
		stateRule := &Rule{
			Title:      types.StringValue(rule.Title),
			Level:      types.StringValue(rule.Level),
			Identifier: types.StringValue(rule.Identifier),
		}

		if rule.Description != "" {
			stateRule.Description = types.StringValue(rule.Description)
		} else {
			stateRule.Description = types.StringNull()
		}

		stateQuery := &Query{
			Combinator: types.StringValue(rule.Query.Combinator),
		}

		stateQuery.Conditions = make([]types.String, len(rule.Query.Conditions))
		for i, u := range rule.Query.Conditions {
			cond, _ := utils.GoObjectToTerraformString(u, r.portClient.JSONEscapeHTML)
			stateQuery.Conditions[i] = cond
		}

		stateRule.Query = stateQuery

		stateRules = append(stateRules, *stateRule)
	}

	// Preserve the original order from state if it exists
	if len(state.Rules) > 0 {
		// Create a map of API rule identifier to API rule for quick lookup
		apiRulesByIdentifier := make(map[string]*cli.Rule)
		for i := range s.Rules {
			apiRulesByIdentifier[s.Rules[i].Identifier] = &s.Rules[i]
		}

		// Process existing rules in their original order, updating with API data
		orderedRules := make([]Rule, 0, len(stateRules))
		processedIdentifiers := make(map[string]bool)

		// First, update existing rules in their original order
		for _, existingRule := range state.Rules {
			identifier := existingRule.Identifier.ValueString()
			if apiRule, exists := apiRulesByIdentifier[identifier]; exists {
				// Update the existing rule with fresh API data while preserving structure
				updatedRule := Rule{
					Title:      types.StringValue(apiRule.Title),
					Level:      types.StringValue(apiRule.Level),
					Identifier: types.StringValue(apiRule.Identifier),
				}

				// Handle description - preserve from state if it exists, otherwise use API value
				if !existingRule.Description.IsNull() && !existingRule.Description.IsUnknown() {
					updatedRule.Description = existingRule.Description
				} else if apiRule.Description != "" {
					updatedRule.Description = types.StringValue(apiRule.Description)
				} else {
					updatedRule.Description = types.StringNull()
				}

				// Update query from API
				stateQuery := &Query{
					Combinator: types.StringValue(apiRule.Query.Combinator),
				}
				stateQuery.Conditions = make([]types.String, len(apiRule.Query.Conditions))
				for i, u := range apiRule.Query.Conditions {
					cond, _ := utils.GoObjectToTerraformString(u, r.portClient.JSONEscapeHTML)
					stateQuery.Conditions[i] = cond
				}
				updatedRule.Query = stateQuery

				orderedRules = append(orderedRules, updatedRule)
				processedIdentifiers[identifier] = true
			}
		}

		// Then append any new rules from API that weren't in the original state
		for _, apiRule := range s.Rules {
			if !processedIdentifiers[apiRule.Identifier] {
				newRule := Rule{
					Title:      types.StringValue(apiRule.Title),
					Level:      types.StringValue(apiRule.Level),
					Identifier: types.StringValue(apiRule.Identifier),
				}

				if apiRule.Description != "" {
					newRule.Description = types.StringValue(apiRule.Description)
				} else {
					newRule.Description = types.StringNull()
				}

				stateQuery := &Query{
					Combinator: types.StringValue(apiRule.Query.Combinator),
				}
				stateQuery.Conditions = make([]types.String, len(apiRule.Query.Conditions))
				for i, u := range apiRule.Query.Conditions {
					cond, _ := utils.GoObjectToTerraformString(u, r.portClient.JSONEscapeHTML)
					stateQuery.Conditions[i] = cond
				}
				newRule.Query = stateQuery

				orderedRules = append(orderedRules, newRule)
			}
		}

		state.Rules = orderedRules
	} else {
		// No existing state, use API order
		state.Rules = stateRules
	}
	if shouldRefreshLevels(state.Levels, s.Levels) {
		state.Levels = fromCliLevelsToTerraformLevels(s.Levels)
	}
}
