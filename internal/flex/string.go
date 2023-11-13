package flex

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GoArrayStringToTerraformList(array []string) types.List {
	if array == nil {
		return types.ListNull(types.StringType)
	}
	attrs := make([]attr.Value, 0, len(array))
	for _, value := range array {
		attrs = append(attrs, basetypes.NewStringValue(value))
	}
	list, _ := types.ListValue(types.StringType, attrs)
	return list
}

func TerraformStringListToGoArray(list []types.String) []string {
	arr := make([]string, len(list))
	for i, t := range list {
		arr[i] = t.ValueString()
	}
	return arr
}
