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
	state.InstallationType = types.StringPointerValue(a.InstallationType)
	state.InstallationAppType = types.StringPointerValue(a.InstallationAppType)
	state.Version = types.StringPointerValue(a.Version)

	if a.Config != nil {
		config, _ := utils.GoObjectToTerraformStringPreferExisting(state.Config, a.Config, r.portClient.JSONEscapeHTML)
		state.Config = config
	}
	if a.Spec != nil {
		if state.Spec.IsNull() || state.Spec.IsUnknown() {
			spec, err := utils.GoObjectToTerraformString(a.Spec, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			state.Spec = spec
		}
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
