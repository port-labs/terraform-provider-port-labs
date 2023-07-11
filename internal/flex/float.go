package flex

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoFloat64ToFramework(v *float64) types.Float64 {
	if v == nil {
		return types.Float64Null()
	}

	return types.Float64Value(*v)
}
