package aggregation_properties

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func refreshAggregationPropertiesState(state *AggregationPropertiesModel, aggregationProperties map[string]cli.BlueprintAggregationProperty) error {
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
		}

		query, err := utils.GoObjectToTerraformString(aggregationProperty.Query)
		if err != nil {
			return err
		}
		state.Properties[aggregationPropertyIdentifier].Query = query

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
