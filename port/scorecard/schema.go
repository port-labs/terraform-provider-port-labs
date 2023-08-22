package scorecard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ConditionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operator": schema.StringAttribute{
			MarkdownDescription: "The operator of the condition",
			Required:            true,
		},
		"property": schema.StringAttribute{
			MarkdownDescription: "The property of the condition",
			Required:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value of the condition",
			Optional:            true,
		},
	}
}

func RuleSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the rule",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the rule",
			Required:            true,
		},
		"level": schema.StringAttribute{
			MarkdownDescription: "The level of the rule",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("Bronze", "Silver", "Gold"),
			},
		},
		"query": schema.SingleNestedAttribute{
			MarkdownDescription: "The query of the rule",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"combinator": schema.StringAttribute{
					MarkdownDescription: "The combinator of the query",
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
					Required: true,
				},
				"conditions": schema.ListNestedAttribute{
					MarkdownDescription: "The conditions of the query",
					Required:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: ConditionSchema(),
					},
				},
			},
		},
	}
}
func ScorecardSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the scorecard",
			Required:            true,
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint of the scorecard",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the scorecard",
			Required:            true,
		},
		"rules": schema.ListNestedAttribute{
			MarkdownDescription: "The rules of the scorecard",
			Required:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: RuleSchema(),
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the scorecard",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the scorecard",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the scorecard",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the scorecard",
			Computed:            true,
		},
	}
}

func (r *ScorecardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Webhook resource",
		Attributes:          ScorecardSchema(),
	}
}
