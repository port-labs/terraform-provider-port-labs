package integration

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

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

// mergeSpec builds the refreshed spec from the server response while preserving
// user-configured values for sensitive integrationSpec fields that the server
// strips (org secret references resolved at runtime).
//
// systemSpec and privateSpec are always excluded — those are server-managed
// and already stripped by the CLI layer's toClientSpec().
func mergeSpec(stateTF types.String, remote *cli.IntegrationClientSpec, jsonEscapeHTML bool) types.String {
	var userSpec map[string]map[string]any
	if !stateTF.IsNull() && !stateTF.IsUnknown() {
		_ = json.Unmarshal([]byte(stateTF.ValueString()), &userSpec)
	}

	merged := make(map[string]map[string]any)

	// integrationSpec: take server values, but keep user values for any key
	// the server returns as nil/empty (stripped secrets).
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

	// appSpec: take server values as-is (no sensitive fields).
	if remote.AppSpec != nil {
		merged["appSpec"] = remote.AppSpec
	}

	encoded, err := utils.GoObjectToTerraformString(merged, jsonEscapeHTML)
	if err != nil {
		return stateTF
	}
	return encoded
}
