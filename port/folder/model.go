package folder

import "github.com/hashicorp/terraform-plugin-framework/types"

type FolderModel struct {
	FolderIdentifier  types.String `tfsdk:"folder_identifier"`
	SidebarIdentifier types.String `tfsdk:"sidebar_identifier"`
	Title             types.String `tfsdk:"title"`
	After             types.String `tfsdk:"after"`
	Parent            types.String `tfsdk:"parent"`
}
