package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	PortalIcon types.String `tfsdk:"portal_icon"`
}
