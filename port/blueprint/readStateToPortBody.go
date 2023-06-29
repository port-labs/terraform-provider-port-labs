package blueprint

import (
	"context"
	"encoding/json"
	"log"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func stringPropResourceToBody(ctx context.Context, state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.Properties.StringProp {
		property := cli.BlueprintProperty{
			Type: "string",
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			property.Title = &title
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
		}

		if !prop.Format.IsNull() {
			format := prop.Format.ValueString()
			property.Format = &format
		}

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			property.Icon = &icon
		}

		if !prop.MinLength.IsNull() {
			property.MinLength = int(prop.MinLength.ValueInt64())
		}

		if !prop.MaxLength.IsNull() {
			property.MaxLength = int(prop.MaxLength.ValueInt64())
		}

		if !prop.Spec.IsNull() {
			spec := prop.Spec.ValueString()
			property.Spec = &spec
		}

		if prop.SpecAuthentication != nil {
			specAuth := &cli.SpecAuthentication{
				AuthorizationUrl: prop.SpecAuthentication.AuthorizationUrl.ValueString(),
				TokenUrl:         prop.SpecAuthentication.TokenUrl.ValueString(),
				ClientId:         prop.SpecAuthentication.ClientId.ValueString(),
			}
			property.SpecAuthentication = specAuth
		}

		if !prop.Pattern.IsNull() {
			property.Pattern = prop.Pattern.ValueString()
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			property.Description = &description
		}

		if !prop.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "string")
			if err != nil {
				return err
			}
			property.Enum = enumList
		}

		if !prop.EnumColors.IsNull() {
			enumColor := map[string]string{}
			for k, v := range prop.EnumColors.Elements() {
				value, _ := v.ToTerraformValue(ctx)
				var keyValue string
				err := value.As(&keyValue)
				if err != nil {
					return err
				}
				enumColor[k] = keyValue
			}

			property.EnumColors = enumColor
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func numberPropResourceToBody(ctx context.Context, state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.Properties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "number",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueFloat64()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Minimum.IsNull() {
				minimum := prop.Minimum.ValueFloat64()
				property.Minimum = &minimum
			}

			if !prop.Maximum.IsNull() {
				maximum := prop.Maximum.ValueFloat64()
				property.Maximum = &maximum
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Enum.IsNull() {
				enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "float64")
				if err != nil {
					return err
				}
				property.Enum = enumList
			}

			if !prop.EnumColors.IsNull() {
				property.EnumColors = map[string]string{}
				for k, v := range prop.EnumColors.Elements() {
					value, _ := v.ToTerraformValue(ctx)
					var keyValue string
					err := value.As(&keyValue)
					if err != nil {
						return err
					}
					property.EnumColors[k] = keyValue
				}
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func booleanPropResourceToBody(state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range state.Properties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "boolean",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueBool()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func objectPropResourceToBody(state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range state.Properties.ObjectProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "object",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				defaultAsString := prop.Default.ValueString()
				defaultObj := make(map[string]interface{})
				err := json.Unmarshal([]byte(defaultAsString), &defaultObj)
				if err != nil {
					log.Fatal(err)
				} else {
					property.Default = defaultObj
				}
			}

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Spec.IsNull() {
				spec := prop.Spec.ValueString()
				property.Spec = &spec
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func arrayPropResourceToBody(ctx context.Context, state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.Properties.ArrayProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "array",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}
			if !prop.MinItems.IsNull() {
				minItems := int(prop.MinItems.ValueInt64())
				property.MinItems = &minItems
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = &maxItems
			}

			if prop.StringItems != nil {
				items := map[string]interface{}{}
				items["type"] = "string"
				if !prop.StringItems.Format.IsNull() {
					items["format"] = prop.StringItems.Format.ValueString()
				}
				if !prop.StringItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "string")
					if err != nil {
						return err
					}

					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.NumberItems != nil {
				items := map[string]interface{}{}
				items["type"] = "number"
				if !prop.NumberItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.NumberItems.Default, "float64")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.BooleanItems != nil {
				items := map[string]interface{}{}
				items["type"] = "boolean"
				if !prop.BooleanItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.BooleanItems.Default, "bool")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.ObjectItems != nil {
				items := map[string]interface{}{}
				items["type"] = "object"
				if !prop.ObjectItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.ObjectItems.Default, "object")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}

	return nil
}

func readStateToPortBody(ctx context.Context, state *BlueprintModel) (map[string]cli.BlueprintProperty, []string, error) {
	props := map[string]cli.BlueprintProperty{}
	var required []string
	if state.Properties != nil {
		if state.Properties.StringProp != nil {
			err := stringPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.ArrayProp != nil {
			err := arrayPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.NumberProp != nil {
			err := numberPropResourceToBody(ctx, state, props, &required)
			if err != nil {
				return nil, nil, err
			}
		}
		if state.Properties.BooleanProp != nil {
			booleanPropResourceToBody(state, props, &required)
		}

		if state.Properties.ObjectProp != nil {
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
			mirrorProp.Title = prop.Title.ValueString()
		}

		mirrorProperties[identifier] = mirrorProp
	}

	return mirrorProperties
}

func calculationPropertiesToBody(state *BlueprintModel) map[string]cli.BlueprintCalculationProperty {
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	for identifier, prop := range state.CalculationProperties {
		calculationProp := cli.BlueprintCalculationProperty{
			Calculation: prop.Calculation.ValueString(),
			Type:        prop.Type.ValueString(),
		}

		if !prop.Title.IsNull() {
			calculationProp.Title = prop.Title.ValueString()
		}

		if !prop.Description.IsNull() {
			calculationProp.Description = prop.Description.ValueString()
		}

		if !prop.Format.IsNull() {
			calculationProp.Format = prop.Format.ValueString()
		}

		calculationProperties[identifier] = calculationProp
	}

	return calculationProperties
}
