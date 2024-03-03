package action

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
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

func writeVisibleToResource(v cli.ActionProperty) (types.Bool, types.String) {
	if v.Visible == nil {
		return types.BoolNull(), types.StringNull()
	}

	visible := reflect.ValueOf(v.Visible)
	switch visible.Kind() {
	case reflect.Bool:
		boolValue := visible.Interface().(bool)
		return types.BoolValue(boolValue), types.StringNull()
	case reflect.Map:
		jq := visible.Interface().(map[string]any)
		jqQueryValue := jq["jqQuery"].(string)
		return types.BoolNull(), types.StringValue(jqQueryValue)
	}

	return types.BoolNull(), types.StringNull()
}

func writeRequiredToResource(v cli.ActionUserInputs) (types.String, []string) {
	// If required is nil, return an empty string and nil
	if v.Required == nil {
		return types.StringNull(), nil
	}

	required := reflect.ValueOf(v.Required)
	switch required.Kind() {
	// if required is a slice of strings that means that the user has specified which properties are required
	case reflect.Slice:
		slice := required.Interface().([]interface{})
		attrs := make([]string, 0, required.Len())
		for _, value := range slice {
			attrs = append(attrs, value.(string))
		}
		return types.StringNull(), attrs
	// if required is a map, that means that the user has specified a jq query to determine which properties are required
	case reflect.Map:
		jq := required.Interface().(map[string]any)
		jqQueryValue := jq["jqQuery"].(string)
		return types.StringValue(jqQueryValue), nil
	}

	// if required is not a slice or a map, return an empty string and nil
	return types.StringNull(), nil
}

func writeInputsToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if len(a.UserInputs.Properties) > 0 {
		properties := &UserPropertiesModel{}
		requiredJq, required := writeRequiredToResource(a.UserInputs)
		for k, v := range a.UserInputs.Properties {
			switch v.Type {
			case "string":
				if properties.StringProps == nil {
					properties.StringProps = make(map[string]StringPropModel)
				}
				stringProp := addStringPropertiesToResource(ctx, &v)

				if requiredJq.IsNull() && lo.Contains(required, k) {
					stringProp.Required = types.BoolValue(true)
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

				if requiredJq.IsNull() && lo.Contains(required, k) {
					numberProp.Required = types.BoolValue(true)
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

				if requiredJq.IsNull() && lo.Contains(required, k) {
					arrayProp.Required = types.BoolValue(true)
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

				if requiredJq.IsNull() && lo.Contains(required, k) {
					booleanProp.Required = types.BoolValue(true)
				}

				properties.BooleanProps[k] = *booleanProp

			case "object":
				if properties.ObjectProps == nil {
					properties.ObjectProps = make(map[string]ObjectPropModel)
				}

				objectProp := addObjectPropertiesToResource(&v)

				if requiredJq.IsNull() && lo.Contains(required, k) {
					objectProp.Required = types.BoolValue(true)
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
			state.OrderProperties = flex.GoArrayStringToTerraformList(ctx, a.UserInputs.Order)
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

	requiredJq, _ := writeRequiredToResource(a.UserInputs)

	state.RequiredJqQuery = requiredJq

	writeInvocationMethodToResource(a, state)

	err := writeInputsToResource(ctx, a, state)
	if err != nil {
		return err
	}
	return nil
}

func setCommonProperties(ctx context.Context, v cli.ActionProperty, prop interface{}) error {
	properties := []string{"Description", "Icon", "Default", "Title", "DependsOn", "Dataset", "Visible"}
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
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *NumberPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *BooleanPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *ArrayPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
			case *ObjectPropModel:
				p.DependsOn = flex.GoArrayStringToTerraformList(ctx, v.DependsOn)
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
					if p.StringItems != nil {
						p.StringItems.Dataset = dataset
					}
				case *ObjectPropModel:
					p.Dataset = dataset
				}
			}

		case "Visible":
			visible, visibleJq := writeVisibleToResource(v)
			if !visible.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Visible = visible
				case *NumberPropModel:
					p.Visible = visible
				case *BooleanPropModel:
					p.Visible = visible
				case *ArrayPropModel:
					p.Visible = visible
				case *ObjectPropModel:
					p.Visible = visible
				}
			}
			if !visibleJq.IsNull() {
				switch p := prop.(type) {
				case *StringPropModel:
					p.VisibleJqQuery = visibleJq
				case *NumberPropModel:
					p.VisibleJqQuery = visibleJq
				case *BooleanPropModel:
					p.VisibleJqQuery = visibleJq
				case *ArrayPropModel:
					p.VisibleJqQuery = visibleJq
				case *ObjectPropModel:
					p.VisibleJqQuery = visibleJq
				}
			}
		}
	}
	return nil
}
