package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshIntegrationState(state *IntegrationModel, a *cli.Integration, integrationId string) error {
	state.ID = types.StringValue(integrationId)
	state.InstallationId = types.StringValue(integrationId)

	state.Title = types.StringPointerValue(a.Title)
	state.InstallationAppType = types.StringPointerValue(a.InstallationAppType)
	state.Version = types.StringPointerValue(a.Version)

	if a.Config != nil {
		config, _ := utils.GoObjectToTerraformString(a.Config)
		state.Config = config
	}
	if a.ChangelogDestination != nil {
		if a.ChangelogDestination.Type == consts.Kafka {
			state.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
		} else {
			if a.ChangelogDestination.Url != "" {
				state.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
					Url: types.StringValue(a.ChangelogDestination.Url),
				}
				if a.ChangelogDestination.Agent != nil {
					state.WebhookChangelogDestination.Agent = types.BoolValue(*a.ChangelogDestination.Agent)
				}
			}
		}
	}

	return nil
}
