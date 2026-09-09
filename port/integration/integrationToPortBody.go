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
		attrs := state.WebhookChangelogDestination.Attributes()
		url, _ := attrs["url"].(types.String)
		agent, _ := attrs["agent"].(types.Bool)
		integration.ChangelogDestination = &cli.ChangelogDestination{
			Type:  consts.Webhook,
			Url:   url.ValueString(),
			Agent: agent.ValueBoolPointer(),
		}
	case isConfigured(state.KafkaChangelogDestination):
		integration.ChangelogDestination = &cli.ChangelogDestination{Type: consts.Kafka}
	}

	return integration, nil
}

// isConfigured reports whether a computed attribute holds a value to send to
// Port. Unknown means Port owns it and will return it on the next refresh.
func isConfigured(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown()
}
