package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"

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
	plan.Spec = planSpec(plan.Spec, state.Spec)
	if plan.Config.IsNull() || plan.Config.IsUnknown() {
		plan.Config = state.Config
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// planSpec overlays the spec sections the configuration declares onto the ones
// recorded in state. Sections the configuration leaves out — appSpec, which
// Port fills in during provisioning — are inherited, so editing integrationSpec
// alone neither shows up as clearing them nor sends a spec that drops them.
func planSpec(config, state types.String) types.String {
	configSections, stateSections := specSections(config), specSections(state)
	if configSections == nil || stateSections == nil {
		return config
	}

	planned := maps.Clone(stateSections)
	maps.Copy(planned, configSections)
	if maps.EqualFunc(planned, stateSections, func(a, b json.RawMessage) bool { return bytes.Equal(a, b) }) {
		return state
	}

	encoded, err := json.Marshal(planned)
	if err != nil {
		return config
	}
	return types.StringValue(string(encoded))
}

func specSections(spec types.String) map[string]json.RawMessage {
	if spec.IsNull() || spec.IsUnknown() {
		return nil
	}
	var sections map[string]json.RawMessage
	if err := json.Unmarshal([]byte(spec.ValueString()), &sections); err != nil {
		return nil
	}
	return sections
}
