package workflow

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func workflowStateToPortBody(ctx context.Context, state *WorkflowModel) (*cli.Workflow, error) {
	w := &cli.Workflow{
		Identifier:  state.Identifier.ValueString(),
		Title:       state.Title.ValueStringPointer(),
		Icon:        state.Icon.ValueStringPointer(),
		Description: state.Description.ValueStringPointer(),
		Category:    state.Category.ValueStringPointer(),
	}

	nodes, err := nodesToPortBody(state.Nodes)
	if err != nil {
		return nil, err
	}
	w.Nodes = nodes

	w.Connections = connectionsToPortBody(state.Connections)

	return w, nil
}

func nodesToPortBody(nodes []WorkflowNodeModel) ([]cli.WorkflowNode, error) {
	result := make([]cli.WorkflowNode, 0, len(nodes))
	for _, n := range nodes {
		node := cli.WorkflowNode{
			Identifier: n.Identifier.ValueString(),
			Title:      n.Title.ValueStringPointer(),
		}

		switch {
		case n.EventTrigger != nil:
			node.Config.Type = consts.EventTrigger
			node.Config.Event = &cli.WorkflowTriggerEvent{
				Type:                n.EventTrigger.Type.ValueString(),
				BlueprintIdentifier: n.EventTrigger.BlueprintIdentifier.ValueString(),
				PropertyIdentifier:  n.EventTrigger.PropertyIdentifier.ValueStringPointer(),
			}

		case n.CursorAgent != nil:
			node.Config.Type = consts.CursorAgent
			node.Config.ApiKey = n.CursorAgent.ApiKey.ValueStringPointer()
			if n.CursorAgent.Prompt != nil {
				node.Config.Prompt = &cli.CursorAgentPrompt{
					Text: n.CursorAgent.Prompt.Text.ValueStringPointer(),
				}
			}
			if n.CursorAgent.Source != nil {
				node.Config.Source = &cli.CursorAgentSource{
					Repository: n.CursorAgent.Source.Repository.ValueStringPointer(),
					Ref:        n.CursorAgent.Source.Ref.ValueStringPointer(),
					PrUrl:      n.CursorAgent.Source.PrUrl.ValueStringPointer(),
				}
			}

		case n.IntegrationAction != nil:
			node.Config.Type = consts.IntegrationAction
			node.Config.InstallationId = n.IntegrationAction.InstallationId.ValueStringPointer()
			node.Config.IntegrationProvider = n.IntegrationAction.IntegrationProvider.ValueStringPointer()
			node.Config.IntegrationInvocationType = n.IntegrationAction.IntegrationInvocationType.ValueStringPointer()
			node.Config.OnFailure = n.IntegrationAction.OnFailure.ValueStringPointer()
			if !n.IntegrationAction.DeferIntegrationInstallation.IsNull() {
				node.Config.DeferIntegrationInstallation = n.IntegrationAction.DeferIntegrationInstallation.ValueBoolPointer()
			}
			if !n.IntegrationAction.ExecutionProperties.IsNull() {
				props, err := utils.TerraformStringToGoType[any](n.IntegrationAction.ExecutionProperties)
				if err != nil {
					return nil, err
				}
				node.Config.IntegrationActionExecutionProperties = props
			}
		}

		result = append(result, node)
	}
	return result, nil
}

func connectionsToPortBody(connections []ConnectionModel) []cli.WorkflowConnection {
	result := make([]cli.WorkflowConnection, 0, len(connections))
	for _, c := range connections {
		result = append(result, cli.WorkflowConnection{
			SourceIdentifier: c.SourceIdentifier.ValueString(),
			TargetIdentifier: c.TargetIdentifier.ValueString(),
		})
	}
	return result
}
