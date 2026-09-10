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

// FrameworkBoolToTruePointer returns a pointer only when the Terraform bool is true.
// Use when serializing optional API flags that should be omitted unless enabled.
func FrameworkBoolToTruePointer(v types.Bool) *bool {
	if v.IsNull() || !v.ValueBool() {
		return nil
	}

	trueValue := true
	return &trueValue
}
