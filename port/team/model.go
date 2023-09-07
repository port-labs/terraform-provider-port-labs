package team

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Description  types.String   `tfsdk:"description"`
	Users        []types.String `tfsdk:"users"`
	CreatedAt    types.String   `tfsdk:"created_at"`
	UpdatedAt    types.String   `tfsdk:"updated_at"`
	ProviderName types.String   `tfsdk:"provider_name"`
}
