package action

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
	"github.com/samber/lo"
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
	a, statusCode, err := r.portClient.ReadAction(ctx, data.Blueprint.ValueString(), data.Identifier.ValueString())
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
	data.ID = types.StringValue(a.Identifier)
	data.Identifier = types.StringValue(a.Identifier)
	data.Blueprint = types.StringValue(blueprintIdentifier)
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

	writeInputsToResource(ctx, a, data)

}

func writeInvocationMethodToResource(a *cli.Action, data *ActionModel) {
	if a.InvocationMethod.Type == "KAFKA" {
		data.KafkaMethod, _ = types.ObjectValue(nil, nil)
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

func setCommonProperties(v cli.BlueprintProperty, prop interface{}) {
	properties := []string{"Description", "Icon", "Default", "Title"}
	for _, property := range properties {
		switch property {
		case "Description":
			if v.Description != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Description = types.StringValue(*v.Description)
				case *NumberPropModel:
					p.Description = types.StringValue(*v.Description)
				case *BooleanPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ArrayPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ObjectPropModel:
					p.Description = types.StringValue(*v.Description)
				}
			}
		case "Icon":
			if v.Icon != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *NumberPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *BooleanPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ArrayPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ObjectPropModel:
					p.Icon = types.StringValue(*v.Icon)
				}
			}
		case "Title":
			if v.Title != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Title = types.StringValue(*v.Title)
				case *NumberPropModel:
					p.Title = types.StringValue(*v.Title)
				case *BooleanPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ArrayPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ObjectPropModel:
					p.Title = types.StringValue(*v.Title)
				}
			}
		case "Default":
			if v.Default != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Default = types.StringValue(v.Default.(string))
				case *NumberPropModel:
					p.Default = types.Float64Value(v.Default.(float64))
				case *BooleanPropModel:
					p.Default = types.BoolValue(v.Default.(bool))
				case *ObjectPropModel:
					js, _ := json.Marshal(v.Default)
					value := string(js)
					p.Default = types.StringValue(value)
				}
			}
		}
	}
}
func addStingPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewStringValue(value.(string)))
		}

		stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
	}

	if v.Format != nil {
		stringProp.Format = types.StringValue(*v.Format)
	}

	if v.MinLength != 0 {
		stringProp.MinLength = types.Int64Value(int64(v.MinLength))
	}

	if v.MaxLength != 0 {
		stringProp.MaxLength = types.Int64Value(int64(v.MaxLength))
	}

	if v.Pattern != "" {
		stringProp.Pattern = types.StringValue(v.Pattern)
	}

	return stringProp
}

func addNumberPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
	numberProp := &NumberPropModel{}
	if v.Minimum != nil {
		numberProp.Minimum = types.Float64Value(*v.Minimum)
	}

	if v.Maximum != nil {
		numberProp.Maximum = types.Float64Value(*v.Maximum)
	}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
		}

		numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)
	} else {
		numberProp.Enum = types.ListNull(types.Float64Type)
	}

	return numberProp
}

func addObjectPropertiesToResource(v *cli.BlueprintProperty) *ObjectPropModel {
	objectProp := &ObjectPropModel{}

	if v.Spec != nil {
		objectProp.Spec = types.StringValue(*v.Spec)
	}

	return objectProp
}

func addArrayPropertiesToResource(v *cli.BlueprintProperty) *ArrayPropModel {
	arrayProp := &ArrayPropModel{}
	if v.MinItems != nil {
		arrayProp.MinItems = types.Int64Value(int64(*v.MinItems))
	}
	if v.MaxItems != nil {
		arrayProp.MaxItems = types.Int64Value(int64(*v.MaxItems))
	}
	if v.Items != nil {
		if v.Items["type"] != "" {
			switch v.Items["type"] {
			case "string":
				arrayProp.StringItems = &StringItems{}
				if v.Default != nil {
					stringArray := make([]string, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						stringArray[i] = v.(string)
					}
					attrs := make([]attr.Value, 0, len(stringArray))
					for _, value := range stringArray {
						attrs = append(attrs, basetypes.NewStringValue(value))
					}
					arrayProp.StringItems.Default, _ = types.ListValue(types.StringType, attrs)
				} else {
					arrayProp.StringItems.Default = types.ListNull(types.StringType)
				}
				if value, ok := v.Items["format"]; ok && value != nil {
					arrayProp.StringItems.Format = types.StringValue(v.Items["format"].(string))
				}
			case "number":
				arrayProp.NumberItems = &NumberItems{}
				if v.Default != nil {
					numberArray := make([]float64, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(numberArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
					}
					arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
				}

			case "boolean":
				arrayProp.BooleanItems = &BooleanItems{}
				if v.Default != nil {
					booleanArray := make([]bool, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(booleanArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewBoolValue(value.(bool)))
					}
					arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
				}
			}
		}
	}

	return arrayProp
}

func writeInputsToResource(ctx context.Context, a *cli.Action, data *ActionModel) {
	if len(a.UserInputs.Properties) > 0 {
		properties := &UserPropertiesModel{}
		for k, v := range a.UserInputs.Properties {
			switch v.Type {
			case "string":
				if properties.StringProp == nil {
					properties.StringProp = make(map[string]StringPropModel)
				}
				stringProp := addStingPropertiesToResource(ctx, &v)

				if lo.Contains(a.UserInputs.Required, k) {
					stringProp.Required = types.BoolValue(true)
				} else {
					stringProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, stringProp)

				properties.StringProp[k] = *stringProp

			case "number":
				if properties.NumberProp == nil {
					properties.NumberProp = make(map[string]NumberPropModel)
				}

				numberProp := addNumberPropertiesToResource(ctx, &v)

				if lo.Contains(a.UserInputs.Required, k) {
					numberProp.Required = types.BoolValue(true)
				} else {
					numberProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, numberProp)

				properties.NumberProp[k] = *numberProp

			case "array":
				if properties.ArrayProp == nil {
					properties.ArrayProp = make(map[string]ArrayPropModel)
				}

				arrayProp := addArrayPropertiesToResource(&v)

				if !data.UserProperties.ArrayProp[k].Required.IsNull() {
					if lo.Contains(a.UserInputs.Required, k) {
						arrayProp.Required = types.BoolValue(true)
					} else {
						arrayProp.Required = types.BoolValue(false)
					}
				}

				setCommonProperties(v, arrayProp)

				properties.ArrayProp[k] = *arrayProp

			case "boolean":
				if properties.BooleanProp == nil {
					properties.BooleanProp = make(map[string]BooleanPropModel)
				}

				booleanProp := &BooleanPropModel{}

				setCommonProperties(v, booleanProp)

				if !data.UserProperties.BooleanProp[k].Required.IsNull() {
					if lo.Contains(a.UserInputs.Required, k) {
						booleanProp.Required = types.BoolValue(true)
					} else {
						booleanProp.Required = types.BoolValue(false)
					}
				}

				properties.BooleanProp[k] = *booleanProp

			case "object":
				if properties.ObjectProp == nil {
					properties.ObjectProp = make(map[string]ObjectPropModel)
				}

				objectProp := addObjectPropertiesToResource(&v)

				if !data.UserProperties.ObjectProp[k].Required.IsNull() {
					if lo.Contains(a.UserInputs.Required, k) {
						objectProp.Required = types.BoolValue(true)
					} else {
						objectProp.Required = types.BoolValue(false)
					}
				}

				setCommonProperties(v, objectProp)

				properties.ObjectProp[k] = *objectProp

			}

		}
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

	a, err := r.portClient.UpdateAction(ctx, bp.Identifier, action.Identifier, action)
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
