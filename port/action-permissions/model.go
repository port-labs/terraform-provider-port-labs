package action_permissions

import "github.com/hashicorp/terraform-plugin-framework/types"

type ExecuteModel struct {
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
	Teams       []types.String `tfsdk:"teams"`
	OwnedByTeam types.Bool     `tfsdk:"owned_by_team"`
	Policy      types.String   `tfsdk:"policy"`
}

type ApproveModel struct {
	Users  []types.String `tfsdk:"users"`
	Roles  []types.String `tfsdk:"roles"`
	Teams  []types.String `tfsdk:"teams"`
	Policy types.String   `tfsdk:"policy"`
}

type PermissionsModel struct {
	Execute *ExecuteModel `tfsdk:"execute"`
	Approve *ApproveModel `tfsdk:"approve"`
}

type ActionPermissionsModel struct {
	ID                  types.String      `tfsdk:"id"`
	ActionIdentifier    types.String      `tfsdk:"action_identifier"`
	BlueprintIdentifier types.String      `tfsdk:"blueprint_identifier"`
	Permissions         *PermissionsModel `tfsdk:"permissions"`
}
