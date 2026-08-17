package workflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func workflowStateToPortBody(ctx context.Context, state *WorkflowModel) (*cli.Workflow, error) {
	w := &cli.Workflow{
		Identifier:            state.Identifier.ValueString(),
		Title:                 state.Title.ValueStringPointer(),
		Icon:                  state.Icon.ValueStringPointer(),
		Description:           state.Description.ValueStringPointer(),
		Category:              state.Category.ValueStringPointer(),
		AllowAnyoneToViewRuns: state.AllowAnyoneToViewRuns.ValueBoolPointer(),
	}

	nodes, err := nodesToPortBody(ctx, state.Nodes)
	if err != nil {
		return nil, err
	}
	w.Nodes = nodes

	w.Connections = connectionsToPortBody(state.Connections)

	return w, nil
}

func nodesToPortBody(ctx context.Context, nodes []WorkflowNodeModel) ([]cli.WorkflowNode, error) {
	result := make([]cli.WorkflowNode, 0, len(nodes))
	for _, n := range nodes {
		node := cli.WorkflowNode{
			Identifier:  n.Identifier.ValueString(),
			Title:       n.Title.ValueStringPointer(),
			Icon:        n.Icon.ValueStringPointer(),
			Description: n.Description.ValueStringPointer(),
			Verbose:     n.Verbose.ValueBoolPointer(),
		}

		links, err := terraformListToStrings(ctx, n.Links)
		if err != nil {
			return nil, err
		}
		node.Links = links

		variables, err := terraformMapToStrings(ctx, n.Variables)
		if err != nil {
			return nil, err
		}
		node.Variables = variables

		config, err := nodeConfigToPortBody(ctx, n)
		if err != nil {
			return nil, err
		}
		node.Config = *config

		result = append(result, node)
	}
	return result, nil
}

func nodeConfigToPortBody(ctx context.Context, n WorkflowNodeModel) (*cli.WorkflowNodeConfig, error) {
	config := &cli.WorkflowNodeConfig{}

	switch {
	case n.SelfServeTrigger != nil:
		config.Type = consts.SelfServeTrigger
		t := n.SelfServeTrigger
		config.Published = t.Published.ValueBoolPointer()
		config.ActionCardButtonText = t.ActionCardButtonText.ValueStringPointer()
		config.ExecuteActionButtonText = t.ExecuteActionButtonText.ValueStringPointer()
		config.Variant = t.Variant.ValueStringPointer()

		if t.UserInputs != nil {
			userInputs, err := selfServeUserInputsToPortBody(ctx, t.UserInputs)
			if err != nil {
				return nil, err
			}
			config.UserInputs = userInputs
		}

		permissions, err := permissionsToPortBody(ctx, t.Permissions)
		if err != nil {
			return nil, err
		}
		config.Permissions = permissions

		for _, c := range t.Contexts {
			config.Contexts = append(config.Contexts, cli.WorkflowTriggerContext{
				On:                  c.On.ValueString(),
				BlueprintIdentifier: c.BlueprintIdentifier.ValueStringPointer(),
				UserInput:           c.UserInput.ValueStringPointer(),
			})
		}

	case n.EventTrigger != nil:
		config.Type = consts.EventTrigger
		t := n.EventTrigger
		config.Published = t.Published.ValueBoolPointer()
		config.Event = &cli.WorkflowTriggerEvent{
			Type:                t.Type.ValueString(),
			BlueprintIdentifier: t.BlueprintIdentifier.ValueString(),
			PropertyIdentifier:  t.PropertyIdentifier.ValueStringPointer(),
		}

		if t.Condition != nil {
			expressions, err := terraformListToStrings(ctx, t.Condition.Expressions)
			if err != nil {
				return nil, err
			}
			config.Condition = &cli.WorkflowNodeCondition{
				Type:        consts.JqCondition,
				Expressions: expressions,
				Combinator:  t.Condition.Combinator.ValueStringPointer(),
			}
		}

	case n.ScheduleTrigger != nil:
		config.Type = consts.ScheduleTrigger
		config.Cron = n.ScheduleTrigger.Cron.ValueStringPointer()
		config.Published = n.ScheduleTrigger.Published.ValueBoolPointer()

	case n.Kafka != nil:
		config.Type = consts.Kafka
		config.OnFailure = n.Kafka.OnFailure.ValueStringPointer()
		payload, err := jsonStringToAny(n.Kafka.Payload)
		if err != nil {
			return nil, err
		}
		config.Payload = payload

	case n.Webhook != nil:
		config.Type = consts.Webhook
		w := n.Webhook
		config.Url = w.Url.ValueStringPointer()
		config.Agent = w.Agent.ValueBoolPointer()
		config.Synchronized = w.Synchronized.ValueBoolPointer()
		config.Method = w.Method.ValueStringPointer()
		config.OnTimeout = w.OnTimeout.ValueStringPointer()
		config.OnFailure = w.OnFailure.ValueStringPointer()

		headers, err := terraformMapToStrings(ctx, w.Headers)
		if err != nil {
			return nil, err
		}
		config.Headers = headers

		body, err := jsonStringToAny(w.Body)
		if err != nil {
			return nil, err
		}
		config.Body = body

	case n.IntegrationAction != nil:
		config.Type = consts.IntegrationAction
		i := n.IntegrationAction
		config.InstallationId = i.InstallationId.ValueStringPointer()
		config.IntegrationProvider = i.IntegrationProvider.ValueStringPointer()
		config.IntegrationInvocationType = i.IntegrationInvocationType.ValueStringPointer()
		config.OnFailure = i.OnFailure.ValueStringPointer()

		executionProperties, err := jsonStringToAny(i.ExecutionProperties)
		if err != nil {
			return nil, err
		}
		config.IntegrationActionExecutionProperties = executionProperties

	case n.UpsertEntity != nil:
		config.Type = consts.UpsertEntity
		u := n.UpsertEntity
		config.BlueprintIdentifier = u.BlueprintIdentifier.ValueStringPointer()
		config.OnFailure = u.OnFailure.ValueStringPointer()

		if u.Mapping != nil {
			mapping, err := upsertMappingToPortBody(ctx, u.Mapping)
			if err != nil {
				return nil, err
			}
			config.Mapping = mapping
		}

	case n.AI != nil:
		config.Type = consts.AI
		a := n.AI
		config.UserPrompt = a.UserPrompt.ValueStringPointer()
		config.SystemPrompt = a.SystemPrompt.ValueStringPointer()
		config.Provider = a.Provider.ValueStringPointer()
		config.Model = a.Model.ValueStringPointer()

		tools, err := terraformListToStrings(ctx, a.Tools)
		if err != nil {
			return nil, err
		}
		config.Tools = tools

		for _, s := range a.McpServers {
			config.McpServers = append(config.McpServers, cli.WorkflowMcpServer{Identifier: s.Identifier.ValueString()})
		}

		outputSchema, err := jsonStringToAny(a.OutputSchema)
		if err != nil {
			return nil, err
		}
		config.OutputSchema = outputSchema

	case n.AIAgent != nil:
		config.Type = consts.AIAgent
		a := n.AIAgent
		config.UserPrompt = a.UserPrompt.ValueStringPointer()
		config.AgentIdentifier = a.AgentIdentifier.ValueStringPointer()
		config.Provider = a.Provider.ValueStringPointer()
		config.Model = a.Model.ValueStringPointer()

		outputSchema, err := jsonStringToAny(a.OutputSchema)
		if err != nil {
			return nil, err
		}
		config.OutputSchema = outputSchema

	case n.Condition != nil:
		config.Type = consts.ConditionNode
		outlets := make([]cli.WorkflowOutlet, 0, len(n.Condition.Outlets))
		for _, o := range n.Condition.Outlets {
			outlets = append(outlets, cli.WorkflowOutlet{
				Identifier:          o.Identifier.ValueString(),
				Title:               o.Title.ValueStringPointer(),
				Expression:          o.Expression.ValueStringPointer(),
				StatusLabel:         statusLabelToPortBody(o.StatusLabel),
				WorkflowStatusLabel: statusLabelToPortBody(o.WorkflowStatusLabel),
			})
		}
		config.Outlets = &outlets

	case n.Input != nil:
		config.Type = consts.InputNode
		i := n.Input
		config.Description = i.Description.ValueStringPointer()

		if i.UserInputs != nil {
			userInputs, err := inputUserInputsToPortBody(ctx, i.UserInputs)
			if err != nil {
				return nil, err
			}
			config.UserInputs = userInputs
		}

		outlets := make([]cli.WorkflowOutlet, 0, len(i.Outlets))
		for _, o := range i.Outlets {
			evaluationMethod := "button"
			numOfResponders := o.NumOfResponders.ValueInt64()
			outlets = append(outlets, cli.WorkflowOutlet{
				Identifier:          o.Identifier.ValueString(),
				Title:               o.Title.ValueStringPointer(),
				EvaluationMethod:    &evaluationMethod,
				NumOfResponders:     &numOfResponders,
				StatusLabel:         statusLabelToPortBody(o.StatusLabel),
				WorkflowStatusLabel: statusLabelToPortBody(o.WorkflowStatusLabel),
			})
		}
		config.Outlets = &outlets

		for _, notification := range i.Notifications {
			converted, err := notificationToPortBody(ctx, notification)
			if err != nil {
				return nil, err
			}
			config.Notifications = append(config.Notifications, *converted)
		}

		responders, err := respondersToPortBody(ctx, i.Responders)
		if err != nil {
			return nil, err
		}
		config.Responders = responders
	}

	return config, nil
}

func selfServeUserInputsToPortBody(ctx context.Context, model *SelfServeUserInputsModel) (*cli.WorkflowUserInputs, error) {
	userInputs, err := userInputsToPortBody(ctx, model.UserProperties, model.Titles, model.RequiredJqQuery, model.OrderProperties)
	if err != nil {
		return nil, err
	}

	userInputs.Steps = userInputStepsToPortBody(model.Steps)
	userInputs.Validations = validationsToPortBody(model.Validations)

	return userInputs, nil
}

func inputUserInputsToPortBody(ctx context.Context, model *InputUserInputsModel) (*cli.WorkflowUserInputs, error) {
	userInputs, err := userInputsToPortBody(ctx, model.UserProperties, model.Titles, model.RequiredJqQuery, model.OrderProperties)
	if err != nil {
		return nil, err
	}

	userInputs.Steps = userInputStepsToPortBody(model.Steps)
	userInputs.Validations = validationsToPortBody(model.Validations)

	buttons := make([]cli.WorkflowInputButton, 0, len(model.Buttons))
	for _, b := range model.Buttons {
		buttons = append(buttons, cli.WorkflowInputButton{
			Identifier: b.Identifier.ValueString(),
			Label:      b.Label.ValueString(),
			Variant:    b.Variant.ValueString(),
			Icon:       b.Icon.ValueStringPointer(),
		})
	}
	userInputs.Buttons = &buttons

	return userInputs, nil
}

func userInputsToPortBody(
	ctx context.Context,
	userProperties *UserPropertiesModel,
	titles map[string]UserInputText,
	requiredJqQuery types.String,
	orderProperties types.List,
) (*cli.WorkflowUserInputs, error) {
	properties, required, err := userPropertiesToBody(ctx, userProperties)
	if err != nil {
		return nil, err
	}

	userInputs := &cli.WorkflowUserInputs{
		Properties: properties,
		Titles:     userInputTitlesToPortBody(titles),
	}

	if !requiredJqQuery.IsNull() {
		userInputs.Required = map[string]string{"jqQuery": requiredJqQuery.ValueString()}
	} else if len(required) > 0 {
		userInputs.Required = required
	}

	if !orderProperties.IsNull() {
		order, err := terraformListToStrings(ctx, orderProperties)
		if err != nil {
			return nil, err
		}
		userInputs.Order = order
	}

	return userInputs, nil
}

func userInputTitlesToPortBody(titles map[string]UserInputText) map[string]cli.ActionTitle {
	if titles == nil {
		return nil
	}

	result := make(map[string]cli.ActionTitle, len(titles))
	for identifier, title := range titles {
		converted := cli.ActionTitle{
			Title:       title.Title.ValueString(),
			Description: title.Description.ValueStringPointer(),
		}

		if !title.VisibleJqQuery.IsNull() {
			converted.Visible = map[string]string{"jqQuery": title.VisibleJqQuery.ValueString()}
		} else if !title.Visible.IsNull() {
			converted.Visible = title.Visible.ValueBool()
		}

		result[identifier] = converted
	}

	return result
}

func userInputStepsToPortBody(steps []UserInputsStepModel) []cli.WorkflowUserInputsStep {
	if steps == nil {
		return nil
	}

	result := make([]cli.WorkflowUserInputsStep, 0, len(steps))
	for _, s := range steps {
		order := make([]string, 0, len(s.Order))
		for _, p := range s.Order {
			order = append(order, p.ValueString())
		}

		step := cli.WorkflowUserInputsStep{
			Title:       s.Title.ValueString(),
			Order:       order,
			Validations: validationsToPortBody(s.Validations),
		}

		if !s.VisibleJqQuery.IsNull() {
			step.Visible = map[string]string{"jqQuery": s.VisibleJqQuery.ValueString()}
		} else if !s.Visible.IsNull() {
			step.Visible = s.Visible.ValueBool()
		}

		result = append(result, step)
	}

	return result
}

func validationsToPortBody(validations []InputValidationModel) []cli.WorkflowInputValidation {
	if validations == nil {
		return nil
	}

	result := make([]cli.WorkflowInputValidation, 0, len(validations))
	for _, v := range validations {
		result = append(result, cli.WorkflowInputValidation{
			Constraint: v.Constraint.ValueString(),
			Message:    v.Message.ValueString(),
		})
	}

	return result
}

func upsertMappingToPortBody(ctx context.Context, model *UpsertMappingModel) (*cli.WorkflowUpsertMapping, error) {
	mapping := &cli.WorkflowUpsertMapping{
		Identifier: model.Identifier.ValueStringPointer(),
		Title:      model.Title.ValueStringPointer(),
		Icon:       model.Icon.ValueStringPointer(),
	}

	if !model.Teams.IsNull() {
		teams, err := terraformListToStrings(ctx, model.Teams)
		if err != nil {
			return nil, err
		}
		mapping.Team = teams
	}

	if !model.Properties.IsNull() {
		properties, err := utils.TerraformStringToGoType[map[string]any](model.Properties)
		if err != nil {
			return nil, err
		}
		mapping.Properties = properties
	}

	if !model.Relations.IsNull() {
		relations, err := utils.TerraformStringToGoType[map[string]any](model.Relations)
		if err != nil {
			return nil, err
		}
		mapping.Relations = relations
	}

	return mapping, nil
}

func notificationToPortBody(ctx context.Context, model NotificationModel) (*cli.WorkflowInputNotification, error) {
	notification := &cli.WorkflowInputNotification{
		Target: model.Target.ValueString(),
		Url:    model.Url.ValueStringPointer(),
		Method: model.Method.ValueStringPointer(),
		Agent:  model.Agent.ValueBoolPointer(),
	}

	for _, f := range model.Fields {
		notification.Fields = append(notification.Fields, cli.WorkflowInputNotificationField{
			Label: f.Label.ValueString(),
			Value: f.Value.ValueString(),
		})
	}

	headers, err := terraformMapToStrings(ctx, model.Headers)
	if err != nil {
		return nil, err
	}
	notification.Headers = headers

	if !model.Body.IsNull() {
		body, err := utils.TerraformStringToGoType[map[string]any](model.Body)
		if err != nil {
			return nil, err
		}
		notification.Body = body
	}

	return notification, nil
}

func principalsToPortBody(ctx context.Context, users, roles, teams types.List) (*cli.WorkflowNodePermissions, error) {
	principals := &cli.WorkflowNodePermissions{}

	var err error
	if principals.Users, err = terraformListToStrings(ctx, users); err != nil {
		return nil, err
	}
	if principals.Roles, err = terraformListToStrings(ctx, roles); err != nil {
		return nil, err
	}
	if principals.Teams, err = terraformListToStrings(ctx, teams); err != nil {
		return nil, err
	}

	return principals, nil
}

func permissionsToPortBody(ctx context.Context, model *PermissionsModel) (*cli.WorkflowNodePermissions, error) {
	if model == nil {
		return nil, nil
	}

	permissions, err := principalsToPortBody(ctx, model.Users, model.Roles, model.Teams)
	if err != nil {
		return nil, err
	}

	if !model.Policy.IsNull() {
		policy, err := utils.TerraformStringToGoType[any](model.Policy)
		if err != nil {
			return nil, err
		}
		permissions.Policy = policy
	}

	return permissions, nil
}

func respondersToPortBody(ctx context.Context, model *RespondersModel) (*cli.WorkflowNodePermissions, error) {
	if model == nil {
		return nil, nil
	}

	responders, err := principalsToPortBody(ctx, model.Users, model.Roles, model.Teams)
	if err != nil {
		return nil, err
	}

	if !model.UsersQuery.IsNull() {
		usersQuery, err := utils.TerraformStringToGoType[any](model.UsersQuery)
		if err != nil {
			return nil, err
		}
		responders.UsersQuery = usersQuery
	}

	return responders, nil
}

func statusLabelToPortBody(model *StatusLabelModel) *cli.WorkflowStatusLabel {
	if model == nil {
		return nil
	}

	return &cli.WorkflowStatusLabel{
		Text:    model.Text.ValueString(),
		Variant: model.Variant.ValueStringPointer(),
	}
}

func connectionsToPortBody(connections []ConnectionModel) []cli.WorkflowConnection {
	result := make([]cli.WorkflowConnection, 0, len(connections))
	for _, c := range connections {
		result = append(result, cli.WorkflowConnection{
			SourceIdentifier:       c.SourceIdentifier.ValueString(),
			TargetIdentifier:       c.TargetIdentifier.ValueString(),
			Description:            c.Description.ValueStringPointer(),
			SourceOutletIdentifier: c.SourceOutletIdentifier.ValueStringPointer(),
			Fallback:               c.Fallback.ValueBoolPointer(),
		})
	}
	return result
}

func jsonStringToAny(value types.String) (any, error) {
	if value.IsNull() {
		return nil, nil
	}

	return utils.TerraformStringToGoType[any](value)
}

func terraformListToStrings(ctx context.Context, list types.List) ([]string, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	values, err := utils.TerraformListToGoArray(ctx, list, "string")
	if err != nil {
		return nil, err
	}

	return utils.InterfaceToStringArray(values), nil
}

func terraformMapToStrings(ctx context.Context, m types.Map) (map[string]string, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}

	result := make(map[string]string, len(m.Elements()))
	if diags := m.ElementsAs(ctx, &result, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert map to string map: %s", diags.Errors()[0].Detail())
	}

	return result, nil
}
