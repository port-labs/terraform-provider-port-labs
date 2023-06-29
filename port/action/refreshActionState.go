package action

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/samber/lo"
)

func writeInvocationMethodToResource(a *cli.Action, data *ActionModel) {
	if a.InvocationMethod.Type == consts.Kafka {
		data.KafkaMethod, _ = types.ObjectValue(nil, nil)
	}

	if a.InvocationMethod.Type == consts.Webhook {
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

		if a.InvocationMethod.ReportWorkflowStatus != nil {
			data.GithubMethod.ReportWorkflowStatus = types.BoolValue(*a.InvocationMethod.ReportWorkflowStatus)
		}
	}

	if a.InvocationMethod.Type == "AZURE-DEVOPS" {
		data.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
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
		data.UserProperties = properties
	}
}

func refreshActionState(ctx context.Context, data *ActionModel, a *cli.Action, blueprintIdentifier string) {
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
