package virtual_mcp_tool

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshVirtualMcpToolState(_ context.Context, state *VirtualMcpToolModel, entity *cli.Entity) error {
	state.ID = types.StringValue(entity.Identifier)
	state.Identifier = types.StringValue(entity.Identifier)
	state.Title = types.StringValue(entity.Title)

	if entity.Icon != "" {
		state.Icon = types.StringValue(entity.Icon)
	} else {
		state.Icon = types.StringNull()
	}

	if entity.CreatedAt != nil {
		state.CreatedAt = types.StringValue(entity.CreatedAt.String())
	}

	if entity.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(entity.UpdatedAt.String())
	}

	state.CreatedBy = types.StringValue(entity.CreatedBy)
	state.UpdatedBy = types.StringValue(entity.UpdatedBy)

	if entity.Team != nil {
		teams := make([]types.String, len(entity.Team))
		for i, team := range entity.Team {
			teams[i] = types.StringValue(team)
		}
		state.Teams = teams
	} else {
		state.Teams = nil
	}

	if name, ok := entity.Properties["name"].(string); ok {
		state.Name = types.StringValue(name)
	} else {
		state.Name = types.StringNull()
	}

	if tool, ok := entity.Properties["tool"].(string); ok {
		state.Tool = types.StringValue(tool)
	} else {
		state.Tool = types.StringNull()
	}

	if description, ok := entity.Properties["description"].(string); ok {
		state.Description = types.StringValue(description)
	} else {
		state.Description = types.StringNull()
	}

	if exposeArguments, ok := entity.Properties["expose_arguments"].(string); ok && exposeArguments != "" {
		state.ExposeArguments = types.StringValue(exposeArguments)
	} else {
		state.ExposeArguments = types.StringValue(ExposeArgumentsAll)
	}

	if argumentsRaw, ok := entity.Properties["arguments"].(map[string]any); ok && len(argumentsRaw) > 0 {
		encoded, err := json.Marshal(argumentsRaw)
		if err != nil {
			return err
		}
		state.Arguments = types.StringValue(string(encoded))
	} else {
		state.Arguments = types.StringNull()
	}

	return nil
}