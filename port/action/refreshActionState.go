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

func writeInvocationMethodToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if a.InvocationMethod.Type == consts.Kafka {
		payload, err := utils.GoObjectToTerraformString(a.InvocationMethod.Payload)
		if err != nil {
			return err
		}

		state.KafkaMethod = &KafkaMethodModel{
			Payload: payload,
		}
	}

	if a.InvocationMethod.Type == consts.Webhook {
		agent, err := utils.GoObjectToTerraformString(a.InvocationMethod.Agent)
		if err != nil {
			return err
		}
		synchronized, err := utils.GoObjectToTerraformString(a.InvocationMethod.Synchronized)
		if err != nil {
			return err
		}
		headers, _ := types.MapValueFrom(ctx, types.StringType, a.InvocationMethod.Headers)
		body, err := utils.GoObjectToTerraformString(a.InvocationMethod.Body)
		if err != nil {
			return err
		}

		state.WebhookMethod = &WebhookMethodModel{
			Url:          types.StringValue(*a.InvocationMethod.Url),
			Agent:        agent,
			Synchronized: synchronized,
			Method:       flex.GoStringToFramework(a.InvocationMethod.Method),
			Headers:      headers,
			Body:         body,
		}
	}

	if a.InvocationMethod.Type == consts.Github {
		workflowInputs, err := utils.GoObjectToTerraformString(a.InvocationMethod.WorkflowInputs)
		if err != nil {
			return err
		}
		reportWorkflowStatus, err := utils.GoObjectToTerraformString(a.InvocationMethod.ReportWorkflowStatus)
		if err != nil {
			return err
		}

		state.GithubMethod = &GithubMethodModel{
			Org:                  types.StringValue(*a.InvocationMethod.Org),
			Repo:                 types.StringValue(*a.InvocationMethod.Repo),
			Workflow:             types.StringValue(*a.InvocationMethod.Workflow),
			WorkflowInputs:       workflowInputs,
			ReportWorkflowStatus: reportWorkflowStatus,
		}
	}

	if a.InvocationMethod.Type == consts.Gitlab {
		pipelineVariables, err := utils.GoObjectToTerraformString(a.InvocationMethod.PipelineVariables)
		if err != nil {
			return err
		}

		state.GitlabMethod = &GitlabMethodModel{
			ProjectName:       types.StringValue(*a.InvocationMethod.ProjectName),
			GroupName:         types.StringValue(*a.InvocationMethod.GroupName),
			DefaultRef:        flex.GoStringToFramework(a.InvocationMethod.DefaultRef),
			PipelineVariables: pipelineVariables,
		}
	}

	if a.InvocationMethod.Type == consts.AzureDevops {
		payload, err := utils.GoObjectToTerraformString(a.InvocationMethod.Payload)
		if err != nil {
			return err
		}

		state.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
			Payload: payload,
		}
	}

	// TODO: return when frontend for upsert entity is ready
	//if a.InvocationMethod.Type == consts.UpsertEntity {
	//	var teams []types.String
	//	switch team := a.InvocationMethod.Team.(type) {
	//	case string:
	//		teams = append(teams, types.StringValue(team))
	//	case []interface{}:
	//		teams = make([]types.String, 0)
	//		for _, t := range team {
	//			teams = append(teams, types.StringValue(t.(string)))
	//		}
	//	}
	//	properties, err := utils.GoObjectToTerraformString(a.InvocationMethod.Properties)
	//	if err != nil {
	//		return err
	//	}
	//	relations, err := utils.GoObjectToTerraformString(a.InvocationMethod.Relations)
	//	if err != nil {
	//		return err
	//	}
	//
	//	state.UpsertEntityMethod = &UpsertEntityMethodModel{
	//		Identifier:          types.StringValue(*a.InvocationMethod.Identifier),
	//		Title:               flex.GoStringToFramework(a.InvocationMethod.Title),
	//		BlueprintIdentifier: types.StringValue(*a.InvocationMethod.BlueprintIdentifier),
	//		Teams:               teams,
	//		Icon:                flex.GoStringToFramework(a.InvocationMethod.Icon),
	//		Properties:          properties,
	//		Relations:           relations,
	//	}
	//}

	return nil
}

func writeDatasetToResource(ds *cli.Dataset) *DatasetModel {
	if ds == nil {
		return nil
	}

	datasetModel := &DatasetModel{
		Combinator: types.StringValue(ds.Combinator),
	}

	for _, v := range ds.Rules {
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

func buildRequired(v *cli.ActionUserInputs) (types.String, []string) {
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

func buildUserProperties(ctx context.Context, a *cli.Action) (*UserPropertiesModel, error) {
	properties := &UserPropertiesModel{}
	if len(a.Trigger.UserInputs.Properties) > 0 {
		requiredJq, required := buildRequired(a.Trigger.UserInputs)
		for k, v := range a.Trigger.UserInputs.Properties {
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
					return nil, err
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
					return nil, err
				}

				properties.NumberProps[k] = *numberProp

			case "array":
				if properties.ArrayProps == nil {
					properties.ArrayProps = make(map[string]ArrayPropModel)
				}

				arrayProp, err := addArrayPropertiesToResource(&v)
				if err != nil {
					return nil, err
				}

				if requiredJq.IsNull() && lo.Contains(required, k) {
					arrayProp.Required = types.BoolValue(true)
				}

				err = setCommonProperties(ctx, v, arrayProp)
				if err != nil {
					return nil, err
				}

				properties.ArrayProps[k] = *arrayProp

			case "boolean":
				if properties.BooleanProps == nil {
					properties.BooleanProps = make(map[string]BooleanPropModel)
				}

				booleanProp := &BooleanPropModel{}

				err := setCommonProperties(ctx, v, booleanProp)
				if err != nil {
					return nil, err
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
					return nil, err
				}

				properties.ObjectProps[k] = *objectProp

			}
		}
	}

	return properties, nil
}

func writeTriggerToResource(ctx context.Context, a *cli.Action, state *ActionModel) error {
	if a.Trigger.Type == consts.SelfService {
		userProperties, err := buildUserProperties(ctx, a)
		if err != nil {
			return err
		}
		requiredJqQuery, _ := buildRequired(a.Trigger.UserInputs)
		orderProperties := types.ListNull(types.StringType)
		if len(a.Trigger.UserInputs.Order) > 0 {
			orderProperties = flex.GoArrayStringToTerraformList(ctx, a.Trigger.UserInputs.Order)
		}

		state.SelfServiceTrigger = &SelfServiceTriggerModel{
			BlueprintIdentifier: flex.GoStringToFramework(a.Trigger.BlueprintIdentifier),
			Operation:           types.StringValue(*a.Trigger.Operation),
			UserProperties:      userProperties,
			RequiredJqQuery:     requiredJqQuery,
			OrderProperties:     orderProperties,
		}
	}

	// TODO: return when frontend for automations is ready
	//if a.Trigger.Type == consts.Automation {
	//	automationTrigger := &AutomationTriggerModel{}
	//
	//	var expressions []types.String
	//	if a.Trigger.Condition != nil {
	//		for _, e := range a.Trigger.Condition.Expressions {
	//			expressions = append(expressions, types.StringValue(e))
	//		}
	//		automationTrigger.JqCondition = &JqConditionModel{
	//			Expressions: expressions,
	//			Combinator:  flex.GoStringToFramework(a.Trigger.Condition.Combinator),
	//		}
	//	}
	//
	//	if a.Trigger.Event.Type == consts.EntityCreated {
	//		automationTrigger.EntityCreatedEvent = &EntityCreatedEventModel{
	//			BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
	//		}
	//	}
	//
	//	if a.Trigger.Event.Type == consts.EntityUpdated {
	//		automationTrigger.EntityUpdatedEvent = &EntityUpdatedEventModel{
	//			BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
	//		}
	//	}
	//
	//	if a.Trigger.Event.Type == consts.EntityDeleted {
	//		automationTrigger.EntityDeletedEvent = &EntityDeletedEventModel{
	//			BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
	//		}
	//	}
	//
	//	if a.Trigger.Event.Type == consts.AnyEntityChange {
	//		automationTrigger.AnyEntityChangeEvent = &AnyEntityChangeEventModel{
	//			BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
	//		}
	//	}
	//
	//	if a.Trigger.Event.Type == consts.TimerPropertyExpired {
	//		automationTrigger.TimerPropertyExpiredEvent = &TimerPropertyExpiredEventModel{
	//			BlueprintIdentifier: types.StringValue(*a.Trigger.Event.BlueprintIdentifier),
	//			PropertyIdentifier:  types.StringValue(*a.Trigger.Event.PropertyIdentifier),
	//		}
	//	}
	//
	//	state.AutomationTrigger = automationTrigger
	//}

	return nil
}

func refreshActionState(ctx context.Context, state *ActionModel, a *cli.Action) error {
	state.ID = types.StringValue(a.Identifier)
	state.Identifier = types.StringValue(a.Identifier)
	state.Blueprint = types.StringNull()
	state.Title = flex.GoStringToFramework(a.Title)
	state.Icon = flex.GoStringToFramework(a.Icon)
	state.Description = flex.GoStringToFramework(a.Description)

	err := writeTriggerToResource(ctx, a, state)
	if err != nil {
		return err
	}

	err = writeInvocationMethodToResource(ctx, a, state)
	if err != nil {
		return err
	}

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
	state.Publish = flex.GoBoolToFramework(a.Publish)

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
