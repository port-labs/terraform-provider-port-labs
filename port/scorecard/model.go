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
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	Rules      []Rule       `tfsdk:"rules"`
}
