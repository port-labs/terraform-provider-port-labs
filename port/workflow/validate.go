package workflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

var _ resource.ResourceWithValidateConfig = &WorkflowResource{}

// Node config attributes cannot be marked Required in the schema: Terraform
// enforces required attributes of a SingleNestedBlock even when the block itself
// is absent, which would make every node config block mandatory.
func (r *WorkflowResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data WorkflowModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeTypes := make(map[string]string, len(data.Nodes))
	outlets := make(map[string]map[string]bool, len(data.Nodes))
	for i, node := range data.Nodes {
		nodePath := path.Root("node").AtListIndex(i)
		identifier := node.Identifier.ValueString()

		if _, duplicate := nodeTypes[identifier]; duplicate && !node.Identifier.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				nodePath.AtName("identifier"),
				"Duplicate node identifier",
				fmt.Sprintf("Node identifier %q is used more than once. Each node must have a unique identifier.", identifier),
			)
		}

		nodeTypes[identifier] = validateNode(resp, nodePath, node)
		outlets[identifier] = nodeOutlets(node)
	}

	validateTriggerPresence(resp, nodeTypes)
	validateConnections(resp, data.Connections, nodeTypes, outlets)
}

var triggerTypes = map[string]bool{
	consts.SelfServeTrigger: true,
	consts.EventTrigger:     true,
	consts.ScheduleTrigger:  true,
}

func validateTriggerPresence(resp *resource.ValidateConfigResponse, nodeTypes map[string]string) {
	for _, nodeType := range nodeTypes {
		if triggerTypes[nodeType] {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("node"),
		"Missing trigger node",
		"A workflow must define at least one trigger node: `self_serve_trigger`, `event_trigger` or `schedule_trigger`.",
	)
}

func nodeOutlets(node WorkflowNodeModel) map[string]bool {
	identifiers := map[string]bool{}
	if node.Condition != nil {
		for _, outlet := range node.Condition.Outlets {
			identifiers[outlet.Identifier.ValueString()] = true
		}
	}
	if node.Input != nil {
		for _, outlet := range node.Input.Outlets {
			identifiers[outlet.Identifier.ValueString()] = true
		}
	}
	return identifiers
}

func validateNode(resp *resource.ValidateConfigResponse, nodePath path.Path, node WorkflowNodeModel) string {
	switch {
	case node.SelfServeTrigger != nil:
		blockPath := nodePath.AtName("self_serve_trigger")
		for i, context := range node.SelfServeTrigger.Contexts {
			contextPath := blockPath.AtName("contexts").AtListIndex(i)
			switch context.On.ValueString() {
			case consts.CreateEntityContext:
				requireSet(resp, contextPath.AtName("blueprint_identifier"), context.BlueprintIdentifier,
					"`blueprint_identifier` is required when `on` is `CREATE_ENTITY`.")
				rejectSet(resp, contextPath.AtName("user_input"), context.UserInput,
					"`user_input` cannot be set when `on` is `CREATE_ENTITY`.")
			case consts.EntityContext:
				requireNonEmpty(resp, contextPath.AtName("user_input"), context.UserInput,
					"`user_input` is required when `on` is `ENTITY`.")
				rejectSet(resp, contextPath.AtName("blueprint_identifier"), context.BlueprintIdentifier,
					"`blueprint_identifier` cannot be set when `on` is `ENTITY`.")
			}
		}
		return consts.SelfServeTrigger

	case node.EventTrigger != nil:
		blockPath := nodePath.AtName("event_trigger")
		requireSet(resp, blockPath.AtName("type"), node.EventTrigger.Type, "`type` is required for an `event_trigger` node.")
		requireSet(resp, blockPath.AtName("blueprint_identifier"), node.EventTrigger.BlueprintIdentifier,
			"`blueprint_identifier` is required for an `event_trigger` node.")

		if node.EventTrigger.Type.ValueString() == consts.WorkflowTimerExpired {
			requireSet(resp, blockPath.AtName("property_identifier"), node.EventTrigger.PropertyIdentifier,
				"`property_identifier` is required when `type` is `TIMER_EXPIRED`.")
		}

		if node.EventTrigger.Condition != nil && node.EventTrigger.Condition.Expressions.IsNull() {
			resp.Diagnostics.AddAttributeError(
				blockPath.AtName("condition").AtName("expressions"),
				"Missing required attribute",
				"`expressions` is required when a `condition` block is set.",
			)
		}
		return consts.EventTrigger

	case node.ScheduleTrigger != nil:
		requireSet(resp, nodePath.AtName("schedule_trigger").AtName("cron"), node.ScheduleTrigger.Cron,
			"`cron` is required for a `schedule_trigger` node.")
		return consts.ScheduleTrigger

	case node.Kafka != nil:
		return consts.Kafka

	case node.Webhook != nil:
		requireSet(resp, nodePath.AtName("webhook").AtName("url"), node.Webhook.Url,
			"`url` is required for a `webhook` node.")
		return consts.Webhook

	case node.IntegrationAction != nil:
		blockPath := nodePath.AtName("integration_action")
		requireSet(resp, blockPath.AtName("installation_id"), node.IntegrationAction.InstallationId,
			"`installation_id` is required for an `integration_action` node.")
		requireSet(resp, blockPath.AtName("integration_provider"), node.IntegrationAction.IntegrationProvider,
			"`integration_provider` is required for an `integration_action` node.")
		requireSet(resp, blockPath.AtName("integration_invocation_type"), node.IntegrationAction.IntegrationInvocationType,
			"`integration_invocation_type` is required for an `integration_action` node.")
		return consts.IntegrationAction

	case node.UpsertEntity != nil:
		requireSet(resp, nodePath.AtName("upsert_entity").AtName("blueprint_identifier"), node.UpsertEntity.BlueprintIdentifier,
			"`blueprint_identifier` is required for an `upsert_entity` node.")
		return consts.UpsertEntity

	case node.AI != nil:
		blockPath := nodePath.AtName("ai")
		requireSet(resp, blockPath.AtName("user_prompt"), node.AI.UserPrompt,
			"`user_prompt` is required for an `ai` node.")
		requirePaired(resp, blockPath, "ai", node.AI.Provider, node.AI.Model)
		return consts.AI

	case node.AIAgent != nil:
		blockPath := nodePath.AtName("ai_agent")
		requireSet(resp, blockPath.AtName("user_prompt"), node.AIAgent.UserPrompt,
			"`user_prompt` is required for an `ai_agent` node.")
		requireSet(resp, blockPath.AtName("agent_identifier"), node.AIAgent.AgentIdentifier,
			"`agent_identifier` is required for an `ai_agent` node.")
		requirePaired(resp, blockPath, "ai_agent", node.AIAgent.Provider, node.AIAgent.Model)
		return consts.AIAgent

	case node.Condition != nil:
		blockPath := nodePath.AtName("condition")
		// The API requires the `outlets` key, and the request body omits empty
		// lists, so a node without outlets is rejected on save either way.
		if len(node.Condition.Outlets) == 0 {
			resp.Diagnostics.AddAttributeError(
				blockPath.AtName("outlets"),
				"Missing required block",
				"A `condition` node must define at least one `outlets` block.",
			)
		}

		seenOutlets := map[string]bool{}
		for i, outlet := range node.Condition.Outlets {
			outletPath := blockPath.AtName("outlets").AtListIndex(i)
			requireUnique(resp, outletPath.AtName("identifier"), outlet.Identifier, seenOutlets, "outlet")
			requireNonEmpty(resp, outletPath.AtName("expression"), outlet.Expression,
				"`expression` cannot be empty.")
			validateStatusLabels(resp, outletPath, outlet.StatusLabel, outlet.WorkflowStatusLabel)
		}
		return consts.ConditionNode

	case node.Input != nil:
		blockPath := nodePath.AtName("input")
		buttons := map[string]bool{}
		// As with condition outlets, the API requires the `buttons` and `outlets`
		// keys while the request body omits empty lists, so an empty one is
		// rejected on save regardless.
		if node.Input.UserInputs == nil || len(node.Input.UserInputs.Buttons) == 0 {
			resp.Diagnostics.AddAttributeError(
				blockPath.AtName("user_inputs").AtName("buttons"),
				"Missing required attribute",
				"An `input` node must define `user_inputs.buttons`.",
			)
		} else {
			seenButtons := map[string]bool{}
			for i, button := range node.Input.UserInputs.Buttons {
				buttonPath := blockPath.AtName("user_inputs").AtName("buttons").AtListIndex(i)
				requireUnique(resp, buttonPath.AtName("identifier"), button.Identifier, seenButtons, "button")
				buttons[button.Identifier.ValueString()] = true
			}
		}

		if len(node.Input.Outlets) == 0 {
			resp.Diagnostics.AddAttributeError(
				blockPath.AtName("outlets"),
				"Missing required block",
				"An `input` node must define at least one `outlets` block.",
			)
		}

		seenOutlets := map[string]bool{}
		for i, outlet := range node.Input.Outlets {
			outletPath := blockPath.AtName("outlets").AtListIndex(i)
			identifier := outlet.Identifier.ValueString()
			requireUnique(resp, outletPath.AtName("identifier"), outlet.Identifier, seenOutlets, "outlet")
			if len(buttons) > 0 && !outlet.Identifier.IsUnknown() && !buttons[identifier] {
				resp.Diagnostics.AddAttributeError(
					outletPath.AtName("identifier"),
					"Unknown button identifier",
					fmt.Sprintf("Outlet %q does not match any button defined in `user_inputs.buttons`.", identifier),
				)
			}
			validateStatusLabels(resp, outletPath, outlet.StatusLabel, outlet.WorkflowStatusLabel)
		}

		for i, notification := range node.Input.Notifications {
			notificationPath := blockPath.AtName("notifications").AtListIndex(i)
			switch notification.Target.ValueString() {
			case "webhook":
				requireSet(resp, notificationPath.AtName("url"), notification.Url,
					"`url` is required when `target` is `webhook`.")
			case "email":
				if len(notification.Fields) == 0 {
					resp.Diagnostics.AddAttributeError(
						notificationPath.AtName("fields"),
						"Missing required block",
						"At least one `fields` block is required when `target` is `email`.",
					)
				}
			}
		}
		return consts.InputNode
	}

	return ""
}

func validateStatusLabels(resp *resource.ValidateConfigResponse, outletPath path.Path, labels ...*StatusLabelModel) {
	names := []string{"status_label", "workflow_status_label"}
	for i, label := range labels {
		if label != nil {
			requireNonEmpty(resp, outletPath.AtName(names[i]).AtName("text"), label.Text,
				fmt.Sprintf("`text` is required when a `%s` block is set.", names[i]))
		}
	}
}

func validateConnections(resp *resource.ValidateConfigResponse, connections []ConnectionModel, nodeTypes map[string]string, outlets map[string]map[string]bool) {
	multipleOutletTypes := map[string]bool{
		consts.ConditionNode: true,
		consts.InputNode:     true,
	}

	// Only condition and input nodes may have several outgoing connections, and
	// each must leave through a distinct outlet. Every other node is limited to a
	// single outgoing connection, so fan-out has to go through a branching node.
	seenSources := map[string]bool{}

	for i, connection := range connections {
		connectionPath := path.Root("connections").AtListIndex(i)
		source := connection.SourceIdentifier.ValueString()
		target := connection.TargetIdentifier.ValueString()

		if source == target && !connection.SourceIdentifier.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				connectionPath,
				"Invalid connection",
				fmt.Sprintf("Node %q cannot connect to itself. Self-connections are not supported.", source),
			)
		}

		validateEndpoint(resp, connectionPath, "source_identifier", connection.SourceIdentifier, nodeTypes)
		validateEndpoint(resp, connectionPath, "target_identifier", connection.TargetIdentifier, nodeTypes)

		if targetType, known := nodeTypes[target]; known && triggerTypes[targetType] {
			resp.Diagnostics.AddAttributeError(
				connectionPath.AtName("target_identifier"),
				"Invalid connection",
				fmt.Sprintf("Node %q is a trigger. Trigger nodes start a workflow, so they cannot have incoming connections.", target),
			)
		}

		sourceType, known := nodeTypes[source]
		if !known {
			continue
		}

		hasOutlet := !connection.SourceOutletIdentifier.IsNull() && !connection.SourceOutletIdentifier.IsUnknown()
		hasFallback := connection.Fallback.ValueBool()

		if hasOutlet && len(outlets[source]) > 0 && !outlets[source][connection.SourceOutletIdentifier.ValueString()] {
			resp.Diagnostics.AddAttributeError(
				connectionPath.AtName("source_outlet_identifier"),
				"Unknown outlet",
				fmt.Sprintf("Node %q has no outlet named %q.", source, connection.SourceOutletIdentifier.ValueString()),
			)
		}

		sourceKey := source
		if multipleOutletTypes[sourceType] {
			sourceKey = fmt.Sprintf("%s\x00%s\x00%t", source, connection.SourceOutletIdentifier.ValueString(), hasFallback)
		}
		if seenSources[sourceKey] {
			detail := fmt.Sprintf("Node %q already has an outgoing connection. A %q node supports only one.", source, sourceType)
			if multipleOutletTypes[sourceType] {
				detail = fmt.Sprintf("Node %q already has an outgoing connection using this outlet. Each connection leaving a %q node must use a distinct outlet.", source, sourceType)
			}
			resp.Diagnostics.AddAttributeError(connectionPath, "Duplicate connection source", detail)
		}
		seenSources[sourceKey] = true

		if multipleOutletTypes[sourceType] {
			if !hasOutlet && !hasFallback {
				resp.Diagnostics.AddAttributeError(
					connectionPath.AtName("source_outlet_identifier"),
					"Missing outlet",
					fmt.Sprintf("Connections leaving a %q node must set `source_outlet_identifier` or `fallback`.", sourceType),
				)
			}
			if hasFallback && sourceType == consts.InputNode {
				resp.Diagnostics.AddAttributeError(
					connectionPath.AtName("fallback"),
					"Invalid fallback",
					"`fallback` is only supported for connections leaving a `condition` node.",
				)
			}
			continue
		}

		if hasOutlet || hasFallback {
			resp.Diagnostics.AddAttributeError(
				connectionPath,
				"Invalid connection",
				fmt.Sprintf("Connections leaving a %q node cannot set `source_outlet_identifier` or `fallback`.", sourceType),
			)
		}
	}

	validateNoCycles(resp, connections)
}

func validateEndpoint(resp *resource.ValidateConfigResponse, connectionPath path.Path, attribute string, value types.String, nodeTypes map[string]string) {
	if value.IsNull() || value.IsUnknown() {
		return
	}

	if _, known := nodeTypes[value.ValueString()]; !known {
		resp.Diagnostics.AddAttributeError(
			connectionPath.AtName(attribute),
			"Unknown node identifier",
			fmt.Sprintf("No node is declared with the identifier %q.", value.ValueString()),
		)
	}
}

// Mirrors the service's topological sort over the connection graph. Self-loops
// are skipped because they are reported with a more specific message.
func validateNoCycles(resp *resource.ValidateConfigResponse, connections []ConnectionModel) {
	adjacency := map[string][]string{}
	incoming := map[string]int{}

	for _, connection := range connections {
		if connection.SourceIdentifier.IsUnknown() || connection.TargetIdentifier.IsUnknown() {
			return
		}

		source := connection.SourceIdentifier.ValueString()
		target := connection.TargetIdentifier.ValueString()
		if source == target {
			continue
		}

		if _, seen := incoming[source]; !seen {
			incoming[source] = 0
		}
		incoming[target]++
		adjacency[source] = append(adjacency[source], target)
	}

	frontier := make([]string, 0, len(incoming))
	for node, degree := range incoming {
		if degree == 0 {
			frontier = append(frontier, node)
		}
	}

	processed := 0
	for len(frontier) > 0 {
		node := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		processed++

		for _, target := range adjacency[node] {
			incoming[target]--
			if incoming[target] == 0 {
				frontier = append(frontier, target)
			}
		}
	}

	if processed < len(incoming) {
		resp.Diagnostics.AddAttributeError(
			path.Root("connections"),
			"Cyclic connections",
			"The connections form a cycle. A workflow graph must be acyclic, so a node cannot be reachable from itself.",
		)
	}
}

// The service falls back to defaults for both or neither, so a lone value is rejected.
func requirePaired(resp *resource.ValidateConfigResponse, blockPath path.Path, blockName string, provider, model types.String) {
	providerSet := !provider.IsNull() && !provider.IsUnknown()
	modelSet := !model.IsNull() && !model.IsUnknown()
	if providerSet == modelSet {
		return
	}

	missing := "model"
	if modelSet {
		missing = "provider"
	}
	resp.Diagnostics.AddAttributeError(
		blockPath.AtName(missing),
		"Missing required attribute",
		fmt.Sprintf("`provider` and `model` must be set together on an `%s` node, or both omitted to use the defaults.", blockName),
	)
}

func requireSet(resp *resource.ValidateConfigResponse, attributePath path.Path, value types.String, detail string) {
	if value.IsNull() && !value.IsUnknown() {
		resp.Diagnostics.AddAttributeError(attributePath, "Missing required attribute", detail)
	}
}

// The API trims these values and rejects the result when it is empty, so a
// blank string is as invalid as an absent one.
func requireNonEmpty(resp *resource.ValidateConfigResponse, attributePath path.Path, value types.String, detail string) {
	if value.IsUnknown() {
		return
	}
	if value.IsNull() || strings.TrimSpace(value.ValueString()) == "" {
		resp.Diagnostics.AddAttributeError(attributePath, "Missing required attribute", detail)
	}
}

func requireUnique(resp *resource.ValidateConfigResponse, attributePath path.Path, value types.String, seen map[string]bool, noun string) {
	if value.IsUnknown() {
		return
	}

	identifier := value.ValueString()
	if seen[identifier] {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			fmt.Sprintf("Duplicate %s identifier", noun),
			fmt.Sprintf("The %s identifier %q is used more than once. Each %s of a node must have a unique identifier.", noun, identifier, noun),
		)
	}
	seen[identifier] = true
}

func rejectSet(resp *resource.ValidateConfigResponse, attributePath path.Path, value types.String, detail string) {
	if !value.IsNull() {
		resp.Diagnostics.AddAttributeError(attributePath, "Invalid attribute combination", detail)
	}
}
