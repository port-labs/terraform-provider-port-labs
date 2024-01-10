package aggregation_property

import (
	"encoding/json"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func aggregationPropertyToBody(state *AggregationPropertyModel) (*cli.BlueprintAggregationProperty, error) {
	if state == nil {
		return nil, nil
	}

	aggregationProperty := cli.BlueprintAggregationProperty{
		Title:       state.Title.ValueStringPointer(),
		Icon:        state.Icon.ValueStringPointer(),
		Description: state.Description.ValueStringPointer(),
		Target:      state.TargetBlueprintIdentifier.ValueString(),
	}

	if !state.Method.CountEntities.IsNull() {
		aggregationProperty.CalculationSpec = map[string]string{
			"func":          "count",
			"calculationBy": "entities",
		}
	} else if state.Method.AverageEntities != nil {
		aggregationProperty.CalculationSpec = map[string]string{
			"func":          "average",
			"calculationBy": "entities",
			"averageOf":     state.Method.AverageEntities.AverageOf.ValueString(),
			"measureTimeBy": state.Method.AverageEntities.MeasureTimeBy.ValueString(),
		}
	} else if state.Method.AverageByProperty != nil {
		aggregationProperty.CalculationSpec = map[string]string{
			"func":          "average",
			"calculationBy": "property",
			"property":      state.Method.AverageByProperty.Property.ValueString(),
			"averageOf":     state.Method.AverageByProperty.AverageOf.ValueString(),
			"measureTimeBy": state.Method.AverageByProperty.MeasureTimeBy.ValueString(),
		}
	} else if state.Method.AggregateByProperty != nil {
		aggregationProperty.CalculationSpec = map[string]string{
			"func":          state.Method.AggregateByProperty.Func.ValueString(),
			"calculationBy": "property",
			"property":      state.Method.AggregateByProperty.Property.ValueString(),
		}
	}

	query, err := queryToPortBody(state.Query.ValueStringPointer())

	if err != nil {
		return nil, err
	}

	// don't set query, if it wasn't set in the state, as the backend only supports setting to an object with
	// the search format, and not empty map or nil
	if query != nil {
		aggregationProperty.Query = *query
	}

	return &aggregationProperty, nil
}

func queryToPortBody(query *string) (*map[string]any, error) {
	if query == nil || *query == "" {
		return nil, nil
	}

	queryMap := make(map[string]any)
	if err := json.Unmarshal([]byte(*query), &queryMap); err != nil {
		return nil, err
	}

	return &queryMap, nil
}
