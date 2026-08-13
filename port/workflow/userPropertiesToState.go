package workflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/samber/lo"
)

// userPropertiesToState rebuilds the declared inputs from the API response.
// configured reports whether a user_properties block was declared, which decides
// whether a form without inputs becomes an empty block or no block at all.
func (r *WorkflowResource) userPropertiesToState(ctx context.Context, userInputs *cli.WorkflowUserInputs, configured bool) (*UserPropertiesModel, error) {
	model := &UserPropertiesModel{}
	_, required := requiredToState(userInputs.Required)

	for identifier, property := range userInputs.Properties {
		isRequired := lo.Contains(required, identifier)

		switch property.Type {
		case "string":
			prop, err := stringPropToState(ctx, property, r.portClient.JSONEscapeHTML)
			if err != nil {
				return nil, fmt.Errorf("string input %q: %w", identifier, err)
			}
			commonPropToState(ctx, property, prop.common(), isRequired)
			if model.StringProps == nil {
				model.StringProps = map[string]StringPropModel{}
			}
			model.StringProps[identifier] = *prop

		case "number":
			prop, err := numberPropToState(ctx, property)
			if err != nil {
				return nil, fmt.Errorf("number input %q: %w", identifier, err)
			}
			commonPropToState(ctx, property, prop.common(), isRequired)
			if model.NumberProps == nil {
				model.NumberProps = map[string]NumberPropModel{}
			}
			model.NumberProps[identifier] = *prop

		case "boolean":
			prop := booleanPropToState(property)
			commonPropToState(ctx, property, prop.common(), isRequired)
			if model.BooleanProps == nil {
				model.BooleanProps = map[string]BooleanPropModel{}
			}
			model.BooleanProps[identifier] = *prop

		case "object":
			prop, err := objectPropToState(property, r.portClient.JSONEscapeHTML)
			if err != nil {
				return nil, fmt.Errorf("object input %q: %w", identifier, err)
			}
			commonPropToState(ctx, property, prop.common(), isRequired)
			if model.ObjectProps == nil {
				model.ObjectProps = map[string]ObjectPropModel{}
			}
			model.ObjectProps[identifier] = *prop

		case "array":
			prop, err := arrayPropToState(ctx, property, r.portClient.JSONEscapeHTML)
			if err != nil {
				return nil, fmt.Errorf("array input %q: %w", identifier, err)
			}
			commonPropToState(ctx, property, prop.common(), isRequired)
			if model.ArrayProps == nil {
				model.ArrayProps = map[string]ArrayPropModel{}
			}
			model.ArrayProps[identifier] = *prop
		}
	}

	if model.StringProps == nil && model.NumberProps == nil && model.BooleanProps == nil &&
		model.ObjectProps == nil && model.ArrayProps == nil && !configured {
		return nil, nil
	}

	return model, nil
}

func commonPropToState(ctx context.Context, property cli.WorkflowInputProperty, common propCommon, required bool) {
	*common.Title = flex.GoStringToFramework(property.Title)
	*common.Icon = flex.GoStringToFramework(property.Icon)
	*common.Description = flex.GoStringToFramework(property.Description)
	*common.DependsOn = flex.GoArrayStringToTerraformList(ctx, property.DependsOn)

	if required {
		*common.Required = types.BoolValue(true)
	}

	*common.Visible, *common.VisibleJqQuery = boolOrJqToState(property.Visible)
	*common.ReadOnly, *common.ReadOnlyJqQuery = boolOrJqToState(property.ReadOnly)
	*common.Disabled, *common.DisabledJqQuery = boolOrJqToState(property.Disabled)
}

func boolOrJqToState(value any) (types.Bool, types.String) {
	switch value := value.(type) {
	case bool:
		return types.BoolValue(value), types.StringNull()
	case map[string]any:
		if jqQuery, ok := value["jqQuery"].(string); ok {
			return types.BoolNull(), types.StringValue(jqQuery)
		}
	case map[string]string:
		if jqQuery, ok := value["jqQuery"]; ok {
			return types.BoolNull(), types.StringValue(jqQuery)
		}
	}
	return types.BoolNull(), types.StringNull()
}

// jqQueryToState reports the jq query of a `{jqQuery}` envelope.
func jqQueryToState(value any) (string, bool) {
	switch value := value.(type) {
	case map[string]any:
		jqQuery, ok := value["jqQuery"].(string)
		return jqQuery, ok
	case map[string]string:
		jqQuery, ok := value["jqQuery"]
		return jqQuery, ok
	}
	return "", false
}

func requiredToState(required any) (types.String, []string) {
	switch required := required.(type) {
	case []any:
		values := make([]string, 0, len(required))
		for _, value := range required {
			if valueString, ok := value.(string); ok {
				values = append(values, valueString)
			}
		}
		return types.StringNull(), values
	case []string:
		return types.StringNull(), required
	case map[string]any:
		if jqQuery, ok := required["jqQuery"].(string); ok {
			return types.StringValue(jqQuery), nil
		}
	}
	return types.StringNull(), nil
}

func stringPropToState(ctx context.Context, property cli.WorkflowInputProperty, jsonEscapeHTML bool) (*StringPropModel, error) {
	prop := &StringPropModel{
		Format:     flex.GoStringToFramework(property.Format),
		Blueprint:  flex.GoStringToFramework(property.Blueprint),
		MinLength:  flex.GoInt64ToFramework(property.MinLength),
		MaxLength:  flex.GoInt64ToFramework(property.MaxLength),
		Default:    types.StringNull(),
		Pattern:    types.StringNull(),
		Enum:       types.ListNull(types.StringType),
		EnumColors: types.MapNull(types.StringType),
	}

	if jqQuery, ok := jqQueryToState(property.Default); ok {
		prop.DefaultJqQuery = types.StringValue(jqQuery)
	} else if value, ok := property.Default.(string); ok {
		prop.Default = types.StringValue(value)
	}

	if jqQuery, ok := jqQueryToState(property.Pattern); ok {
		prop.PatternJqQuery = types.StringValue(jqQuery)
	} else if value, ok := property.Pattern.(string); ok && value != "" {
		prop.Pattern = types.StringValue(value)
	}

	enum, enumJqQuery, err := enumToState(property.Enum, types.StringType, func(value any) (attr.Value, error) {
		valueString, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected a string enum value, got %T", value)
		}
		return types.StringValue(valueString), nil
	})
	if err != nil {
		return nil, err
	}
	prop.Enum, prop.EnumJqQuery = enum, enumJqQuery

	if property.EnumColors != nil {
		prop.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, property.EnumColors)
	}

	prop.Dataset = datasetToState(property.Dataset, jsonEscapeHTML)
	prop.Sort = sortToState(property.Sort)

	return prop, nil
}

func numberPropToState(_ context.Context, property cli.WorkflowInputProperty) (*NumberPropModel, error) {
	prop := &NumberPropModel{
		Default:          types.Float64Null(),
		Minimum:          flex.GoFloat64ToFramework(property.Minimum),
		Maximum:          flex.GoFloat64ToFramework(property.Maximum),
		ExclusiveMinimum: flex.GoFloat64ToFramework(property.ExclusiveMinimum),
		ExclusiveMaximum: flex.GoFloat64ToFramework(property.ExclusiveMaximum),
		Enum:             types.ListNull(types.Float64Type),
	}

	if jqQuery, ok := jqQueryToState(property.Default); ok {
		prop.DefaultJqQuery = types.StringValue(jqQuery)
	} else if value, ok := property.Default.(float64); ok {
		prop.Default = types.Float64Value(value)
	}

	enum, enumJqQuery, err := enumToState(property.Enum, types.Float64Type, func(value any) (attr.Value, error) {
		valueFloat, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("expected a number enum value, got %T", value)
		}
		return types.Float64Value(valueFloat), nil
	})
	if err != nil {
		return nil, err
	}
	prop.Enum, prop.EnumJqQuery = enum, enumJqQuery

	return prop, nil
}

func booleanPropToState(property cli.WorkflowInputProperty) *BooleanPropModel {
	prop := &BooleanPropModel{Default: types.BoolNull()}

	if jqQuery, ok := jqQueryToState(property.Default); ok {
		prop.DefaultJqQuery = types.StringValue(jqQuery)
	} else if value, ok := property.Default.(bool); ok {
		prop.Default = types.BoolValue(value)
	}

	return prop
}

func objectPropToState(property cli.WorkflowInputProperty, jsonEscapeHTML bool) (*ObjectPropModel, error) {
	prop := &ObjectPropModel{
		Format:  flex.GoStringToFramework(property.Format),
		Default: types.StringNull(),
	}

	if jqQuery, ok := jqQueryToState(property.Default); ok {
		prop.DefaultJqQuery = types.StringValue(jqQuery)
	} else if property.Default != nil {
		defaultValue, err := utils.GoObjectToTerraformString(property.Default, jsonEscapeHTML)
		if err != nil {
			return nil, fmt.Errorf("converting `default` to a string: %w", err)
		}
		prop.Default = defaultValue
	}

	return prop, nil
}

func arrayPropToState(ctx context.Context, property cli.WorkflowInputProperty, jsonEscapeHTML bool) (*ArrayPropModel, error) {
	prop := &ArrayPropModel{
		UniqueItems: flex.GoBoolToFramework(property.UniqueItems),
		Sort:        sortToState(property.Sort),
	}

	minItems, minItemsJqQuery, err := intOrJqToState(property.MinItems, "minItems")
	if err != nil {
		return nil, err
	}
	prop.MinItems, prop.MinItemsJqQuery = minItems, minItemsJqQuery

	maxItems, maxItemsJqQuery, err := intOrJqToState(property.MaxItems, "maxItems")
	if err != nil {
		return nil, err
	}
	prop.MaxItems, prop.MaxItemsJqQuery = maxItems, maxItemsJqQuery

	if jqQuery, ok := jqQueryToState(property.Default); ok {
		prop.DefaultJqQuery = types.StringValue(jqQuery)
	}

	if property.Items == nil {
		return prop, nil
	}

	// The default of an array lives on the property while the rest of the item
	// configuration lives on items, so it is only read when it is not a jq query.
	defaults, _ := property.Default.([]any)
	if !prop.DefaultJqQuery.IsNull() {
		defaults = nil
	}

	switch property.Items["type"] {
	case "string":
		items := &StringItemsModel{
			Default:    types.ListNull(types.StringType),
			Enum:       types.ListNull(types.StringType),
			EnumColors: types.MapNull(types.StringType),
		}

		if value, ok := property.Items["format"].(string); ok {
			items.Format = types.StringValue(value)
		}
		if value, ok := property.Items["blueprint"].(string); ok {
			items.Blueprint = types.StringValue(value)
		}
		if value, ok := property.Items["dataset"]; ok && value != nil {
			dataset, err := utils.GoObjectToTerraformString(value, jsonEscapeHTML)
			if err != nil {
				return nil, fmt.Errorf("converting `string_items.dataset` to a string: %w", err)
			}
			items.Dataset = dataset
		}

		enum, enumJqQuery, err := enumToState(property.Items["enum"], types.StringType, func(value any) (attr.Value, error) {
			valueString, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("expected a string enum value, got %T", value)
			}
			return types.StringValue(valueString), nil
		})
		if err != nil {
			return nil, err
		}
		items.Enum, items.EnumJqQuery = enum, enumJqQuery

		if colors, ok := property.Items["enumColors"].(map[string]any); ok {
			items.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, toStringMap(colors))
		}

		if defaults != nil {
			values := make([]attr.Value, 0, len(defaults))
			for _, value := range defaults {
				valueString, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string default item, got %T", value)
				}
				values = append(values, types.StringValue(valueString))
			}
			items.Default, _ = types.ListValue(types.StringType, values)
		}

		prop.StringItems = items

	case "number":
		items := &NumberItemsModel{
			Default:    types.ListNull(types.Float64Type),
			Enum:       types.ListNull(types.Float64Type),
			EnumColors: types.MapNull(types.StringType),
		}

		enum, enumJqQuery, err := enumToState(property.Items["enum"], types.Float64Type, func(value any) (attr.Value, error) {
			valueFloat, ok := value.(float64)
			if !ok {
				return nil, fmt.Errorf("expected a number enum value, got %T", value)
			}
			return types.Float64Value(valueFloat), nil
		})
		if err != nil {
			return nil, err
		}
		items.Enum, items.EnumJqQuery = enum, enumJqQuery

		if colors, ok := property.Items["enumColors"].(map[string]any); ok {
			items.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, toStringMap(colors))
		}

		if defaults != nil {
			values := make([]attr.Value, 0, len(defaults))
			for _, value := range defaults {
				valueFloat, ok := value.(float64)
				if !ok {
					return nil, fmt.Errorf("expected a number default item, got %T", value)
				}
				values = append(values, types.Float64Value(valueFloat))
			}
			items.Default, _ = types.ListValue(types.Float64Type, values)
		}

		prop.NumberItems = items

	case "object":
		items := &ObjectItemsModel{
			Default: types.ListNull(types.MapType{ElemType: types.StringType}),
		}

		if value, ok := property.Items["format"].(string); ok {
			items.Format = types.StringValue(value)
		}

		if defaults != nil {
			values := make([]attr.Value, 0, len(defaults))
			for _, value := range defaults {
				valueMap, ok := value.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("expected an object default item, got %T", value)
				}
				entries := map[string]attr.Value{}
				for key, entry := range valueMap {
					entries[key] = types.StringValue(fmt.Sprintf("%v", entry))
				}
				mapValue, _ := types.MapValue(types.StringType, entries)
				values = append(values, mapValue)
			}
			items.Default, _ = types.ListValue(types.MapType{ElemType: types.StringType}, values)
		}

		prop.ObjectItems = items
	}

	return prop, nil
}

func intOrJqToState(value any, name string) (types.Int64, types.String, error) {
	switch value := value.(type) {
	case nil:
		return types.Int64Null(), types.StringNull(), nil
	case float64:
		return types.Int64Value(int64(value)), types.StringNull(), nil
	case int:
		return types.Int64Value(int64(value)), types.StringNull(), nil
	case map[string]any:
		if jqQuery, ok := value["jqQuery"].(string); ok {
			return types.Int64Null(), types.StringValue(jqQuery), nil
		}
	}
	return types.Int64Null(), types.StringNull(), fmt.Errorf("`%s` must be a number or a jq query, got %T", name, value)
}

func enumToState(value any, elementType attr.Type, convert func(any) (attr.Value, error)) (types.List, types.String, error) {
	switch value := value.(type) {
	case nil:
		return types.ListNull(elementType), types.StringNull(), nil
	case []any:
		values := make([]attr.Value, 0, len(value))
		for _, entry := range value {
			converted, err := convert(entry)
			if err != nil {
				return types.ListNull(elementType), types.StringNull(), err
			}
			values = append(values, converted)
		}
		list, diags := types.ListValue(elementType, values)
		if diags.HasError() {
			return types.ListNull(elementType), types.StringNull(), fmt.Errorf("building the enum list: %v", diags.Errors())
		}
		return list, types.StringNull(), nil
	}

	if jqQuery, ok := jqQueryToState(value); ok {
		return types.ListNull(elementType), types.StringValue(jqQuery), nil
	}

	return types.ListNull(elementType), types.StringNull(), nil
}

func toStringMap(value map[string]any) map[string]string {
	result := make(map[string]string, len(value))
	for key, entry := range value {
		result[key] = fmt.Sprintf("%v", entry)
	}
	return result
}

func sortToState(sort *cli.EntitiesSortModel) *EntitiesSortModel {
	if sort == nil {
		return nil
	}
	return &EntitiesSortModel{
		Property: types.StringValue(sort.Property),
		Order:    types.StringValue(sort.Order),
	}
}

func datasetToState(dataset *cli.WorkflowDataset, jsonEscapeHTML bool) *DatasetModel {
	if dataset == nil {
		return nil
	}

	model := &DatasetModel{Combinator: types.StringValue(dataset.Combinator)}
	for _, rule := range dataset.Rules {
		model.Rules = append(model.Rules, datasetRuleToState(rule, jsonEscapeHTML))
	}
	return model
}

func datasetRuleToState(rule cli.WorkflowDatasetRule, jsonEscapeHTML bool) DatasetRuleModel {
	model := DatasetRuleModel{}

	if rule.Combinator != nil && *rule.Combinator != "" {
		model.Combinator = types.StringValue(*rule.Combinator)
		for _, nested := range rule.Rules {
			model.Rules = append(model.Rules, datasetRuleToState(nested, jsonEscapeHTML))
		}
		return model
	}

	model.Blueprint = flex.GoStringToFramework(rule.Blueprint)
	model.Property = flex.GoStringToFramework(rule.Property)
	model.Operator = flex.GoStringToFramework(&rule.Operator)

	if jqQuery, ok := jqQueryToState(rule.Value); ok {
		model.Value = &ValueModel{JqQuery: types.StringValue(jqQuery)}
	} else if rule.Value != nil {
		if value, err := utils.GoObjectToTerraformString(rule.Value, jsonEscapeHTML); err == nil {
			model.ValueJson = value
		}
	}

	return model
}
