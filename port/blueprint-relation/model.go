package blueprint_relation

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BlueprintRelationModel struct {
	ID          types.String `tfsdk:"id"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Identifier  types.String `tfsdk:"identifier"`
	Target      types.String `tfsdk:"target"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Required    types.Bool   `tfsdk:"required"`
	Many        types.Bool   `tfsdk:"many"`
}
