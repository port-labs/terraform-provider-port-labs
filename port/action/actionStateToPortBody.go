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
		Identifier: data.Identifier.ValueString(),
		Title:      data.Title.ValueString(),
		Trigger:    data.Trigger.ValueString(),
	}

	if !data.Icon.IsNull() {
		icon := data.Icon.ValueString()
		action.Icon = &icon
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		action.Description = &description
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
		org := data.GithubMethod.Org.ValueString()
		repo := data.GithubMethod.Repo.ValueString()
		workflow := data.GithubMethod.Workflow.ValueString()
		githubInvocation := &cli.InvocationMethod{
			Type:     consts.Github,
			Org:      &org,
			Repo:     &repo,
			Workflow: &workflow,
		}

		if !data.GithubMethod.OmitPayload.IsNull() {
			omitPayload := data.GithubMethod.OmitPayload.ValueBool()
			githubInvocation.OmitPayload = &omitPayload
		}

		if !data.GithubMethod.OmitUserInputs.IsNull() {
			omitUserInputs := data.GithubMethod.OmitUserInputs.ValueBool()
			githubInvocation.OmitUserInputs = &omitUserInputs
		}

		if !data.GithubMethod.ReportWorkflowStatus.IsNull() {
			reportWorkflowStatus := data.GithubMethod.ReportWorkflowStatus.ValueBool()
			githubInvocation.ReportWorkflowStatus = &reportWorkflowStatus
		}
		return githubInvocation
	}

	if !data.KafkaMethod.IsNull() {
		return &cli.InvocationMethod{
			Type: consts.Kafka,
		}
	}

	if data.WebhookMethod != nil {
		url := data.WebhookMethod.Url.ValueString()
		webhookInvocation := &cli.InvocationMethod{
			Type: consts.Webhook,
			Url:  &url,
		}
		if !data.WebhookMethod.Agent.IsNull() {
			agent := data.WebhookMethod.Agent.ValueBool()
			webhookInvocation.Agent = &agent
		}
		if !data.WebhookMethod.Synchronized.IsNull() {
			synchronized := data.WebhookMethod.Synchronized.ValueBool()
			webhookInvocation.Synchronized = &synchronized
		}
		if !data.WebhookMethod.Method.IsNull() {
			method := data.WebhookMethod.Method.ValueString()
			webhookInvocation.Method = &method
		}

		return webhookInvocation
	}

	if data.GitlabMethod != nil {
		projectName := data.GitlabMethod.ProjectName.ValueString()
		groupName := data.GitlabMethod.GroupName.ValueString()
		gitlabInvocation := &cli.InvocationMethod{
			Type:        consts.Gitlab,
			ProjectName: &projectName,
			GroupName:   &groupName,
		}

		if !data.GitlabMethod.OmitPayload.IsNull() {
			omitPayload := data.GitlabMethod.OmitPayload.ValueBool()
			gitlabInvocation.OmitPayload = &omitPayload
		}

		if !data.GitlabMethod.OmitUserInputs.IsNull() {
			omitUserInputs := data.GitlabMethod.OmitUserInputs.ValueBool()
			gitlabInvocation.OmitUserInputs = &omitUserInputs
		}

		if !data.GitlabMethod.DefaultRef.IsNull() {
			defaultRef := data.GitlabMethod.DefaultRef.ValueString()
			gitlabInvocation.DefaultRef = &defaultRef
		}

		if !data.GitlabMethod.Agent.IsNull() {
			agent := data.GitlabMethod.Agent.ValueBool()
			gitlabInvocation.Agent = &agent
		}

		return gitlabInvocation
	}

	return nil
}
func actionPermissionsToPortBody(state *PermissionsModel) *cli.ActionPermissions {
	if state == nil {
		return nil
	}

	return &cli.ActionPermissions{
		Execute: cli.ActionExecutePermissions{
			Users:       flex.TerraformStringListToGoArray(state.Execute.Users),
			Roles:       flex.TerraformStringListToGoArray(state.Execute.Roles),
			Teams:       flex.TerraformStringListToGoArray(state.Execute.Teams),
			OwnedByTeam: state.Execute.OwnedByTeam.ValueBoolPointer(),
		},
		Approve: cli.ActionApprovePermissions{
			Users: flex.TerraformStringListToGoArray(state.Approve.Users),
			Roles: flex.TerraformStringListToGoArray(state.Approve.Roles),
			Teams: flex.TerraformStringListToGoArray(state.Approve.Teams),
		},
	}
}
