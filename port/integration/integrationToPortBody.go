package integration

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func integrationToPortBody(state *IntegrationModel) (*cli.Integration, error) {
	if state == nil {
		return nil, nil
	}

	integration := &cli.Integration{
		InstallationId: state.InstallationId.ValueString(),
	}

	integration.Title = state.Title.ValueStringPointer()
	integration.Version = state.Version.ValueStringPointer()
	integration.InstallationAppType = state.InstallationAppType.ValueStringPointer()

	if !state.Config.IsNull() {
		configStr := state.Config.ValueString()
		config, err := utils.TerraformJsonStringToGoObject(&configStr)
		if err != nil {
			return nil, err
		}
		integration.Config = config
	}
	if !state.KafkaChangelogDestination.IsNull() {
		integration.ChangelogDestination = &cli.ChangelogDestination{
			Type: consts.Kafka,
		}
	}
	if state.WebhookChangelogDestination != nil {
		integration.ChangelogDestination = &cli.ChangelogDestination{
			Type:  consts.Webhook,
			Url:   state.WebhookChangelogDestination.Url.ValueString(),
			Agent: state.WebhookChangelogDestination.Agent.ValueBoolPointer(),
		}
	}

	return integration, nil
}
