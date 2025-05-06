package scorecard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Query struct {
	Combinator types.String   `tfsdk:"combinator"`
	Conditions []types.String `tfsdk:"conditions"`
}

type Rule struct {
	Identifier  types.String `tfsdk:"identifier"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Level       types.String `tfsdk:"level"`
	Query       *Query       `tfsdk:"query"`
}

type Level struct {
	Title types.String `tfsdk:"title"`
	Color types.String `tfsdk:"color"`
}

type ScorecardModel struct {
	ID         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Blueprint  types.String `tfsdk:"blueprint"`
	Title      types.String `tfsdk:"title"`
	Filter     *Query       `tfsdk:"filter"`
	Levels     []Level      `tfsdk:"levels"`
	Rules      []Rule       `tfsdk:"rules"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	UpdatedBy  types.String `tfsdk:"updated_by"`
}
