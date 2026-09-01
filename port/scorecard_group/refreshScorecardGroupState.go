package scorecard_group

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func configuredPropertyKeys(stateProperties types.String) map[string]struct{} {
	if stateProperties.IsNull() || stateProperties.IsUnknown() {
		return nil
	}

	var props map[string]any
	if err := json.Unmarshal([]byte(stateProperties.ValueString()), &props); err != nil || len(props) == 0 {
		return nil
	}

	keys := make(map[string]struct{}, len(props))
	for key := range props {
		keys[key] = struct{}{}
	}
	return keys
}

func propertiesFromAPIForRead(stateProperties types.String, apiProperties map[string]any, jsonEscapeHTML bool) types.String {
	if stateProperties.IsNull() || stateProperties.IsUnknown() {
		return types.StringNull()
	}

	configuredKeys := configuredPropertyKeys(stateProperties)
	if len(configuredKeys) == 0 {
		return types.StringNull()
	}

	var stateMap map[string]any
	if err := json.Unmarshal([]byte(stateProperties.ValueString()), &stateMap); err != nil {
		return stateProperties
	}

	merged := make(map[string]any, len(stateMap))
	for key := range stateMap {
		if apiValue, ok := apiProperties[key]; ok {
			merged[key] = apiValue
		} else {
			merged[key] = nil
		}
	}

	properties, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return stateProperties
	}
	return properties
}

func shouldRefreshGroupLevels(stateLevels []scorecard.Level, cliLevels []cli.Level) bool {
	if len(stateLevels) == 0 && reflect.DeepEqual(cliLevels, scorecard.DefaultCliLevels()) {
		return false
	}
	if len(stateLevels) > 0 || (len(stateLevels) == 0 && !reflect.DeepEqual(cliLevels, scorecard.DefaultCliLevels())) {
		return true
	}
	return false
}

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
		stateRules = append(stateRules, ruleFromCLI(rule, scorecard.Rule{}, jsonEscapeHTML))
	}
	return stateRules
}

func ruleFromCLI(rule cli.Rule, existingRule scorecard.Rule, jsonEscapeHTML bool) scorecard.Rule {
	stateRule := scorecard.Rule{
		Title:      types.StringValue(rule.Title),
		Level:      types.StringValue(rule.Level),
		Identifier: types.StringValue(rule.Identifier),
	}
	if !existingRule.Description.IsNull() && !existingRule.Description.IsUnknown() {
		stateRule.Description = existingRule.Description
	} else if rule.Description != "" {
		stateRule.Description = types.StringValue(rule.Description)
	} else {
		stateRule.Description = types.StringNull()
	}
	stateRule.Query = queryFromCLI(&rule.Query, jsonEscapeHTML)
	return stateRule
}

func rulesFromStateAndCLI(stateRules []scorecard.Rule, apiRules []cli.Rule, jsonEscapeHTML bool) []scorecard.Rule {
	apiStateRules := rulesFromCLI(apiRules, jsonEscapeHTML)
	if len(stateRules) == 0 {
		return apiStateRules
	}

	apiRulesByIdentifier := make(map[string]scorecard.Rule, len(apiStateRules))
	for _, rule := range apiStateRules {
		apiRulesByIdentifier[rule.Identifier.ValueString()] = rule
	}

	orderedRules := make([]scorecard.Rule, 0, len(apiStateRules))
	processedIdentifiers := make(map[string]bool)

	for _, existingRule := range stateRules {
		identifier := existingRule.Identifier.ValueString()
		apiRule, exists := apiRulesByIdentifier[identifier]
		if !exists {
			continue
		}
		updatedRule := apiRule
		if !existingRule.Description.IsNull() && !existingRule.Description.IsUnknown() {
			updatedRule.Description = existingRule.Description
		}
		orderedRules = append(orderedRules, updatedRule)
		processedIdentifiers[identifier] = true
	}

	for _, apiRule := range apiStateRules {
		if !processedIdentifiers[apiRule.Identifier.ValueString()] {
			orderedRules = append(orderedRules, apiRule)
		}
	}

	return orderedRules
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

func memberSpecFromCLI(spec cli.ScorecardGroupMemberSpec, stateSpec MemberSpecModel, jsonEscapeHTML bool) MemberSpecModel {
	return MemberSpecModel{
		Filter: queryFromCLI(spec.Filter, jsonEscapeHTML),
		Rules:  rulesFromStateAndCLI(stateSpec.Rules, spec.Rules, jsonEscapeHTML),
	}
}

func (r *ScorecardGroupResource) jsonEscapeHTML() bool {
	if r.portClient == nil {
		return false
	}
	return r.portClient.JSONEscapeHTML
}

func refreshPerBlueprintState(state *ScorecardGroupModel, group *cli.ScorecardGroup, jsonEscapeHTML bool) {
	state.Scorecards = memberSpecsFromGroup(state.Scorecards, group, jsonEscapeHTML)
	state.Blueprints = nil
	state.Rules = nil
	state.Filters = nil
}

func memberSpecsFromGroup(previous map[string]MemberSpecModel, group *cli.ScorecardGroup, jsonEscapeHTML bool) map[string]MemberSpecModel {
	if len(group.Scorecards) > 0 {
		scorecards := make(map[string]MemberSpecModel, len(group.Scorecards))
		for blueprintID, memberSpec := range group.Scorecards {
			existingSpec := MemberSpecModel{}
			if previous != nil {
				existingSpec = previous[blueprintID]
			}
			scorecards[blueprintID] = memberSpecFromCLI(memberSpec, existingSpec, jsonEscapeHTML)
		}
		return scorecards
	}

	blueprints := group.Blueprints
	if len(blueprints) == 0 {
		blueprints = make([]string, 0, len(previous))
		for blueprintID := range previous {
			blueprints = append(blueprints, blueprintID)
		}
	}

	scorecards := make(map[string]MemberSpecModel, len(blueprints))
	for _, blueprintID := range blueprints {
		existingSpec := MemberSpecModel{}
		if previous != nil {
			existingSpec = previous[blueprintID]
		}
		var filter *cli.Query
		if group.Filters != nil {
			filter = group.Filters[blueprintID]
		}
		scorecards[blueprintID] = MemberSpecModel{
			Filter: queryFromCLI(filter, jsonEscapeHTML),
			Rules:  rulesFromStateAndCLI(existingSpec.Rules, group.Rules, jsonEscapeHTML),
		}
	}
	return scorecards
}

func refreshSharedRulesState(state *ScorecardGroupModel, group *cli.ScorecardGroup, jsonEscapeHTML bool) {
	state.Scorecards = nil
	if len(group.Blueprints) > 0 {
		state.Blueprints = make([]types.String, len(group.Blueprints))
		for i, blueprint := range group.Blueprints {
			state.Blueprints[i] = types.StringValue(blueprint)
		}
	} else {
		state.Blueprints = nil
	}

	state.Rules = rulesFromStateAndCLI(state.Rules, group.Rules, jsonEscapeHTML)
	if len(group.Filters) > 0 {
		state.Filters = make(map[string]*scorecard.Query, len(group.Filters))
		for blueprintID, filter := range group.Filters {
			state.Filters[blueprintID] = queryFromCLI(filter, jsonEscapeHTML)
		}
	} else {
		state.Filters = nil
	}
}

func (r *ScorecardGroupResource) refreshScorecardGroupState(ctx context.Context, state *ScorecardGroupModel, group *cli.ScorecardGroup, syncPropertiesFromAPI bool) {
	state.ID = types.StringValue(group.Identifier)
	state.Identifier = types.StringValue(group.Identifier)
	state.Title = types.StringValue(group.Title)
	if group.CreatedAt != nil {
		state.CreatedAt = types.StringValue(group.CreatedAt.String())
	} else {
		state.CreatedAt = types.StringNull()
	}
	if group.CreatedBy != "" {
		state.CreatedBy = types.StringValue(group.CreatedBy)
	} else {
		state.CreatedBy = types.StringNull()
	}
	if group.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(group.UpdatedAt.String())
	} else {
		state.UpdatedAt = types.StringNull()
	}
	if group.UpdatedBy != "" {
		state.UpdatedBy = types.StringValue(group.UpdatedBy)
	} else {
		state.UpdatedBy = types.StringNull()
	}

	if shouldRefreshGroupLevels(state.Levels, group.Levels) {
		if len(group.Levels) > 0 {
			state.Levels = levelsFromCLI(group.Levels)
		} else {
			state.Levels = nil
		}
	}

	jsonEscapeHTML := r.jsonEscapeHTML()
	if syncPropertiesFromAPI {
		state.Properties = propertiesFromAPIForRead(state.Properties, group.Properties, jsonEscapeHTML)
	}
	switch {
	case len(state.Scorecards) > 0:
		// Keep per-blueprint config even when the API canonicalizes to shared-rules.
		refreshPerBlueprintState(state, group, jsonEscapeHTML)
	case len(state.Blueprints) > 0 || len(state.Rules) > 0 || len(state.Filters) > 0:
		refreshSharedRulesState(state, group, jsonEscapeHTML)
	case len(group.Scorecards) > 0:
		refreshPerBlueprintState(state, group, jsonEscapeHTML)
	default:
		refreshSharedRulesState(state, group, jsonEscapeHTML)
	}
}
