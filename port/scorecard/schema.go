package scorecard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			Required:            true,
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
		},
		"query": schema.SingleNestedAttribute{
			MarkdownDescription: "The query of the rule",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"combinator": schema.StringAttribute{
					MarkdownDescription: "The combinator of the query",
					Required:            true,
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
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the scorecard",
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
	}
}

func (r *ScorecardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Webhook resource",
		Attributes:          ScorecardSchema(),
	}
}
