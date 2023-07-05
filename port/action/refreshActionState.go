package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/samber/lo"
)

func writeInvocationMethodToResource(a *cli.Action, state *ActionModel) {
	if a.InvocationMethod.Type == consts.Kafka {
		state.KafkaMethod, _ = types.ObjectValue(nil, nil)
	}

	if a.InvocationMethod.Type == consts.Webhook {
		state.WebhookMethod = &WebhookMethodModel{
			Url:   types.StringValue(*a.InvocationMethod.Url),
			Agent: flex.GoBoolToFramework(a.InvocationMethod.Agent),
		}
	}

	if a.InvocationMethod.Type == consts.Github {
		state.GithubMethod = &GithubMethodModel{
			Repo:                 types.StringValue(*a.InvocationMethod.Repo),
			Org:                  types.StringValue(*a.InvocationMethod.Org),
			OmitPayload:          flex.GoBoolToFramework(a.InvocationMethod.OmitPayload),
			OmitUserInputs:       flex.GoBoolToFramework(a.InvocationMethod.OmitUserInputs),
			Workflow:             flex.GoStringToFramework(a.InvocationMethod.Workflow),
			ReportWorkflowStatus: flex.GoBoolToFramework(a.InvocationMethod.ReportWorkflowStatus),
		}
	}

	if a.InvocationMethod.Type == consts.AzureDevops {
		state.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
		}
	}
}

func addStingPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{
		MinLength: flex.GoInt64ToFramework(v.MinLength),
		MaxLength: flex.GoInt64ToFramework(v.MaxLength),
		Pattern:   flex.GoStringToFramework(v.Pattern),
	}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewStringValue(value.(string)))
		}

		stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
	}

	return stringProp
}

func addNumberPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
	numberProp := &NumberPropModel{
		Minimum: flex.GoFloat64ToFramework(v.Minimum),
		Maximum: flex.GoFloat64ToFramework(v.Maximum),
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
	objectProp := &ObjectPropModel{
		Spec: flex.GoStringToFramework(v.Spec),
	}

	return objectProp
}

func addArrayPropertiesToResource(v *cli.BlueprintProperty) *ArrayPropModel {
	arrayProp := &ArrayPropModel{
		MinItems: flex.GoInt64ToFramework(v.MinItems),
		MaxItems: flex.GoInt64ToFramework(v.MaxItems),
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

func writeInputsToResource(ctx context.Context, a *cli.Action, state *ActionModel) {
	if len(a.UserInputs.Properties) > 0 {
		properties := &UserPropertiesModel{}
		for k, v := range a.UserInputs.Properties {
			switch v.Type {
			case "string":
				if properties.StringProps == nil {
					properties.StringProps = make(map[string]StringPropModel)
				}
				stringProp := addStingPropertiesToResource(ctx, &v)

				if lo.Contains(a.UserInputs.Required, k) {
					stringProp.Required = types.BoolValue(true)
				} else {
					stringProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, stringProp)

				properties.StringProps[k] = *stringProp

			case "number":
				if properties.NumberProps == nil {
					properties.NumberProps = make(map[string]NumberPropModel)
				}

				numberProp := addNumberPropertiesToResource(ctx, &v)

				if lo.Contains(a.UserInputs.Required, k) {
					numberProp.Required = types.BoolValue(true)
				} else {
					numberProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, numberProp)

				properties.NumberProps[k] = *numberProp

			case "array":
				if properties.ArrayProps == nil {
					properties.ArrayProps = make(map[string]ArrayPropModel)
				}

				arrayProp := addArrayPropertiesToResource(&v)

				if lo.Contains(a.UserInputs.Required, k) {
					arrayProp.Required = types.BoolValue(true)
				} else {
					arrayProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, arrayProp)

				properties.ArrayProps[k] = *arrayProp

			case "boolean":
				if properties.BooleanProps == nil {
					properties.BooleanProps = make(map[string]BooleanPropModel)
				}

				booleanProp := &BooleanPropModel{}

				setCommonProperties(v, booleanProp)

				if lo.Contains(a.UserInputs.Required, k) {
					booleanProp.Required = types.BoolValue(true)
				} else {
					booleanProp.Required = types.BoolValue(false)
				}

				properties.BooleanProps[k] = *booleanProp

			case "object":
				if properties.ObjectProps == nil {
					properties.ObjectProps = make(map[string]ObjectPropModel)
				}

				objectProp := addObjectPropertiesToResource(&v)

				if lo.Contains(a.UserInputs.Required, k) {
					objectProp.Required = types.BoolValue(true)
				} else {
					objectProp.Required = types.BoolValue(false)
				}

				setCommonProperties(v, objectProp)

				properties.ObjectProps[k] = *objectProp

			}
		}
		state.UserProperties = properties
	}
}

func refreshActionState(ctx context.Context, state *ActionModel, a *cli.Action, blueprintIdentifier string) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintIdentifier, a.Identifier))
	state.Identifier = types.StringValue(a.Identifier)
	state.Blueprint = types.StringValue(blueprintIdentifier)
	state.Title = types.StringValue(a.Title)
	state.Trigger = types.StringValue(a.Trigger)

	state.Icon = flex.GoStringToFramework(a.Icon)
	state.Description = flex.GoStringToFramework(a.Description)
	state.RequiredApproval = flex.GoBoolToFramework(a.RequiredApproval)

	if a.ApprovalNotification != nil {
		if a.ApprovalNotification.Type == "email" {
			state.ApprovalEmailNotification, _ = types.ObjectValue(nil, nil)
		} else {
			state.ApprovalWebhookNotification = &ApprovalWebhookNotificationModel{
				Url: types.StringValue(a.ApprovalNotification.Url),
			}

			if a.ApprovalNotification.Format != nil {
				state.ApprovalWebhookNotification.Format = types.StringValue(*a.ApprovalNotification.Format)
			}

		}
	}

	writeInvocationMethodToResource(a, state)

	writeInputsToResource(ctx, a, state)

}

func setCommonProperties(v cli.BlueprintProperty, prop interface{}) {
	properties := []string{"Description", "Icon", "Default", "Title", "Format", "Blueprint"}
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
		case "Default":
			switch p := prop.(type) {
			case *StringPropModel:
				if v.Default == nil {
					p.Default = types.StringNull()
				} else {
					p.Default = types.StringValue(v.Default.(string))
				}
			case *NumberPropModel:
				if v.Default == nil {
					p.Default = types.Float64Null()
				} else {
					p.Default = types.Float64Value(v.Default.(float64))
				}
			case *BooleanPropModel:
				if v.Default == nil {
					p.Default = types.BoolNull()
				} else {
					p.Default = types.BoolValue(v.Default.(bool))
				}
			case *ObjectPropModel:
				if v.Default == nil {
					p.Default = types.StringNull()
				} else {
					js, _ := json.Marshal(v.Default)
					value := string(js)
					p.Default = types.StringValue(value)
				}
			}

		case "Blueprint":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Blueprint = flex.GoStringToFramework(v.Blueprint)
			case *NumberPropModel:
				p.Blueprint = flex.GoStringToFramework(v.Blueprint)
			case *BooleanPropModel:
				p.Blueprint = flex.GoStringToFramework(v.Blueprint)
			case *ArrayPropModel:
				p.Blueprint = flex.GoStringToFramework(v.Blueprint)
			case *ObjectPropModel:
				p.Blueprint = flex.GoStringToFramework(v.Blueprint)
			}

		case "Format":
			if v.Format != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Format = flex.GoStringToFramework(v.Format)
				case *NumberPropModel:
					p.Format = flex.GoStringToFramework(v.Format)
				case *BooleanPropModel:
					p.Format = flex.GoStringToFramework(v.Format)
				case *ArrayPropModel:
					p.Format = flex.GoStringToFramework(v.Format)
				case *ObjectPropModel:
					p.Format = flex.GoStringToFramework(v.Format)
				}
			}
		}
	}
}
