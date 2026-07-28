package workflow

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *WorkflowResource) refreshWorkflowState(ctx context.Context, state *WorkflowModel, w *cli.Workflow) error {
	// Preserve prior nodes/connections so we can keep ordering stable and
	// re-use secret values that the API does not return.
	priorNodes := make(map[string]WorkflowNodeModel, len(state.Nodes))
	for _, n := range state.Nodes {
		priorNodes[n.Identifier.ValueString()] = n
	}

	state.ID = types.StringValue(w.Identifier)
	state.Identifier = types.StringValue(w.Identifier)
	state.Title = flex.GoStringToFramework(w.Title)
	state.Icon = flex.GoStringToFramework(w.Icon)
	state.Description = flex.GoStringToFramework(w.Description)

	if w.Tags == nil {
		state.Tags = types.ListNull(types.StringType)
	} else {
		state.Tags = flex.GoArrayStringToTerraformList(ctx, *w.Tags)
	}

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

	switch apiNode.Type {
	case consts.CursorAgent:
		apiKey := flex.GoStringToFramework(apiNode.ApiKey)
		// The API may not return the api_key (secret); keep the prior value.
		if apiKey.IsNull() && prior.CursorAgent != nil {
			apiKey = prior.CursorAgent.ApiKey
		}
		cursorAgent := &CursorAgentModel{ApiKey: apiKey}
		if apiNode.Prompt != nil {
			cursorAgent.Prompt = &CursorPromptModel{Text: flex.GoStringToFramework(apiNode.Prompt.Text)}
		}
		if apiNode.Source != nil {
			cursorAgent.Source = &CursorSourceModel{PrUrl: flex.GoStringToFramework(apiNode.Source.PrUrl)}
		}
		node.CursorAgent = cursorAgent

	case consts.Delay:
		node.Delay = &DelayModel{Seconds: secondsToModel(apiNode.Seconds)}

	case consts.IntegrationAction:
		integrationAction := &IntegrationActionModel{
			InstallationId:        flex.GoStringToFramework(apiNode.InstallationId),
			IntegrationActionType: flex.GoStringToFramework(apiNode.IntegrationActionType),
			WorkflowInputs:        types.StringNull(),
			ReportWorkflowStatus:  types.StringNull(),
		}
		if props := apiNode.IntegrationActionExecutionProperties; props != nil {
			integrationAction.Org = flex.GoStringToFramework(props.Org)
			integrationAction.Repo = flex.GoStringToFramework(props.Repo)
			integrationAction.Workflow = flex.GoStringToFramework(props.Workflow)
			workflowInputs, err := utils.GoObjectToTerraformString(props.WorkflowInputs, r.portClient.JSONEscapeHTML)
			if err != nil {
				return node, err
			}
			integrationAction.WorkflowInputs = workflowInputs
			reportWorkflowStatus, err := utils.GoObjectToTerraformString(props.ReportWorkflowStatus, r.portClient.JSONEscapeHTML)
			if err != nil {
				return node, err
			}
			integrationAction.ReportWorkflowStatus = reportWorkflowStatus
		}
		node.IntegrationAction = integrationAction

	default:
		// Any other type is treated as an event trigger, whose discriminator is
		// the event type itself.
		node.EventTrigger = &EventTriggerModel{
			Type:                types.StringValue(apiNode.Type),
			BlueprintIdentifier: flex.GoStringToFramework(apiNode.BlueprintIdentifier),
		}
	}

	return node, nil
}

// reorderNodes returns the API nodes ordered to match the prior state order
// (matched by identifier), with any new nodes appended at the end. This keeps
// the resource stable against non-deterministic API ordering.
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

// secondsToModel converts the API seconds value (a number or a dynamic
// expression string) to a string.
func secondsToModel(seconds any) types.String {
	switch v := seconds.(type) {
	case nil:
		return types.StringNull()
	case string:
		return types.StringValue(v)
	case float64:
		return types.StringValue(strconv.FormatInt(int64(v), 10))
	case int64:
		return types.StringValue(strconv.FormatInt(v, 10))
	case int:
		return types.StringValue(strconv.Itoa(v))
	default:
		return types.StringNull()
	}
}
