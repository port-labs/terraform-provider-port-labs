package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func PropsResourceToBody(ctx context.Context, state *PropertiesModel) (map[string]cli.BlueprintProperty, []string, error) {
	props := map[string]cli.BlueprintProperty{}
	var required []string
	if state != nil {
		if state.StringProps != nil {
			err := stringPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.ArrayProps != nil {
			err := arrayPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.NumberProps != nil {
			err := numberPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.BooleanProps != nil {
			booleanPropResourceToBody(state, props, &required)
		}

		if state.ObjectProps != nil {
			objectPropResourceToBody(state, props, &required)
		}

	}
	return props, required, nil
}

func RelationsResourceToBody(state map[string]RelationModel) map[string]cli.Relation {
	relations := map[string]cli.Relation{}

	for identifier, prop := range state {
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

func MirrorPropertiesToBody(state map[string]MirrorPropertyModel) map[string]cli.BlueprintMirrorProperty {
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}

	for identifier, prop := range state {
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

func CalculationPropertiesToBody(ctx context.Context, state map[string]CalculationPropertyModel) map[string]cli.BlueprintCalculationProperty {
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	for identifier, prop := range state {
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

		if !prop.Spec.IsNull() {
			spec := prop.Spec.ValueString()
			calculationProp.Spec = &spec
		}

		if prop.SpecAuthentication != nil {
			specAuth := &cli.SpecAuthentication{
				AuthorizationUrl: prop.SpecAuthentication.AuthorizationUrl.ValueString(),
				TokenUrl:         prop.SpecAuthentication.TokenUrl.ValueString(),
				ClientId:         prop.SpecAuthentication.ClientId.ValueString(),
			}
			calculationProp.SpecAuthentication = specAuth
		}

		calculationProperties[identifier] = calculationProp
	}

	return calculationProperties
}
