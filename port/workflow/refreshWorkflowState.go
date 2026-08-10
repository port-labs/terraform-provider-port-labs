package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *WorkflowResource) refreshWorkflowState(ctx context.Context, state *WorkflowModel, w *cli.Workflow) error {
	// Keep prior node order stable and reuse secret values the API does not return.
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

	orderedNodes := reorderNodes(state.Nodes, w.Nodes)
	nodes := make([]WorkflowNodeModel, 0, len(orderedNodes))
	for _, apiNode := range orderedNodes {
		node, err := r.nodeToModel(apiNode, priorNodes[apiNode.Identifier])
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}
	state.Nodes = nodes

	state.Connections = reorderConnections(state.Connections, w.Connections)

	return nil
}

func (r *WorkflowResource) nodeToModel(apiNode cli.WorkflowNode, prior WorkflowNodeModel) (WorkflowNodeModel, error) {
	node := WorkflowNodeModel{
		Identifier: types.StringValue(apiNode.Identifier),
		Title:      flex.GoStringToFramework(apiNode.Title),
	}

	config := apiNode.Config
	switch config.Type {
	case consts.CursorAgent:
		apiKey := flex.GoStringToFramework(config.ApiKey)
		// api_key is a secret the API does not echo back; keep the prior value.
		if apiKey.IsNull() && prior.CursorAgent != nil {
			apiKey = prior.CursorAgent.ApiKey
		}
		cursorAgent := &CursorAgentModel{ApiKey: apiKey}
		if config.Prompt != nil {
			cursorAgent.Prompt = &CursorPromptModel{Text: flex.GoStringToFramework(config.Prompt.Text)}
		}
		if config.Source != nil {
			cursorAgent.Source = &CursorSourceModel{
				Repository: flex.GoStringToFramework(config.Source.Repository),
				Ref:        flex.GoStringToFramework(config.Source.Ref),
				PrUrl:      flex.GoStringToFramework(config.Source.PrUrl),
			}
		}
		node.CursorAgent = cursorAgent

	case consts.IntegrationAction:
		integrationAction := &IntegrationActionModel{
			InstallationId:            flex.GoStringToFramework(config.InstallationId),
			IntegrationProvider:       flex.GoStringToFramework(config.IntegrationProvider),
			IntegrationInvocationType: flex.GoStringToFramework(config.IntegrationInvocationType),
			OnFailure:                 flex.GoStringToFramework(config.OnFailure),
			ExecutionProperties:       types.StringNull(),
			DeferIntegrationInstallation: func() types.Bool {
				if config.DeferIntegrationInstallation == nil {
					return types.BoolNull()
				}
				return types.BoolValue(*config.DeferIntegrationInstallation)
			}(),
		}
		if config.IntegrationActionExecutionProperties != nil {
			executionProperties, err := utils.GoObjectToTerraformString(config.IntegrationActionExecutionProperties, r.portClient.JSONEscapeHTML)
			if err != nil {
				return node, err
			}
			integrationAction.ExecutionProperties = executionProperties
		}
		node.IntegrationAction = integrationAction

	case consts.EventTrigger:
		eventTrigger := &EventTriggerModel{}
		if config.Event != nil {
			eventTrigger.Type = types.StringValue(config.Event.Type)
			eventTrigger.BlueprintIdentifier = types.StringValue(config.Event.BlueprintIdentifier)
			eventTrigger.PropertyIdentifier = flex.GoStringToFramework(config.Event.PropertyIdentifier)
		}
		node.EventTrigger = eventTrigger
	}

	return node, nil
}

// reorderNodes orders API nodes by the prior state order, appending new nodes last.
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
	key := func(source, target string) string { return source + "\x00" + target }

	byKey := make(map[string]cli.WorkflowConnection, len(apiConnections))
	for _, c := range apiConnections {
		byKey[key(c.SourceIdentifier, c.TargetIdentifier)] = c
	}

	ordered := make([]ConnectionModel, 0, len(apiConnections))
	seen := make(map[string]bool, len(apiConnections))
	appendConnection := func(c cli.WorkflowConnection) {
		ordered = append(ordered, ConnectionModel{
			SourceIdentifier: types.StringValue(c.SourceIdentifier),
			TargetIdentifier: types.StringValue(c.TargetIdentifier),
		})
	}

	for _, p := range prior {
		k := key(p.SourceIdentifier.ValueString(), p.TargetIdentifier.ValueString())
		if c, ok := byKey[k]; ok && !seen[k] {
			appendConnection(c)
			seen[k] = true
		}
	}
	for _, c := range apiConnections {
		k := key(c.SourceIdentifier, c.TargetIdentifier)
		if !seen[k] {
			appendConnection(c)
			seen[k] = true
		}
	}
	return ordered
}
