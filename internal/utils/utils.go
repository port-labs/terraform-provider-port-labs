package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SpreadMaps(target map[string]schema.Attribute, source map[string]schema.Attribute) {
	for key, value := range source {
		target[key] = value
	}
}

func GenID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[:18])
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
		}
	}
	return elems, nil

}
