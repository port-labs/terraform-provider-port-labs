package virtual_mcp_server

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func virtualMcpServerResourceToPortBody(state *VirtualMcpServerModel) (*cli.Entity, error) {
	e := &cli.Entity{
		Blueprint:  VirtualMcpServerBlueprintIdentifier,
		Title:      state.Title.ValueString(),
		Properties: map[string]any{},
		Relations:  map[string]any{},
	}

	if !state.Identifier.IsUnknown() && !state.Identifier.IsNull() {
		e.Identifier = state.Identifier.ValueString()
	}

	if !state.Icon.IsNull() {
		e.Icon = state.Icon.ValueString()
	}

	if state.Teams != nil {
		e.Team = make([]string, len(state.Teams))
		for i, team := range state.Teams {
			e.Team[i] = team.ValueString()
		}
	}

	if !state.Description.IsNull() {
		e.Properties["description"] = state.Description.ValueString()
	}

	if !state.Instructions.IsNull() {
		e.Properties["instructions"] = state.Instructions.ValueString()
	}

	if state.PassthroughTools != nil {
		patterns := make([]string, len(state.PassthroughTools))
		for i, pattern := range state.PassthroughTools {
			patterns[i] = pattern.ValueString()
		}
		e.Properties["passthrough_tools"] = patterns
	}

	if state.VirtualMcpTools != nil {
		toolIDs := make([]string, len(state.VirtualMcpTools))
		for i, toolID := range state.VirtualMcpTools {
			toolIDs[i] = toolID.ValueString()
		}
		e.Relations[VirtualMcpToolsRelation] = toolIDs
	}

	return e, nil
}