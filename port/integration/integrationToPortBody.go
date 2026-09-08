package integration

import (
	"fmt"
	"regexp"

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
		return nil, fmt.Errorf(
			"installation_id must match the pattern %s: got %q",
			installationIdPattern, installationId,
		)
	}

	installationType := state.installationType()
	integration := &cli.Integration{
		InstallationId:      installationId,
		Title:               state.Title.ValueStringPointer(),
		Version:             state.Version.ValueStringPointer(),
		InstallationAppType: state.InstallationAppType.ValueStringPointer(),
		InstallationType:    &installationType,
	}

	if state.isSaas() {
		spec, err := parseSpecFromConfig(state.Spec)
		if err != nil {
			return nil, err
		}
		integration.Spec = spec
	}

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
