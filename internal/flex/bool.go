package flex

import "github.com/hashicorp/terraform-plugin-framework/types"

func GoBoolToFramework(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*v)
}

func GoBoolToFrameworkDefaultFalse(v *bool) types.Bool {
	if v == nil {
		return types.BoolValue(false)
	}

	return types.BoolValue(*v)
}
