package virtual_mcp_tool

import (
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func virtualMcpToolResourceToPortBody(state *VirtualMcpToolModel) (*cli.Entity, error) {
	e := &cli.Entity{
		Blueprint:  VirtualMcpToolBlueprintIdentifier,
		Title:      state.Title.ValueString(),
		Properties: map[string]any{},
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

	e.Properties["name"] = state.Name.ValueString()
	e.Properties["tool"] = state.Tool.ValueString()

	if !state.Description.IsNull() {
		e.Properties["description"] = state.Description.ValueString()
	}

	exposeArguments := ExposeArgumentsAll
	if !state.ExposeArguments.IsNull() {
		exposeArguments = state.ExposeArguments.ValueString()
	}
	e.Properties["expose_arguments"] = exposeArguments

	if !state.Arguments.IsNull() && state.Arguments.ValueString() != "" {
		var arguments map[string]any
		if err := json.Unmarshal([]byte(state.Arguments.ValueString()), &arguments); err != nil {
			return nil, err
		}
		e.Properties["arguments"] = arguments
	}

	return e, nil
}