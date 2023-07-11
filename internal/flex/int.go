package flex

import "github.com/hashicorp/terraform-plugin-framework/types"

func GoInt64ToFramework(v *int) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*v))
}
