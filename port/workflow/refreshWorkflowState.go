package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/action"
)

func (r *WorkflowResource) refreshWorkflowState(ctx context.Context, state *WorkflowModel, w *cli.Workflow) error {
	// The API returns nodes topologically sorted, so restore the prior order.
	priorNodes := make(map[string]WorkflowNodeModel, len(state.Nodes))
	for _, n := range state.Nodes {
		priorNodes[n.Identifier.ValueString()] = n
	}

	state.ID = types.StringValue(w.Identifier)
	state.Identifier = types.StringValue(w.Identifier)
	state.Title = flex.GoStringToFramework(w.Title)
	state.Icon = flex.GoStringToFramework(w.Icon)
	state.Description = flex.GoStringToFramework(w.Description)
	state.Category = flex.GoStringToFramework(w.Category)
	state.AllowAnyoneToViewRuns = flex.GoBoolToFramework(w.AllowAnyoneToViewRuns)

	orderedNodes := reorderNodes(state.Nodes, w.Nodes)
	nodes := make([]WorkflowNodeModel, 0, len(orderedNodes))
	for _, apiNode := range orderedNodes {
		node, err := r.nodeToModel(ctx, apiNode, priorNodes[apiNode.Identifier])
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}
	state.Nodes = nodes

	state.Connections = reorderConnections(state.Connections, w.Connections)

	return nil
}

func derefSlice[T any](values *[]T) []T {
	if values == nil {
		return nil
	}
	return *values
}

// The API defaults these fields to an empty string and the schema declares the
// same default, so an absent value has to refresh as "" rather than null to keep
// the state consistent with the plan.
func stringOrEmpty(value *string) types.String {
	if value == nil {
		return types.StringValue("")
	}
	return types.StringValue(*value)
}

func (r *WorkflowResource) nodeToModel(ctx context.Context, apiNode cli.WorkflowNode, prior WorkflowNodeModel) (WorkflowNodeModel, error) {
	node := WorkflowNodeModel{
		Identifier:  types.StringValue(apiNode.Identifier),
		Title:       flex.GoStringToFramework(apiNode.Title),
		Icon:        flex.GoStringToFramework(apiNode.Icon),
		Description: flex.GoStringToFramework(apiNode.Description),
		Verbose:     flex.GoBoolToFramework(apiNode.Verbose),
		Links:       stringsToList(ctx, apiNode.Links, prior.Links),
		Variables:   stringsToMap(ctx, apiNode.Variables, prior.Variables),
	}

	config := apiNode.Config
	switch config.Type {
	case consts.SelfServeTrigger:
		trigger := &SelfServeTriggerModel{
			ActionCardButtonText:    flex.GoStringToFramework(config.ActionCardButtonText),
			ExecuteActionButtonText: flex.GoStringToFramework(config.ExecuteActionButtonText),
			Variant:                 flex.GoStringToFramework(config.Variant),
			Published:               flex.GoBoolToFramework(config.Published),
			Permissions:             permissionsToModel(ctx, config.Permissions, r.portClient.JSONEscapeHTML),
		}

		if config.UserInputs != nil {
			configured := prior.SelfServeTrigger != nil && prior.SelfServeTrigger.UserInputs != nil
			userInputs, err := r.userInputsToModel(ctx, config.UserInputs, configured)
			if err != nil {
				return node, err
			}
			trigger.UserInputs = &SelfServeUserInputsModel{
				UserProperties:  userInputs.UserProperties,
				Titles:          userInputs.Titles,
				RequiredJqQuery: userInputs.RequiredJqQuery,
				OrderProperties: userInputs.OrderProperties,
				Steps:           userInputs.Steps,
			}
		}

		for _, c := range config.Contexts {
			trigger.Contexts = append(trigger.Contexts, TriggerContextModel{
				On:                  types.StringValue(c.On),
				BlueprintIdentifier: flex.GoStringToFramework(c.BlueprintIdentifier),
				UserInput:           flex.GoStringToFramework(c.UserInput),
			})
		}

		node.SelfServeTrigger = trigger

	case consts.EventTrigger:
		trigger := &EventTriggerModel{
			Published: flex.GoBoolToFramework(config.Published),
		}
		if config.Event != nil {
			trigger.Type = types.StringValue(config.Event.Type)
			trigger.BlueprintIdentifier = types.StringValue(config.Event.BlueprintIdentifier)
			trigger.PropertyIdentifier = flex.GoStringToFramework(config.Event.PropertyIdentifier)
		}
		if config.Condition != nil {
			trigger.Condition = &NodeConditionModel{
				Expressions: flex.GoArrayStringToTerraformList(ctx, config.Condition.Expressions),
				Combinator:  flex.GoStringToFramework(config.Condition.Combinator),
			}
		}
		node.EventTrigger = trigger

	case consts.ScheduleTrigger:
		node.ScheduleTrigger = &ScheduleTriggerModel{
			Cron:      flex.GoStringToFramework(config.Cron),
			Published: flex.GoBoolToFramework(config.Published),
		}

	case consts.Kafka:
		payload, err := anyToJSONString(config.Payload, r.portClient.JSONEscapeHTML)
		if err != nil {
			return node, err
		}
		node.Kafka = &KafkaModel{
			Payload:   payload,
			OnFailure: flex.GoStringToFramework(config.OnFailure),
		}

	case consts.Webhook:
		body, err := anyToJSONString(config.Body, r.portClient.JSONEscapeHTML)
		if err != nil {
			return node, err
		}
		var priorHeaders types.Map
		if prior.Webhook != nil {
			priorHeaders = prior.Webhook.Headers
		}
		node.Webhook = &WebhookModel{
			Url:          flex.GoStringToFramework(config.Url),
			Agent:        flex.GoBoolToFramework(config.Agent),
			Synchronized: flex.GoBoolToFramework(config.Synchronized),
			Method:       flex.GoStringToFramework(config.Method),
			Headers:      stringsToMap(ctx, config.Headers, priorHeaders),
			Body:         body,
			OnTimeout:    flex.GoStringToFramework(config.OnTimeout),
			OnFailure:    flex.GoStringToFramework(config.OnFailure),
		}

	case consts.IntegrationAction:
		executionProperties, err := anyToJSONString(config.IntegrationActionExecutionProperties, r.portClient.JSONEscapeHTML)
		if err != nil {
			return node, err
		}
		node.IntegrationAction = &IntegrationActionModel{
			InstallationId:            flex.GoStringToFramework(config.InstallationId),
			IntegrationProvider:       flex.GoStringToFramework(config.IntegrationProvider),
			IntegrationInvocationType: flex.GoStringToFramework(config.IntegrationInvocationType),
			ExecutionProperties:       executionProperties,
			OnFailure:                 flex.GoStringToFramework(config.OnFailure),
		}

	case consts.UpsertEntity:
		upsert := &UpsertEntityModel{
			BlueprintIdentifier: flex.GoStringToFramework(config.BlueprintIdentifier),
			OnFailure:           flex.GoStringToFramework(config.OnFailure),
		}
		if config.Mapping != nil {
			mapping, err := upsertMappingToModel(ctx, config.Mapping, r.portClient.JSONEscapeHTML)
			if err != nil {
				return node, err
			}
			upsert.Mapping = mapping
		}
		node.UpsertEntity = upsert

	case consts.AI:
		outputSchema, err := anyToJSONString(config.OutputSchema, r.portClient.JSONEscapeHTML)
		if err != nil {
			return node, err
		}
		var priorTools types.List
		if prior.AI != nil {
			priorTools = prior.AI.Tools
		}
		ai := &AIModel{
			UserPrompt:   flex.GoStringToFramework(config.UserPrompt),
			SystemPrompt: stringOrEmpty(config.SystemPrompt),
			Provider:     flex.GoStringToFramework(config.Provider),
			Model:        flex.GoStringToFramework(config.Model),
			Tools:        stringsToList(ctx, config.Tools, priorTools),
			OutputSchema: outputSchema,
		}
		for _, s := range config.McpServers {
			ai.McpServers = append(ai.McpServers, McpServerModel{Identifier: types.StringValue(s.Identifier)})
		}
		node.AI = ai

	case consts.AIAgent:
		outputSchema, err := anyToJSONString(config.OutputSchema, r.portClient.JSONEscapeHTML)
		if err != nil {
			return node, err
		}
		node.AIAgent = &AIAgentModel{
			UserPrompt:      flex.GoStringToFramework(config.UserPrompt),
			AgentIdentifier: flex.GoStringToFramework(config.AgentIdentifier),
			Provider:        flex.GoStringToFramework(config.Provider),
			Model:           flex.GoStringToFramework(config.Model),
			OutputSchema:    outputSchema,
		}

	case consts.ConditionNode:
		condition := &ConditionModel{}
		for _, o := range derefSlice(config.Outlets) {
			condition.Outlets = append(condition.Outlets, ConditionOutletModel{
				Identifier:          types.StringValue(o.Identifier),
				Title:               stringOrEmpty(o.Title),
				Expression:          flex.GoStringToFramework(o.Expression),
				StatusLabel:         statusLabelToModel(o.StatusLabel),
				WorkflowStatusLabel: statusLabelToModel(o.WorkflowStatusLabel),
			})
		}
		node.Condition = condition

	case consts.InputNode:
		input := &InputModel{
			Description: flex.GoStringToFramework(config.Description),
			Responders:  respondersToModel(ctx, config.Responders, r.portClient.JSONEscapeHTML),
		}

		if config.UserInputs != nil {
			configured := prior.Input != nil && prior.Input.UserInputs != nil
			userInputs, err := r.userInputsToModel(ctx, config.UserInputs, configured)
			if err != nil {
				return node, err
			}
			input.UserInputs = &InputUserInputsModel{
				UserProperties:  userInputs.UserProperties,
				Titles:          userInputs.Titles,
				RequiredJqQuery: userInputs.RequiredJqQuery,
				OrderProperties: userInputs.OrderProperties,
				Steps:           userInputs.Steps,
			}
			for _, b := range derefSlice(config.UserInputs.Buttons) {
				input.UserInputs.Buttons = append(input.UserInputs.Buttons, InputButtonModel{
					Identifier: types.StringValue(b.Identifier),
					Label:      types.StringValue(b.Label),
					Variant:    types.StringValue(b.Variant),
					Icon:       flex.GoStringToFramework(b.Icon),
				})
			}
		}

		for _, o := range derefSlice(config.Outlets) {
			outlet := InputOutletModel{
				Identifier:          types.StringValue(o.Identifier),
				Title:               stringOrEmpty(o.Title),
				StatusLabel:         statusLabelToModel(o.StatusLabel),
				WorkflowStatusLabel: statusLabelToModel(o.WorkflowStatusLabel),
			}
			if o.NumOfResponders != nil {
				outlet.NumOfResponders = types.Int64Value(*o.NumOfResponders)
			}
			input.Outlets = append(input.Outlets, outlet)
		}

		for i, n := range config.Notifications {
			var priorNotification *NotificationModel
			if prior.Input != nil && i < len(prior.Input.Notifications) {
				priorNotification = &prior.Input.Notifications[i]
			}
			notification, err := notificationToModel(ctx, n, priorNotification, r.portClient.JSONEscapeHTML)
			if err != nil {
				return node, err
			}
			input.Notifications = append(input.Notifications, *notification)
		}

		node.Input = input
	}

	return node, nil
}

type userInputsModel struct {
	UserProperties  *action.UserPropertiesModel
	Titles          map[string]action.ActionTitle
	RequiredJqQuery types.String
	OrderProperties types.List
	Steps           []action.Step
}

func (r *WorkflowResource) userInputsToModel(ctx context.Context, userInputs *cli.WorkflowUserInputs, configured bool) (*userInputsModel, error) {
	apiUserInputs := &cli.ActionUserInputs{
		Properties: userInputs.Properties,
		Required:   userInputs.Required,
		Order:      userInputs.Order,
		Steps:      userInputs.Steps,
		Titles:     userInputs.Titles,
	}

	mapper := action.NewUserInputsMapper(r.portClient)

	userProperties, err := mapper.UserPropertiesToState(ctx, apiUserInputs, configured)
	if err != nil {
		return nil, err
	}

	titles, err := mapper.UserInputTitlesToState(apiUserInputs)
	if err != nil {
		return nil, err
	}

	return &userInputsModel{
		UserProperties:  userProperties,
		Titles:          titles,
		RequiredJqQuery: action.UserInputsRequiredToState(apiUserInputs),
		OrderProperties: action.UserInputsOrderToState(ctx, apiUserInputs.Order),
		Steps:           action.UserInputStepsToState(apiUserInputs.Steps),
	}, nil
}

func upsertMappingToModel(ctx context.Context, mapping *cli.WorkflowUpsertMapping, jsonEscapeHTML bool) (*UpsertMappingModel, error) {
	model := &UpsertMappingModel{
		Identifier: flex.GoStringToFramework(mapping.Identifier),
		Title:      flex.GoStringToFramework(mapping.Title),
		Icon:       flex.GoStringToFramework(mapping.Icon),
		Teams:      types.ListNull(types.StringType),
	}

	switch team := mapping.Team.(type) {
	case string:
		model.Teams = flex.GoArrayStringToTerraformList(ctx, []string{team})
	case []any:
		model.Teams = flex.GoArrayStringToTerraformList(ctx, utils.InterfaceToStringArray(team))
	case []string:
		model.Teams = flex.GoArrayStringToTerraformList(ctx, team)
	}

	if mapping.Properties != nil {
		properties, err := utils.GoObjectToTerraformString(mapping.Properties, jsonEscapeHTML)
		if err != nil {
			return nil, err
		}
		model.Properties = properties
	}

	if mapping.Relations != nil {
		relations, err := utils.GoObjectToTerraformString(mapping.Relations, jsonEscapeHTML)
		if err != nil {
			return nil, err
		}
		model.Relations = relations
	}

	return model, nil
}

func notificationToModel(ctx context.Context, notification cli.WorkflowInputNotification, prior *NotificationModel, jsonEscapeHTML bool) (*NotificationModel, error) {
	var priorHeaders types.Map
	if prior != nil {
		priorHeaders = prior.Headers
	}

	model := &NotificationModel{
		Target:  types.StringValue(notification.Target),
		Url:     flex.GoStringToFramework(notification.Url),
		Method:  flex.GoStringToFramework(notification.Method),
		Agent:   flex.GoBoolToFramework(notification.Agent),
		Headers: stringsToMap(ctx, notification.Headers, priorHeaders),
	}

	for _, f := range notification.Fields {
		model.Fields = append(model.Fields, NotificationFieldModel{
			Label: types.StringValue(f.Label),
			Value: types.StringValue(f.Value),
		})
	}

	if notification.Body != nil {
		body, err := utils.GoObjectToTerraformString(notification.Body, jsonEscapeHTML)
		if err != nil {
			return nil, err
		}
		model.Body = body
	}

	return model, nil
}

// An empty permissions object is equivalent to an omitted block.
func hasNoPrincipals(permissions *cli.WorkflowNodePermissions) bool {
	return permissions == nil ||
		(len(permissions.Users) == 0 && len(permissions.Roles) == 0 && len(permissions.Teams) == 0 &&
			permissions.Policy == nil && permissions.UsersQuery == nil)
}

func permissionsToModel(ctx context.Context, permissions *cli.WorkflowNodePermissions, jsonEscapeHTML bool) *PermissionsModel {
	if hasNoPrincipals(permissions) {
		return nil
	}

	model := &PermissionsModel{
		Users: stringsToList(ctx, permissions.Users, types.ListNull(types.StringType)),
		Roles: stringsToList(ctx, permissions.Roles, types.ListNull(types.StringType)),
		Teams: stringsToList(ctx, permissions.Teams, types.ListNull(types.StringType)),
	}

	if permissions.Policy != nil {
		if policy, err := utils.GoObjectToTerraformString(permissions.Policy, jsonEscapeHTML); err == nil {
			model.Policy = policy
		}
	}

	return model
}

func respondersToModel(ctx context.Context, responders *cli.WorkflowNodePermissions, jsonEscapeHTML bool) *RespondersModel {
	if hasNoPrincipals(responders) {
		return nil
	}

	model := &RespondersModel{
		Users: stringsToList(ctx, responders.Users, types.ListNull(types.StringType)),
		Roles: stringsToList(ctx, responders.Roles, types.ListNull(types.StringType)),
		Teams: stringsToList(ctx, responders.Teams, types.ListNull(types.StringType)),
	}

	if responders.UsersQuery != nil {
		if usersQuery, err := utils.GoObjectToTerraformString(responders.UsersQuery, jsonEscapeHTML); err == nil {
			model.UsersQuery = usersQuery
		}
	}

	return model
}

func statusLabelToModel(label *cli.WorkflowStatusLabel) *StatusLabelModel {
	if label == nil {
		return nil
	}

	return &StatusLabelModel{
		Text:    types.StringValue(label.Text),
		Variant: flex.GoStringToFramework(label.Variant),
	}
}

func anyToJSONString(value any, jsonEscapeHTML bool) (types.String, error) {
	if value == nil {
		return types.StringNull(), nil
	}

	return utils.GoObjectToTerraformString(value, jsonEscapeHTML)
}

// stringsToList keeps the prior value when the API returns an empty collection,
// so that an unset optional attribute is not refreshed into an empty list.
func stringsToList(ctx context.Context, values []string, prior types.List) types.List {
	if len(values) == 0 {
		if !prior.IsNull() && !prior.IsUnknown() {
			return prior
		}
		return types.ListNull(types.StringType)
	}

	return flex.GoArrayStringToTerraformList(ctx, values)
}

func stringsToMap(ctx context.Context, values map[string]string, prior types.Map) types.Map {
	if len(values) == 0 {
		if !prior.IsNull() && !prior.IsUnknown() {
			return prior
		}
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(values))
	for k, v := range values {
		elements[k] = types.StringValue(v)
	}

	result, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}

	return result
}

func reorderNodes(prior []WorkflowNodeModel, apiNodes []cli.WorkflowNode) []cli.WorkflowNode {
	byID := make(map[string]cli.WorkflowNode, len(apiNodes))
	for _, n := range apiNodes {
		byID[n.Identifier] = n
	}

	ordered := make([]cli.WorkflowNode, 0, len(apiNodes))
	seen := make(map[string]bool, len(apiNodes))
	for _, p := range prior {
		id := p.Identifier.ValueString()
		if n, ok := byID[id]; ok && !seen[id] {
			ordered = append(ordered, n)
			seen[id] = true
		}
	}
	for _, n := range apiNodes {
		if !seen[n.Identifier] {
			ordered = append(ordered, n)
			seen[n.Identifier] = true
		}
	}
	return ordered
}

func reorderConnections(prior []ConnectionModel, apiConnections []cli.WorkflowConnection) []ConnectionModel {
	// A node with multiple outlets can connect to the same target more than
	// once, so the outlet and fallback are part of a connection's identity.
	key := func(source, target string, outlet *string, fallback *bool) string {
		outletKey := ""
		if outlet != nil {
			outletKey = *outlet
		}
		fallbackKey := ""
		if fallback != nil && *fallback {
			fallbackKey = "fallback"
		}
		return source + "\x00" + target + "\x00" + outletKey + "\x00" + fallbackKey
	}

	byKey := make(map[string]cli.WorkflowConnection, len(apiConnections))
	for _, c := range apiConnections {
		byKey[key(c.SourceIdentifier, c.TargetIdentifier, c.SourceOutletIdentifier, c.Fallback)] = c
	}

	ordered := make([]ConnectionModel, 0, len(apiConnections))
	seen := make(map[string]bool, len(apiConnections))
	appendConnection := func(c cli.WorkflowConnection) {
		ordered = append(ordered, ConnectionModel{
			SourceIdentifier:       types.StringValue(c.SourceIdentifier),
			TargetIdentifier:       types.StringValue(c.TargetIdentifier),
			Description:            flex.GoStringToFramework(c.Description),
			SourceOutletIdentifier: flex.GoStringToFramework(c.SourceOutletIdentifier),
			Fallback:               flex.GoBoolToFramework(c.Fallback),
		})
	}

	for _, p := range prior {
		k := key(
			p.SourceIdentifier.ValueString(),
			p.TargetIdentifier.ValueString(),
			p.SourceOutletIdentifier.ValueStringPointer(),
			p.Fallback.ValueBoolPointer(),
		)
		if c, ok := byKey[k]; ok && !seen[k] {
			appendConnection(c)
			seen[k] = true
		}
	}
	for _, c := range apiConnections {
		k := key(c.SourceIdentifier, c.TargetIdentifier, c.SourceOutletIdentifier, c.Fallback)
		if !seen[k] {
			appendConnection(c)
			seen[k] = true
		}
	}
	return ordered
}
