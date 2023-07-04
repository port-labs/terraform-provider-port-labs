package blueprint

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/samber/lo"
)

func addStringPropertiesToState(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewStringValue(value.(string)))
		}

		stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
	}

	if v.EnumColors != nil {
		stringProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		stringProp.EnumColors = types.MapNull(types.StringType)
	}

	if v.Format != nil {
		stringProp.Format = types.StringValue(*v.Format)
	}

	if v.Spec != nil {
		stringProp.Spec = types.StringValue(*v.Spec)
	}

	if v.MinLength != 0 {
		stringProp.MinLength = types.Int64Value(int64(v.MinLength))
	}

	if v.MaxLength != 0 {
		stringProp.MaxLength = types.Int64Value(int64(v.MaxLength))
	}

	if v.Pattern != "" {
		stringProp.Pattern = types.StringValue(v.Pattern)
	}

	if v.SpecAuthentication != nil {
		stringProp.SpecAuthentication = &SpecAuthenticationModel{
			AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
			TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
			ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
		}
	}

	return stringProp
}

func addNumberPropertiesToState(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
	numberProp := &NumberPropModel{}
	if v.Minimum != nil {
		numberProp.Minimum = types.Float64Value(*v.Minimum)
	}

	if v.Maximum != nil {
		numberProp.Maximum = types.Float64Value(*v.Maximum)
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

func addObjectPropertiesToState(v *cli.BlueprintProperty) *ObjectPropModel {
	objectProp := &ObjectPropModel{}

	if v.Spec != nil {
		objectProp.Spec = types.StringValue(*v.Spec)
	}

	return objectProp
}

func addArrayPropertiesToState(v *cli.BlueprintProperty) *ArrayPropModel {
	arrayProp := &ArrayPropModel{}
	if v.MinItems != nil {
		arrayProp.MinItems = types.Int64Value(int64(*v.MinItems))
	}
	if v.MaxItems != nil {
		arrayProp.MaxItems = types.Int64Value(int64(*v.MaxItems))
	}
	if v.Items != nil {
		if v.Items["type"] != "" {
			switch v.Items["type"] {
			case "string":
				arrayProp.StringItems = &StringItems{}
				if v.Default != nil {
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
			case "number":
				arrayProp.NumberItems = &NumberItems{}
				if v.Default != nil {
					numberArray := make([]float64, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(numberArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
					}
					arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
				} else {
					arrayProp.NumberItems.Default = types.ListNull(types.Float64Type)
				}

			case "boolean":
				arrayProp.BooleanItems = &BooleanItems{}
				if v.Default != nil {
					booleanArray := make([]bool, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(booleanArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewBoolValue(value.(bool)))
					}
					arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
				} else {
					arrayProp.BooleanItems.Default = types.ListNull(types.BoolType)
				}

			case "object":
				arrayProp.ObjectItems = &ObjectItems{}
				if v.Default != nil {
					objectArray := make([]map[string]interface{}, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						objectArray[i] = v.(map[string]interface{})
					}
					attrs := make([]attr.Value, 0, len(objectArray))
					for _, value := range objectArray {
						js, _ := json.Marshal(&value)
						stringValue := string(js)
						attrs = append(attrs, basetypes.NewStringValue(stringValue))
					}
					arrayProp.ObjectItems.Default, _ = types.ListValue(types.StringType, attrs)
				} else {
					arrayProp.ObjectItems.Default = types.ListNull(types.StringType)
				}
			}
		}
	}

	return arrayProp
}

func setCommonProperties(v cli.BlueprintProperty, prop interface{}) {
	properties := []string{"Description", "Icon", "Default", "Title"}
	for _, property := range properties {
		switch property {
		case "Description":
			if v.Description != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Description = types.StringValue(*v.Description)
				case *NumberPropModel:
					p.Description = types.StringValue(*v.Description)
				case *BooleanPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ArrayPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ObjectPropModel:
					p.Description = types.StringValue(*v.Description)
				}
			}
		case "Icon":
			if v.Icon != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *NumberPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *BooleanPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ArrayPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ObjectPropModel:
					p.Icon = types.StringValue(*v.Icon)
				}
			}
		case "Title":
			if v.Title != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Title = types.StringValue(*v.Title)
				case *NumberPropModel:
					p.Title = types.StringValue(*v.Title)
				case *BooleanPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ArrayPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ObjectPropModel:
					p.Title = types.StringValue(*v.Title)
				}
			}
		case "Default":
			if v.Default != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Default = types.StringValue(v.Default.(string))
				case *NumberPropModel:
					p.Default = types.Float64Value(v.Default.(float64))
				case *BooleanPropModel:
					p.Default = types.BoolValue(v.Default.(bool))
				case *ObjectPropModel:
					js, _ := json.Marshal(v.Default)
					value := string(js)
					p.Default = types.StringValue(value)
				}
			}
		}
	}
}

func updatePropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) error {
	properties := &PropertiesModel{}

	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]StringPropModel)
			}
			stringProp := addStringPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, stringProp)

			properties.StringProp[k] = *stringProp

		case "number":
			if properties.NumberProp == nil {
				properties.NumberProp = make(map[string]NumberPropModel)
			}

			numberProp := addNumberPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, numberProp)

			properties.NumberProp[k] = *numberProp

		case "array":
			if properties.ArrayProp == nil {
				properties.ArrayProp = make(map[string]ArrayPropModel)
			}

			arrayProp := addArrayPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, arrayProp)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, booleanProp)

			if lo.Contains(b.Schema.Required, k) {
				booleanProp.Required = types.BoolValue(true)
			} else {
				booleanProp.Required = types.BoolValue(false)
			}

			properties.BooleanProp[k] = *booleanProp

		case "object":
			if properties.ObjectProp == nil {
				properties.ObjectProp = make(map[string]ObjectPropModel)
			}

			objectProp := addObjectPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, objectProp)

			properties.ObjectProp[k] = *objectProp

		}

	}

	bm.Properties = properties

	return nil
}

func addRelationsToState(b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.Relations {
		if bm.Relations == nil {
			bm.Relations = make(map[string]RelationModel)
		}

		relationModel := &RelationModel{
			Target: types.StringValue(*v.Target),
		}

		if v.Title != nil {
			relationModel.Title = types.StringValue(*v.Title)
		}

		if v.Many != nil {
			relationModel.Many = types.BoolValue(*v.Many)
		}

		if v.Required != nil {
			relationModel.Required = types.BoolValue(*v.Required)
		}

		bm.Relations[k] = *relationModel

	}
}

func addMirrorPropertiesToState(b *cli.Blueprint, bm *BlueprintModel) {
	if b.MirrorProperties != nil {
		for k, v := range b.MirrorProperties {
			if bm.MirrorProperties == nil {
				bm.MirrorProperties = make(map[string]MirrorPropertyModel)
			}

			mirrorPropertyModel := &MirrorPropertyModel{
				Path: types.StringValue(v.Path),
			}
			if v.Title != "" {
				mirrorPropertyModel.Title = types.StringValue(v.Title)
			}

			bm.MirrorProperties[k] = *mirrorPropertyModel

		}
	}
}

func addCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.CalculationProperties {
		if bm.CalculationProperties == nil {
			bm.CalculationProperties = make(map[string]CalculationPropertyModel)
		}

		calculationPropertyModel := &CalculationPropertyModel{
			Calculation: types.StringValue(v.Calculation),
			Type:        types.StringValue(v.Type),
		}
		if v.Title != nil {
			calculationPropertyModel.Title = types.StringValue(*v.Title)
		}

		if v.Description != nil {
			calculationPropertyModel.Description = types.StringValue(*v.Description)
		}

		if v.Format != nil {
			calculationPropertyModel.Format = types.StringValue(*v.Format)
		}

		if v.Colorized != nil {
			calculationPropertyModel.Colorized = types.BoolValue(*v.Colorized)
		}

		if v.Colors != nil {
			calculationPropertyModel.Colors, _ = types.MapValueFrom(ctx, types.StringType, Colors)
		} else {
			calculationPropertyModel.Colors = types.MapNull(types.StringType)
		}

		bm.CalculationProperties[k] = *calculationPropertyModel

	}
}
