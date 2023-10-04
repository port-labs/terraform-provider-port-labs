package actionpermissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExecuteModel struct {
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
	Teams       []types.String `tfsdk:"teams"`
	OwnedByTeam types.Bool     `tfsdk:"ownedByTeam"`
}

type ApproveModel struct {
	Users []types.String `tfsdk:"users"`
	Roles []types.String `tfsdk:"roles"`
	Teams []types.String `tfsdk:"teams"`
}

type ActionPermissionsModel struct {
	Execute *ExecuteModel `tfsdk:"execute"`
	Approve *ApproveModel `tfsdk:"approve"`
}
