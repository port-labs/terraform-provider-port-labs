package folder

import "github.com/hashicorp/terraform-plugin-framework/types"

type FolderModel struct {
	ID         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	After      types.String `tfsdk:"after"`
	Parent     types.String `tfsdk:"parent"`
}
