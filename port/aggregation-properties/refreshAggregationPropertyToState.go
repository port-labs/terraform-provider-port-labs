package aggregation_properties

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func refreshAggregationPropertyState(state *AggregationPropertiesModel, aggregationProperties map[string]cli.BlueprintAggregationProperty) error {
	state.ID = state.BlueprintIdentifier

	state.Properties = map[string]*AggregationPropertyModel{}

	for aggregationPropertyIdentifier, aggregationProperty := range aggregationProperties {

		state.Properties[aggregationPropertyIdentifier] = &AggregationPropertyModel{
			types.StringPointerValue(aggregationProperty.Title),
			types.StringPointerValue(aggregationProperty.Icon),
			types.StringPointerValue(aggregationProperty.Description),
			types.StringValue(aggregationProperty.Target),
			nil,
			types.StringPointerValue(nil),
		}

		if aggregationProperty.Query != nil {
			query, err := json.Marshal(aggregationProperty.Query)
			if err != nil {
				return err
			}
			state.Properties[aggregationPropertyIdentifier].Query = types.StringValue(string(query))
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