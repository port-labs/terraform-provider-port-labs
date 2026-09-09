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

	if !a.Spec.IsEmpty() {
		state.Spec = mergeSpec(state.Spec, a.Spec, r.portClient.JSONEscapeHTML)
	}
	if a.Config != nil {
		state.Config, _ = utils.GoObjectToTerraformStringPreferExisting(state.Config, a.Config, r.portClient.JSONEscapeHTML)
	}
}

// applyWriteResult syncs the server response after Create or Update while
// keeping every Computed attribute at its planned value. Port changes status
// and version during writes, and fills spec/config with server defaults — none
// of that may differ from what Terraform planned. Read picks it all up next.
func applyWriteResult(plan *IntegrationModel, a *cli.Integration, integrationId string) {
	savedSpec := plan.Spec
	savedConfig := plan.Config
	savedStatus := plan.Status
	savedVersion := plan.Version

	applyServerFields(plan, a, integrationId)

	plan.Spec = nullIfUnknown(savedSpec)
	plan.Config = nullIfUnknown(savedConfig)
	plan.Status = nullIfUnknown(savedStatus)
	plan.Version = nullIfUnknown(savedVersion)
}

func nullIfUnknown(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

// mergeSpec renders Port's spec as JSON, restoring values that Port strips or
// overwrites: sensitive integrationSpec values (org secret references) and
// user-provided appSpec fields (e.g. liveEventsEnabled).
func mergeSpec(state types.String, remote *cli.IntegrationClientSpec, jsonEscapeHTML bool) types.String {
	prior := priorSpecSections(state)

	merged := make(map[string]any, 2)
	switch {
	case remote.IntegrationSpec != nil:
		merged["integrationSpec"] = withPriorSecrets(remote.IntegrationSpec, prior["integrationSpec"])
	case prior["integrationSpec"] != nil:
		merged["integrationSpec"] = prior["integrationSpec"]
	}

	// Server always returns appSpec. Merge with user's values so fields the
	// user explicitly set (e.g. liveEventsEnabled) survive a refresh.
	switch {
	case remote.AppSpec != nil && prior["appSpec"] != nil:
		merged["appSpec"] = mergeAppSpec(remote.AppSpec, prior["appSpec"])
	case remote.AppSpec != nil:
		merged["appSpec"] = remote.AppSpec
	case prior["appSpec"] != nil:
		merged["appSpec"] = prior["appSpec"]
	}

	encoded, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return state
	}
	return encoded
}

// mergeAppSpec merges server appSpec with user's prior values. Server wins for
// fields it manages, but user-set fields are preserved when the server returns
// the same key with a different value (user explicitly controls it).
func mergeAppSpec(remote, prior map[string]any) map[string]any {
	merged := make(map[string]any, len(remote))
	for k, v := range remote {
		merged[k] = v
	}
	// Keep user values for keys they explicitly set — the server returned
	// value is the "default" but the user override should win on read so
	// that the next plan doesn't show a diff.
	for k, v := range prior {
		if _, exists := merged[k]; !exists {
			merged[k] = v
		}
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

// withPriorSecrets fills the blanks Port leaves behind with the values already
// held in state, so secret references configured in HCL survive a refresh.
func withPriorSecrets(remote, prior map[string]any) map[string]any {
	merged := make(map[string]any, len(remote))
	for k, v := range remote {
		merged[k] = v
	}
	for k, v := range prior {
		if isBlank(merged[k]) && !isBlank(v) {
			merged[k] = v
		}
	}
	return merged
}

func isBlank(v any) bool {
	return v == nil || v == ""
}
