package action

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithImportState = &ActionResource{}

func NewActionResource() resource.Resource {
	return &ActionResource{}
}

type ActionResource struct {
	portClient *cli.PortClient
}

func (r *ActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action"
}

func (r *ActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_id>:<action_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}

func (r *ActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ActionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := data.Blueprint.ValueString()
	a, statusCode, err := r.portClient.ReadAction(ctx, data.Identifier.ValueString(), data.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading action", err.Error())
		return
	}

	writeActionFieldsToResource(ctx, data, a, blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeActionFieldsToResource(ctx context.Context, data *ActionModel, a *cli.Action, blueprintIdentifier string) {
	data.Title = types.StringValue(a.Title)
	data.Trigger = types.StringValue(a.Trigger)
	if a.Icon != nil {
		data.Icon = types.StringValue(*a.Icon)
	}
	if a.Description != nil {
		data.Description = types.StringValue(*a.Description)
	}

	if a.RequiredApproval != nil {
		data.RequiredApproval = types.BoolValue(*a.RequiredApproval)
	}

	writeInvocationMethodToResource(a, data)

	writeInputsToResource(a, data)

}

func writeInvocationMethodToResource(a *cli.Action, data *ActionModel) {
	if a.InvocationMethod.Type == "KAFKA" {
		data.KafkaMethod = types.MapNull(types.StringType)
	}

	if a.InvocationMethod.Type == "WEBHOOK" {
		data.WebhookMethod = &WebhookMethodModel{
			Url: types.StringValue(*a.InvocationMethod.Url),
		}
		if a.InvocationMethod.Agent != nil {
			data.WebhookMethod.Agent = types.BoolValue(*a.InvocationMethod.Agent)
		}
	}

	if a.InvocationMethod.Type == "GITHUB" {
		data.GithubMethod = &GithubMethodModel{
			Repo: types.StringValue(*a.InvocationMethod.Repo),
			Org:  types.StringValue(*a.InvocationMethod.Org),
		}

		if a.InvocationMethod.OmitPayload != nil {
			data.GithubMethod.OmitPayload = types.BoolValue(*a.InvocationMethod.OmitPayload)
		}

		if a.InvocationMethod.OmitUserInputs != nil {
			data.GithubMethod.OmitUserInputs = types.BoolValue(*a.InvocationMethod.OmitUserInputs)
		}

		if a.InvocationMethod.Workflow != nil {
			data.GithubMethod.Workflow = types.StringValue(*a.InvocationMethod.Workflow)
		}

		if a.InvocationMethod.Branch != nil {
			data.GithubMethod.Branch = types.StringValue(*a.InvocationMethod.Branch)
		}
	}

	if a.InvocationMethod.Type == "AZURE-DEVOPS" {
		data.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
		}
	}
}

func writeInputsToResource(a *cli.Action, data *ActionModel) {
	if len(a.UserInputs.Properties) > 0 {

	}

}
func (r *ActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteAction(ctx, data.Identifier.ValueString(), data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete action", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	action, err := actionResourceToBody(ctx, data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	a, err := r.portClient.CreateAction(ctx, bp.Identifier, action)
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	data.ID = types.StringValue(a.Identifier)
	data.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	action, err := actionResourceToBody(ctx, data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	a, err := r.portClient.UpdateAction(ctx, bp.Identifier, action.ID, action)
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	data.ID = types.StringValue(a.Identifier)
	data.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func actionResourceToBody(ctx context.Context, data *ActionModel, bp *cli.Blueprint) (*cli.Action, error) {
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

	action.InvocationMethod = invocationMethodToBody(data)

	if data.UserProperties != nil {
		actionPropertiesToBody(ctx, action, data)
	}

	return action, nil
}

func stringPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
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

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			property.Icon = &icon
		}

		if !prop.MinLength.IsNull() {
			property.MinLength = int(prop.MinLength.ValueInt64())
		}

		if !prop.MaxLength.IsNull() {
			property.MaxLength = int(prop.MaxLength.ValueInt64())
		}

		if !prop.Pattern.IsNull() {
			property.Pattern = prop.Pattern.ValueString()
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			property.Description = &description
		}

		if !prop.Enum.IsNull() {
			enumList := []interface{}{}
			for _, enum := range prop.Enum.Elements() {
				v, _ := enum.ToTerraformValue(ctx)
				var keyValue string
				v.As(&keyValue)
				enumList = append(enumList, keyValue)
			}
			property.Enum = enumList
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func numberPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.UserProperties.NumberProp {
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
				property.Enum = []interface{}{}
				for _, e := range prop.Enum.Elements() {
					v, _ := e.ToTerraformValue(ctx)
					var keyValue big.Float
					v.As(&keyValue)
					floatValue, _ := keyValue.Float64()
					property.Enum = append(property.Enum, floatValue)
				}
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
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

func objectPropResourceToBody(d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
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
					log.Fatal(err)
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
}

func arrayPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
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

			if prop.StringItems != nil {
				items := map[string]interface{}{}
				items["type"] = "string"
				if !prop.StringItems.Format.IsNull() {
					items["format"] = prop.StringItems.Format.ValueString()
				}
				if !prop.StringItems.Default.IsNull() {
					defaultList := []interface{}{}
					for _, e := range prop.StringItems.Default.Elements() {
						v, _ := e.ToTerraformValue(ctx)
						var keyValue string
						v.As(&keyValue)
						defaultList = append(defaultList, keyValue)
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.NumberItems != nil {
				items := map[string]interface{}{}
				items["type"] = "number"
				if !prop.NumberItems.Default.IsNull() {
					items["default"] = prop.NumberItems.Default
				}
				property.Items = items
			}

			if prop.BooleanItems != nil {
				items := map[string]interface{}{}
				items["type"] = "boolean"
				if !prop.BooleanItems.Default.IsNull() {
					items["default"] = prop.BooleanItems.Default
				}
				property.Items = items
			}

			if prop.ObjectItems != nil {
				items := map[string]interface{}{}
				items["type"] = "object"
				if !prop.ObjectItems.Default.IsNull() {
					items["default"] = prop.ObjectItems.Default
				}
				property.Items = items
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func actionPropertiesToBody(ctx context.Context, action *cli.Action, data *ActionModel) {
	required := []string{}
	props := map[string]cli.BlueprintProperty{}

	if data.UserProperties.StringProp != nil {
		stringPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.ArrayProp != nil {
		arrayPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.NumberProp != nil {
		numberPropResourceToBody(ctx, data, props, &required)
	}
	if data.UserProperties.BooleanProp != nil {
		booleanPropResourceToBody(data, props, &required)
	}

	if data.UserProperties.ObjectProp != nil {
		objectPropResourceToBody(data, props, &required)
	}

	action.UserInputs.Properties = props
	action.UserInputs.Required = required

}
func invocationMethodToBody(data *ActionModel) *cli.InvocationMethod {
	if data.AzureMethod != nil {
		org := data.AzureMethod.Org.ValueString()
		webhook := data.AzureMethod.Webhook.ValueString()
		return &cli.InvocationMethod{
			Type:    "AZURE-DEVOPS",
			Org:     &org,
			Webhook: &webhook,
		}
	}

	if data.GithubMethod != nil {
		org := data.GithubMethod.Org.ValueString()
		repo := data.GithubMethod.Repo.ValueString()
		githubInvocation := &cli.InvocationMethod{
			Type: "GITHUB",
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
			Type: "KAFKA",
		}
	}

	if data.WebhookMethod != nil {
		url := data.WebhookMethod.Url.ValueString()
		webhookInvocation := &cli.InvocationMethod{
			Type: "WEBHOOK",
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
