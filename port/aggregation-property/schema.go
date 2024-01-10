package aggregation_property

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AggregationPropertySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"aggregation_identifier": schema.StringAttribute{
			Description: "The identifier of the aggregation property in the blueprint",
			Required:    true,
		},
		"blueprint_identifier": schema.StringAttribute{
			Description: "The identifier of the blueprint the aggregation property will be added to",
			Required:    true,
		},
		"target_blueprint_identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint to perform the aggregation on",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the aggregation property",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the aggregation property",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the aggregation property",
			Optional:            true,
		},
		"method": schema.SingleNestedAttribute{
			MarkdownDescription: "The aggregation method to perform on the target blueprint, one of count_entities, average_entities, average_by_property, aggregate_by_property",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"count_entities": schema.BoolAttribute{
					MarkdownDescription: "Function to count the entities of the target entities",
					Optional:            true,
					Validators: []validator.Bool{
						boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("average_entities")),
						boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("average_by_property")),
						boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("aggregate_by_property")),
					},
				},
				"average_entities": schema.SingleNestedAttribute{
					MarkdownDescription: "Function to average the entities of the target entities",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"average_of": schema.StringAttribute{
							MarkdownDescription: "The time periods to calculate the average of, e.g. hour, day, week, month",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("day"),
							Validators: []validator.String{
								stringvalidator.OneOf("hour", "day", "week", "month"),
							},
						},
						"measure_time_by": schema.StringAttribute{
							MarkdownDescription: "The property name on which to calculate the the time periods, e.g. $createdAt, $updated_at or any other date property",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("$createdAt"),
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("count_entities"),
							path.MatchRelative().AtParent().AtName("average_by_property"),
							path.MatchRelative().AtParent().AtName("aggregate_by_property"),
						),
					},
				},
				"average_by_property": schema.SingleNestedAttribute{
					MarkdownDescription: "Function to calculate the average by property value of the target entities",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"average_of": schema.StringAttribute{
							MarkdownDescription: "The time periods to calculate the average by, e.g. hour, day, week, month",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("hour", "day", "week", "month", "total"),
							},
						},
						"measure_time_by": schema.StringAttribute{
							MarkdownDescription: "The property name on which to calculate the the time periods, e.g. $createdAt, $updated_at or any other date property",
							Required:            true,
						},
						"property": schema.StringAttribute{
							MarkdownDescription: "The property name on which to calculate the average by",
							Required:            true,
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("count_entities"),
							path.MatchRelative().AtParent().AtName("average_entities"),
							path.MatchRelative().AtParent().AtName("aggregate_by_property"),
						),
					},
				},
				"aggregate_by_property": schema.SingleNestedAttribute{
					MarkdownDescription: "Function to calculate the aggregate by property value of the target entities, such as sum, min, max, median",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"func": schema.StringAttribute{
							MarkdownDescription: "The func of the aggregate by property",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("sum", "min", "max", "median"),
							},
						},
						"property": schema.StringAttribute{
							MarkdownDescription: "The property of the aggregate by property",
							Required:            true,
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("count_entities"),
							path.MatchRelative().AtParent().AtName("average_entities"),
							path.MatchRelative().AtParent().AtName("average_by_property"),
						),
					},
				},
			},
		},
		"query": schema.StringAttribute{
			MarkdownDescription: "Query to filter the target entities",
			Optional:            true,
		},
	}
}

func (r *AggregationPropertyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: AggregationPropertyResourceMarkdownDescription,
		Attributes:          AggregationPropertySchema(),
	}
}

var AggregationPropertyResourceMarkdownDescription = `

# Aggregation Property

This resource allows you to manage an aggregation property.
`
