package user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserModel struct {
	ID                types.String   `tfsdk:"id"`
	Email             types.String   `tfsdk:"email"`
	Roles             []types.String `tfsdk:"roles"`
	Teams             []types.String `tfsdk:"teams"`
	InactivityTimeout types.Int64    `tfsdk:"inactivity_timeout"`
	Status            types.String   `tfsdk:"status"`
	ManagedByScim     types.Bool     `tfsdk:"managed_by_scim"`
}
