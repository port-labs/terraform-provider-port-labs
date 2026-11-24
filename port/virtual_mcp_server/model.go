package virtual_mcp_server

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VirtualMcpServerModel struct {
	ID               types.String   `tfsdk:"id"`
	Identifier       types.String   `tfsdk:"identifier"`
	Title            types.String   `tfsdk:"title"`
	Icon             types.String   `tfsdk:"icon"`
	Description      types.String   `tfsdk:"description"`
	Instructions     types.String   `tfsdk:"instructions"`
	PassthroughTools []types.String `tfsdk:"passthrough_tools"`
	VirtualMcpTools  []types.String `tfsdk:"virtual_mcp_tools"`
	URL              types.String   `tfsdk:"url"`
	Teams            []types.String `tfsdk:"teams"`
	CreatedAt        types.String   `tfsdk:"created_at"`
	CreatedBy        types.String   `tfsdk:"created_by"`
	UpdatedAt        types.String   `tfsdk:"updated_at"`
	UpdatedBy        types.String   `tfsdk:"updated_by"`
}