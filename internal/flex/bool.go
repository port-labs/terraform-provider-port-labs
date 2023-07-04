package flex

import "github.com/hashicorp/terraform-plugin-framework/types"

func GoBoolToFramework(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*v)
}
