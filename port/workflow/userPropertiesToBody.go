package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func userPropertiesToBody(ctx context.Context, model *UserPropertiesModel) (map[string]cli.WorkflowInputProperty, []string, error) {
	properties := map[string]cli.WorkflowInputProperty{}
	var required []string

	if model == nil {
		return properties, nil, nil
	}

	for identifier, prop := range model.StringProps {
		property, err := stringPropToBody(ctx, prop)
		if err != nil {
			return nil, nil, fmt.Errorf("string input %q: %w", identifier, err)
		}
		properties[identifier] = *property
		required = appendIfRequired(required, identifier, prop.Required)
	}

	for identifier, prop := range model.NumberProps {
		property, err := numberPropToBody(ctx, prop)
		if err != nil {
			return nil, nil, fmt.Errorf("number input %q: %w", identifier, err)
		}
		properties[identifier] = *property
		required = appendIfRequired(required, identifier, prop.Required)
	}

	for identifier, prop := range model.BooleanProps {
		property, err := booleanPropToBody(ctx, prop)
		if err != nil {
			return nil, nil, fmt.Errorf("boolean input %q: %w", identifier, err)
		}
		properties[identifier] = *property
		required = appendIfRequired(required, identifier, prop.Required)
	}

	for identifier, prop := range model.ObjectProps {
		property, err := objectPropToBody(ctx, prop)
		if err != nil {
			return nil, nil, fmt.Errorf("object input %q: %w", identifier, err)
		}
		properties[identifier] = *property
		required = appendIfRequired(required, identifier, prop.Required)
	}

	for identifier, prop := range model.ArrayProps {
		property, err := arrayPropToBody(ctx, prop)
		if err != nil {
			return nil, nil, fmt.Errorf("array input %q: %w", identifier, err)
		}
		properties[identifier] = *property
		required = appendIfRequired(required, identifier, prop.Required)
	}

	sort.Strings(required)

	return properties, required, nil
}

func appendIfRequired(required []string, identifier string, value types.Bool) []string {
	if value.ValueBool() {
		return append(required, identifier)
	}
	return required
}

func commonPropToBody(ctx context.Context, property *cli.WorkflowInputProperty, common propCommon) error {
	property.Title = common.Title.ValueStringPointer()
	property.Icon = common.Icon.ValueStringPointer()
	property.Description = common.Description.ValueStringPointer()

	if !common.DependsOn.IsNull() {
		dependsOn, err := utils.TerraformListToGoArray(ctx, *common.DependsOn, "string")
		if err != nil {
			return err
		}
		property.DependsOn = utils.InterfaceToStringArray(dependsOn)
	}

	property.Visible = boolOrJqToBody(*common.Visible, *common.VisibleJqQuery)
	property.ReadOnly = boolOrJqToBody(*common.ReadOnly, *common.ReadOnlyJqQuery)
	property.Disabled = boolOrJqToBody(*common.Disabled, *common.DisabledJqQuery)

	return nil
}

func boolOrJqToBody(value types.Bool, jqQuery types.String) any {
	if !jqQuery.IsNull() {
		return jqQueryBody(jqQuery)
	}
	if !value.IsNull() {
		return value.ValueBool()
	}
	return nil
}

func jqQueryBody(jqQuery types.String) map[string]string {
	return map[string]string{"jqQuery": jqQuery.ValueString()}
}

func stringPropToBody(ctx context.Context, prop StringPropModel) (*cli.WorkflowInputProperty, error) {
	property := &cli.WorkflowInputProperty{Type: "string"}
	if err := commonPropToBody(ctx, property, prop.common()); err != nil {
		return nil, err
	}

	property.Format = prop.Format.ValueStringPointer()
	property.Blueprint = prop.Blueprint.ValueStringPointer()

	if !prop.DefaultJqQuery.IsNull() {
		property.Default = jqQueryBody(prop.DefaultJqQuery)
	} else if !prop.Default.IsNull() {
		property.Default = prop.Default.ValueString()
	}

	if !prop.MinLength.IsNull() {
		minLength := int(prop.MinLength.ValueInt64())
		property.MinLength = &minLength
	}

	if !prop.MaxLength.IsNull() {
		maxLength := int(prop.MaxLength.ValueInt64())
		property.MaxLength = &maxLength
	}

	if !prop.PatternJqQuery.IsNull() {
		property.Pattern = jqQueryBody(prop.PatternJqQuery)
	} else if !prop.Pattern.IsNull() {
		property.Pattern = prop.Pattern.ValueString()
	}

	enum, err := enumToBody(ctx, prop.Enum, prop.EnumJqQuery, "string")
	if err != nil {
		return nil, err
	}
	property.Enum = enum

	enumColors, err := enumColorsToBody(ctx, prop.EnumColors)
	if err != nil {
		return nil, err
	}
	property.EnumColors = enumColors

	dataset, err := datasetToBody(prop.Dataset)
	if err != nil {
		return nil, err
	}
	property.Dataset = dataset
	property.Sort = sortToBody(prop.Sort)

	if encryption := encryptionToBody(prop.Encryption, prop.ClientSideEncryption); encryption != nil {
		property.Encryption = encryption
	}

	return property, nil
}

func numberPropToBody(ctx context.Context, prop NumberPropModel) (*cli.WorkflowInputProperty, error) {
	property := &cli.WorkflowInputProperty{Type: "number"}
	if err := commonPropToBody(ctx, property, prop.common()); err != nil {
		return nil, err
	}

	if !prop.DefaultJqQuery.IsNull() {
		property.Default = jqQueryBody(prop.DefaultJqQuery)
	} else if !prop.Default.IsNull() {
		property.Default = prop.Default.ValueFloat64()
	}

	property.Minimum = prop.Minimum.ValueFloat64Pointer()
	property.Maximum = prop.Maximum.ValueFloat64Pointer()
	property.ExclusiveMinimum = prop.ExclusiveMinimum.ValueFloat64Pointer()
	property.ExclusiveMaximum = prop.ExclusiveMaximum.ValueFloat64Pointer()

	enum, err := enumToBody(ctx, prop.Enum, prop.EnumJqQuery, "float64")
	if err != nil {
		return nil, err
	}
	property.Enum = enum

	return property, nil
}

func booleanPropToBody(ctx context.Context, prop BooleanPropModel) (*cli.WorkflowInputProperty, error) {
	property := &cli.WorkflowInputProperty{Type: "boolean"}
	if err := commonPropToBody(ctx, property, prop.common()); err != nil {
		return nil, err
	}

	if !prop.DefaultJqQuery.IsNull() {
		property.Default = jqQueryBody(prop.DefaultJqQuery)
	} else if !prop.Default.IsNull() {
		property.Default = prop.Default.ValueBool()
	}

	return property, nil
}

func objectPropToBody(ctx context.Context, prop ObjectPropModel) (*cli.WorkflowInputProperty, error) {
	property := &cli.WorkflowInputProperty{Type: "object"}
	if err := commonPropToBody(ctx, property, prop.common()); err != nil {
		return nil, err
	}

	property.Format = prop.Format.ValueStringPointer()

	if !prop.DefaultJqQuery.IsNull() {
		property.Default = jqQueryBody(prop.DefaultJqQuery)
	} else if !prop.Default.IsNull() {
		defaultValue := map[string]any{}
		if err := json.Unmarshal([]byte(prop.Default.ValueString()), &defaultValue); err != nil {
			return nil, fmt.Errorf("`default` must be a JSON object: %w", err)
		}
		property.Default = defaultValue
	}

	if encryption := encryptionToBody(prop.Encryption, prop.ClientSideEncryption); encryption != nil {
		property.Encryption = encryption
	}

	return property, nil
}

func arrayPropToBody(ctx context.Context, prop ArrayPropModel) (*cli.WorkflowInputProperty, error) {
	property := &cli.WorkflowInputProperty{Type: "array"}
	if err := commonPropToBody(ctx, property, prop.common()); err != nil {
		return nil, err
	}

	if !prop.DefaultJqQuery.IsNull() {
		property.Default = jqQueryBody(prop.DefaultJqQuery)
	}

	if !prop.MinItemsJqQuery.IsNull() {
		property.MinItems = jqQueryBody(prop.MinItemsJqQuery)
	} else if !prop.MinItems.IsNull() {
		property.MinItems = int(prop.MinItems.ValueInt64())
	}

	if !prop.MaxItemsJqQuery.IsNull() {
		property.MaxItems = jqQueryBody(prop.MaxItemsJqQuery)
	} else if !prop.MaxItems.IsNull() {
		property.MaxItems = int(prop.MaxItems.ValueInt64())
	}

	property.UniqueItems = prop.UniqueItems.ValueBoolPointer()
	property.Sort = sortToBody(prop.Sort)

	if err := arrayItemsToBody(ctx, property, prop); err != nil {
		return nil, err
	}

	return property, nil
}

func arrayItemsToBody(ctx context.Context, property *cli.WorkflowInputProperty, prop ArrayPropModel) error {
	switch {
	case prop.StringItems != nil:
		items := map[string]any{"type": "string"}
		if !prop.StringItems.Format.IsNull() {
			items["format"] = prop.StringItems.Format.ValueString()
		}
		if !prop.StringItems.Blueprint.IsNull() {
			items["blueprint"] = prop.StringItems.Blueprint.ValueString()
		}
		if !prop.StringItems.Dataset.IsNull() {
			dataset, err := utils.TerraformJsonStringToGoObject(prop.StringItems.Dataset.ValueStringPointer())
			if err != nil {
				return fmt.Errorf("`string_items.dataset` must be a JSON object: %w", err)
			}
			items["dataset"] = dataset
		}

		enum, err := enumToBody(ctx, prop.StringItems.Enum, prop.StringItems.EnumJqQuery, "string")
		if err != nil {
			return err
		}
		if enum != nil {
			items["enum"] = enum
		}

		enumColors, err := enumColorsToBody(ctx, prop.StringItems.EnumColors)
		if err != nil {
			return err
		}
		if enumColors != nil {
			items["enumColors"] = enumColors
		}

		if !prop.StringItems.Pattern.IsNull() {
			items["pattern"] = prop.StringItems.Pattern.ValueString()
		}
		if !prop.StringItems.MinLength.IsNull() {
			items["minLength"] = int(prop.StringItems.MinLength.ValueInt64())
		}
		if !prop.StringItems.MaxLength.IsNull() {
			items["maxLength"] = int(prop.StringItems.MaxLength.ValueInt64())
		}

		if !prop.StringItems.Default.IsNull() {
			defaults, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "string")
			if err != nil {
				return err
			}
			property.Default = defaults
		}

		property.Items = items

	case prop.NumberItems != nil:
		items := map[string]any{"type": "number"}

		enum, err := enumToBody(ctx, prop.NumberItems.Enum, prop.NumberItems.EnumJqQuery, "float64")
		if err != nil {
			return err
		}
		if enum != nil {
			items["enum"] = enum
		}

		enumColors, err := enumColorsToBody(ctx, prop.NumberItems.EnumColors)
		if err != nil {
			return err
		}
		if enumColors != nil {
			items["enumColors"] = enumColors
		}

		if !prop.NumberItems.Default.IsNull() {
			defaults, err := utils.TerraformListToGoArray(ctx, prop.NumberItems.Default, "float64")
			if err != nil {
				return err
			}
			property.Default = defaults
		}

		property.Items = items

	case prop.ObjectItems != nil:
		items := map[string]any{"type": "object"}
		if !prop.ObjectItems.Format.IsNull() {
			items["format"] = prop.ObjectItems.Format.ValueString()
		}

		if !prop.ObjectItems.Default.IsNull() {
			defaults := make([]map[string]any, 0, len(prop.ObjectItems.Default.Elements()))
			for _, element := range prop.ObjectItems.Default.Elements() {
				elementMap, ok := element.(types.Map)
				if !ok {
					return fmt.Errorf("`object_items.default` must hold maps, got %T", element)
				}
				value := map[string]any{}
				for key, entry := range elementMap.Elements() {
					entryString, ok := entry.(types.String)
					if !ok {
						return fmt.Errorf("`object_items.default` must hold string values, got %T", entry)
					}
					value[key] = entryString.ValueString()
				}
				defaults = append(defaults, value)
			}
			property.Default = defaults
		}

		property.Items = items
	}

	return nil
}

func enumToBody(ctx context.Context, enum types.List, enumJqQuery types.String, elementType string) (any, error) {
	if !enumJqQuery.IsNull() {
		return jqQueryBody(enumJqQuery), nil
	}
	if enum.IsNull() {
		return nil, nil
	}

	values, err := utils.TerraformListToGoArray(ctx, enum, elementType)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func enumColorsToBody(ctx context.Context, enumColors types.Map) (map[string]string, error) {
	if enumColors.IsNull() {
		return nil, nil
	}

	colors := map[string]string{}
	for key, value := range enumColors.Elements() {
		valueString, ok := value.(types.String)
		if !ok {
			return nil, fmt.Errorf("`enum_colors` must hold strings, got %T", value)
		}
		colors[key] = valueString.ValueString()
	}
	return colors, nil
}

func sortToBody(model *EntitiesSortModel) *cli.EntitiesSortModel {
	if model == nil {
		return nil
	}
	return &cli.EntitiesSortModel{
		Property: model.Property.ValueString(),
		Order:    model.Order.ValueString(),
	}
}

func datasetToBody(model *DatasetModel) (*cli.WorkflowDataset, error) {
	if model == nil {
		return nil, nil
	}

	dataset := &cli.WorkflowDataset{Combinator: model.Combinator.ValueString()}
	dataset.Rules = make([]cli.WorkflowDatasetRule, 0, len(model.Rules))
	for _, rule := range model.Rules {
		converted, err := datasetRuleToBody(rule)
		if err != nil {
			return nil, err
		}
		dataset.Rules = append(dataset.Rules, converted)
	}
	return dataset, nil
}

func datasetRuleToBody(model DatasetRuleModel) (cli.WorkflowDatasetRule, error) {
	rule := cli.WorkflowDatasetRule{}

	if !model.Combinator.IsNull() && model.Combinator.ValueString() != "" {
		combinator := model.Combinator.ValueString()
		rule.Combinator = &combinator
		rule.Rules = make([]cli.WorkflowDatasetRule, 0, len(model.Rules))
		for _, nested := range model.Rules {
			converted, err := datasetRuleToBody(nested)
			if err != nil {
				return rule, err
			}
			rule.Rules = append(rule.Rules, converted)
		}
		return rule, nil
	}

	rule.Operator = model.Operator.ValueString()
	rule.Blueprint = model.Blueprint.ValueStringPointer()
	rule.Property = model.Property.ValueStringPointer()

	switch {
	case model.Value != nil && !model.Value.JqQuery.IsNull():
		rule.Value = map[string]string{"jqQuery": model.Value.JqQuery.ValueString()}
	case !model.ValueJson.IsNull():
		var value any
		if err := json.Unmarshal([]byte(model.ValueJson.ValueString()), &value); err != nil {
			return rule, fmt.Errorf("`dataset` rule `value_json` must be valid JSON: %w", err)
		}
		rule.Value = value
	}

	return rule, nil
}
