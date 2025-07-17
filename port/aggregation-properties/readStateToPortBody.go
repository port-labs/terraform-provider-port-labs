package aggregation_properties

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func aggregationPropertiesToBody(state *AggregationPropertiesModel) (*map[string]cli.BlueprintAggregationProperty, error) {
	if state == nil {
		return nil, nil
	}

	aggregationProperties := make(map[string]cli.BlueprintAggregationProperty)

	for aggregationPropertyIdentifier, aggregationProperty := range state.Properties {

		newAggregationProperty := cli.BlueprintAggregationProperty{
			Title:       aggregationProperty.Title.ValueStringPointer(),
			Icon:        aggregationProperty.Icon.ValueStringPointer(),
			Description: aggregationProperty.Description.ValueStringPointer(),
			Target:      aggregationProperty.TargetBlueprintIdentifier.ValueString(),
		}

		if !aggregationProperty.Method.CountEntities.IsNull() {
			newAggregationProperty.CalculationSpec = map[string]string{
				"func":          "count",
				"calculationBy": "entities",
			}
		} else if aggregationProperty.Method.AverageEntities != nil {
			newAggregationProperty.CalculationSpec = map[string]string{
				"func":          "average",
				"calculationBy": "entities",
				"averageOf":     aggregationProperty.Method.AverageEntities.AverageOf.ValueString(),
				"measureTimeBy": aggregationProperty.Method.AverageEntities.MeasureTimeBy.ValueString(),
			}
		} else if aggregationProperty.Method.AverageByProperty != nil {
			newAggregationProperty.CalculationSpec = map[string]string{
				"func":          "average",
				"calculationBy": "property",
				"property":      aggregationProperty.Method.AverageByProperty.Property.ValueString(),
				"averageOf":     aggregationProperty.Method.AverageByProperty.AverageOf.ValueString(),
				"measureTimeBy": aggregationProperty.Method.AverageByProperty.MeasureTimeBy.ValueString(),
			}
		} else if aggregationProperty.Method.AggregateByProperty != nil {
			newAggregationProperty.CalculationSpec = map[string]string{
				"func":          aggregationProperty.Method.AggregateByProperty.Func.ValueString(),
				"calculationBy": "property",
				"property":      aggregationProperty.Method.AggregateByProperty.Property.ValueString(),
			}
		}

		query, err := utils.TerraformJsonStringToGoObject(aggregationProperty.Query.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		// don't set query, if it wasn't set in the state, as the backend only supports setting to an object with
		// the search format, and not empty map or nil
		if query != nil {
			newAggregationProperty.Query = *query
		}

		if len(aggregationProperty.PathFilter) > 0 {
			pathFilter := make([]cli.AggregationPropertyPathFilter, len(aggregationProperty.PathFilter))
			for i, pf := range aggregationProperty.PathFilter {
				pathFilter[i] = cli.AggregationPropertyPathFilter{}

				// Only set FromBlueprint if it's not null
				if !pf.FromBlueprint.IsNull() {
					pathFilter[i].FromBlueprint = pf.FromBlueprint.ValueString()
				}

				var pathStrings []string
				for _, pathElement := range pf.Path.Elements() {
					if pathStr, ok := pathElement.(types.String); ok {
						pathStrings = append(pathStrings, pathStr.ValueString())
					}
				}
				pathFilter[i].Path = pathStrings
			}
			newAggregationProperty.PathFilter = pathFilter
		}

		aggregationProperties[aggregationPropertyIdentifier] = newAggregationProperty
	}

	return &aggregationProperties, nil
}
