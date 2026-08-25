package scorecard_group

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

type MemberSpecModel struct {
	Filter *scorecard.Query `tfsdk:"filter"`
	Rules  []scorecard.Rule `tfsdk:"rules"`
}

type ScorecardGroupModel struct {
	ID         types.String              `tfsdk:"id"`
	Identifier types.String              `tfsdk:"identifier"`
	Title      types.String              `tfsdk:"title"`
	Levels     []scorecard.Level         `tfsdk:"levels"`
	Blueprints []types.String            `tfsdk:"blueprints"`
	Rules      []scorecard.Rule          `tfsdk:"rules"`
	Filter     *scorecard.Query          `tfsdk:"filter"`
	Scorecards map[string]MemberSpecModel `tfsdk:"scorecards"`
	Properties types.Dynamic              `tfsdk:"properties"`
	CreatedAt  types.String              `tfsdk:"created_at"`
	CreatedBy  types.String              `tfsdk:"created_by"`
	UpdatedAt  types.String              `tfsdk:"updated_at"`
	UpdatedBy  types.String              `tfsdk:"updated_by"`
}
