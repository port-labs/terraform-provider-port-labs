package flex

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GoStringToFramework(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}

	return types.StringValue(*v)
}

func GoArrayStringToTerraformList(ctx context.Context, array []string) types.List {
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
