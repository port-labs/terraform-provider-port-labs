package aggregation_properties

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *AggregationPropertiesResource) refreshAggregationPropertiesState(state *AggregationPropertiesModel, aggregationProperties map[string]cli.BlueprintAggregationProperty) error {
	state.ID = state.BlueprintIdentifier

	state.Properties = map[string]*AggregationPropertyModel{}

	for aggregationPropertyIdentifier, aggregationProperty := range aggregationProperties {

		state.Properties[aggregationPropertyIdentifier] = &AggregationPropertyModel{
			Title:                     types.StringPointerValue(aggregationProperty.Title),
			Icon:                      types.StringPointerValue(aggregationProperty.Icon),
			Description:               types.StringPointerValue(aggregationProperty.Description),
			TargetBlueprintIdentifier: types.StringValue(aggregationProperty.Target),
			Method:                    nil,
			Query:                     types.StringPointerValue(nil),
			PathFilter:                nil,
		}

		query, err := utils.GoObjectToTerraformString(aggregationProperty.Query, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}
		state.Properties[aggregationPropertyIdentifier].Query = query

		// Handle PathFilter conversion from Port API response to Terraform state
		if len(aggregationProperty.PathFilter) > 0 {
			pathFilter := make([]AggregationPropertyPathFilterModel, len(aggregationProperty.PathFilter))
			for i, pf := range aggregationProperty.PathFilter {
				pathFilter[i] = AggregationPropertyPathFilterModel{}
				
				// Set FromBlueprint, using null value if empty
				if pf.FromBlueprint != "" {
					pathFilter[i].FromBlueprint = types.StringValue(pf.FromBlueprint)
				} else {
					pathFilter[i].FromBlueprint = types.StringNull()
				}

				// Convert path from []string to types.List
				pathElements := make([]attr.Value, len(pf.Path))
				for j, pathStr := range pf.Path {
					pathElements[j] = types.StringValue(pathStr)
				}
				pathList, _ := types.ListValue(types.StringType, pathElements)
				pathFilter[i].Path = pathList
			}
			state.Properties[aggregationPropertyIdentifier].PathFilter = pathFilter
		}

		if aggregationProperty.CalculationSpec != nil {
			if calculationBy, ok := aggregationProperty.CalculationSpec["calculationBy"]; ok {
				if calculationBy == "entities" {
					if entitiesFunc, ok := aggregationProperty.CalculationSpec["func"]; ok {
						if entitiesFunc == "count" {
							state.Properties[aggregationPropertyIdentifier].Method = &AggregationMethodsModel{
								CountEntities: types.BoolValue(true),
							}
						} else if entitiesFunc == "average" {
							state.Properties[aggregationPropertyIdentifier].Method = &AggregationMethodsModel{
								AverageEntities: &AverageEntitiesModel{
									AverageOf:     types.StringValue(aggregationProperty.CalculationSpec["averageOf"]),
									MeasureTimeBy: types.StringValue(aggregationProperty.CalculationSpec["measureTimeBy"]),
								},
							}
						}
					}
				} else if calculationBy == "property" {
					if propertyFunc, ok := aggregationProperty.CalculationSpec["func"]; ok {
						if propertyFunc == "average" {
							state.Properties[aggregationPropertyIdentifier].Method = &AggregationMethodsModel{
								AverageByProperty: &AverageByProperty{
									MeasureTimeBy: types.StringValue(aggregationProperty.CalculationSpec["measureTimeBy"]),
									AverageOf:     types.StringValue(aggregationProperty.CalculationSpec["averageOf"]),
									Property:      types.StringValue(aggregationProperty.CalculationSpec["property"]),
								},
							}
						} else {
							state.Properties[aggregationPropertyIdentifier].Method = &AggregationMethodsModel{
								AggregateByProperty: &AggregateByPropertyModel{
									Func:     types.StringValue(aggregationProperty.CalculationSpec["func"]),
									Property: types.StringValue(aggregationProperty.CalculationSpec["property"]),
								},
							}
						}
					}
				}
			}
		}
	}
	return nil
}
