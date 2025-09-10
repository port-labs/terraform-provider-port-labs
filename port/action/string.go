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

func stringPropResourceToBody(ctx context.Context, d *SelfServiceTriggerModel, props map[string]cli.ActionProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.StringProps {
		property := cli.ActionProperty{
			Type: "string",
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			property.Title = &title
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
		}

		if !prop.DefaultJqQuery.IsNull() {
			defaultJqQuery := prop.DefaultJqQuery.ValueString()
			jqQueryMap := map[string]string{
				"jqQuery": defaultJqQuery,
			}
			property.Default = jqQueryMap
		}

		if !prop.Format.IsNull() {
			format := prop.Format.ValueString()
			property.Format = &format
		}

		if !prop.Blueprint.IsNull() {
			blueprint := prop.Blueprint.ValueString()
			property.Blueprint = &blueprint
		}

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			property.Icon = &icon
		}

		if !prop.MinLength.IsNull() {
			minLength := int(prop.MinLength.ValueInt64())
			property.MinLength = &minLength
		}

		if !prop.MaxLength.IsNull() {
			maxLength := int(prop.MaxLength.ValueInt64())
			property.MaxLength = &maxLength
		}

		if !prop.Pattern.IsNull() {
			pattern := prop.Pattern.ValueString()
			if pattern != "" {
				property.Pattern = &pattern
			}
		}

		if !prop.PatternJqQuery.IsNull() {
			patternJqQuery := prop.PatternJqQuery.ValueString()
			if patternJqQuery != "" {
				patternJqQueryMap := map[string]string{
					"jqQuery": patternJqQuery,
				}
				property.Pattern = patternJqQueryMap
			}
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

		if !prop.Encryption.IsNull() {
			encryption := prop.Encryption.ValueString()
			property.Encryption = &encryption
		}

		if prop.Dataset != nil {
			property.Dataset = actionDataSetToPortBody(prop.Dataset)
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

		if !prop.Disabled.IsNull() {
			val := prop.Disabled.ValueBool()
			property.Disabled = &val
		}

		if !prop.DisabledJqQuery.IsNull() {
			DisabledJqQuery := map[string]string{
				"jqQuery": prop.DisabledJqQuery.ValueString(),
			}
			property.Disabled = DisabledJqQuery 
		}

		if prop.Sort != nil {
			property.Sort = &cli.EntitiesSortModel{
				Property: prop.Sort.Property.ValueString(),
				Order:    prop.Sort.Order.ValueString(),
			}
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addStringPropertiesToResource(ctx context.Context, v *cli.ActionProperty) *StringPropModel {
	stringProp := &StringPropModel{
		MinLength:  flex.GoInt64ToFramework(v.MinLength),
		MaxLength:  flex.GoInt64ToFramework(v.MaxLength),
		Format:     flex.GoStringToFramework(v.Format),
		Blueprint:  flex.GoStringToFramework(v.Blueprint),
		Encryption: flex.GoStringToFramework(v.Encryption),
		Dataset:    writeDatasetToResource(v.Dataset),
	}

	stringProp.Pattern = types.StringNull()
	stringProp.PatternJqQuery = types.StringNull()

	if v.Pattern != nil {
		vPattern := reflect.ValueOf(v.Pattern)

		if vPattern.Kind() == reflect.String {
			// Regular pattern
			patternValue := v.Pattern.(string)
			if patternValue != "" {
				stringProp.Pattern = types.StringValue(patternValue)
			}
		} else if vPattern.Kind() == reflect.Map {
			// JQ Query pattern
			patternMap, ok := v.Pattern.(map[string]interface{})
			if ok && patternMap != nil {
				if jqQuery, ok := patternMap["jqQuery"]; ok && jqQuery != nil {
					jqQueryStr, isString := jqQuery.(string)
					if isString && jqQueryStr != "" {
						stringProp.PatternJqQuery = types.StringValue(jqQueryStr)
					}
				}
			}
		}
	}

	if v.Enum != nil {
		v := reflect.ValueOf(v.Enum)
		switch v.Kind() {
		case reflect.Slice:
			slice := v.Interface().([]interface{})
			attrs := make([]attr.Value, 0, v.Len())
			for _, value := range slice {
				attrs = append(attrs, basetypes.NewStringValue(value.(string)))
			}
			stringProp.Enum, _ = types.ListValue(types.StringType, attrs)

		case reflect.Map:
			v := v.Interface().(map[string]interface{})
			jqQueryValue := v["jqQuery"].(string)
			stringProp.EnumJqQuery = flex.GoStringToFramework(&jqQueryValue)
			stringProp.Enum = types.ListNull(types.StringType)

		}
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
		stringProp.EnumJqQuery = types.StringNull()
	}

	if v.EnumColors != nil {
		stringProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		stringProp.EnumColors = types.MapNull(types.StringType)
	}

	if v.Sort != nil {
		stringProp.Sort = &EntitiesSortModel{
			Property: types.StringValue(v.Sort.Property),
			Order:    types.StringValue(v.Sort.Order),
		}
	}

	return stringProp
}
