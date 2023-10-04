package action

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
	"github.com/samber/lo"
)

func writeInvocationMethodToResource(a *cli.Action, state *ActionModel) {
	if a.InvocationMethod.Type == consts.Kafka {
		state.KafkaMethod, _ = types.ObjectValue(nil, nil)
	}

	if a.InvocationMethod.Type == consts.Webhook {
		state.WebhookMethod = &WebhookMethodModel{
			Url:          types.StringValue(*a.InvocationMethod.Url),
			Agent:        flex.GoBoolToFramework(a.InvocationMethod.Agent),
			Synchronized: flex.GoBoolToFramework(a.InvocationMethod.Synchronized),
			Method:       flex.GoStringToFramework(a.InvocationMethod.Method),
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

	if a.InvocationMethod.Type == consts.Gitlab {
		state.GitlabMethod = &GitlabMethodModel{
			ProjectName:    types.StringValue(*a.InvocationMethod.ProjectName),
			GroupName:      types.StringValue(*a.InvocationMethod.GroupName),
			OmitPayload:    flex.GoBoolToFramework(a.InvocationMethod.OmitPayload),
			OmitUserInputs: flex.GoBoolToFramework(a.InvocationMethod.OmitUserInputs),
			DefaultRef:     types.StringValue(*a.InvocationMethod.DefaultRef),
			Agent:          flex.GoBoolToFramework(a.InvocationMethod.Agent),
		}
	}
}

func writeDatasetToResource(v cli.ActionProperty) *DatasetModel {
	if v.Dataset == nil {
		return nil
	}

	dataset := v.Dataset

	datasetModel := &DatasetModel{
		Combinator: types.StringValue(dataset.Combinator),
	}

	for _, v := range dataset.Rules {
		rule := &Rule{
			Blueprint: flex.GoStringToFramework(v.Blueprint),
			Property:  flex.GoStringToFramework(v.Property),
			Operator:  flex.GoStringToFramework(&v.Operator),
			Value: &Value{
				JqQuery: flex.GoStringToFramework(&v.Value.JqQuery),
			},
		}
		datasetModel.Rules = append(datasetModel.Rules, *rule)
	}

	return datasetModel

}
func writeInputsToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if len(a.UserInputs.Properties) > 0 {
		properties := &UserPropertiesModel{}
		for k, v := range a.UserInputs.Properties {
			switch v.Type {
			case "string":
				if properties.StringProps == nil {
					properties.StringProps = make(map[string]StringPropModel)
				}
				stringProp := addStringPropertiesToResource(ctx, &v)

				if lo.Contains(a.UserInputs.Required, k) {
					stringProp.Required = types.BoolValue(true)
				} else {
					stringProp.Required = types.BoolValue(false)
				}

				err := setCommonProperties(ctx, v, stringProp)
				if err != nil {
					return err
				}

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

				err := setCommonProperties(ctx, v, numberProp)
				if err != nil {
					return err
				}

				properties.NumberProps[k] = *numberProp

			case "array":
				if properties.ArrayProps == nil {
					properties.ArrayProps = make(map[string]ArrayPropModel)
				}

				arrayProp, err := addArrayPropertiesToResource(&v)
				if err != nil {
					return err
				}

				if lo.Contains(a.UserInputs.Required, k) {
					arrayProp.Required = types.BoolValue(true)
				} else {
					arrayProp.Required = types.BoolValue(false)
				}

				err = setCommonProperties(ctx, v, arrayProp)
				if err != nil {
					return err
				}

				properties.ArrayProps[k] = *arrayProp

			case "boolean":
				if properties.BooleanProps == nil {
					properties.BooleanProps = make(map[string]BooleanPropModel)
				}

				booleanProp := &BooleanPropModel{}

				err := setCommonProperties(ctx, v, booleanProp)
				if err != nil {
					return err
				}

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

				err := setCommonProperties(ctx, v, objectProp)
				if err != nil {
					return err
				}

				properties.ObjectProps[k] = *objectProp

			}
		}
		state.UserProperties = properties
		if len(a.UserInputs.Order) > 0 {
			state.OrderProperties = flex.GoArrayStringToTerraformList(a.UserInputs.Order)
		}
	}
	return nil
}

func refreshActionState(ctx context.Context, state *ActionModel, a *cli.Action, blueprintIdentifier string) error {
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

	err := writeInputsToResource(ctx, a, state)
	if err != nil {
		return err
	}
	return nil
}

func setCommonProperties(ctx context.Context, v cli.ActionProperty, prop interface{}) error {
	properties := []string{"Description", "Icon", "Default", "Title", "DependsOn", "Dataset"}
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
			case *ObjectPropModel:
				if v, ok := v.Default.(map[string]interface{}); ok {
					if v["jqQuery"] != nil {
						p.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
					} else {
						defaultValue, err := utils.GoObjectToTerraformString(v)
						if err != nil {
							return fmt.Errorf("error converting default value to terraform string: %s", err.Error())
						}
						if defaultValue.IsNull() {
							p.Default = types.StringNull()
							p.DefaultJqQuery = types.StringNull()
						}
						p.Default = defaultValue
					}
				}
			}
		case "DependsOn":
			switch p := prop.(type) {
			case *StringPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(v.DependsOn)
			case *NumberPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(v.DependsOn)
			case *BooleanPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(v.DependsOn)
			case *ArrayPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(v.DependsOn)
			case *ObjectPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(v.DependsOn)
			}

		case "Dataset":
			dataset := writeDatasetToResource(v)
			if dataset != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Dataset = dataset
				case *NumberPropModel:
					p.Dataset = dataset
				case *BooleanPropModel:
					p.Dataset = dataset
				case *ArrayPropModel:
					p.Dataset = dataset
				case *ObjectPropModel:
					p.Dataset = dataset
				}
			}
		}
	}
	return nil
}
