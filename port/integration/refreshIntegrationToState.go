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

	if a.ChangelogDestination != nil {
		if a.ChangelogDestination.Type == consts.Kafka {
			m.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
			m.WebhookChangelogDestination = nil
		} else {
			if a.ChangelogDestination.Url != "" {
				m.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
					Url: types.StringValue(a.ChangelogDestination.Url),
				}
				if a.ChangelogDestination.Agent != nil {
					m.WebhookChangelogDestination.Agent = types.BoolValue(*a.ChangelogDestination.Agent)
				}
				m.KafkaChangelogDestination = types.ObjectNull(map[string]attr.Type{})
			}
		}
	} else {
		m.KafkaChangelogDestination = types.ObjectNull(map[string]attr.Type{})
		m.WebhookChangelogDestination = nil
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
// holding spec and config at their planned values, which Terraform compares
// against. A Computed attribute the user left out of the config arrives here
// unknown, and every attribute must be known once apply returns.
func applyWriteResult(plan *IntegrationModel, a *cli.Integration, integrationId string) {
	applyServerFields(plan, a, integrationId)

	plan.Spec = nullIfUnknown(plan.Spec)
	plan.Config = nullIfUnknown(plan.Config)
}

func nullIfUnknown(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

// mergeSpec renders Port's spec as JSON, restoring the sensitive
// integrationSpec values (organization secret references) that Port blanks out
// in read responses.
func mergeSpec(state types.String, remote *cli.IntegrationClientSpec, jsonEscapeHTML bool) types.String {
	prior := priorIntegrationSpec(state)

	merged := make(map[string]any, 2)
	switch {
	case remote.IntegrationSpec != nil:
		merged["integrationSpec"] = withPriorSecrets(remote.IntegrationSpec, prior)
	case prior != nil:
		merged["integrationSpec"] = prior
	}
	if remote.AppSpec != nil {
		merged["appSpec"] = remote.AppSpec
	}

	encoded, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return state
	}
	return encoded
}

func priorIntegrationSpec(state types.String) map[string]any {
	if state.IsNull() || state.IsUnknown() {
		return nil
	}
	var spec struct {
		IntegrationSpec map[string]any `json:"integrationSpec"`
	}
	if err := json.Unmarshal([]byte(state.ValueString()), &spec); err != nil {
		return nil
	}
	return spec.IntegrationSpec
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
