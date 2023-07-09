package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func numberPropResourceToBody(ctx context.Context, state *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.UserProperties.NumberProps {
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

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
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

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addNumberPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
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

	return numberProp
}
