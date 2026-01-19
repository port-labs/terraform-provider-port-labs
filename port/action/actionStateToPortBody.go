package action

import (
	"context"
	"reflect"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func actionDataSetToPortBody(dataSet *DatasetModel) *cli.Dataset {
	cliDateSet := &cli.Dataset{
		Combinator: dataSet.Combinator.ValueString(),
	}
	rules := make([]cli.DatasetRule, 0, len(dataSet.Rules))
	for _, rule := range dataSet.Rules {
		dataSetRule := cli.DatasetRule{
			Operator: rule.Operator.ValueString(),
		}

		if rule.Value != nil && !rule.Value.JqQuery.IsNull() {
			dataSetRule.Value = &cli.DatasetValue{
				JqQuery: rule.Value.JqQuery.ValueString(),
			}
		}

		if !rule.Blueprint.IsNull() {
			blueprint := rule.Blueprint.ValueString()
			dataSetRule.Blueprint = &blueprint
		}
		if !rule.Property.IsNull() {
			rule := rule.Property.ValueString()
			dataSetRule.Property = &rule
		}

		rules = append(rules, dataSetRule)
	}
	cliDateSet.Rules = rules
	return cliDateSet
}

func actionStateToPortBody(ctx context.Context, data *ActionModel) (*cli.Action, error) {
	var err error
	action := &cli.Action{
		Identifier:  data.Identifier.ValueString(),
		Title:       data.Title.ValueStringPointer(),
		Icon:        data.Icon.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Publish:     data.Publish.ValueBoolPointer(),
	}

	action.RequiredApproval = utils.TerraformStringToBooleanOrString(data.RequiredApproval)
	if action.RequiredApproval != nil && reflect.TypeOf(action.RequiredApproval).Kind() == reflect.String {
		action.RequiredApproval = map[string]interface{}{"type": action.RequiredApproval}
	}

	action.Trigger, err = triggerToBody(ctx, data)
	if err != nil {
		return nil, err
	}

	action.InvocationMethod, err = invocationMethodToBody(ctx, data)
	if err != nil {
		return nil, err
	}

	if !data.ApprovalEmailNotification.IsNull() {
		action.ApprovalNotification = &cli.ApprovalNotification{
			Type: "email",
		}
	}
	if data.ApprovalWebhookNotification != nil {
		action.ApprovalNotification = &cli.ApprovalNotification{
			Type:   "webhook",
			Url:    data.ApprovalWebhookNotification.Url.ValueString(),
			Format: data.ApprovalWebhookNotification.Format.ValueStringPointer(),
		}
	}

	if !data.AllowAnyoneToViewRuns.IsNull() {
		action.AllowAnyoneToViewRuns = data.AllowAnyoneToViewRuns.ValueBoolPointer()
	}

	return action, nil
}

func triggerToBody(ctx context.Context, data *ActionModel) (*cli.Trigger, error) {
	if data.SelfServiceTrigger != nil {
		selfServiceTrigger := &cli.Trigger{
			Type:                consts.SelfService,
			BlueprintIdentifier: data.SelfServiceTrigger.BlueprintIdentifier.ValueStringPointer(),
			Operation:           data.SelfServiceTrigger.Operation.ValueStringPointer(),
			UserInputs: &cli.ActionUserInputs{
				Properties: make(map[string]cli.ActionProperty),
			},
		}

		if data.SelfServiceTrigger.UserProperties != nil {
			err := actionPropertiesToBody(ctx, selfServiceTrigger, data.SelfServiceTrigger)
			if err != nil {
				return nil, err
			}
		}

		if data.SelfServiceTrigger.Titles != nil {
			err := actionTitlesToBody(ctx, selfServiceTrigger, data.SelfServiceTrigger)
			if err != nil {
				return nil, err
			}
		}

		if !data.SelfServiceTrigger.RequiredJqQuery.IsNull() {
			selfServiceTrigger.UserInputs.Required = map[string]string{
				"jqQuery": data.SelfServiceTrigger.RequiredJqQuery.ValueString(),
			}
		}

		if !data.SelfServiceTrigger.OrderProperties.IsNull() {
			order, err := utils.TerraformListToGoArray(ctx, data.SelfServiceTrigger.OrderProperties, "string")
			if err != nil {
				return nil, err
			}
			orderString := utils.InterfaceToStringArray(order)
			selfServiceTrigger.UserInputs.Order = orderString
		}

		if data.SelfServiceTrigger.Steps != nil {
			steps := make([]cli.Step, 0, len(data.SelfServiceTrigger.Steps))

			for _, s := range data.SelfServiceTrigger.Steps {
				o := make([]string, 0, len(s.Order))
				for _, p := range s.Order {
					o = append(o, p.ValueString())
				}
				stepObj := cli.Step{
					Title: s.Title.ValueString(),
					Order: o,
				}

				if !s.VisibleJqQuery.IsNull() {
					stepObj.Visible = map[string]string{
						"jqQuery": s.VisibleJqQuery.ValueString(),
					}
				} else if !s.Visible.IsNull() {
					stepObj.Visible = s.Visible.ValueBool()
				}

				steps = append(steps, stepObj)
			}

			selfServiceTrigger.UserInputs.Steps = steps
		}

		if !data.SelfServiceTrigger.Condition.IsNull() {
			condition, err := utils.TerraformStringToGoType[cli.TriggerCondition](data.SelfServiceTrigger.Condition)
			if err != nil {
				return nil, err
			}
			selfServiceTrigger.Condition = &condition
		}

		if !data.SelfServiceTrigger.ActionCardButtonText.IsNull() {
			selfServiceTrigger.ActionCardButtonText = data.SelfServiceTrigger.ActionCardButtonText.ValueStringPointer()
		}

		if !data.SelfServiceTrigger.ExecuteActionButtonText.IsNull() {
			selfServiceTrigger.ExecuteActionButtonText = data.SelfServiceTrigger.ExecuteActionButtonText.ValueStringPointer()
		}

		return selfServiceTrigger, nil
	}

	if data.AutomationTrigger != nil {
		automationTrigger := &cli.Trigger{
			Type: consts.Automation,
		}

		if data.AutomationTrigger.JqCondition != nil {
			automationTrigger.Condition = &cli.TriggerCondition{
				Type:        consts.JqCondition,
				Expressions: flex.TerraformStringListToGoArray(data.AutomationTrigger.JqCondition.Expressions),
				Combinator:  data.AutomationTrigger.JqCondition.Combinator.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.EntityCreatedEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:                consts.EntityCreated,
				BlueprintIdentifier: data.AutomationTrigger.EntityCreatedEvent.BlueprintIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.EntityUpdatedEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:                consts.EntityUpdated,
				BlueprintIdentifier: data.AutomationTrigger.EntityUpdatedEvent.BlueprintIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.EntityDeletedEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:                consts.EntityDeleted,
				BlueprintIdentifier: data.AutomationTrigger.EntityDeletedEvent.BlueprintIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.AnyEntityChangeEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:                consts.AnyEntityChange,
				BlueprintIdentifier: data.AutomationTrigger.AnyEntityChangeEvent.BlueprintIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.TimerPropertyExpiredEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:                consts.TimerPropertyExpired,
				BlueprintIdentifier: data.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.ValueStringPointer(),
				PropertyIdentifier:  data.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.RunCreatedEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:             consts.RunCreated,
				ActionIdentifier: data.AutomationTrigger.RunCreatedEvent.ActionIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.RunUpdatedEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:             consts.RunUpdated,
				ActionIdentifier: data.AutomationTrigger.RunUpdatedEvent.ActionIdentifier.ValueStringPointer(),
			}
		}

		if data.AutomationTrigger.AnyRunChangeEvent != nil {
			automationTrigger.Event = &cli.TriggerEvent{
				Type:             consts.AnyRunChange,
				ActionIdentifier: data.AutomationTrigger.AnyRunChangeEvent.ActionIdentifier.ValueStringPointer(),
			}
		}

		return automationTrigger, nil
	}

	return nil, nil
}

func actionPropertiesToBody(ctx context.Context, actionTrigger *cli.Trigger, data *SelfServiceTriggerModel) error {
	required := []string{}
	props := map[string]cli.ActionProperty{}
	var err error
	if data.UserProperties.StringProps != nil {
		err = stringPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.ArrayProps != nil {
		err = arrayPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.NumberProps != nil {
		err = numberPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.BooleanProps != nil {
		err = booleanPropResourceToBody(ctx, data, props, &required)
	}

	if data.UserProperties.ObjectProps != nil {
		err = objectPropResourceToBody(ctx, data, props, &required)
	}

	if err != nil {
		return err
	}

	actionTrigger.UserInputs.Properties = props

	// if requiredJqQuery is set, required shouldn't be set and vice versa
	if !data.RequiredJqQuery.IsNull() {
		RequiredJqQueryMap := map[string]string{
			"jqQuery": data.RequiredJqQuery.ValueString(),
		}
		actionTrigger.UserInputs.Required = RequiredJqQueryMap
	} else if len(required) > 0 {
		actionTrigger.UserInputs.Required = required
	}

	return nil
}

func actionTitlesToBody(ctx context.Context, actionTrigger *cli.Trigger, data *SelfServiceTriggerModel) error {
	actionTitles := map[string]cli.ActionTitle{}
	var err error
	if data.Titles != nil {
		for key, actionTitle := range data.Titles {

			cliTitle := cli.ActionTitle{
				Title:       actionTitle.Title.ValueString(),
				Description: actionTitle.Description.ValueStringPointer(),
			}

			if !actionTitle.Visible.IsNull() {
				cliTitle.Visible = actionTitle.Visible.ValueBool()
			}

			if !actionTitle.VisibleJqQuery.IsNull() {
				VisibleJqQueryMap := map[string]string{
					"jqQuery": actionTitle.VisibleJqQuery.ValueString(),
				}
				cliTitle.Visible = VisibleJqQueryMap
			}

			actionTitles[key] = cliTitle
		}
	}

	if err != nil {
		return err
	}

	actionTrigger.UserInputs.Titles = actionTitles

	return nil
}

func invocationMethodToBody(ctx context.Context, data *ActionModel) (*cli.InvocationMethod, error) {
	if data.KafkaMethod != nil {
		payload, err := utils.TerraformStringToGoType[interface{}](data.KafkaMethod.Payload)
		if err != nil {
			return nil, err
		}

		return &cli.InvocationMethod{Type: consts.Kafka, Payload: payload}, nil
	}

	if data.WebhookMethod != nil {
		agent, err := utils.TerraformStringToGoType[interface{}](data.WebhookMethod.Agent)
		if err != nil {
			return nil, err
		}
		synchronized, err := utils.TerraformStringToGoType[interface{}](data.WebhookMethod.Synchronized)
		if err != nil {
			return nil, err
		}
		headers := make(map[string]string)
		for key, value := range data.WebhookMethod.Headers.Elements() {
			tv, _ := value.ToTerraformValue(ctx)
			var keyValue string
			err = tv.As(&keyValue)
			if err != nil {
				return nil, err
			}
			headers[key] = keyValue
		}
		body, err := utils.TerraformStringToGoType[interface{}](data.WebhookMethod.Body)
		if err != nil {
			return nil, err
		}

		webhookInvocation := &cli.InvocationMethod{
			Type:         consts.Webhook,
			Url:          data.WebhookMethod.Url.ValueStringPointer(),
			Agent:        agent,
			Synchronized: synchronized,
			Method:       data.WebhookMethod.Method.ValueStringPointer(),
			Headers:      headers,
			Body:         body,
		}

		return webhookInvocation, nil
	}

	if data.GithubMethod != nil {
		reportWorkflowStatus, err := utils.TerraformStringToGoType[interface{}](data.GithubMethod.ReportWorkflowStatus)
		if err != nil {
			return nil, err
		}
		wi, err := utils.TerraformStringToGoType[interface{}](data.GithubMethod.WorkflowInputs)
		if err != nil {
			return nil, err
		}
		workflowInputs, _ := wi.(map[string]interface{})

		githubInvocation := &cli.InvocationMethod{
			Type:                 consts.Github,
			Org:                  data.GithubMethod.Org.ValueStringPointer(),
			Repo:                 data.GithubMethod.Repo.ValueStringPointer(),
			Workflow:             data.GithubMethod.Workflow.ValueStringPointer(),
			WorkflowInputs:       workflowInputs,
			ReportWorkflowStatus: reportWorkflowStatus,
		}

		return githubInvocation, nil
	}

	if data.GitlabMethod != nil {
		pv, err := utils.TerraformStringToGoType[interface{}](data.GitlabMethod.PipelineVariables)
		if err != nil {
			return nil, err
		}
		pipelineVariables, _ := pv.(map[string]interface{})

		gitlabInvocation := &cli.InvocationMethod{
			Type:              consts.Gitlab,
			ProjectName:       data.GitlabMethod.ProjectName.ValueStringPointer(),
			GroupName:         data.GitlabMethod.GroupName.ValueStringPointer(),
			DefaultRef:        data.GitlabMethod.DefaultRef.ValueStringPointer(),
			PipelineVariables: pipelineVariables,
		}

		return gitlabInvocation, nil
	}

	if data.AzureMethod != nil {
		payload, err := utils.TerraformStringToGoType[interface{}](data.AzureMethod.Payload)
		if err != nil {
			return nil, err
		}

		azureInvocation := &cli.InvocationMethod{
			Type:    consts.AzureDevops,
			Org:     data.AzureMethod.Org.ValueStringPointer(),
			Webhook: data.AzureMethod.Webhook.ValueStringPointer(),
			Payload: payload,
		}

		return azureInvocation, nil
	}

	if data.UpsertEntityMethod != nil {
		var mapping cli.MappingSchema
		if data.UpsertEntityMethod.Mapping != nil {
			var team interface{}
			if data.UpsertEntityMethod.Mapping.Teams != nil {
				team = flex.TerraformStringListToGoArray(data.UpsertEntityMethod.Mapping.Teams)
			}
			if !data.UpsertEntityMethod.Mapping.TeamsJQ.IsNull() {
				team = data.UpsertEntityMethod.Mapping.TeamsJQ.ValueString()
			}
			properties, err := utils.TerraformJsonStringToGoObject(data.UpsertEntityMethod.Mapping.Properties.ValueStringPointer())
			if err != nil {
				return nil, err
			}

			relations, err := utils.TerraformJsonStringToGoObject(data.UpsertEntityMethod.Mapping.Relations.ValueStringPointer())
			if err != nil {
				return nil, err
			}

			mapping = cli.MappingSchema{
				Team:       team,
				Identifier: data.UpsertEntityMethod.Mapping.Identifier.ValueStringPointer(),
				Title:      data.UpsertEntityMethod.Title.ValueStringPointer(),
				Icon:       data.UpsertEntityMethod.Mapping.Icon.ValueStringPointer(),
			}

			if properties != nil {
				mapping.Properties = *properties
			}

			if relations != nil {
				mapping.Relations = *relations
			}
		}

		upsertEntityInvocation := &cli.InvocationMethod{
			Type:                consts.UpsertEntity,
			BlueprintIdentifier: data.UpsertEntityMethod.BlueprintIdentifier.ValueStringPointer(),
			Mapping:             &mapping,
		}

		return upsertEntityInvocation, nil
	}

	return nil, nil
}
