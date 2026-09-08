package integration

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *IntegrationResource) refreshIntegrationState(state *IntegrationModel, remote *cli.Integration) error {
	state.ID = types.StringValue(remote.InstallationId)
	state.InstallationId = types.StringValue(remote.InstallationId)
	state.Title = types.StringPointerValue(remote.Title)
	state.InstallationAppType = types.StringPointerValue(remote.InstallationAppType)
	state.InstallationType = types.StringPointerValue(remote.InstallationType)
	state.Version = types.StringPointerValue(remote.Version)

	if remote.StatusInfo != nil {
		state.Status = types.StringValue(remote.StatusInfo.IntegrationStatus.Status)
	} else {
		state.Status = types.StringNull()
	}

	if remote.Spec != nil {
		state.Spec = mergeSpec(state.Spec, remote.Spec, r.portClient.JSONEscapeHTML)
	}

	if remote.Config != nil {
		config, _ := utils.GoObjectToTerraformStringPreferExisting(state.Config, remote.Config, r.portClient.JSONEscapeHTML)
		state.Config = config
	}

	if remote.ChangelogDestination != nil {
		switch remote.ChangelogDestination.Type {
		case consts.Kafka:
			state.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
		case consts.Webhook:
			if remote.ChangelogDestination.Url != "" {
				state.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
					Url: types.StringValue(remote.ChangelogDestination.Url),
				}
				if remote.ChangelogDestination.Agent != nil {
					state.WebhookChangelogDestination.Agent = types.BoolValue(*remote.ChangelogDestination.Agent)
				}
			}
		}
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
	// Parse user's current spec from state so we can preserve secret refs.
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
