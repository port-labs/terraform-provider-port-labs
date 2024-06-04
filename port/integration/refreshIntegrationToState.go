package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshIntegrationState(state *IntegrationModel, a *cli.Integration, integrationId string) error {
	state.ID = types.StringValue(integrationId)
	state.InstallationId = types.StringValue(integrationId)
	if a.InstallationAppType != nil && len(*a.InstallationAppType) != 0 {
		state.InstallationAppType = flex.GoStringToFramework(a.InstallationAppType)
	}
	state.Title = types.StringValue(a.Title)
	state.Version = types.StringValue(a.Version)

	if a.Config != nil {
		config, _ := utils.GoObjectToTerraformString(a.Config)
		state.Config = config
	}
	if a.ChangelogDestination != nil {
		if a.ChangelogDestination.Type == consts.Kafka {
			state.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
		} else {
			state.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
				Url: types.StringValue(a.ChangelogDestination.Url),
			}
			if a.ChangelogDestination.Agent != nil {
				state.WebhookChangelogDestination.Agent = types.BoolValue(*a.ChangelogDestination.Agent)
			}
		}
	}

	return nil
}
