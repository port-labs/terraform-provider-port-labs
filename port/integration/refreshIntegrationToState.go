package integration

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// Port enriches an integration after every write: it adds spec.appSpec and
// fills config with the integration's default mappings. Terraform requires
// Create and Update to return exactly the planned value, so those two
// attributes are refreshed from Port on Read only, and ModifyPlan keeps the
// resulting server-owned additions out of the next diff.

// applyServerFields copies the attributes Port owns onto the model, leaving
// spec and config untouched.
func applyServerFields(m *IntegrationModel, a *cli.Integration, integrationId string) {
	m.ID = types.StringValue(integrationId)
	m.InstallationId = types.StringValue(integrationId)
	m.Title = types.StringPointerValue(a.Title)
	m.InstallationAppType = types.StringPointerValue(a.InstallationAppType)
	m.InstallationType = types.StringPointerValue(a.InstallationType)
	m.Version = types.StringPointerValue(a.Version)

	if a.StatusInfo != nil {
		m.Status = types.StringValue(a.StatusInfo.IntegrationStatus.Status)
	} else {
		m.Status = types.StringNull()
	}

	m.KafkaChangelogDestination = types.ObjectNull(kafkaChangelogDestinationType)
	m.WebhookChangelogDestination = types.ObjectNull(webhookChangelogDestinationType)

	switch dest := a.ChangelogDestination; {
	case dest == nil:
	case dest.Type == consts.Kafka:
		m.KafkaChangelogDestination = types.ObjectValueMust(kafkaChangelogDestinationType, map[string]attr.Value{})
	case dest.Type == consts.Webhook && dest.Url != "":
		m.WebhookChangelogDestination = types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
			"url":   types.StringValue(dest.Url),
			"agent": types.BoolPointerValue(dest.Agent),
		})
	}
}

// refreshIntegrationState syncs the full server view onto state. Used by Read,
// where picking up Port's own additions is what surfaces drift.
func (r *IntegrationResource) refreshIntegrationState(state *IntegrationModel, a *cli.Integration, integrationId string) {
	applyServerFields(state, a, integrationId)

	prior, remote := state.Spec, a.Spec
	if state.isSaasOAuth2() {
		prior, remote = extractAppSpec(prior), remote.AppSpecOnly()
	}
	if !remote.IsEmpty() {
		state.Spec = mergeSpec(prior, remote, r.portClient.JSONEscapeHTML)
	}
	if a.Config != nil {
		state.Config, _ = utils.GoObjectToTerraformStringPreferExisting(state.Config, a.Config, r.portClient.JSONEscapeHTML)
	}
}

// applyWriteResult stores the plan. The plan is already the pending change
// computed from state and config. The create/update response is an in-flight
// operation and omits configured fields, so it must not be copied onto the plan.
// Unknown computed values become null; Read fills server-owned values later.
func applyWriteResult(plan *IntegrationModel, installationId string) {
	plan.ID = types.StringValue(installationId)
	plan.InstallationAppType = nullIfUnknown(plan.InstallationAppType)
	plan.InstallationType = nullIfUnknown(plan.InstallationType)
	plan.Title = nullIfUnknown(plan.Title)
	plan.Version = nullIfUnknown(plan.Version)
	plan.Config = nullIfUnknown(plan.Config)
	plan.Spec = nullIfUnknown(plan.Spec)
	plan.CreatePortResourcesOrigin = nullIfUnknown(plan.CreatePortResourcesOrigin)
	plan.Status = nullIfUnknown(plan.Status)
	plan.KafkaChangelogDestination = nullObjectIfUnknown(plan.KafkaChangelogDestination, kafkaChangelogDestinationType)
	plan.WebhookChangelogDestination = nullObjectIfUnknown(plan.WebhookChangelogDestination, webhookChangelogDestinationType)
}

func nullIfUnknown(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

func nullObjectIfUnknown(v types.Object, attrTypes map[string]attr.Type) types.Object {
	if v.IsUnknown() {
		return types.ObjectNull(attrTypes)
	}
	return v
}

// mergeSpec renders Port's spec as JSON, restoring values that Port strips or
// overwrites. Explicit spec sections keep only keys declared in HCL; omitted
// sections inherit Port defaults on Read.
func mergeSpec(state types.String, remote *cli.IntegrationClientSpec, jsonEscapeHTML bool) types.String {
	prior := priorSpecSections(state)

	merged := make(map[string]any, 2)
	mergeSpecSectionInto(merged, "integrationSpec", remote.IntegrationSpec, prior["integrationSpec"], nil)
	mergeSpecSectionInto(merged, "appSpec", remote.AppSpec, prior["appSpec"], consts.IsServerManagedAppSpecKey)

	encoded, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return state
	}
	return encoded
}

// mergeSpecSectionInto merges one spec section on Read. Explicit sections keep
// only keys declared in HCL; omitted sections inherit Port defaults. skipKey
// drops server-managed fields (appSpec only today).
func mergeSpecSectionInto(merged map[string]any, key string, remote, prior map[string]any, skipKey func(string) bool) {
	switch {
	case remote != nil && len(prior) > 0:
		merged[key] = mergeSpecSectionExplicit(remote, prior, skipKey)
	case remote != nil:
		merged[key] = mergeSpecSectionImplicit(remote, skipKey)
	case len(prior) > 0:
		merged[key] = prior
	}
}

func mergeSpecSectionExplicit(_remote, prior map[string]any, skipKey func(string) bool) map[string]any {
	if len(prior) == 0 {
		return nil
	}
	merged := make(map[string]any, len(prior))
	for k, v := range prior {
		if skipKey != nil && skipKey(k) {
			continue
		}
		merged[k] = v
	}
	return merged
}

func mergeSpecSectionImplicit(remote map[string]any, skipKey func(string) bool) map[string]any {
	if len(remote) == 0 {
		return nil
	}
	merged := make(map[string]any, len(remote))
	for k, v := range remote {
		if skipKey != nil && skipKey(k) {
			continue
		}
		merged[k] = v
	}
	return merged
}

func priorSpecSections(state types.String) map[string]map[string]any {
	if state.IsNull() || state.IsUnknown() {
		return nil
	}
	var spec struct {
		IntegrationSpec map[string]any `json:"integrationSpec"`
		AppSpec         map[string]any `json:"appSpec"`
	}
	if err := json.Unmarshal([]byte(state.ValueString()), &spec); err != nil {
		return nil
	}
	return map[string]map[string]any{
		"integrationSpec": spec.IntegrationSpec,
		"appSpec":         spec.AppSpec,
	}
}
