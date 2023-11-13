package action

import (
	"context"

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

func actionStateToPortBody(ctx context.Context, data *ActionModel, bp *cli.Blueprint) (*cli.Action, error) {
	action := &cli.Action{
		Identifier:       data.Identifier.ValueString(),
		Title:            data.Title.ValueString(),
		Trigger:          data.Trigger.ValueString(),
		Icon:             data.Icon.ValueStringPointer(),
		Description:      data.Description.ValueStringPointer(),
		RequiredApproval: data.RequiredApproval.ValueBoolPointer(),
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

	action.InvocationMethod = invocationMethodToBody(data)

	if data.UserProperties != nil {
		err := actionPropertiesToBody(ctx, action, data)
		if err != nil {
			return nil, err
		}
	} else {
		action.UserInputs.Properties = make(map[string]cli.ActionProperty)
	}

	if !data.OrderProperties.IsNull() {
		order, err := utils.TerraformListToGoArray(ctx, data.OrderProperties, "string")
		if err != nil {
			return nil, err
		}
		orderString := utils.InterfaceToStringArray(order)
		action.UserInputs.Order = orderString
	}

	return action, nil
}

func actionPropertiesToBody(ctx context.Context, action *cli.Action, data *ActionModel) error {
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
	action.UserInputs.Required = required

	return nil
}
func invocationMethodToBody(data *ActionModel) *cli.InvocationMethod {
	if data.AzureMethod != nil {
		org := data.AzureMethod.Org.ValueString()
		webhook := data.AzureMethod.Webhook.ValueString()
		return &cli.InvocationMethod{
			Type:    consts.AzureDevops,
			Org:     &org,
			Webhook: &webhook,
		}
	}

	if data.GithubMethod != nil {
		return &cli.InvocationMethod{
			Type:                 consts.Github,
			Org:                  data.GithubMethod.Org.ValueStringPointer(),
			Repo:                 data.GithubMethod.Repo.ValueStringPointer(),
			Workflow:             data.GithubMethod.Workflow.ValueStringPointer(),
			OmitPayload:          data.GithubMethod.OmitPayload.ValueBoolPointer(),
			OmitUserInputs:       data.GithubMethod.OmitUserInputs.ValueBoolPointer(),
			ReportWorkflowStatus: data.GithubMethod.ReportWorkflowStatus.ValueBoolPointer(),
		}
	}

	if !data.KafkaMethod.IsNull() {
		return &cli.InvocationMethod{
			Type: consts.Kafka,
		}
	}

	if data.WebhookMethod != nil {
		return &cli.InvocationMethod{
			Type:         consts.Webhook,
			Url:          data.WebhookMethod.Url.ValueStringPointer(),
			Agent:        data.WebhookMethod.Agent.ValueBoolPointer(),
			Synchronized: data.WebhookMethod.Synchronized.ValueBoolPointer(),
			Method:       data.WebhookMethod.Method.ValueStringPointer(),
		}
	}

	if data.GitlabMethod != nil {
		return &cli.InvocationMethod{
			Type:           consts.Gitlab,
			ProjectName:    data.GitlabMethod.ProjectName.ValueStringPointer(),
			GroupName:      data.GitlabMethod.GroupName.ValueStringPointer(),
			OmitPayload:    data.GitlabMethod.OmitPayload.ValueBoolPointer(),
			OmitUserInputs: data.GitlabMethod.OmitUserInputs.ValueBoolPointer(),
			DefaultRef:     data.GitlabMethod.DefaultRef.ValueStringPointer(),
			Agent:          data.GitlabMethod.Agent.ValueBoolPointer(),
		}
	}

	return nil
}
