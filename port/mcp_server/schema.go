package mcp_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func McpServerSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the MCP server entity",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display name of the MCP server",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the MCP server",
			Optional:            true,
		},
		"url": schema.StringAttribute{
			MarkdownDescription: "The remote MCP server endpoint URL",
			Required:            true,
		},
		"auth_method": schema.StringAttribute{
			MarkdownDescription: "The authentication method used to connect to the MCP server (for example, `oauth`, `headers`, or `none`)",
			Optional:            true,
		},
		"oauth_config": schema.SingleNestedAttribute{
			MarkdownDescription: "OAuth configuration for MCP servers that use OAuth authentication",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"client_id": schema.StringAttribute{
					MarkdownDescription: "OAuth client ID for manual OAuth registration",
					Optional:            true,
				},
				"client_secret": schema.StringAttribute{
					MarkdownDescription: "OAuth client secret for manual OAuth registration. Port may not return this value on subsequent reads, so Terraform retains the value from create or update in state.",
					Optional:            true,
					Sensitive:           true,
				},
				"scope": schema.StringAttribute{
					MarkdownDescription: "OAuth scope to request during authorization",
					Optional:            true,
				},
				"grant_type": schema.StringAttribute{
					MarkdownDescription: "OAuth grant type used when connecting to the MCP server. Use `authorization_code` for interactive user consent flows, or `client_credentials` for machine-to-machine authentication without a browser redirect. Defaults to `authorization_code` when omitted.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(GrantTypeAuthorizationCode, GrantTypeClientCredentials),
					},
				},
			},
		},
		"headers": schema.MapAttribute{
			MarkdownDescription: "HTTP headers sent when connecting to the MCP server (for API key or token authentication)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: "The teams the MCP server belongs to",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the MCP server",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the MCP server",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the MCP server",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the MCP server",
			Computed:            true,
		},
	}
}

func (r *McpServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          McpServerSchema(),
	}
}

var ResourceMarkdownDescription = `
# MCP Server resource

Manages an external MCP (Model Context Protocol) server connected to Port under the ` + "`_mcp_server`" + ` system blueprint.

Use ` + "`oauth_config.grant_type`" + ` to choose between interactive OAuth (` + "`authorization_code`" + `) and machine-to-machine OAuth (` + "`client_credentials`" + `).

## Example Usage

` + "```hcl" + `
resource "port_mcp_server" "example" {
  title       = "Example MCP Server"
  url         = "https://mcp.example.com/v1"
  auth_method = "oauth"
  oauth_config = {
    client_id     = "my-client-id"
    client_secret = var.mcp_client_secret
    scope         = "read write"
    grant_type    = "client_credentials"
  }
}
` + "\n```"
