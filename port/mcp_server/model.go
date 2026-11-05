package mcp_server

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OAuthConfigModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Scope        types.String `tfsdk:"scope"`
	GrantType    types.String `tfsdk:"grant_type"`
}

type McpServerModel struct {
	ID          types.String        `tfsdk:"id"`
	Identifier  types.String        `tfsdk:"identifier"`
	Title       types.String        `tfsdk:"title"`
	Icon        types.String        `tfsdk:"icon"`
	URL         types.String        `tfsdk:"url"`
	AuthMethod  types.String        `tfsdk:"auth_method"`
	OAuthConfig *OAuthConfigModel   `tfsdk:"oauth_config"`
	Headers     map[string]types.String `tfsdk:"headers"`
	Teams       []types.String      `tfsdk:"teams"`
	CreatedAt   types.String        `tfsdk:"created_at"`
	CreatedBy   types.String        `tfsdk:"created_by"`
	UpdatedAt   types.String        `tfsdk:"updated_at"`
	UpdatedBy   types.String        `tfsdk:"updated_by"`
}
