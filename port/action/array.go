package action

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
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

		if prop.StringItems.Dataset != nil {
			items["dataset"] = actionDataSetToPortBody(prop.StringItems.Dataset)
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

func arrayPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.ActionProperty, required *[]string) error {
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
				property.MinItems = &minItems
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = &maxItems
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

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addArrayPropertiesToResource(v *cli.ActionProperty) (*ArrayPropModel, error) {
	arrayProp := &ArrayPropModel{
		MinItems: flex.GoInt64ToFramework(v.MinItems),
		MaxItems: flex.GoInt64ToFramework(v.MaxItems),
	}

	if v.Default != nil {
		switch v := v.Default.(type) {
		// We only test for map[string]interface{} ATM
		case map[string]interface{}:
			arrayProp.DefaultJqQuery = types.StringValue(v["jqQuery"].(string))
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
					arrayProp.StringItems.Dataset = writeDatasetToResource(v.Items["dataset"].(*cli.Dataset))
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
						stringfiyValue, err := json.Marshal(value)
						if err != nil {
							return nil, err
						}
						stringValue := string(stringfiyValue)
						attrs = append(attrs, basetypes.NewStringValue(stringValue))
					}
					arrayProp.ObjectItems.Default, _ = types.ListValue(types.StringType, attrs)
				}

			}
		}
	}

	return arrayProp, nil
}
