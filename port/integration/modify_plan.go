package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithModifyPlan = &IntegrationResource{}

type immutableField struct {
	name    string
	changed func(state, plan *IntegrationModel) bool
}

var immutableFields = []immutableField{
	{
		name:    "installation_id",
		changed: func(s, p *IntegrationModel) bool { return !p.InstallationId.Equal(s.InstallationId) },
	},
	{
		name:    "installation_app_type",
		changed: func(s, p *IntegrationModel) bool { return !p.InstallationAppType.Equal(s.InstallationAppType) },
	},
	{
		name:    "installation_type",
		changed: func(s, p *IntegrationModel) bool { return !p.InstallationType.Equal(s.InstallationType) },
	},
}

func (r *IntegrationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state, plan IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, f := range immutableFields {
		if f.changed(&state, &plan) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("cannot change %s", f.name),
				fmt.Sprintf(
					"The Port API does not support changing `%s` on an existing integration. To use a different value, destroy this resource (which deletes the integration from Port) and create a new `port_integration`.",
					f.name,
				),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Read pulls Port's own additions into state (spec.appSpec, default config
	// mappings). Neither belongs in a diff unless the configuration itself
	// changed, so fall back to state when only Port moved.
	plan.Spec = specPlan(state.Spec, plan.Spec)
	if plan.Config.IsNull() || plan.Config.IsUnknown() {
		plan.Config = state.Config
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// specPlan keeps the spec recorded in state whenever the configured
// integrationSpec is unchanged, so Port's appSpec additions do not read as a
// removal the next time Terraform plans.
func specPlan(state, plan types.String) types.String {
	if state.IsNull() || plan.IsNull() || plan.IsUnknown() {
		return plan
	}

	stateSections, planSections := specSections(state), specSections(plan)
	if stateSections == nil || planSections == nil {
		return plan
	}
	if bytes.Equal(stateSections["integrationSpec"], planSections["integrationSpec"]) {
		return state
	}
	return plan
}

func specSections(v types.String) map[string]json.RawMessage {
	var sections map[string]json.RawMessage
	if err := json.Unmarshal([]byte(v.ValueString()), &sections); err != nil {
		return nil
	}
	return sections
}
