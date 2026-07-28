package workflow

import (
	"context"
	"strconv"

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
	}

	if !state.Tags.IsNull() {
		tagsRaw, err := utils.TerraformListToGoArray(ctx, state.Tags, "string")
		if err != nil {
			return nil, err
		}
		tags := utils.InterfaceToStringArray(tagsRaw)
		// Use a pointer so an explicit empty list clears tags on update.
		w.Tags = &tags
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
			node.Type = n.EventTrigger.Type.ValueString()
			node.BlueprintIdentifier = n.EventTrigger.BlueprintIdentifier.ValueStringPointer()

		case n.CursorAgent != nil:
			node.Type = consts.CursorAgent
			node.ApiKey = n.CursorAgent.ApiKey.ValueStringPointer()
			if n.CursorAgent.Prompt != nil {
				node.Prompt = &cli.CursorAgentPrompt{
					Text: n.CursorAgent.Prompt.Text.ValueStringPointer(),
				}
			}
			if n.CursorAgent.Source != nil {
				node.Source = &cli.CursorAgentSource{
					PrUrl: n.CursorAgent.Source.PrUrl.ValueStringPointer(),
				}
			}

		case n.Delay != nil:
			node.Type = consts.Delay
			node.Seconds = secondsToPortBody(n.Delay.Seconds.ValueString())

		case n.IntegrationAction != nil:
			node.Type = consts.IntegrationAction
			node.InstallationId = n.IntegrationAction.InstallationId.ValueStringPointer()
			node.IntegrationActionType = n.IntegrationAction.IntegrationActionType.ValueStringPointer()

			execProps := &cli.WorkflowIntegrationActionExecutionProperties{
				Org:      n.IntegrationAction.Org.ValueStringPointer(),
				Repo:     n.IntegrationAction.Repo.ValueStringPointer(),
				Workflow: n.IntegrationAction.Workflow.ValueStringPointer(),
			}
			if !n.IntegrationAction.WorkflowInputs.IsNull() {
				wi, err := utils.TerraformStringToGoType[map[string]any](n.IntegrationAction.WorkflowInputs)
				if err != nil {
					return nil, err
				}
				execProps.WorkflowInputs = wi
			}
			if !n.IntegrationAction.ReportWorkflowStatus.IsNull() {
				execProps.ReportWorkflowStatus = utils.TerraformStringToBooleanOrString(n.IntegrationAction.ReportWorkflowStatus)
			}
			node.IntegrationActionExecutionProperties = execProps
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

// secondsToPortBody sends a plain integer as a number and a dynamic expression
// as a string.
func secondsToPortBody(seconds string) any {
	if seconds == "" {
		return nil
	}
	if v, err := strconv.ParseInt(seconds, 10, 64); err == nil {
		return v
	}
	return seconds
}
