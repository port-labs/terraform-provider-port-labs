package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func numberPropResourceToBody(ctx context.Context, state *PropertiesModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.NumberProps {
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

func AddNumberPropertiesToState(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
	numberProp := &NumberPropModel{
		Minimum: flex.GoFloat64ToFramework(v.Minimum),
		Maximum: flex.GoFloat64ToFramework(v.Maximum),
	}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
		}

		numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)
	} else {
		numberProp.Enum = types.ListNull(types.Float64Type)
	}

	if v.EnumColors != nil {
		numberProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		numberProp.EnumColors = types.MapNull(types.StringType)
	}

	return numberProp
}
