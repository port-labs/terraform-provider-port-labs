package integration

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// ──────────────────────────────────────────────────────────────────────────────
// WHY THIS IS COMPLICATED — the "spec / config" consistency problem
// ──────────────────────────────────────────────────────────────────────────────
//
// Terraform enforces a strict contract: the value a provider returns after
// Create/Update MUST exactly match the planned value for every attribute that
// was known in the plan. If it doesn't → "Provider produced inconsistent result
// after apply".
//
// The Port API adds server-managed fields after create/update:
//   - spec: server adds "appSpec" (liveEvents config, etc.)
//   - config: server populates default resource mappings
//
// These fields aren't in the user's HCL, so they aren't in the plan.
// If we naively return them from Create/Update, Terraform errors.
//
// The solution uses two mechanisms:
//
// 1. refreshIntegrationState (this file):
//    Called from Create/Update/Read. Takes a `planSpec` parameter:
//      - planSpec != nil (Create/Update): only include appSpec if it was
//        already in the plan (e.g. from a previous Read → ModifyPlan cycle).
//      - planSpec == nil (Read): always include appSpec from server.
//    This ensures Create/Update return values matching the plan exactly.
//
// 2. ModifyPlan (modify_plan.go):
//    Runs before every plan. On updates (when state exists):
//      - spec: if user's integrationSpec hasn't changed, keeps the state
//        value (which includes appSpec from last Read). This prevents a
//        phantom diff where Terraform wants to remove appSpec.
//      - config: if user didn't declare config, keeps state value.
//
// The full lifecycle for a SaaS integration:
//
//   terraform apply (Create):
//     Plan:   spec = {"integrationSpec": {...}}           ← from HCL
//     Apply:  server adds appSpec, but we exclude it      ← planSpec match
//     State:  spec = {"integrationSpec": {...}}           ✓ matches plan
//
//   terraform plan (Read + ModifyPlan):
//     Read:   spec = {"appSpec": {...}, "integrationSpec": {...}}  ← full server
//     State updated with appSpec
//     ModifyPlan: integrationSpec unchanged → keep state  ← suppress diff
//     Plan output: "No changes"                           ✓ no phantom diff
//
//   terraform apply (Update, user changes integrationSpec):
//     Plan:   spec = {"appSpec": {...}, "integrationSpec": {NEW}}  ← from ModifyPlan
//     Apply:  planSpec has appSpec → we include it         ← planSpec match
//     State:  spec = {"appSpec": {...}, "integrationSpec": {NEW}} ✓ matches plan
// ──────────────────────────────────────────────────────────────────────────────

// refreshIntegrationState updates Terraform state from the server response.
// Always merges all server fields (including appSpec and config).
// On Create/Update, the caller (resource.go) restores planned spec/config
// values afterward to satisfy Terraform's consistency contract.
func (r *IntegrationResource) refreshIntegrationState(state *IntegrationModel, a *cli.Integration, integrationId string) error {
	state.ID = types.StringValue(integrationId)
	state.InstallationId = types.StringValue(integrationId)
	state.Title = types.StringPointerValue(a.Title)
	state.InstallationAppType = types.StringPointerValue(a.InstallationAppType)
	state.InstallationType = types.StringPointerValue(a.InstallationType)
	state.Version = types.StringPointerValue(a.Version)

	if a.StatusInfo != nil {
		state.Status = types.StringValue(a.StatusInfo.IntegrationStatus.Status)
	} else {
		state.Status = types.StringNull()
	}

	if a.Spec != nil {
		state.Spec = mergeSpec(state.Spec, a.Spec, r.portClient.JSONEscapeHTML)
	}

	if a.Config != nil {
		config, _ := utils.GoObjectToTerraformStringPreferExisting(state.Config, a.Config, r.portClient.JSONEscapeHTML)
		state.Config = config
	}

	// Terraform requires all values to be known after apply. Computed+Optional
	// attributes start as "unknown" on Create when the user didn't set them.
	// Resolve to null so Terraform doesn't error with "still indicated an
	// unknown value".
	if state.Spec.IsUnknown() {
		state.Spec = types.StringNull()
	}
	if state.Config.IsUnknown() {
		state.Config = types.StringNull()
	}

	if a.ChangelogDestination != nil {
		if a.ChangelogDestination.Type == consts.Kafka {
			state.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
			state.WebhookChangelogDestination = nil
		} else {
			if a.ChangelogDestination.Url != "" {
				state.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
					Url: types.StringValue(a.ChangelogDestination.Url),
				}
				if a.ChangelogDestination.Agent != nil {
					state.WebhookChangelogDestination.Agent = types.BoolValue(*a.ChangelogDestination.Agent)
				}
				state.KafkaChangelogDestination = types.ObjectNull(map[string]attr.Type{})
			}
		}
	} else {
		state.KafkaChangelogDestination = types.ObjectNull(map[string]attr.Type{})
		state.WebhookChangelogDestination = nil
	}

	return nil
}

// mergeSpec builds the refreshed spec JSON from the server response.
// Always includes all server fields (integrationSpec, appSpec).
// systemSpec / privateSpec are excluded (stripped by CLI layer).
// Preserves user-set secrets that the server strips (org secret references).
func mergeSpec(stateTF types.String, remote *cli.IntegrationClientSpec, jsonEscapeHTML bool) types.String {
	var userSpec map[string]map[string]any
	if !stateTF.IsNull() && !stateTF.IsUnknown() {
		_ = json.Unmarshal([]byte(stateTF.ValueString()), &userSpec)
	}

	merged := make(map[string]map[string]any)

	if remote.IntegrationSpec != nil {
		is := make(map[string]any, len(remote.IntegrationSpec))
		for k, v := range remote.IntegrationSpec {
			is[k] = v
		}
		if userIS := userSpec["integrationSpec"]; userIS != nil {
			for k, userVal := range userIS {
				serverVal, exists := is[k]
				if (!exists || serverVal == nil || serverVal == "") && userVal != nil && userVal != "" {
					is[k] = userVal
				}
			}
		}
		merged["integrationSpec"] = is
	} else if userIS := userSpec["integrationSpec"]; userIS != nil {
		merged["integrationSpec"] = userIS
	}

	if remote.AppSpec != nil {
		merged["appSpec"] = remote.AppSpec
	}

	encoded, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return stateTF
	}
	return encoded
}
