package scorecard

import (
	"context"
	"fmt"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

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

}
