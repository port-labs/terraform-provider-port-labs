package action

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *ActionResource) writeInvocationMethodToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if a.InvocationMethod.Type == consts.Kafka {
		var oldPayload types.String
		if state.KafkaMethod != nil {
			oldPayload = state.KafkaMethod.Payload
		}
		payload, err := utils.GoObjectToTerraformStringPreferExisting(oldPayload, a.InvocationMethod.Payload, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.KafkaMethod = &KafkaMethodModel{
			Payload: payload,
		}
	}

	if a.InvocationMethod.Type == consts.Webhook {
		var oldAgent, oldSynchronized, oldBody types.String
		if state.WebhookMethod != nil {
			oldAgent = state.WebhookMethod.Agent
			oldSynchronized = state.WebhookMethod.Synchronized
			oldBody = state.WebhookMethod.Body
		}
		agent, err := utils.GoObjectToTerraformStringPreferExisting(oldAgent, a.InvocationMethod.Agent, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		synchronized, err := utils.GoObjectToTerraformStringPreferExisting(oldSynchronized, a.InvocationMethod.Synchronized, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		headers, _ := types.MapValueFrom(ctx, types.StringType, a.InvocationMethod.Headers)
		body, err := utils.GoObjectToTerraformStringPreferExisting(oldBody, a.InvocationMethod.Body, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.WebhookMethod = &WebhookMethodModel{
			Url:          types.StringValue(*a.InvocationMethod.Url),
			Agent:        agent,
			Synchronized: synchronized,
			Method:       flex.GoStringToFramework(a.InvocationMethod.Method),
			Headers:      headers,
			Body:         body,
		}
	}

	if a.InvocationMethod.Type == consts.Github {
		var oldWorkflowInputs, oldReportWorkflowStatus types.String
		if state.GithubMethod != nil {
			oldWorkflowInputs = state.GithubMethod.WorkflowInputs
			oldReportWorkflowStatus = state.GithubMethod.ReportWorkflowStatus
		}
		workflowInputs, err := utils.GoObjectToTerraformStringPreferExisting(oldWorkflowInputs, a.InvocationMethod.WorkflowInputs, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		reportWorkflowStatus, err := utils.GoObjectToTerraformStringPreferExisting(oldReportWorkflowStatus, a.InvocationMethod.ReportWorkflowStatus, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.GithubMethod = &GithubMethodModel{
			Org:                  types.StringValue(*a.InvocationMethod.Org),
			Repo:                 types.StringValue(*a.InvocationMethod.Repo),
			Workflow:             types.StringValue(*a.InvocationMethod.Workflow),
			WorkflowInputs:       workflowInputs,
			ReportWorkflowStatus: reportWorkflowStatus,
		}
	}

	if a.InvocationMethod.Type == consts.Gitlab {
		var oldPipelineVariables types.String
		if state.GitlabMethod != nil {
			oldPipelineVariables = state.GitlabMethod.PipelineVariables
		}
		pipelineVariables, err := utils.GoObjectToTerraformStringPreferExisting(oldPipelineVariables, a.InvocationMethod.PipelineVariables, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.GitlabMethod = &GitlabMethodModel{
			ProjectName:       types.StringValue(*a.InvocationMethod.ProjectName),
			GroupName:         types.StringValue(*a.InvocationMethod.GroupName),
			DefaultRef:        flex.GoStringToFramework(a.InvocationMethod.DefaultRef),
			PipelineVariables: pipelineVariables,
		}
	}

	if a.InvocationMethod.Type == consts.AzureDevops {
		var oldPayload types.String
		if state.AzureMethod != nil {
			oldPayload = state.AzureMethod.Payload
		}
		payload, err := utils.GoObjectToTerraformStringPreferExisting(oldPayload, a.InvocationMethod.Payload, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
			Payload: payload,
		}
	}

	if a.InvocationMethod.Type == consts.UpsertEntity {
		var teams []types.String
		var teamsJQ types.String
		switch team := a.InvocationMethod.Mapping.Team.(type) {
		case string:
			teamsJQ = types.StringValue(team)
		case []interface{}:
			teams = make([]types.String, 0)
			for _, t := range team {
				teams = append(teams, types.StringValue(t.(string)))
			}
			teamsJQ = types.StringNull()
		}
		var oldProperties, oldRelations types.String
		if state.UpsertEntityMethod != nil && state.UpsertEntityMethod.Mapping != nil {
			oldProperties = state.UpsertEntityMethod.Mapping.Properties
			oldRelations = state.UpsertEntityMethod.Mapping.Relations
		}
		properties, err := utils.GoObjectToTerraformStringPreferExisting(oldProperties, a.InvocationMethod.Mapping.Properties, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		relations, err := utils.GoObjectToTerraformStringPreferExisting(oldRelations, a.InvocationMethod.Mapping.Relations, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.UpsertEntityMethod = &UpsertEntityMethodModel{
			Title:               flex.GoStringToFramework(a.InvocationMethod.Mapping.Title),
			BlueprintIdentifier: types.StringValue(*a.InvocationMethod.BlueprintIdentifier),
			Mapping: &MappingModel{
				Properties: properties,
				Relations:  relations,
				Icon:       flex.GoStringToFramework(a.InvocationMethod.Mapping.Icon),
				Teams:      teams,
				TeamsJQ:    teamsJQ,
				Identifier: types.StringPointerValue(a.InvocationMethod.Mapping.Identifier),
			},
		}
	}

	if a.InvocationMethod.Type == consts.IntegrationAction {
		var oldWorkflowInputs, oldReportWorkflowStatus types.String
		if state.IntegrationMethod != nil && state.IntegrationMethod.IntegrationActionExecutionProperties != nil {
			oldWorkflowInputs = state.IntegrationMethod.IntegrationActionExecutionProperties.WorkflowInputs
			oldReportWorkflowStatus = state.IntegrationMethod.IntegrationActionExecutionProperties.ReportWorkflowStatus
		}
		workflowInputs, err := utils.GoObjectToTerraformStringPreferExisting(oldWorkflowInputs, a.InvocationMethod.IntegrationActionExecutionProperties.WorkflowInputs, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		reportWorkflowStatus, err := utils.GoObjectToTerraformStringPreferExisting(oldReportWorkflowStatus, a.InvocationMethod.IntegrationActionExecutionProperties.ReportWorkflowStatus, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.IntegrationMethod = &IntegrationMethodModel{
			InstallationId:        types.StringValue(*a.InvocationMethod.InstallationId),
			IntegrationActionType: types.StringValue(*a.InvocationMethod.IntegrationActionType),
			IntegrationActionExecutionProperties: &IntegrationActionExecutionPropertiesModel{
				Org:                  types.StringValue(*a.InvocationMethod.IntegrationActionExecutionProperties.Org),
				Repo:                 types.StringValue(*a.InvocationMethod.IntegrationActionExecutionProperties.Repo),
				Workflow:             types.StringValue(*a.InvocationMethod.IntegrationActionExecutionProperties.Workflow),
				WorkflowInputs:       workflowInputs,
				ReportWorkflowStatus: reportWorkflowStatus,
			},
		}
	}

	return nil
}

func writeDatasetToResource(ds *cli.Dataset) *DatasetModel {
	if ds == nil {
		return nil
	}

	datasetModel := &DatasetModel{
		Combinator: types.StringValue(ds.Combinator),
	}

	for _, v := range ds.Rules {
		rule := convertCliRuleToRule(v)
		datasetModel.Rules = append(datasetModel.Rules, rule)
	}

	return datasetModel
}

// convertCliRuleToRule recursively converts a CLI DatasetRule to a Terraform Rule.
// Handles both leaf rules (with operator) and group rules (with combinator + nested rules).
func convertCliRuleToRule(v cli.DatasetRule) Rule {
	rule := Rule{}

	// Check if this is a group rule (has combinator)
	if v.Combinator != nil && *v.Combinator != "" {
		rule.Combinator = types.StringValue(*v.Combinator)
		// Recursively convert nested rules
		if len(v.Rules) > 0 {
			rule.Rules = make([]Rule, 0, len(v.Rules))
			for _, nestedRule := range v.Rules {
				rule.Rules = append(rule.Rules, convertCliRuleToRule(nestedRule))
			}
		}
		// Group rules don't have operator/property/blueprint/value
		rule.Operator = types.StringNull()
		rule.Blueprint = types.StringNull()
		rule.Property = types.StringNull()
	} else {
		// Leaf rule - has operator
		rule.Blueprint = flex.GoStringToFramework(v.Blueprint)
		rule.Property = flex.GoStringToFramework(v.Property)
		rule.Operator = flex.GoStringToFramework(&v.Operator)
		rule.Combinator = types.StringNull()

		if v.Value != nil {
			rule.Value = &Value{
				JqQuery: flex.GoStringToFramework(&v.Value.JqQuery),
			}
		}
	}

	return rule
}

func buildBoolOrJq(prop any) (types.Bool, types.String) {
	if prop == nil {
		return types.BoolNull(), types.StringNull()
	}

	reflectedProp := reflect.ValueOf(prop)
	switch reflectedProp.Kind() {
	case reflect.Bool:
		boolValue := reflectedProp.Interface().(bool)
		return types.BoolValue(boolValue), types.StringNull()
	case reflect.Map:
		jq := reflectedProp.Interface().(map[string]any)
		jqQueryValue := jq["jqQuery"].(string)
		return types.BoolNull(), types.StringValue(jqQueryValue)
	}

	return types.BoolNull(), types.StringNull()
}

func buildRequired(v *cli.ActionUserInputs) (types.String, []string) {
	// If required is nil, return an empty string and nil
	if v.Required == nil {
		return types.StringNull(), nil
	}

	required := reflect.ValueOf(v.Required)
	switch required.Kind() {
	// if required is a slice of strings that means that the user has specified which properties are required
	case reflect.Slice:
		slice := required.Interface().([]interface{})
		attrs := make([]string, 0, required.Len())
		for _, value := range slice {
			attrs = append(attrs, value.(string))
		}
		return types.StringNull(), attrs
	// if required is a map, that means that the user has specified a jq query to determine which properties are required
	case reflect.Map:
		jq := required.Interface().(map[string]any)
		jqQueryValue := jq["jqQuery"].(string)
		return types.StringValue(jqQueryValue), nil
	}

	// if required is not a slice or a map, return an empty string and nil
	return types.StringNull(), nil
}

func (r *ActionResource) buildUserProperties(ctx context.Context, a *cli.Action, state *ActionModel) (*UserPropertiesModel, error) {
	properties := &UserPropertiesModel{}
	if len(a.Trigger.UserInputs.Properties) > 0 {
		requiredJq, required := buildRequired(a.Trigger.UserInputs)
		for k, v := range a.Trigger.UserInputs.Properties {
			switch v.Type {
			case "string":
				if properties.StringProps == nil {
					properties.StringProps = make(map[string]StringPropModel)
				}
				stringProp := addStringPropertiesToResource(ctx, &v)

				if requiredJq.IsNull() && lo.Contains(required, k) {
					stringProp.Required = types.BoolValue(true)
				}

				err := r.setCommonProperties(ctx, v, stringProp)
				if err != nil {
					return nil, err
				}

				properties.StringProps[k] = *stringProp

			case "number":
				if properties.NumberProps == nil {
					properties.NumberProps = make(map[string]NumberPropModel)
				}

				numberProp := addNumberPropertiesToResource(ctx, &v)

				if requiredJq.IsNull() && lo.Contains(required, k) {
					numberProp.Required = types.BoolValue(true)
				}

				err := r.setCommonProperties(ctx, v, numberProp)
				if err != nil {
					return nil, err
				}

				properties.NumberProps[k] = *numberProp

			case "array":
				if properties.ArrayProps == nil {
					properties.ArrayProps = make(map[string]ArrayPropModel)
				}

				var oldArrayProp *ArrayPropModel
				if state.SelfServiceTrigger != nil && state.SelfServiceTrigger.UserProperties != nil && state.SelfServiceTrigger.UserProperties.ArrayProps != nil {
					if old, ok := state.SelfServiceTrigger.UserProperties.ArrayProps[k]; ok {
						oldArrayProp = &old
					}
				}
				arrayProp, err := r.addArrayPropertiesToResource(&v, oldArrayProp)
				if err != nil {
					return nil, err
				}

				if requiredJq.IsNull() && lo.Contains(required, k) {
					arrayProp.Required = types.BoolValue(true)
				}

				err = r.setCommonProperties(ctx, v, arrayProp)
				if err != nil {
					return nil, err
				}

				properties.ArrayProps[k] = *arrayProp

			case "boolean":
				if properties.BooleanProps == nil {
					properties.BooleanProps = make(map[string]BooleanPropModel)
				}

				booleanProp := &BooleanPropModel{}

				err := r.setCommonProperties(ctx, v, booleanProp)
				if err != nil {
					return nil, err
				}

				if requiredJq.IsNull() && lo.Contains(required, k) {
					booleanProp.Required = types.BoolValue(true)
				}

				properties.BooleanProps[k] = *booleanProp

			case "object":
				if properties.ObjectProps == nil {
					properties.ObjectProps = make(map[string]ObjectPropModel)
				}

				objectProp := addObjectPropertiesToResource(&v)

				if requiredJq.IsNull() && lo.Contains(required, k) {
					objectProp.Required = types.BoolValue(true)
				}

				var oldObjectDefault types.String
				if state.SelfServiceTrigger != nil && state.SelfServiceTrigger.UserProperties != nil && state.SelfServiceTrigger.UserProperties.ObjectProps != nil {
					if oldProp, ok := state.SelfServiceTrigger.UserProperties.ObjectProps[k]; ok {
						oldObjectDefault = oldProp.Default
					}
				}

				err := r.setCommonProperties(ctx, v, objectProp)
				if err != nil {
					return nil, err
				}
				err = r.setObjectDefault(objectProp, v.Default, oldObjectDefault)
				if err != nil {
					return nil, err
				}

				properties.ObjectProps[k] = *objectProp

			}
		}
	}
	if properties.StringProps == nil && properties.NumberProps == nil && properties.ArrayProps == nil && properties.BooleanProps == nil && properties.ObjectProps == nil {
		// this logic is handling default initialization of user properties as there is no option to define default user properties in the action schema
		// if there was a state defined for the user properties, return the initiated properties
		if state.SelfServiceTrigger != nil && state.SelfServiceTrigger.UserProperties != nil {
			return properties, nil
		}
		// if there are no user properties defined, return nil
		return nil, nil
	}
	return properties, nil
}

func (r *ActionResource) buildActionTitles(a *cli.Action) (map[string]ActionTitle, error) {
	if a.Trigger.UserInputs.Titles == nil {
		return nil, nil
	}

	actionTitles := make(map[string]ActionTitle)

	for key, actionTitle := range a.Trigger.UserInputs.Titles {
		stateTitle := ActionTitle{
			Title:       types.StringValue(actionTitle.Title),
			Description: flex.GoStringToFramework(actionTitle.Description),
		}

		if actionTitle.Visible != nil {
			visible := reflect.ValueOf(actionTitle.Visible)
			switch visible.Kind() {
			case reflect.Bool:
				boolValue := visible.Interface().(bool)
				stateTitle.Visible = types.BoolValue(boolValue)
			case reflect.Map:
				jq := visible.Interface().(map[string]any)
				jqQueryValue := jq["jqQuery"].(string)
				stateTitle.VisibleJqQuery = types.StringValue(jqQueryValue)
			}
		}

		actionTitles[key] = stateTitle
	}
	return actionTitles, nil
}

func (r *ActionResource) writeTriggerToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if a.Trigger.Type == consts.SelfService {
		userProperties, err := r.buildUserProperties(ctx, a, state)
		if err != nil {
			return err
		}
		actionTitles, err := r.buildActionTitles(a)
		if err != nil {
			return err
		}
		requiredJqQuery, _ := buildRequired(a.Trigger.UserInputs)
		state.SelfServiceTrigger = &SelfServiceTriggerModel{
			BlueprintIdentifier:     flex.GoStringToFramework(a.Trigger.BlueprintIdentifier),
			Operation:               types.StringValue(*a.Trigger.Operation),
			UserProperties:          userProperties,
			RequiredJqQuery:         requiredJqQuery,
			Titles:                  actionTitles,
			ActionCardButtonText:    flex.GoStringToFramework(a.Trigger.ActionCardButtonText),
			ExecuteActionButtonText: flex.GoStringToFramework(a.Trigger.ExecuteActionButtonText),
		}

		if len(a.Trigger.UserInputs.Order) > 0 {
			state.SelfServiceTrigger.OrderProperties = flex.GoArrayStringToTerraformList(ctx, a.Trigger.UserInputs.Order)
		} else {
			state.SelfServiceTrigger.OrderProperties = types.ListNull(types.StringType)
		}

		if len(a.Trigger.UserInputs.Steps) > 0 {
			steps := make([]Step, 0, len(a.Trigger.UserInputs.Steps))
			for _, step := range a.Trigger.UserInputs.Steps {
				t := basetypes.NewStringValue(step.Title)
				o := make([]types.String, 0, len(step.Order))
				for _, p := range step.Order {
					o = append(o, types.StringValue(p))
				}
				s := Step{
					Title: t,
					Order: o,
				}
				visible, visibleJq := buildBoolOrJq(step.Visible)
				if !visible.IsNull() {
					s.Visible = visible
				}
				if !visibleJq.IsNull() {
					s.VisibleJqQuery = visibleJq
				}
				steps = append(steps, s)
			}

			state.SelfServiceTrigger.Steps = steps
		}

		if a.Trigger.Condition != nil {
			var oldCondition types.String
			if state.SelfServiceTrigger != nil {
				oldCondition = state.SelfServiceTrigger.Condition
			}
			triggerCondition, err := utils.GoObjectToTerraformStringPreferExisting(oldCondition, a.Trigger.Condition, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			state.SelfServiceTrigger.Condition = triggerCondition
		}
	}

	if a.Trigger.Type == consts.Automation {
		automationTrigger := &AutomationTriggerModel{}

		if a.Trigger.Condition != nil {
			// Initialize expressions as empty slice to properly handle empty arrays from API
			expressions := make([]types.String, 0, len(a.Trigger.Condition.Expressions))
			for _, e := range a.Trigger.Condition.Expressions {
				expressions = append(expressions, types.StringValue(e))
			}
			automationTrigger.JqCondition = &JqConditionModel{
				Expressions: expressions,
				Combinator:  flex.GoStringToFramework(a.Trigger.Condition.Combinator),
			}
		} else if state.AutomationTrigger != nil && state.AutomationTrigger.JqCondition != nil && len(state.AutomationTrigger.JqCondition.Expressions) == 0 {
			// Preserve jq_condition with empty expressions from prior state
			// When user configures jq_condition with empty expressions, we don't send
			// condition to API (empty expressions = no filtering). On read, API returns
			// no condition, so we preserve the user's intent to have the jq_condition block.
			automationTrigger.JqCondition = &JqConditionModel{
				Expressions: make([]types.String, 0),
				Combinator:  state.AutomationTrigger.JqCondition.Combinator,
			}
		}

		if a.Trigger.Event.Type == consts.EntityCreated {
			automationTrigger.EntityCreatedEvent = &EntityCreatedEventModel{
				BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.EntityUpdated {
			automationTrigger.EntityUpdatedEvent = &EntityUpdatedEventModel{
				BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.EntityDeleted {
			automationTrigger.EntityDeletedEvent = &EntityDeletedEventModel{
				BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.AnyEntityChange {
			automationTrigger.AnyEntityChangeEvent = &AnyEntityChangeEventModel{
				BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.TimerPropertyExpired {
			automationTrigger.TimerPropertyExpiredEvent = &TimerPropertyExpiredEventModel{
				BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
				PropertyIdentifier:  types.StringValue(*a.Trigger.Event.PropertyIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.RunCreated {
			automationTrigger.RunCreatedEvent = &RunCreatedEvent{
				ActionIdentifier: types.StringValue(*a.Trigger.Event.ActionIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.RunUpdated {
			automationTrigger.RunUpdatedEvent = &RunUpdatedEvent{
				ActionIdentifier: types.StringValue(*a.Trigger.Event.ActionIdentifier),
			}
		}

		if a.Trigger.Event.Type == consts.AnyRunChange {
			automationTrigger.AnyRunChangeEvent = &AnyRunChangeEvent{
				ActionIdentifier: types.StringValue(*a.Trigger.Event.ActionIdentifier),
			}
		}

		state.AutomationTrigger = automationTrigger
	}

	return nil
}

func (r *ActionResource) refreshActionState(ctx context.Context, state *ActionModel, a *cli.Action) error {
	state.ID = types.StringValue(a.Identifier)
	state.Identifier = types.StringValue(a.Identifier)
	state.Blueprint = types.StringNull()
	state.Title = flex.GoStringToFramework(a.Title)
	state.Icon = flex.GoStringToFramework(a.Icon)
	state.Description = flex.GoStringToFramework(a.Description)

	err := r.writeTriggerToResource(ctx, a, state)
	if err != nil {
		return err
	}

	err = r.writeInvocationMethodToResource(ctx, a, state)
	if err != nil {
		return err
	}

	if a.RequiredApproval == nil {
		state.RequiredApproval = types.StringNull()
	} else if reflect.TypeOf(a.RequiredApproval).Kind() == reflect.Map {
		state.RequiredApproval = types.StringValue(a.RequiredApproval.(map[string]interface{})["type"].(string))
	} else {
		state.RequiredApproval = types.StringValue(strconv.FormatBool(a.RequiredApproval.(bool)))
	}

	if a.ApprovalNotification != nil {
		switch a.ApprovalNotification.Type {
		case "email":
			state.ApprovalEmailNotification, _ = types.ObjectValue(nil, nil)
		case "none":
			state.ApprovalNoneNotification, _ = types.ObjectValue(nil, nil)
		default:
			state.ApprovalWebhookNotification = &ApprovalWebhookNotificationModel{
				Url: types.StringValue(a.ApprovalNotification.Url),
			}

			if a.ApprovalNotification.Format != nil {
				state.ApprovalWebhookNotification.Format = types.StringValue(*a.ApprovalNotification.Format)
			}
		}
	}
	state.Publish = flex.GoBoolToFramework(a.Publish)
	state.AllowAnyoneToViewRuns = flex.GoBoolToFramework(a.AllowAnyoneToViewRuns)

	return nil
}

func (r *ActionResource) setObjectDefault(objectProp *ObjectPropModel, defaultValue interface{}, oldObjectDefault types.String) error {
	v, ok := defaultValue.(map[string]interface{})
	if !ok {
		return nil
	}
	if v["jqQuery"] != nil {
		objectProp.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
		return nil
	}
	encoded, err := utils.GoObjectToTerraformStringPreferExisting(oldObjectDefault, v, r.portClient.JSONEscapeHTML)
	if err != nil {
		return fmt.Errorf("error converting default value to terraform string: %s", err.Error())
	}
	if encoded.IsNull() {
		objectProp.Default = types.StringNull()
		objectProp.DefaultJqQuery = types.StringNull()
	}
	objectProp.Default = encoded
	return nil
}

func (r *ActionResource) setCommonProperties(ctx context.Context, v cli.ActionProperty, prop interface{}) error {
	properties := []string{"Description", "Icon", "Default", "Title", "DependsOn", "Dataset", "Visible", "Disabled"}
	for _, property := range properties {
		switch property {
		case "Description":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *NumberPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *BooleanPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ArrayPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ObjectPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			}
		case "Icon":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *NumberPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *BooleanPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ArrayPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ObjectPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			}
		case "Title":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *NumberPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *BooleanPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ArrayPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ObjectPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			}
		// Due to the possibility of an error being raised when converting null to a pointer, we are unable to utilize flex in this scenario.
		case "Default":
			switch p := prop.(type) {
			case *StringPropModel:
				if v.Default == nil {
					p.Default = types.StringNull()
					p.DefaultJqQuery = types.StringNull()
				} else {
					switch v := v.Default.(type) {
					case string:
						p.Default = types.StringValue(v)
					case map[string]interface{}:
						p.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
					}
				}
			case *NumberPropModel:
				if v.Default == nil {
					p.Default = types.Float64Null()
					p.DefaultJqQuery = types.StringNull()
				} else {
					switch v := v.Default.(type) {
					case float64:
						p.Default = types.Float64Value(v)
					case map[string]interface{}:
						p.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
					}
				}
			case *BooleanPropModel:
				if v.Default == nil {
					p.Default = types.BoolNull()
					p.DefaultJqQuery = types.StringNull()
				} else {
					switch v := v.Default.(type) {
					case bool:
						p.Default = types.BoolValue(v)
					case map[string]interface{}:
						p.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
					}
				}
			}
		case "DependsOn":
			switch p := prop.(type) {
			case *StringPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *NumberPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *BooleanPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *ArrayPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *ObjectPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			}

		case "Visible":
			visible, visibleJq := buildBoolOrJq(v.Visible)
			if !visible.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Visible = visible
				case *NumberPropModel:
					p.Visible = visible
				case *BooleanPropModel:
					p.Visible = visible
				case *ArrayPropModel:
					p.Visible = visible
				case *ObjectPropModel:
					p.Visible = visible
				}
			}
			if !visibleJq.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.VisibleJqQuery = visibleJq
				case *NumberPropModel:
					p.VisibleJqQuery = visibleJq
				case *BooleanPropModel:
					p.VisibleJqQuery = visibleJq
				case *ArrayPropModel:
					p.VisibleJqQuery = visibleJq
				case *ObjectPropModel:
					p.VisibleJqQuery = visibleJq
				}
			}

		case "Disabled":
			disabled, disabledJq := buildBoolOrJq(v.Disabled)
			if !disabled.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Disabled = disabled
				case *NumberPropModel:
					p.Disabled = disabled
				case *BooleanPropModel:
					p.Disabled = disabled
				case *ArrayPropModel:
					p.Disabled = disabled
				case *ObjectPropModel:
					p.Disabled = disabled
				}
			}
			if !disabledJq.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.DisabledJqQuery = disabledJq
				case *NumberPropModel:
					p.DisabledJqQuery = disabledJq
				case *BooleanPropModel:
					p.DisabledJqQuery = disabledJq
				case *ArrayPropModel:
					p.DisabledJqQuery = disabledJq
				case *ObjectPropModel:
					p.DisabledJqQuery = disabledJq
				}
			}
		}
	}
	return nil
}
