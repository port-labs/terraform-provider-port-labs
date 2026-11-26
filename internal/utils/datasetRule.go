package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// PropertyToState converts the `property` of an API dataset rule, which is
// either an identifier (string) or an object for contextual and path rules,
// into the `property` and `property_json` attributes.
//
// The prior `property_json` is kept when it is semantically equal, so
// formatting and HTML escaping never cause drift.
func PropertyToState(prior types.String, property any, jsonEscapeHTML bool) (identifier, object types.String, err error) {
	switch property := property.(type) {
	case nil:
		return types.StringNull(), types.StringNull(), nil
	case string:
		return types.StringValue(property), types.StringNull(), nil
	default:
		object, err = GoObjectToTerraformStringPreferExisting(prior, property, jsonEscapeHTML)
		return types.StringNull(), object, err
	}
}

// PriorAt returns the element of a prior list at index i, or nil when the
// prior list is shorter or absent.
func PriorAt[T any](prior []T, i int) *T {
	if i < 0 || i >= len(prior) {
		return nil
	}
	return &prior[i]
}
