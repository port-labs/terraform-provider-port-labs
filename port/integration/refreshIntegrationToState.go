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

	state.GithubExternalPropertiesNamespaceClaims = githubExternalPropertiesNamespaceClaimsToState(a.GithubExternalPropertiesNamespaceClaims)

	return nil
}

func githubExternalPropertiesNamespaceClaimsToState(claims map[string]bool) types.Map {
	if len(claims) == 0 {
		return types.MapNull(types.BoolType)
	}

	elements := make(map[string]attr.Value, len(claims))
	for org, claimed := range claims {
		elements[org] = types.BoolValue(claimed)
	}

	result, diags := types.MapValue(types.BoolType, elements)
	if diags.HasError() {
		return types.MapNull(types.BoolType)
	}

	return result
}
