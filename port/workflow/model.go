package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkflowModel struct {
	ID                    types.String `tfsdk:"id"`
	Identifier            types.String `tfsdk:"identifier"`
	Title                 types.String `tfsdk:"title"`
	Icon                  types.String `tfsdk:"icon"`
	Description           types.String `tfsdk:"description"`
	Tags                  types.List   `tfsdk:"tags"`
	AllowAnyoneToViewRuns types.Bool   `tfsdk:"allow_anyone_to_view_runs"`
	Nodes                 types.String `tfsdk:"nodes"`
	Connections           types.String `tfsdk:"connections"`
	CreatedAt             types.String `tfsdk:"created_at"`
	CreatedBy             types.String `tfsdk:"created_by"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
	UpdatedBy             types.String `tfsdk:"updated_by"`
}
