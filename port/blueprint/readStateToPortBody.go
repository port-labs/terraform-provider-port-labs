package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
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

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			relationProp.Description = &description
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

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			calculationProp.Icon = &icon
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
				if stringValue, ok := value.(basetypes.StringValue); ok {
					colors[key] = stringValue.ValueString()
				}
			}

			calculationProp.Colors = colors
		}

		calculationProperties[identifier] = calculationProp
	}

	return calculationProperties
}

func PropertiesToBody(ctx context.Context, state *PropertiesModel) (map[string]cli.BlueprintProperty, []string, error) {
	bm := &BlueprintModel{Properties: state}
	return propsResourceToBody(ctx, bm)
}

func RelationsToBody(state map[string]RelationModel) map[string]cli.Relation {
	bm := &BlueprintModel{Relations: state}
	return relationsResourceToBody(bm)
}

func MirrorPropertiesToBody(state map[string]MirrorPropertyModel) map[string]cli.BlueprintMirrorProperty {
	bm := &BlueprintModel{MirrorProperties: state}
	return mirrorPropertiesToBody(bm)
}

func CalculationPropertiesToBody(ctx context.Context, state map[string]CalculationPropertyModel) map[string]cli.BlueprintCalculationProperty {
	bm := &BlueprintModel{CalculationProperties: state}
	return calculationPropertiesToBody(ctx, bm)
}
