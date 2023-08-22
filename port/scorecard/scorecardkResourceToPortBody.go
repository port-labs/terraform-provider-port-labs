package scorecard

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func scorecardResourceToPortBody(ctx context.Context, state *ScorecardModel) (*cli.Scorecard, error) {
	s := &cli.Scorecard{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	rules := []cli.Rule{}

	for _, stateRule := range state.Rules {
		rule := &cli.Rule{
			Level:      stateRule.Level.ValueString(),
			Identifier: stateRule.Identifier.ValueString(),
			Title:      stateRule.Title.ValueString(),
		}

		query := &cli.Query{
			Combinator: stateRule.Query.Combinator.ValueString(),
		}

		conditions := []cli.Condition{}
		for _, stateCondition := range stateRule.Query.Conditions {
			condition := &cli.Condition{
				Property: stateCondition.Property.ValueString(),
				Operator: stateCondition.Operator.ValueString(),
			}

			if !stateCondition.Value.IsNull() {
				value := stateCondition.Value.ValueString()
				condition.Value = &value
			}

			conditions = append(conditions, *condition)
		}
		query.Conditions = conditions
		rule.Query = *query
		rules = append(rules, *rule)
	}

	s.Rules = rules

	return s, nil
}
