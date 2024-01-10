package aggregation_property

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func refreshAggregationPropertyState(state *AggregationPropertyModel, aggregationProperty cli.BlueprintAggregationProperty, blueprintIdentifier string, aggregationIdentifier string) error {
	state.ID = types.StringValue(fmt.Sprint("%s:%s", blueprintIdentifier, aggregationIdentifier))
	state.BlueprintIdentifier = types.StringValue(blueprintIdentifier)
	state.AggregationIdentifier = types.StringValue(aggregationIdentifier)
	state.TargetBlueprintIdentifier = types.StringValue(aggregationProperty.Target)
	state.Icon = types.StringPointerValue(aggregationProperty.Icon)
	state.Title = types.StringPointerValue(aggregationProperty.Title)
	state.Description = types.StringPointerValue(aggregationProperty.Description)

	if aggregationProperty.Query != nil {
		query, err := json.Marshal(aggregationProperty.Query)
		if err != nil {
			return err
		}
		state.Query = types.StringValue(string(query))
	}

	if aggregationProperty.CalculationSpec != nil {
		if calculationBy, ok := aggregationProperty.CalculationSpec["calculationBy"]; ok {
			if calculationBy == "entities" {
				if entitiesFunc, ok := aggregationProperty.CalculationSpec["func"]; ok {
					if entitiesFunc == "count" {
						state.Method = &AggregationMethodsModel{
							CountEntities: types.BoolValue(true),
						}
					} else if entitiesFunc == "average" {
						state.Method = &AggregationMethodsModel{
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
						state.Method = &AggregationMethodsModel{
							AverageByProperty: &AverageByProperty{
								MeasureTimeBy: types.StringValue(aggregationProperty.CalculationSpec["measureTimeBy"]),
								AverageOf:     types.StringValue(aggregationProperty.CalculationSpec["averageOf"]),
								Property:      types.StringValue(aggregationProperty.CalculationSpec["property"]),
							},
						}
					} else {
						state.Method = &AggregationMethodsModel{
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
	return nil
}
