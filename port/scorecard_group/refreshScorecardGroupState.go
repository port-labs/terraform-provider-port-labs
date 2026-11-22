package scorecard_group

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

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

func rulesFromCLI(rules []cli.Rule, jsonEscapeHTML bool, syncPropertiesFromAPI bool) ([]scorecard.Rule, error) {
	stateRules := make([]scorecard.Rule, 0, len(rules))
	for _, rule := range rules {
		stateRule, err := ruleFromCLI(rule, scorecard.Rule{}, jsonEscapeHTML, syncPropertiesFromAPI)
		if err != nil {
			return nil, err
		}
		stateRules = append(stateRules, stateRule)
	}
	return stateRules, nil
}

func ruleFromCLI(rule cli.Rule, existingRule scorecard.Rule, jsonEscapeHTML bool, syncPropertiesFromAPI bool) (scorecard.Rule, error) {
	stateRule := scorecard.Rule{
		Title:      types.StringValue(rule.Title),
		Level:      types.StringValue(rule.Level),
		Identifier: types.StringValue(rule.Identifier),
		Properties: existingRule.Properties,
		Relations:  existingRule.Relations,
	}
	if !existingRule.Description.IsNull() && !existingRule.Description.IsUnknown() {
		stateRule.Description = existingRule.Description
	} else if rule.Description != "" {
		stateRule.Description = types.StringValue(rule.Description)
	} else {
		stateRule.Description = types.StringNull()
	}
	stateRule.Query = queryFromCLI(&rule.Query, jsonEscapeHTML)
	if err := scorecard.SyncRuleCustomFieldsFromAPI(&stateRule, rule, jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
		return scorecard.Rule{}, err
	}
	return stateRule, nil
}

func rulesFromStateAndCLI(stateRules []scorecard.Rule, apiRules []cli.Rule, jsonEscapeHTML bool, syncPropertiesFromAPI bool) ([]scorecard.Rule, error) {
	if len(stateRules) == 0 {
		return rulesFromCLI(apiRules, jsonEscapeHTML, syncPropertiesFromAPI)
	}

	apiRulesByIdentifier := make(map[string]cli.Rule, len(apiRules))
	for _, rule := range apiRules {
		apiRulesByIdentifier[rule.Identifier] = rule
	}

	orderedRules := make([]scorecard.Rule, 0, len(apiRules))
	processedIdentifiers := make(map[string]bool)

	for _, existingRule := range stateRules {
		identifier := existingRule.Identifier.ValueString()
		apiRule, exists := apiRulesByIdentifier[identifier]
		if !exists {
			continue
		}
		updatedRule, err := ruleFromCLI(apiRule, existingRule, jsonEscapeHTML, syncPropertiesFromAPI)
		if err != nil {
			return nil, err
		}
		orderedRules = append(orderedRules, updatedRule)
		processedIdentifiers[identifier] = true
	}

	for _, apiRule := range apiRules {
		if !processedIdentifiers[apiRule.Identifier] {
			newRule, err := ruleFromCLI(apiRule, scorecard.Rule{}, jsonEscapeHTML, syncPropertiesFromAPI)
			if err != nil {
				return nil, err
			}
			orderedRules = append(orderedRules, newRule)
		}
	}

	return orderedRules, nil
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

func memberSpecFromCLI(spec cli.ScorecardGroupMemberSpec, stateSpec MemberSpecModel, jsonEscapeHTML bool, syncPropertiesFromAPI bool) (MemberSpecModel, error) {
	rules, err := rulesFromStateAndCLI(stateSpec.Rules, spec.Rules, jsonEscapeHTML, syncPropertiesFromAPI)
	if err != nil {
		return MemberSpecModel{}, err
	}
	return MemberSpecModel{
		Filter: queryFromCLI(spec.Filter, jsonEscapeHTML),
		Rules:  rules,
	}, nil
}

func (r *ScorecardGroupResource) jsonEscapeHTML() bool {
	if r.portClient == nil {
		return false
	}
	return r.portClient.JSONEscapeHTML
}

func refreshPerBlueprintState(state *ScorecardGroupModel, group *cli.ScorecardGroup, jsonEscapeHTML bool, syncPropertiesFromAPI bool) error {
	scorecards, err := memberSpecsFromGroup(state.Scorecards, group, jsonEscapeHTML, syncPropertiesFromAPI)
	if err != nil {
		return err
	}
	state.Scorecards = scorecards
	state.Blueprints = nil
	state.Rules = nil
	state.Filters = nil
	return nil
}

func memberSpecsFromGroup(previous map[string]MemberSpecModel, group *cli.ScorecardGroup, jsonEscapeHTML bool, syncPropertiesFromAPI bool) (map[string]MemberSpecModel, error) {
	if len(group.Scorecards) > 0 {
		scorecards := make(map[string]MemberSpecModel, len(group.Scorecards))
		for blueprintID, memberSpec := range group.Scorecards {
			existingSpec := MemberSpecModel{}
			if previous != nil {
				existingSpec = previous[blueprintID]
			}
			spec, err := memberSpecFromCLI(memberSpec, existingSpec, jsonEscapeHTML, syncPropertiesFromAPI)
			if err != nil {
				return nil, err
			}
			scorecards[blueprintID] = spec
		}
		return scorecards, nil
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
		rules, err := rulesFromStateAndCLI(existingSpec.Rules, group.Rules, jsonEscapeHTML, syncPropertiesFromAPI)
		if err != nil {
			return nil, err
		}
		scorecards[blueprintID] = MemberSpecModel{
			Filter: queryFromCLI(filter, jsonEscapeHTML),
			Rules:  rules,
		}
	}
	return scorecards, nil
}

func refreshSharedRulesState(state *ScorecardGroupModel, group *cli.ScorecardGroup, jsonEscapeHTML bool, syncPropertiesFromAPI bool) error {
	state.Scorecards = nil
	if len(group.Blueprints) > 0 {
		state.Blueprints = make([]types.String, len(group.Blueprints))
		for i, blueprint := range group.Blueprints {
			state.Blueprints[i] = types.StringValue(blueprint)
		}
	} else {
		state.Blueprints = nil
	}

	rules, err := rulesFromStateAndCLI(state.Rules, group.Rules, jsonEscapeHTML, syncPropertiesFromAPI)
	if err != nil {
		return err
	}
	state.Rules = rules
	if len(group.Filters) > 0 {
		state.Filters = make(map[string]*scorecard.Query, len(group.Filters))
		for blueprintID, filter := range group.Filters {
			state.Filters[blueprintID] = queryFromCLI(filter, jsonEscapeHTML)
		}
	} else {
		state.Filters = nil
	}
	return nil
}

func (r *ScorecardGroupResource) refreshScorecardGroupState(ctx context.Context, state *ScorecardGroupModel, group *cli.ScorecardGroup, syncPropertiesFromAPI bool) error {
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
	if err := utils.SyncJSONObjectState(&state.ScorecardProperties, group.ScorecardProperties, "scorecard properties", jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
		return err
	}

	if err := utils.SyncJSONObjectState(&state.GroupProperties, group.GroupProperties, "group properties", jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
		return err
	}
	if err := utils.SyncJSONObjectState(&state.GroupRelations, group.GroupRelations, "group relations", jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
		return err
	}
	if err := utils.SyncJSONObjectState(&state.ScorecardRelations, group.ScorecardRelations, "scorecard relations", jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
		return err
	}
	switch {
	case len(state.Scorecards) > 0:
		// Keep per-blueprint config even when the API canonicalizes to shared-rules.
		if err := refreshPerBlueprintState(state, group, jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
			return err
		}
	case len(state.Blueprints) > 0 || len(state.Rules) > 0 || len(state.Filters) > 0:
		if err := refreshSharedRulesState(state, group, jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
			return err
		}
	case len(group.Scorecards) > 0:
		if err := refreshPerBlueprintState(state, group, jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
			return err
		}
	default:
		if err := refreshSharedRulesState(state, group, jsonEscapeHTML, syncPropertiesFromAPI); err != nil {
			return err
		}
	}
	return nil
}
