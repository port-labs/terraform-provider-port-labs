package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.String = (*jsonIgnoreDiff)(nil)
var _ planmodifier.List = (*jsonIgnoreDiff)(nil)

func JsonIgnoreDiffPlanModifier() planmodifier.String {
	return jsonIgnoreDiff{}
}

func JsonIgnoreDiffPlanModifierList() planmodifier.List {
	return jsonIgnoreDiff{}
}

type jsonIgnoreDiff struct {
}

// Description implements planmodifier.String
func (jsonIgnoreDiff) Description(context.Context) string {
	return "Compares json for object equality to ignore formatting changes"
}

// MarkdownDescription implements planmodifier.String
func (j jsonIgnoreDiff) MarkdownDescription(ctx context.Context) string {
	return j.Description(ctx)
}

// PlanModifyString implements planmodifier.String
func (j jsonIgnoreDiff) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
}

// PlanModifyList implements planmodifier.List
func (j jsonIgnoreDiff) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() {
		resp.PlanValue = req.PlanValue

		return
	}

	// If the current value is semantically equivalent to the planned value
	// then return the current value, else return the planned value.

	planned, diags := req.PlanValue.ToListValue(ctx)

	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)

		return
	}

	current, diags := req.StateValue.ToListValue(ctx)

	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)

		return
	}

	if len(planned.Elements()) != len(current.Elements()) {
		resp.PlanValue = req.PlanValue

		return
	}

	currentVals := make([]attr.Value, len(current.Elements()))
	copy(currentVals, current.Elements())

	for _, plannedVal := range planned.Elements() {
		for i, currentVal := range currentVals {
			if currentVal.Equal(plannedVal) {
				// Remove from the slice.
				currentVals = append(currentVals[:i], currentVals[i+1:]...)

				break
			}
		}
	}

	if len(currentVals) == 0 {
		// Every planned value is equal to a current value.
		resp.PlanValue = req.StateValue
	} else {
		resp.PlanValue = req.PlanValue
	}
}
