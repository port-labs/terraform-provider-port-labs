// Package planmodifiers holds plan modifiers shared by several resources.
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseStateForUnknownUnlessChanged is UseStateForUnknown for computed values that are stable across updates except when
// one of the given string attributes changes (e.g. an id derived from identifier/blueprint, or a created_at that is
// reset when the object is recreated under a new blueprint). When every trigger attribute is known and equal to its
// prior state the prior value is kept, otherwise the value stays unknown so Terraform accepts the new one.
func UseStateForUnknownUnlessChanged(triggers ...path.Path) planmodifier.String {
	return useStateForUnknownUnlessChanged{triggers: triggers}
}

type useStateForUnknownUnlessChanged struct {
	triggers []path.Path
}

func (m useStateForUnknownUnlessChanged) Description(_ context.Context) string {
	return fmt.Sprintf("Keeps the prior state value unless any of %v changes.", m.triggers)
}

func (m useStateForUnknownUnlessChanged) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useStateForUnknownUnlessChanged) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Create (or replacement, which Terraform plans as a create with null prior state), destroy, a value that is
	// already known, or a value that comes from unknown config: nothing to keep.
	if req.StateValue.IsNull() || req.Plan.Raw.IsNull() || !req.PlanValue.IsUnknown() || req.ConfigValue.IsUnknown() {
		return
	}

	for _, p := range m.triggers {
		var planned, prior types.String
		if diags := req.Plan.GetAttribute(ctx, p, &planned); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		if diags := req.State.GetAttribute(ctx, p, &prior); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		if planned.IsUnknown() || !planned.Equal(prior) {
			return
		}
	}

	resp.PlanValue = req.StateValue
}
