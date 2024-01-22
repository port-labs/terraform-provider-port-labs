package page

import "github.com/hashicorp/terraform-plugin-framework/types"

type PageModel struct {
	ID                  types.String   `tfsdk:"id"`
	Identifier          types.String   `tfsdk:"identifier"`
	Type                types.String   `tfsdk:"type"`
	ShowInSidebar       types.Bool     `tfsdk:"show_in_sidebar"`
	Section             types.String   `tfsdk:"section"`
	Icon                types.String   `tfsdk:"icon"`
	Title               types.String   `tfsdk:"title"`
	Locked              types.Bool     `tfsdk:"locked"`
	Blueprint           types.String   `tfsdk:"blueprint"`
	RequiredQueryParams []types.String `tfsdk:"required_query_params"`
	Widgets             []types.String `tfsdk:"widgets"`
	CreatedAt           types.String   `tfsdk:"created_at"`
	CreatedBy           types.String   `tfsdk:"created_by"`
	UpdatedAt           types.String   `tfsdk:"updated_at"`
	UpdatedBy           types.String   `tfsdk:"updated_by"`
}
