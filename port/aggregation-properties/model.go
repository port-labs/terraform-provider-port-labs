package aggregation_properties

import "github.com/hashicorp/terraform-plugin-framework/types"

type AggregationPropertiesModel struct {
	ID                  types.String                         `tfsdk:"id"`
	BlueprintIdentifier types.String                         `tfsdk:"blueprint_identifier"`
	Properties          map[string]*AggregationPropertyModel `tfsdk:"properties"`
}

type AggregationPropertyModel struct {
	Title                     types.String             `tfsdk:"title"`
	Icon                      types.String             `tfsdk:"icon"`
	Description               types.String             `tfsdk:"description"`
	TargetBlueprintIdentifier types.String             `tfsdk:"target_blueprint_identifier"`
	Method                    *AggregationMethodsModel `tfsdk:"method"`
	Query                     types.String             `tfsdk:"query"`
}

type AggregationMethodsModel struct {
	CountEntities       types.Bool                `tfsdk:"count_entities"`
	AverageEntities     *AverageEntitiesModel     `tfsdk:"average_entities"`
	AverageByProperty   *AverageByProperty        `tfsdk:"average_by_property"`
	AggregateByProperty *AggregateByPropertyModel `tfsdk:"aggregate_by_property"`
}

type AverageEntitiesModel struct {
	AverageOf     types.String `tfsdk:"average_of"`
	MeasureTimeBy types.String `tfsdk:"measure_time_by"`
}

type AverageByProperty struct {
	MeasureTimeBy types.String `tfsdk:"measure_time_by"`
	AverageOf     types.String `tfsdk:"average_of"`
	Property      types.String `tfsdk:"property"`
}

type AggregateByPropertyModel struct {
	Func     types.String `tfsdk:"func"`
	Property types.String `tfsdk:"property"`
}
