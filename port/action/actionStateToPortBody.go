package action

import (
	"context"
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

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
		action.UserInputs.Properties = make(map[string]cli.BlueprintProperty)
	}

	return action, nil
}

func stringPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.StringProp {
		property := cli.BlueprintProperty{
			Type: "string",
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			property.Title = &title
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
		}

		if !prop.Format.IsNull() {
			format := prop.Format.ValueString()
			property.Format = &format
		}

		if !prop.Blueprint.IsNull() {
			blueprint := prop.Blueprint.ValueString()
			property.Blueprint = &blueprint
		}

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			property.Icon = &icon
		}

		if !prop.MinLength.IsNull() {
			minLength := int(prop.MinLength.ValueInt64())
			property.MinLength = &minLength
		}

		if !prop.MaxLength.IsNull() {
			maxLength := int(prop.MaxLength.ValueInt64())
			property.MaxLength = &maxLength
		}

		if !prop.Pattern.IsNull() {
			pattern := prop.Pattern.ValueString()
			property.Pattern = &pattern
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			property.Description = &description
		}

		if !prop.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "string")
			if err != nil {
				return err
			}

			property.Enum = enumList
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func numberPropResourceToBody(ctx context.Context, state *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.UserProperties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "number",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueFloat64()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Minimum.IsNull() {
				minimum := prop.Minimum.ValueFloat64()
				property.Minimum = &minimum
			}

			if !prop.Maximum.IsNull() {
				maximum := prop.Maximum.ValueFloat64()
				property.Maximum = &maximum
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Enum.IsNull() {

				enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "float64")
				if err != nil {
					return err
				}
				property.Enum = enumList
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func booleanPropResourceToBody(d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.UserProperties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "boolean",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueBool()
			}

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func objectPropResourceToBody(d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ObjectProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "object",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				defaultAsString := prop.Default.ValueString()
				defaultObj := make(map[string]interface{})
				err := json.Unmarshal([]byte(defaultAsString), &defaultObj)
				if err != nil {
					return err
				} else {
					property.Default = defaultObj
				}
			}

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Spec.IsNull() {
				spec := prop.Spec.ValueString()
				property.Spec = &spec
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func handleArrayItemsToBody(ctx context.Context, property *cli.BlueprintProperty, prop ArrayPropModel, required *[]string) error {
	if prop.StringItems != nil {
		items := map[string]interface{}{}
		items["type"] = "string"
		if !prop.StringItems.Format.IsNull() {
			items["format"] = prop.StringItems.Format.ValueString()
		}

		if !prop.StringItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "string")
			if err != nil {
				return err
			}

			property.Default = defaultList
		}
		property.Items = items
	}

	if prop.NumberItems != nil {
		items := map[string]interface{}{}
		items["type"] = "number"
		if !prop.NumberItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "float64")
			if err != nil {
				return err
			}

			items["default"] = defaultList
		}
		property.Items = items
	}

	if prop.BooleanItems != nil {
		items := map[string]interface{}{}
		items["type"] = "boolean"
		if !prop.BooleanItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "bool")
			if err != nil {
				return err
			}

			items["default"] = defaultList
		}
		property.Items = items
	}

	if prop.ObjectItems != nil {
		items := map[string]interface{}{}
		items["type"] = "object"
		if !prop.ObjectItems.Default.IsNull() {
			defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "object")
			if err != nil {
				return err
			}
			items["default"] = defaultList
		}
		property.Items = items
	}
	return nil
}
func arrayPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ArrayProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "array",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}
			if !prop.MinItems.IsNull() {
				minItems := int(prop.MinItems.ValueInt64())
				property.MinItems = &minItems
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = &maxItems
			}

			err := handleArrayItemsToBody(ctx, &property, prop, required)
			if err != nil {
				return err
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func actionPropertiesToBody(ctx context.Context, action *cli.Action, data *ActionModel) error {
	required := []string{}
	props := map[string]cli.BlueprintProperty{}
	var err error
	if data.UserProperties.StringProp != nil {
		err = stringPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.ArrayProp != nil {
		err = arrayPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.NumberProp != nil {
		err = numberPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.BooleanProp != nil {
		booleanPropResourceToBody(data, props, &required)
	}

	if data.UserProperties.ObjectProp != nil {
		err = objectPropResourceToBody(data, props, &required)
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
		githubInvocation := &cli.InvocationMethod{
			Type: consts.Github,
			Org:  &org,
			Repo: &repo,
		}
		if !data.GithubMethod.Workflow.IsNull() {
			workflow := data.GithubMethod.Workflow.ValueString()
			githubInvocation.Workflow = &workflow
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
		return webhookInvocation
	}
	return nil
}
