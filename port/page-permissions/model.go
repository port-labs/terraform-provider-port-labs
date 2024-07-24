package page_permissions

import "github.com/hashicorp/terraform-plugin-framework/types"

type ReadPagePermissionsModel struct {
	Users []types.String `tfsdk:"users"`
	Roles []types.String `tfsdk:"roles"`
	Teams []types.String `tfsdk:"teams"`
}

type PagePermissionsModel struct {
	ID             types.String              `tfsdk:"id"`
	PageIdentifier types.String              `tfsdk:"page_identifier"`
	Read           *ReadPagePermissionsModel `tfsdk:"read"`
}
