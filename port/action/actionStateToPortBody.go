package action

import (
	"context"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func actionDataSetToPortBody(dataSet *DatasetModel) *cli.Dataset {
	cliDateSet := &cli.Dataset{
		Combinator: dataSet.Combinator.ValueString(),
	}
	rules := make([]cli.DatasetRule, 0, len(dataSet.Rules))
	for _, rule := range dataSet.Rules {
		dataSetRule := cli.DatasetRule{
			Operator: rule.Operator.ValueString(),
			Value: &cli.DatasetValue{
				JqQuery: rule.Value.JqQuery.ValueString(),
			},
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
		Identifier: data.Identifier.ValueString(),
	}

	if !data.Title.IsNull() {
		title := data.Title.ValueString()
		action.Title = &title
	}

	if !data.Icon.IsNull() {
		icon := data.Icon.ValueString()
		action.Icon = &icon
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		action.Description = &description
	}

	action.Trigger, err = triggerToBody(ctx, data)
	if err != nil {
		return nil, err
	}

	action.InvocationMethod, err = invocationMethodToBody(ctx, data)
	if err != nil {
		return nil, err
	}

	if !data.RequiredApproval.IsNull() {
		requiredApproval := data.RequiredApproval.ValueBool()
		action.RequiredApproval = &requiredApproval
	}

	if !data.ApprovalEmailNotification.IsNull() {
		action.ApprovalNotification = &cli.ApprovalNotification{
			Type: "email",
		}
	}
	if data.ApprovalWebhookNotification != nil {
		action.ApprovalNotification = &cli.ApprovalNotification{
			Type: "webhook",
			Url:  data.ApprovalWebhookNotification.Url.ValueString(),
		}

		if !data.ApprovalWebhookNotification.Format.IsNull() {
			format := data.ApprovalWebhookNotification.Format.ValueString()
			action.ApprovalNotification.Format = &format
		}
	}

	if !data.Publish.IsNull() {
		publish := data.Publish.ValueBool()
		action.Publish = &publish
	}

	return action, nil
}

func triggerToBody(ctx context.Context, data *ActionModel) (*cli.Trigger, error) {
	if data.SelfServiceTrigger != nil {
		selfServiceTrigger := &cli.Trigger{
			Type:                consts.SelfService,
			BlueprintIdentifier: data.SelfServiceTrigger.BlueprintIdentifier.ValueStringPointer(),
			Operation:           data.SelfServiceTrigger.Operation.ValueStringPointer(),
		}

		if data.SelfServiceTrigger.UserProperties != nil {
			err := actionPropertiesToBody(ctx, selfServiceTrigger, data.SelfServiceTrigger)
			if err != nil {
				return nil, err
			}
		} else {
			selfServiceTrigger.UserInputs.Properties = make(map[string]cli.ActionProperty)
		}

		if !data.SelfServiceTrigger.OrderProperties.IsNull() {
			order, err := utils.TerraformListToGoArray(ctx, data.SelfServiceTrigger.OrderProperties, "string")
			if err != nil {
				return nil, err
			}
			orderString := utils.InterfaceToStringArray(order)
			selfServiceTrigger.UserInputs.Order = orderString
		}

		return selfServiceTrigger, nil
	}

	if data.AutomationTrigger != nil {
		automationTrigger := &cli.Trigger{
			Type: consts.Automation,
			Condition: &cli.TriggerCondition{
				Expressions: flex.TerraformStringListToGoArray(data.AutomationTrigger.JqCondition.Expressions),
				Combinator:  data.AutomationTrigger.JqCondition.Combinator.ValueStringPointer(),
			},
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

		return automationTrigger, nil
	}

	return nil, nil
}

func actionPropertiesToBody(ctx context.Context, action *cli.Trigger, data *SelfServiceTriggerModel) error {
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

	action.UserInputs.Properties = props

	// if requiredJqQuery is set, required shouldn't be set and vice versa
	if !data.RequiredJqQuery.IsNull() {
		RequiredJqQueryMap := map[string]string{
			"jqQuery": data.RequiredJqQuery.ValueString(),
		}
		action.UserInputs.Required = RequiredJqQueryMap
	} else {
		action.UserInputs.Required = required
	}

	return nil
}

func invocationMethodToBody(ctx context.Context, data *ActionModel) (*cli.InvocationMethod, error) {
	if data.KafkaMethod != nil {
		payload, err := utils.TerraformStringToGoObject(data.KafkaMethod.Payload)
		if err != nil {
			return nil, err
		}

		return &cli.InvocationMethod{Type: consts.Kafka, Payload: payload}, nil
	}

	if data.WebhookMethod != nil {
		agent, err := utils.TerraformStringToGoObject(data.WebhookMethod.Agent)
		if err != nil {
			return nil, err
		}
		synchronized, err := utils.TerraformStringToGoObject(data.WebhookMethod.Synchronized)
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
		body, err := utils.TerraformStringToGoObject(data.WebhookMethod.Body)
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
		reportWorkflowStatus, err := utils.TerraformStringToGoObject(data.GithubMethod.ReportWorkflowStatus)
		if err != nil {
			return nil, err
		}
		workflowInputs, err := utils.TerraformStringToGoObject(data.GithubMethod.WorkflowInputs)
		if err != nil {
			return nil, err
		}

		githubInvocation := &cli.InvocationMethod{
			Type:                 consts.Github,
			Org:                  data.GithubMethod.Org.ValueStringPointer(),
			Repo:                 data.GithubMethod.Repo.ValueStringPointer(),
			Workflow:             data.GithubMethod.Workflow.ValueStringPointer(),
			WorkflowInputs:       workflowInputs.(map[string]interface{}),
			ReportWorkflowStatus: reportWorkflowStatus,
		}

		return githubInvocation, nil
	}

	if data.GitlabMethod != nil {
		pipelineVariables, err := utils.TerraformStringToGoObject(data.GitlabMethod.PipelineVariables)
		if err != nil {
			return nil, err
		}

		gitlabInvocation := &cli.InvocationMethod{
			Type:              consts.Gitlab,
			ProjectName:       data.GitlabMethod.ProjectName.ValueStringPointer(),
			GroupName:         data.GitlabMethod.GroupName.ValueStringPointer(),
			DefaultRef:        data.GitlabMethod.DefaultRef.ValueStringPointer(),
			PipelineVariables: pipelineVariables.(map[string]interface{}),
		}

		return gitlabInvocation, nil
	}

	if data.AzureMethod != nil {
		payload, err := utils.TerraformStringToGoObject(data.AzureMethod.Payload)
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
		upsertEntityInvocation := &cli.InvocationMethod{
			Type:                consts.UpsertEntity,
			Identifier:          data.UpsertEntityMethod.Identifier.ValueStringPointer(),
			Title:               data.UpsertEntityMethod.Title.ValueStringPointer(),
			BlueprintIdentifier: data.UpsertEntityMethod.BlueprintIdentifier.ValueStringPointer(),
			Team:                flex.TerraformStringListToGoArray(data.UpsertEntityMethod.Teams),
			Icon:                data.UpsertEntityMethod.Icon.ValueStringPointer(),
		}

		return upsertEntityInvocation, nil
	}

	return nil, nil
}
