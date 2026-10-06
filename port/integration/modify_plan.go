package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
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
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan IntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config IntegrationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	specConfiguredInHCL := specIsConfigured(config.Spec)

	isCreate := req.State.Raw.IsNull()

	if isCreate {
		validatePlanRules(&plan, nil, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		r.validateSpecAtPlan(ctx, &plan, specConfiguredInHCL, &resp.Diagnostics)
		return
	}

	var state IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	validatePlanRules(&plan, &state, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read pulls Port's own additions into state (spec.appSpec, default config
	// mappings). Neither belongs in a diff unless the configuration itself
	// changed, so fall back to state when only Port moved.
	plan.Spec = planSpec(plan.Spec, state.Spec)
	if plan.isSaasOAuth2() {
		plan.Spec = extractAppSpec(plan.Spec)
	}
	if plan.Config.IsNull() || plan.Config.IsUnknown() {
		plan.Config = state.Config
	}

	r.validateSpecAtPlan(ctx, &plan, specConfiguredInHCL, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// planSpec overlays the spec sections the configuration declares onto the ones
// recorded in state. Sections the configuration leaves out inherit Port defaults
// from state so partial edits do not clear undeclared keys on apply.
func planSpec(config, state types.String) types.String {
	configSections, stateSections := specSections(config), specSections(state)
	if configSections == nil || stateSections == nil {
		return config
	}

	planned := maps.Clone(stateSections)
	maps.Copy(planned, configSections)

	if _, configDeclaresIntegrationSpec := configSections["integrationSpec"]; configDeclaresIntegrationSpec {
		planned["integrationSpec"] = configSections["integrationSpec"]
	}
	if _, configDeclaresAppSpec := configSections["appSpec"]; configDeclaresAppSpec {
		planned["appSpec"] = configSections["appSpec"]
	} else if inherited, ok := planned["appSpec"]; ok {
		if filtered := filterServerManagedAppSpecJSON(inherited); len(filtered) > 0 {
			planned["appSpec"] = filtered
		} else {
			delete(planned, "appSpec")
		}
	}

	if maps.EqualFunc(planned, stateSections, func(a, b json.RawMessage) bool { return bytes.Equal(a, b) }) {
		return state
	}

	encoded, err := json.Marshal(planned)
	if err != nil {
		return config
	}
	return types.StringValue(string(encoded))
}

func filterServerManagedAppSpecJSON(raw json.RawMessage) json.RawMessage {
	var appSpec map[string]any
	if err := json.Unmarshal(raw, &appSpec); err != nil {
		return raw
	}
	filtered := mergeSpecSectionImplicit(appSpec, consts.IsServerManagedAppSpecKey)
	if len(filtered) == 0 {
		return nil
	}
	out, err := json.Marshal(filtered)
	if err != nil {
		return raw
	}
	return out
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

func extractAppSpec(spec types.String) types.String {
	if spec.IsUnknown() {
		return spec
	}
	app := priorSpecSections(spec)["appSpec"]
	if app == nil {
		return types.StringNull()
	}
	out, err := json.Marshal(map[string]any{"appSpec": app})
	if err != nil {
		return types.StringNull()
	}
	return types.StringValue(string(out))
}
