package scorecard_group

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var scorecardEntityReservedPropertyKeys = map[string]struct{}{
	"blueprint": {},
	"levels":    {},
	"filter":    {},
}

func validateScorecardGroupProperties(properties map[string]any) error {
	for key := range properties {
		if _, reserved := scorecardEntityReservedPropertyKeys[key]; reserved {
			return fmt.Errorf("property %q is reserved and cannot be set in `properties`", key)
		}
	}
	return nil
}

func propertiesToCLI(ctx context.Context, properties types.Dynamic) (map[string]any, error) {
	if properties.IsNull() || properties.IsUnknown() {
		return nil, nil
	}

	underlying := properties.UnderlyingValue()
	mapValue, ok := underlying.(types.Map)
	if !ok {
		return nil, fmt.Errorf("properties must be a map, got %T", underlying)
	}

	result := make(map[string]any, len(mapValue.Elements()))
	for key, value := range mapValue.Elements() {
		goValue, err := attrValueToGoValue(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert property %q: %w", key, err)
		}
		result[key] = goValue
	}

	if err := validateScorecardGroupProperties(result); err != nil {
		return nil, err
	}

	return result, nil
}

func propertiesFromCLI(ctx context.Context, properties map[string]any) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(properties) == 0 {
		return types.DynamicNull(), diags
	}

	elements := make(map[string]attr.Value, len(properties))
	for key, value := range properties {
		attrValue, err := goValueToAttrValue(ctx, value)
		if err != nil {
			diags.AddError("failed to convert scorecard group properties", fmt.Sprintf("property %q: %s", key, err))
			return types.DynamicUnknown(), diags
		}
		elements[key] = attrValue
	}

	mapValue, d := types.MapValue(types.DynamicType, elements)
	diags.Append(d...)
	if diags.HasError() {
		return types.DynamicUnknown(), diags
	}

	return types.DynamicValue(mapValue), diags
}

func attrValueToGoValue(ctx context.Context, value attr.Value) (any, error) {
	switch v := value.(type) {
	case types.Dynamic:
		if v.IsNull() || v.IsUnknown() {
			return nil, nil
		}
		return attrValueToGoValue(ctx, v.UnderlyingValue())
	case types.String:
		if v.IsNull() {
			return nil, nil
		}
		return v.ValueString(), nil
	case types.Bool:
		if v.IsNull() {
			return nil, nil
		}
		return v.ValueBool(), nil
	case types.Int64:
		if v.IsNull() {
			return nil, nil
		}
		return v.ValueInt64(), nil
	case types.Float64:
		if v.IsNull() {
			return nil, nil
		}
		return v.ValueFloat64(), nil
	case types.Number:
		if v.IsNull() {
			return nil, nil
		}
		f, _ := v.ValueBigFloat().Float64()
		return f, nil
	case types.List:
		if v.IsNull() {
			return nil, nil
		}
		items := make([]any, 0, len(v.Elements()))
		for _, element := range v.Elements() {
			item, err := attrValueToGoValue(ctx, element)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	case types.Map:
		if v.IsNull() {
			return nil, nil
		}
		items := make(map[string]any, len(v.Elements()))
		for key, element := range v.Elements() {
			item, err := attrValueToGoValue(ctx, element)
			if err != nil {
				return nil, err
			}
			items[key] = item
		}
		return items, nil
	case types.Object:
		if v.IsNull() {
			return nil, nil
		}
		items := make(map[string]any, len(v.Attributes()))
		for key, element := range v.Attributes() {
			item, err := attrValueToGoValue(ctx, element)
			if err != nil {
				return nil, err
			}
			items[key] = item
		}
		return items, nil
	default:
		return nil, fmt.Errorf("unsupported property value type %T", value)
	}
}

func goValueToAttrValue(ctx context.Context, value any) (attr.Value, error) {
	underlying, err := goValueToUnderlyingAttrValue(ctx, value)
	if err != nil {
		return nil, err
	}
	return types.DynamicValue(underlying), nil
}

func goValueToUnderlyingAttrValue(ctx context.Context, value any) (attr.Value, error) {
	switch val := value.(type) {
	case nil:
		return types.StringNull(), nil
	case string:
		return types.StringValue(val), nil
	case bool:
		return types.BoolValue(val), nil
	case float64:
		return types.Float64Value(val), nil
	case float32:
		return types.Float64Value(float64(val)), nil
	case int:
		return types.Int64Value(int64(val)), nil
	case int64:
		return types.Int64Value(val), nil
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return types.Int64Value(i), nil
		}
		f, err := val.Float64()
		if err != nil {
			return nil, err
		}
		return types.Float64Value(f), nil
	case []any:
		elements := make([]attr.Value, len(val))
		for i, item := range val {
			element, err := goValueToAttrValue(ctx, item)
			if err != nil {
				return nil, err
			}
			elements[i] = element
		}
		listValue, diags := types.ListValue(types.DynamicType, elements)
		if diags.HasError() {
			return nil, fmt.Errorf("%s", diags.Errors()[0].Summary())
		}
		return listValue, nil
	case map[string]any:
		elements := make(map[string]attr.Value, len(val))
		for key, item := range val {
			element, err := goValueToAttrValue(ctx, item)
			if err != nil {
				return nil, err
			}
			elements[key] = element
		}
		mapValue, diags := types.MapValue(types.DynamicType, elements)
		if diags.HasError() {
			return nil, fmt.Errorf("%s", diags.Errors()[0].Summary())
		}
		return mapValue, nil
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("unsupported value type %T: %w", value, err)
		}
		var decoded any
		if err := json.Unmarshal(bytes, &decoded); err != nil {
			return nil, err
		}
		return goValueToUnderlyingAttrValue(ctx, decoded)
	}
}
