package scorecard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Condition struct {
	Operator types.String `tfsdk:"operator"`
	Property types.String `tfsdk:"property"`
	Value    types.String `tfsdk:"value"`
}

type Query struct {
	Combinator types.String `tfsdk:"combinator"`
	Conditions []Condition  `tfsdk:"conditions"`
}

type Rule struct {
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	Level      types.String `tfsdk:"level"`
	Query      *Query       `tfsdk:"query"`
}

type ScorecardModel struct {
	ID         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Blueprint  types.String `tfsdk:"blueprint"`
	Title      types.String `tfsdk:"title"`
	Rules      []Rule       `tfsdk:"rules"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	UpdatedBy  types.String `tfsdk:"updated_by"`
}
