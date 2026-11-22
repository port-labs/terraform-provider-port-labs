package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	HiddenGlobalSearchTabs types.Set    `tfsdk:"hidden_global_search_tabs"`
}
