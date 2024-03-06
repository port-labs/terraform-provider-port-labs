package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CopyMaps(target map[string]schema.Attribute, source map[string]schema.Attribute) {
	for key, value := range source {
		target[key] = value
	}
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

func GoObjectToTerraformString(v interface{}) (types.String, error) {
	if v == nil {
		return types.StringNull(), nil
	}
	js, err := json.Marshal(v)
	if err != nil {
		return types.StringNull(), err
	}

	value := string(js)
	return types.StringValue(value), nil
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
