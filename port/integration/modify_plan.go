package integration

import (
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

	// Suppress spec diffs caused only by server-managed fields (appSpec).
	// If the user's integrationSpec hasn't changed, keep the state value.
	plan.Spec = suppressServerOnlySpecDiff(state.Spec, plan.Spec)

	// Suppress config diffs when user didn't declare config (null in HCL).
	if plan.Config.IsNull() || plan.Config.IsUnknown() {
		plan.Config = state.Config
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// suppressServerOnlySpecDiff keeps the state spec when the user's integrationSpec
// hasn't changed — preventing diffs caused by server-added appSpec.
func suppressServerOnlySpecDiff(stateSpec, planSpec types.String) types.String {
	if stateSpec.IsNull() || planSpec.IsNull() || planSpec.IsUnknown() {
		return planSpec
	}

	var stateMap, planMap map[string]json.RawMessage
	if err := json.Unmarshal([]byte(stateSpec.ValueString()), &stateMap); err != nil {
		return planSpec
	}
	if err := json.Unmarshal([]byte(planSpec.ValueString()), &planMap); err != nil {
		return planSpec
	}

	// Compare only integrationSpec — if it matches, the diff is server-only.
	stateIS, _ := json.Marshal(stateMap["integrationSpec"])
	planIS, _ := json.Marshal(planMap["integrationSpec"])
	if string(stateIS) == string(planIS) {
		return stateSpec
	}

	return planSpec
}
