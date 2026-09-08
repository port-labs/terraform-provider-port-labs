package integration

import (
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

	// Spec is deliberately NOT refreshed from the server response.
	// The server strips sensitive integrationSpec values (they're org secret
	// references resolved at runtime) and injects appSpec defaults that the
	// user never configured. Overwriting would cause permanent drift.
	// We keep whatever the user configured in their HCL.

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
