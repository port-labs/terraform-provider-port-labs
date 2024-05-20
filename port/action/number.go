package action

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func numberPropResourceToBody(ctx context.Context, state *SelfServiceTriggerModel, props map[string]cli.ActionProperty, required *[]string) error {
	for propIdentifier, prop := range state.UserProperties.NumberProps {
		props[propIdentifier] = cli.ActionProperty{
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

			if !prop.DefaultJqQuery.IsNull() {
				defaultJqQuery := prop.DefaultJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": defaultJqQuery,
				}
				property.Default = jqQueryMap
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

			if !prop.EnumJqQuery.IsNull() {
				enumJqQueryMap := map[string]string{
					"jqQuery": prop.EnumJqQuery.ValueString(),
				}
				property.Enum = enumJqQueryMap
			}

			if !prop.DependsOn.IsNull() {
				dependsOn, err := utils.TerraformListToGoArray(ctx, prop.DependsOn, "string")
				if err != nil {
					return err
				}
				property.DependsOn = utils.InterfaceToStringArray(dependsOn)

			}

			if !prop.Visible.IsNull() {
				property.Visible = prop.Visible.ValueBoolPointer()
			}

			if !prop.VisibleJqQuery.IsNull() {
				VisibleJqQueryMap := map[string]string{
					"jqQuery": prop.VisibleJqQuery.ValueString(),
				}
				property.Visible = VisibleJqQueryMap
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addNumberPropertiesToResource(ctx context.Context, v *cli.ActionProperty) *NumberPropModel {
	numberProp := &NumberPropModel{
		Minimum: flex.GoFloat64ToFramework(v.Minimum),
		Maximum: flex.GoFloat64ToFramework(v.Maximum),
	}

	if v.Enum != nil {
		v := reflect.ValueOf(v.Enum)
		switch v.Kind() {
		case reflect.Slice:
			slice := v.Interface().([]interface{})
			attrs := make([]attr.Value, 0, v.Len())
			for _, value := range slice {
				attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
			}

			numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)

		case reflect.Map:
			v := v.Interface().(map[string]interface{})
			jqQueryValue := v["jqQuery"].(string)
			numberProp.EnumJqQuery = flex.GoStringToFramework(&jqQueryValue)
			numberProp.Enum = types.ListNull(types.Float64Type)
		}
	} else {
		numberProp.Enum = types.ListNull(types.Float64Type)
		numberProp.EnumJqQuery = types.StringNull()
	}

	return numberProp
}
