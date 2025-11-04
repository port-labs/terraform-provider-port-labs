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

func handleArrayItemsToBody(ctx context.Context, property *cli.ActionProperty, prop ArrayPropModel, required *[]string) error {
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

		if !prop.StringItems.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Enum, "string")
			if err != nil {
				return err
			}
			items["enum"] = enumList
		}

		if !prop.StringItems.Dataset.IsNull() {
			v, err := utils.TerraformJsonStringToGoObject(prop.StringItems.Dataset.ValueStringPointer())
			if err != nil {
				return err
			}

			items["dataset"] = v
		}

		if !prop.StringItems.Format.IsNull() {
			items["format"] = prop.StringItems.Format.ValueString()
		}

		if !prop.StringItems.Blueprint.IsNull() {
			items["blueprint"] = prop.StringItems.Blueprint.ValueString()
		}

		if !prop.StringItems.EnumJqQuery.IsNull() {
			enumJqQueryMap := map[string]string{
				"jqQuery": prop.StringItems.EnumJqQuery.ValueString(),
			}
			items["enum"] = enumJqQueryMap
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

			items["default"] = defaultList
		}

		if !prop.NumberItems.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.NumberItems.Enum, "float64")
			if err != nil {
				return err
			}
			items["enum"] = enumList
		}

		if !prop.NumberItems.EnumJqQuery.IsNull() {
			enumJqQueryMap := map[string]string{
				"jqQuery": prop.NumberItems.EnumJqQuery.ValueString(),
			}
			items["enum"] = enumJqQueryMap
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

			items["default"] = defaultList
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
			items["default"] = defaultList
		}

		property.Items = items
	}
	return nil
}

func arrayPropResourceToBody(ctx context.Context, d *SelfServiceTriggerModel, props map[string]cli.ActionProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ArrayProps {
		props[propIdentifier] = cli.ActionProperty{
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

			if !prop.DefaultJqQuery.IsNull() {
				defaultJqQuery := prop.DefaultJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": defaultJqQuery,
				}
				property.Default = jqQueryMap
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}
			if !prop.MinItems.IsNull() {
				minItems := int(prop.MinItems.ValueInt64())
				property.MinItems = minItems
			}

			if !prop.MinItemsJqQuery.IsNull() {
				minItemsJqQuery := prop.MinItemsJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": minItemsJqQuery,
				}
				property.MinItems = jqQueryMap
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = maxItems
			}

			if !prop.MaxItemsJqQuery.IsNull() {
				maxItemsJqQuery := prop.MaxItemsJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": maxItemsJqQuery,
				}
				property.MaxItems = jqQueryMap
			}

			if !prop.DependsOn.IsNull() {
				dependsOn, err := utils.TerraformListToGoArray(ctx, prop.DependsOn, "string")
				if err != nil {
					return err
				}
				property.DependsOn = utils.InterfaceToStringArray(dependsOn)

			}

			err := handleArrayItemsToBody(ctx, &property, prop, required)
			if err != nil {
				return err
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
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func (r *ActionResource) addArrayPropertiesToResource(v *cli.ActionProperty) (*ArrayPropModel, error) {
	arrayProp := &ArrayPropModel{}

	// Handle MinItems - can be int or JQ query
	if v.MinItems != nil {
		switch minItems := v.MinItems.(type) {
		case float64:
			arrayProp.MinItems = types.Int64Value(int64(minItems))
		case int:
			arrayProp.MinItems = types.Int64Value(int64(minItems))
		case map[string]interface{}:
			if jqQuery, ok := minItems["jqQuery"].(string); ok {
				arrayProp.MinItemsJqQuery = types.StringValue(jqQuery)
			}
		}
	}

	// Handle MaxItems - can be int or JQ query
	if v.MaxItems != nil {
		switch maxItems := v.MaxItems.(type) {
		case float64:
			arrayProp.MaxItems = types.Int64Value(int64(maxItems))
		case int:
			arrayProp.MaxItems = types.Int64Value(int64(maxItems))
		case map[string]interface{}:
			if jqQuery, ok := maxItems["jqQuery"].(string); ok {
				arrayProp.MaxItemsJqQuery = types.StringValue(jqQuery)
			}
		}
	}

	if v.Default != nil {
		switch v := v.Default.(type) {
		// We only test for map[string]interface{} ATM
		case map[string]interface{}:
			arrayProp.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
		}
	}

	if v.Sort != nil {
		arrayProp.Sort = &EntitiesSortModel{
			Property: types.StringValue(v.Sort.Property),
			Order:    types.StringValue(v.Sort.Order),
		}
	}

	if v.Items != nil {
		if v.Items["type"] != "" {
			switch v.Items["type"] {
			case "string":
				arrayProp.StringItems = &StringItems{}
				if v.Default != nil && arrayProp.DefaultJqQuery.IsNull() {
					stringArray := make([]string, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						stringArray[i] = v.(string)
					}
					attrs := make([]attr.Value, 0, len(stringArray))
					for _, value := range stringArray {
						attrs = append(attrs, basetypes.NewStringValue(value))
					}
					arrayProp.StringItems.Default, _ = types.ListValue(types.StringType, attrs)
				} else {
					arrayProp.StringItems.Default = types.ListNull(types.StringType)
				}
				if value, ok := v.Items["format"]; ok && value != nil {
					arrayProp.StringItems.Format = types.StringValue(v.Items["format"].(string))
				}
				if value, ok := v.Items["blueprint"]; ok && value != nil {
					arrayProp.StringItems.Blueprint = types.StringValue(v.Items["blueprint"].(string))
				}
				if value, ok := v.Items["dataset"]; ok && value != nil {
					ds, err := utils.GoObjectToTerraformString(v.Items["dataset"], r.portClient.JSONEscapeHTML)
					if err != nil {
						return nil, err
					}
					arrayProp.StringItems.Dataset = ds
				}

				if value, ok := v.Items["enum"]; ok && value != nil {
					v := reflect.ValueOf(value)
					switch v.Kind() {
					case reflect.Slice:
						slice := v.Interface().([]interface{})
						attrs := make([]attr.Value, 0, v.Len())
						for _, value := range slice {
							attrs = append(attrs, basetypes.NewStringValue(value.(string)))
						}
						arrayProp.StringItems.Enum, _ = types.ListValue(types.StringType, attrs)
					case reflect.Map:
						v := v.Interface().(map[string]interface{})
						jqQueryValue := v["jqQuery"].(string)
						arrayProp.StringItems.EnumJqQuery = flex.GoStringToFramework(&jqQueryValue)
						arrayProp.StringItems.Enum = types.ListNull(types.StringType)
					}
				} else {
					arrayProp.StringItems.Enum = types.ListNull(types.StringType)
				}

			case "number":
				arrayProp.NumberItems = &NumberItems{}
				if v.Default != nil && arrayProp.DefaultJqQuery.IsNull() {
					numberArray := make([]float64, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(numberArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
					}
					arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
				} else {
					arrayProp.NumberItems.Default = types.ListNull(types.Float64Type)
				}

				if value, ok := v.Items["enum"]; ok && value != nil {
					v := reflect.ValueOf(value)
					switch v.Kind() {
					case reflect.Slice:
						slice := v.Interface().([]interface{})
						attrs := make([]attr.Value, 0, v.Len())
						for _, value := range slice {
							attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
						}
						arrayProp.NumberItems.Enum, _ = types.ListValue(types.Float64Type, attrs)
					case reflect.Map:
						v := v.Interface().(map[string]interface{})
						jqQueryValue := v["jqQuery"].(string)
						arrayProp.NumberItems.EnumJqQuery = flex.GoStringToFramework(&jqQueryValue)
						arrayProp.NumberItems.Enum = types.ListNull(types.Float64Type)
					}
				} else {
					arrayProp.NumberItems.Enum = types.ListNull(types.Float64Type)
				}

			case "boolean":
				arrayProp.BooleanItems = &BooleanItems{}
				if v.Default != nil && arrayProp.DefaultJqQuery.IsNull() {
					booleanArray := make([]bool, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(booleanArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewBoolValue(value.(bool)))
					}
					arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
				}

			case "object":
				arrayProp.ObjectItems = &ObjectItems{}
				if v.Default != nil && arrayProp.DefaultJqQuery.IsNull() {
					objectArray := make([]map[string]interface{}, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(objectArray))
					for _, value := range v.Default.([]interface{}) {
						stringValue, err := utils.GoObjectToTerraformString(value, r.portClient.JSONEscapeHTML)
						if err != nil {
							return nil, err
						}
						attrs = append(attrs, stringValue)
					}
					arrayProp.ObjectItems.Default, _ = types.ListValue(types.StringType, attrs)
				}

			}
		}
	}

	return arrayProp, nil
}
