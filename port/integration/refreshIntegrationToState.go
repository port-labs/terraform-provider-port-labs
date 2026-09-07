package integration

import (
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
	state.Version = types.StringPointerValue(a.Version)

	if a.Config != nil {
		config, _ := utils.GoObjectToTerraformStringPreferExisting(state.Config, a.Config, r.portClient.JSONEscapeHTML)
		state.Config = config
	}

	state.KafkaChangelogDestination = types.ObjectNull(kafkaChangelogDestinationType)
	state.WebhookChangelogDestination = types.ObjectNull(webhookChangelogDestinationType)

	switch dest := a.ChangelogDestination; {
	case dest == nil:
	case dest.Type == consts.Kafka:
		state.KafkaChangelogDestination = types.ObjectValueMust(kafkaChangelogDestinationType, map[string]attr.Value{})
	case dest.Type == consts.Webhook && dest.Url != "":
		state.WebhookChangelogDestination = types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
			"url":   types.StringValue(dest.Url),
			"agent": types.BoolPointerValue(dest.Agent),
		})
	}

	return nil
}
