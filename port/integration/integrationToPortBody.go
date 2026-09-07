package integration

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const installationIdPattern = `^[a-z0-9-]+$`

var installationIdRegex = regexp.MustCompile(installationIdPattern)

func integrationToPortBody(state *IntegrationModel) (*cli.Integration, error) {
	if state == nil {
		return nil, nil
	}

	installationId := state.InstallationId.ValueString()

	if !installationIdRegex.MatchString(installationId) {
		return nil, fmt.Errorf("installation_id must match the pattern %s: must contain only lowercase letters, numbers, and dashes. Got: %q", installationIdPattern, installationId)
	}

	integration := &cli.Integration{
		InstallationId: installationId,
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
	switch {
	case isConfigured(state.WebhookChangelogDestination):
		integration.ChangelogDestination = webhookChangelogDestinationToPortBody(state.WebhookChangelogDestination)
	case isConfigured(state.KafkaChangelogDestination):
		integration.ChangelogDestination = &cli.ChangelogDestination{Type: consts.Kafka}
	}

	return integration, nil
}

func isConfigured(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown()
}

func webhookChangelogDestinationToPortBody(obj types.Object) *cli.ChangelogDestination {
	url, ok := obj.Attributes()["url"].(types.String)
	if !ok || url.IsNull() {
		return nil
	}
	// A missing agent leaves the zero value, which is null.
	agent, _ := obj.Attributes()["agent"].(types.Bool)

	return &cli.ChangelogDestination{
		Type:  consts.Webhook,
		Url:   url.ValueString(),
		Agent: agent.ValueBoolPointer(),
	}
}
