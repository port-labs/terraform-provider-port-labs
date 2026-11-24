package virtual_mcp_tool

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func VirtualMcpToolSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the virtual tool entity",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display title of the virtual tool",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the virtual tool",
			Optional:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The tool name exposed to MCP clients (up to 48 characters)",
			Required:            true,
		},
		"tool": schema.StringAttribute{
			MarkdownDescription: "The underlying Port or MCP connector tool name (for example `list_entities` or `github_list_commits`)",
			Required:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description shown to agents when the tool is listed",
			Optional:            true,
		},
		"expose_arguments": schema.StringAttribute{
			MarkdownDescription: "Whether agents see all underlying tool arguments (`all`) or only arguments listed in `arguments` (`listed`)",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(ExposeArgumentsAll),
			Validators: []validator.String{
				stringvalidator.OneOf(ExposeArgumentsAll, ExposeArgumentsListed),
			},
		},
		"arguments": schema.StringAttribute{
			MarkdownDescription: "Per-argument rules for the underlying tool, as a JSON object keyed by dotted argument paths (for example `{\"blueprintIdentifier\":{\"value\":\"service\"}}`)",
			Optional:            true,
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: "The teams the virtual tool belongs to",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the virtual tool",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the virtual tool",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the virtual tool",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the virtual tool",
			Computed:            true,
		},
	}
}