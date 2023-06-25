package entity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ArrayPropModel struct {
	StringItems  types.Map `tfsdk:"string_items"`
	NumberItems  types.Map `tfsdk:"number_items"`
	BooleanItems types.Map `tfsdk:"boolean_items"`
}

type EntityPropertiesModel struct {
	StringProp  map[string]string  `tfsdk:"string_prop"`
	NumberProp  map[string]float64 `tfsdk:"number_prop"`
	BooleanProp map[string]bool    `tfsdk:"boolean_prop"`
	ObjectProp  map[string]string  `tfsdk:"object_prop"`
	ArrayProp   *ArrayPropModel    `tfsdk:"array_prop"`
}

type RelationModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	Identifiers types.List   `tfsdk:"identifiers"`
}

type EntityModel struct {
	ID         types.String           `tfsdk:"id"`
	Identifier types.String           `tfsdk:"identifier"`
	Blueprint  types.String           `tfsdk:"blueprint"`
	Title      types.String           `tfsdk:"title"`
	Icon       types.String           `tfsdk:"icon"`
	RunID      types.String           `tfsdk:"run_id"`
	CreatedAt  types.String           `tfsdk:"created_at"`
	CreatedBy  types.String           `tfsdk:"created_by"`
	UpdatedAt  types.String           `tfsdk:"updated_at"`
	UpdatedBy  types.String           `tfsdk:"updated_by"`
	Properties *EntityPropertiesModel `tfsdk:"properties"`
	Teams      []types.String         `tfsdk:"teams"`
	Relations  types.Map              `tfsdk:"relations"`
}
