package scorecard

import (
	"context"
	"encoding/json"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)


func DefaultLevels() []cli.Level {
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

func scorecardResourceToPortBody(ctx context.Context, state *ScorecardModel) (*cli.Scorecard, error) {
	s := &cli.Scorecard{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	var rules []cli.Rule

	for _, stateRule := range state.Rules {
		rule := &cli.Rule{
			Level:      stateRule.Level.ValueString(),
			Identifier: stateRule.Identifier.ValueString(),
			Title:      stateRule.Title.ValueString(),
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
	
	if len(state.Levels) == 0 {
		s.Levels = DefaultLevels()
	} else {
		var levels []cli.Level
		for _, stateLevel := range state.Levels {
			level := &cli.Level{
				Color: stateLevel.Color.ValueString(),
				Title: stateLevel.Title.ValueString(),
			}
			levels = append(levels, *level)
		}
		s.Levels = levels
	}
	return s, nil
}
