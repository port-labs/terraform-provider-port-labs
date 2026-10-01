package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	IdleTimeMS types.Int64  `tfsdk:"idle_time_ms"`
}
