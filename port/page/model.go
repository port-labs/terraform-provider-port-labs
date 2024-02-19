package page

import "github.com/hashicorp/terraform-plugin-framework/types"

type PageModel struct {
	ID         types.String   `tfsdk:"id"`
	Identifier types.String   `tfsdk:"identifier"`
	Title      types.String   `tfsdk:"title"`
	Type       types.String   `tfsdk:"type"`
	Parent     types.String   `tfsdk:"parent"`
	After      types.String   `tfsdk:"after"`
	Icon       types.String   `tfsdk:"icon"`
	Locked     types.Bool     `tfsdk:"locked"`
	Blueprint  types.String   `tfsdk:"blueprint"`
	Widgets    []types.String `tfsdk:"widgets"`
	CreatedAt  types.String   `tfsdk:"created_at"`
	CreatedBy  types.String   `tfsdk:"created_by"`
	UpdatedAt  types.String   `tfsdk:"updated_at"`
	UpdatedBy  types.String   `tfsdk:"updated_by"`
}
