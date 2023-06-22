package entity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringPropModel struct {
	Value types.String `tfsdk:"value"`
}

type NumberPropModel struct {
	Value types.Float64 `tfsdk:"value"`
}

type EntityPropertiesModel struct {
	StringProp  map[string]string  `tfsdk:"string_prop"`
	NumberProp  map[string]float64 `tfsdk:"number_prop"`
	BooleanProp map[string]bool    `tfsdk:"boolean_prop"`
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
}
