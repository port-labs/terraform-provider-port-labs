package folder

import "github.com/hashicorp/terraform-plugin-framework/types"

type FolderModel struct {
	ID         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	Sidebar    types.String `tfsdk:"sidebar"`
	After      types.String `tfsdk:"after"`
	Parent     types.String `tfsdk:"parent"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	UpdatedBy  types.String `tfsdk:"updated_by"`
}
