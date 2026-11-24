package virtual_mcp_tool

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VirtualMcpToolModel struct {
	ID              types.String   `tfsdk:"id"`
	Identifier      types.String   `tfsdk:"identifier"`
	Title           types.String   `tfsdk:"title"`
	Icon            types.String   `tfsdk:"icon"`
	Name            types.String   `tfsdk:"name"`
	Tool            types.String   `tfsdk:"tool"`
	Description     types.String   `tfsdk:"description"`
	ExposeArguments types.String   `tfsdk:"expose_arguments"`
	Arguments       types.String   `tfsdk:"arguments"`
	Teams           []types.String `tfsdk:"teams"`
	CreatedAt       types.String   `tfsdk:"created_at"`
	CreatedBy       types.String   `tfsdk:"created_by"`
	UpdatedAt       types.String   `tfsdk:"updated_at"`
	UpdatedBy       types.String   `tfsdk:"updated_by"`
}