package blueprint

import (
	"context"
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func propsResourceToBody(ctx context.Context, state *BlueprintModel) (map[string]cli.BlueprintProperty, []string, error) {
	props := map[string]cli.BlueprintProperty{}
	var required []string
	if state.Properties != nil {
		if state.Properties.StringProps != nil {
			err := stringPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.ArrayProps != nil {
			err := arrayPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.NumberProps != nil {
			err := numberPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.BooleanProps != nil {
			booleanPropResourceToBody(state, props, &required)
		}

		if state.Properties.ObjectProps != nil {
			objectPropResourceToBody(state, props, &required)
		}

	}
	return props, required, nil
}

func relationsResourceToBody(state *BlueprintModel) map[string]cli.Relation {
	relations := map[string]cli.Relation{}

	for identifier, prop := range state.Relations {
		target := prop.Target.ValueString()
		relationProp := cli.Relation{
			Target: &target,
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			relationProp.Title = &title
		}
		if !prop.Many.IsNull() {
			many := prop.Many.ValueBool()
			relationProp.Many = &many
		}

		if !prop.Required.IsNull() {
			required := prop.Required.ValueBool()
			relationProp.Required = &required
		}

		relations[identifier] = relationProp
	}

	return relations
}

func mirrorPropertiesToBody(state *BlueprintModel) map[string]cli.BlueprintMirrorProperty {
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}

	for identifier, prop := range state.MirrorProperties {
		mirrorProp := cli.BlueprintMirrorProperty{
			Path: prop.Path.ValueString(),
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			mirrorProp.Title = &title
		}

		mirrorProperties[identifier] = mirrorProp
	}

	return mirrorProperties
}

func calculationPropertiesToBody(ctx context.Context, state *BlueprintModel) map[string]cli.BlueprintCalculationProperty {
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	for identifier, prop := range state.CalculationProperties {
		calculationProp := cli.BlueprintCalculationProperty{
			Calculation: prop.Calculation.ValueString(),
			Type:        prop.Type.ValueString(),
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			calculationProp.Title = &title
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			calculationProp.Description = &description
		}

		if !prop.Format.IsNull() {
			format := prop.Format.ValueString()
			calculationProp.Format = &format
		}

		if !prop.Colorized.IsNull() {
			colorized := prop.Colorized.ValueBool()
			calculationProp.Colorized = &colorized
		}

		if !prop.Colors.IsNull() {
			colors := make(map[string]string)
			for key, value := range prop.Colors.Elements() {
				colors[key] = value.String()
			}

			calculationProp.Colors = colors
		}

		calculationProperties[identifier] = calculationProp
	}

	return calculationProperties
}

func aggregationPropertiesToBody(ctx context.Context, state *BlueprintModel) (map[string]cli.BlueprintAggregationProperty, error) {
	aggregationProperties := map[string]cli.BlueprintAggregationProperty{}

	for identifier, prop := range state.AggregationProperties {
		aggregationProp := cli.BlueprintAggregationProperty{
			Target: prop.Target.ValueString(),
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			aggregationProp.Title = &title
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			aggregationProp.Description = &description
		}

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			aggregationProp.Icon = &icon
		}

		if !prop.Method.CountEntities.IsNull() {
			aggregationProp.CalculationSpec = map[string]string{
				"func":          "count",
				"calculationBy": "entities",
			}
		} else if prop.Method.AverageEntities != nil {
			aggregationProp.CalculationSpec = map[string]string{
				"func":          "average",
				"calculationBy": "entities",
				"averageOf":     prop.Method.AverageEntities.AverageOf.ValueString(),
				"measureTimeBy": prop.Method.AverageEntities.MeasureTimeBy.ValueString(),
			}
		} else if prop.Method.AverageByProperty != nil {
			aggregationProp.CalculationSpec = map[string]string{
				"func":          "average",
				"calculationBy": "property",
				"property":      prop.Method.AverageByProperty.Property.ValueString(),
				"averageOf":     prop.Method.AverageByProperty.AverageOf.ValueString(),
				"measureTimeBy": prop.Method.AverageByProperty.MeasureTimeBy.ValueString(),
			}
		} else if prop.Method.AggregateByProperty != nil {
			aggregationProp.CalculationSpec = map[string]string{
				"func":          prop.Method.AggregateByProperty.Func.ValueString(),
				"calculationBy": "property",
				"property":      prop.Method.AggregateByProperty.Property.ValueString(),
			}
		}

		query, err := queryToPortBody(prop.Query.ValueStringPointer())
		if err != nil {
			return nil, err
		}

		// don't set query, if it wasn't set in the state, as the backend only supports setting to an object with
		// the search format, and not empty map or nil
		if query != nil {
			aggregationProp.Query = *query
		}

		aggregationProperties[identifier] = aggregationProp
	}

	return aggregationProperties, nil
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
