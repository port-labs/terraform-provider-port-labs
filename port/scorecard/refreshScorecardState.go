package scorecard

import (
	"context"
	"fmt"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"reflect"
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

func refreshScorecardState(ctx context.Context, state *ScorecardModel, s *cli.Scorecard, blueprintIdentifier string) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintIdentifier, s.Identifier))
	state.Identifier = types.StringValue(s.Identifier)
	state.Blueprint = types.StringValue(blueprintIdentifier)
	state.Title = types.StringValue(s.Title)
	state.CreatedAt = types.StringValue(s.CreatedAt.String())
	state.CreatedBy = types.StringValue(s.CreatedBy)
	state.UpdatedAt = types.StringValue(s.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(s.UpdatedBy)

	stateRules := []Rule{}
	for _, rule := range s.Rules {
		stateRule := &Rule{
			Title:      types.StringValue(rule.Title),
			Level:      types.StringValue(rule.Level),
			Identifier: types.StringValue(rule.Identifier),
		}
		stateQuery := &Query{
			Combinator: types.StringValue(rule.Query.Combinator),
		}

		stateQuery.Conditions = make([]types.String, len(rule.Query.Conditions))
		for i, u := range rule.Query.Conditions {
			cond, _ := utils.GoObjectToTerraformString(u)
			stateQuery.Conditions[i] = cond
		}

		stateRule.Query = stateQuery

		stateRules = append(stateRules, *stateRule)
	}

	state.Rules = stateRules
	if shouldRefreshLevels(state.Levels, s.Levels) {
		state.Levels = fromCliLevelsToTerraformLevels(s.Levels)
	}
}
