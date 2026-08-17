package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CopyGenericMaps[T any](target map[string]T, source map[string]T) {
	for key, value := range source {
		target[key] = value
	}
}

func CopyMaps(target map[string]schema.Attribute, source map[string]schema.Attribute) {
	CopyGenericMaps(target, source)
}

func GenID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[len(id)-18:])
}

func TerraformListToGoArray(ctx context.Context, list types.List, arrayType string) ([]interface{}, error) {
	elems := []interface{}{}
	for _, elem := range list.Elements() {
		v, _ := elem.ToTerraformValue(ctx)
		switch arrayType {
		case "string":
			var stringValue string
			err := v.As(&stringValue)
			if err != nil {
				return nil, err
			}
			elems = append(elems, stringValue)
		case "float64":
			var keyValue big.Float
			err := v.As(&keyValue)
			if err != nil {
				return nil, err
			}
			floatValue, _ := keyValue.Float64()
			elems = append(elems, floatValue)

		case "bool":
			var boolValue bool
			err := v.As(&boolValue)
			if err != nil {
				return nil, err
			}
			elems = append(elems, boolValue)

		case "object":
			var stringValue string
			err := v.As(&stringValue)
			if err != nil {
				return nil, err
			}
			defaultObject := map[string]interface{}{}
			err = json.Unmarshal([]byte(stringValue), &defaultObject)
			if err != nil {
				return nil, err
			}
			elems = append(elems, defaultObject)
		}
	}
	return elems, nil

}

var nillableKinds = []reflect.Kind{reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice}

func GoObjectToTerraformString(v interface{}, jsonEscapeHTML bool) (types.String, error) {
	if v == nil {
		return types.StringNull(), nil
	}

	isNillable := slices.Contains(nillableKinds, reflect.TypeOf(v).Kind())
	if isNillable && reflect.ValueOf(v).IsNil() {
		return types.StringNull(), nil
	}

	jsonBuilder := new(strings.Builder)
	jsonEncoder := json.NewEncoder(jsonBuilder)
	jsonEncoder.SetEscapeHTML(jsonEscapeHTML)
	err := jsonEncoder.Encode(v)
	if err != nil {
		return types.StringNull(), err
	}

	jsonStr, _ := strings.CutSuffix(jsonBuilder.String(), "\n")
	return types.StringValue(jsonStr), nil
}

func normalizeJSONString(s string, jsonEscapeHTML bool) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return "", err
	}
	encoded, err := GoObjectToTerraformString(v, jsonEscapeHTML)
	if err != nil {
		return "", err
	}
	return encoded.ValueString(), nil
}

func JSONStringsSemanticallyEqual(a, b string, jsonEscapeHTML bool) (bool, error) {
	if a == b {
		return true, nil
	}
	normalizedA, err := normalizeJSONString(a, jsonEscapeHTML)
	if err != nil {
		return false, err
	}
	normalizedB, err := normalizeJSONString(b, jsonEscapeHTML)
	if err != nil {
		return false, err
	}
	return normalizedA == normalizedB, nil
}

func GoObjectToTerraformStringPreferExisting(preferred types.String, v interface{}, jsonEscapeHTML bool) (types.String, error) {
	encoded, err := GoObjectToTerraformString(v, jsonEscapeHTML)
	if err != nil {
		return types.StringNull(), err
	}

	if preferred.IsNull() || preferred.IsUnknown() {
		return encoded, nil
	}

	equal, err := JSONStringsSemanticallyEqual(preferred.ValueString(), encoded.ValueString(), jsonEscapeHTML)
	if err != nil || !equal {
		return encoded, nil
	}

	return preferred, nil
}

func TerraformStringAt[T attr.Value](elements []T, i int) types.String {
	if i < 0 || i >= len(elements) {
		return types.StringNull()
	}
	s, ok := any(elements[i]).(types.String)
	if !ok {
		return types.StringNull()
	}
	return s
}

func TerraformStringAtList(list types.List, i int) types.String {
	if list.IsNull() || list.IsUnknown() {
		return types.StringNull()
	}
	return TerraformStringAt(list.Elements(), i)
}

func TerraformStringToGoType[T any](s types.String) (T, error) {
	var obj T

	if s.IsNull() {
		return obj, nil
	}

	if err := json.Unmarshal([]byte(s.ValueString()), &obj); err != nil {
		return obj, err
	}

	return obj, nil
}

func TerraformJsonStringToGoObject(v *string) (*map[string]any, error) {
	if v == nil || *v == "" {
		return nil, nil
	}

	vMap := make(map[string]any)
	if err := json.Unmarshal([]byte(*v), &vMap); err != nil {
		return nil, err
	}

	return &vMap, nil
}

func InterfaceToStringArray(o interface{}) []string {
	items := o.([]interface{})
	res := make([]string, len(items))
	for i, item := range items {
		res[i] = item.(string)
	}

	return res
}

func TFStringListToStringArray(list []types.String) []string {
	res := make([]string, len(list))
	for i, item := range list {
		res[i] = item.ValueString()
	}

	return res
}

func TerraformStringToBooleanOrString(s types.String) interface{} {
	var obj interface{}

	if s.IsNull() {
		return obj
	}

	value := s.ValueString()
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	return value
}
