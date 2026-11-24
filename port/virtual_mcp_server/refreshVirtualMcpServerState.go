package virtual_mcp_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshVirtualMcpServerState(_ context.Context, state *VirtualMcpServerModel, entity *cli.Entity) error {
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

	if description, ok := entity.Properties["description"].(string); ok {
		state.Description = types.StringValue(description)
	} else {
		state.Description = types.StringNull()
	}

	if instructions, ok := entity.Properties["instructions"].(string); ok {
		state.Instructions = types.StringValue(instructions)
	} else {
		state.Instructions = types.StringNull()
	}

	if patternsRaw, ok := entity.Properties["passthrough_tools"].([]any); ok {
		patterns := make([]types.String, 0, len(patternsRaw))
		for _, patternRaw := range patternsRaw {
			if pattern, ok := patternRaw.(string); ok {
				patterns = append(patterns, types.StringValue(pattern))
			}
		}
		state.PassthroughTools = patterns
	} else {
		state.PassthroughTools = nil
	}

	if relationsRaw, ok := entity.Relations[VirtualMcpToolsRelation].([]any); ok {
		toolIDs := make([]types.String, 0, len(relationsRaw))
		for _, toolRaw := range relationsRaw {
			if toolID, ok := toolRaw.(string); ok {
				toolIDs = append(toolIDs, types.StringValue(toolID))
			}
		}
		state.VirtualMcpTools = toolIDs
	} else if relationsRaw, ok := entity.Relations[VirtualMcpToolsRelation].([]string); ok {
		toolIDs := make([]types.String, len(relationsRaw))
		for i, toolID := range relationsRaw {
			toolIDs[i] = types.StringValue(toolID)
		}
		state.VirtualMcpTools = toolIDs
	} else {
		state.VirtualMcpTools = nil
	}

	if url, ok := entity.Properties["url"].(string); ok && url != "" {
		state.URL = types.StringValue(url)
	} else {
		state.URL = types.StringNull()
	}

	return nil
}