package entity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ArrayPropsModel struct {
	StringItems  types.Map `tfsdk:"string_items"`
	NumberItems  types.Map `tfsdk:"number_items"`
	BooleanItems types.Map `tfsdk:"boolean_items"`
	ObjectItems  types.Map `tfsdk:"object_items"`
}

type EntityPropertiesModel struct {
	StringProps  map[string]types.String  `tfsdk:"string_props"`
	NumberProps  map[string]types.Float64 `tfsdk:"number_props"`
	BooleanProps map[string]types.Bool    `tfsdk:"boolean_props"`
	ObjectProps  map[string]types.String  `tfsdk:"object_props"`
	ArrayProps   *ArrayPropsModel         `tfsdk:"array_props"`
}

type RelationModel struct {
	SingleRelation map[string]string   `tfsdk:"single_relations"`
	ManyRelations  map[string][]string `tfsdk:"many_relations"`
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
	Relations  *RelationModel         `tfsdk:"relations"`
}
