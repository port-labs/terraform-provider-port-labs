package scorecard_group

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceQuerySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"combinator": schema.StringAttribute{
			MarkdownDescription: "The combinator of the query.",
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("and", "or"),
			},
		},
		"conditions": schema.ListAttribute{
			MarkdownDescription: "The conditions of the query. Each condition object should be encoded to a string.",
			Computed:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
	}
}

func dataSourceMemberSpecSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"filter": schema.SingleNestedAttribute{
			MarkdownDescription: "An optional set of conditions to filter entities evaluated by this member scorecard.",
			Computed:            true,
			Attributes:          dataSourceQuerySchema(),
		},
		"rules": schema.ListNestedAttribute{
			MarkdownDescription: "The rules that define this member scorecard.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: dataSourceRuleSchema(),
			},
		},
	}
}

func dataSourceRuleSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the rule",
			Computed:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the rule",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the rule",
			Computed:            true,
		},
		"level": schema.StringAttribute{
			MarkdownDescription: "The level of the rule",
			Computed:            true,
		},
		"query": schema.SingleNestedAttribute{
			MarkdownDescription: "The query of the rule",
			Computed:            true,
			Attributes:          dataSourceQuerySchema(),
		},
	}
}

func dataSourceLevelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"color": schema.StringAttribute{
			MarkdownDescription: "The color of the level",
			Computed:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the level",
			Computed:            true,
		},
	}
}

func scorecardGroupDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the scorecard group to read.",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the scorecard group (applied to all member scorecards).",
			Computed:            true,
		},
		"levels": schema.ListNestedAttribute{
			MarkdownDescription: "The available levels of the scorecard group, shared by all members.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: dataSourceLevelSchema(),
			},
		},
		"properties": schema.StringAttribute{
			MarkdownDescription: "Additional `_scorecard` blueprint properties applied to every member scorecard in the group, as a JSON encoded string.",
			Computed:            true,
		},
		"blueprints": schema.SetAttribute{
			MarkdownDescription: "Blueprint identifiers in shared-rules mode.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"rules": schema.ListNestedAttribute{
			MarkdownDescription: "The rules applied to every blueprint in shared-rules mode.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: dataSourceRuleSchema(),
			},
		},
		"filters": schema.MapNestedAttribute{
			MarkdownDescription: "Optional filters per blueprint in shared-rules mode, keyed by blueprint identifier.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: dataSourceQuerySchema(),
			},
		},
		"scorecards": schema.MapNestedAttribute{
			MarkdownDescription: "Map of blueprint identifier to member scorecard filter/rules in per-blueprint mode.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: dataSourceMemberSpecSchema(),
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the scorecard group.",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the scorecard group.",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the scorecard group.",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the scorecard group.",
			Computed:            true,
		},
	}
}
