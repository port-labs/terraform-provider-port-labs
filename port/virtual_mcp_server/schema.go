package virtual_mcp_server

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func VirtualMcpServerSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the virtual MCP server entity (used in the server URL path)",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display title of the virtual MCP server",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the virtual MCP server",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "What the virtual MCP server is for",
			Optional:            true,
		},
		"instructions": schema.StringAttribute{
			MarkdownDescription: "Instructions sent to agents when they connect to this server",
			Optional:            true,
		},
		"passthrough_tools": schema.ListAttribute{
			MarkdownDescription: "Regex patterns of Port or connector tools to expose unchanged (patterns are anchored to the full tool name)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"virtual_mcp_tools": schema.ListAttribute{
			MarkdownDescription: "Identifiers of `port_virtual_mcp_tool` entities exposed by this server",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"url": schema.StringAttribute{
			MarkdownDescription: "The server's MCP endpoint URL (calculated by Port from the entity identifier)",
			Computed:            true,
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: "The teams the virtual MCP server belongs to",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the virtual MCP server",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the virtual MCP server",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the virtual MCP server",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the virtual MCP server",
			Computed:            true,
		},
	}
}