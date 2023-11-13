package action

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func stringPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.ActionProperty, required *[]string) error {
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
			property.Pattern = &pattern
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
		Pattern:    types.StringPointerValue(v.Pattern),
		Format:     types.StringPointerValue(v.Format),
		Blueprint:  types.StringPointerValue(v.Blueprint),
		Encryption: types.StringPointerValue(v.Encryption),
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
			stringProp.EnumJqQuery = types.StringPointerValue(&jqQueryValue)
			stringProp.Enum = types.ListNull(types.StringType)

		}
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
		stringProp.EnumJqQuery = types.StringNull()
	}

	return stringProp
}
