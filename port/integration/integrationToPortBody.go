package integration

import (
	"fmt"
	"net/url"
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

	if !state.OauthBrokerUrl.IsNull() {
		oauthBrokerUrl := state.OauthBrokerUrl.ValueString()
		parsed, err := url.Parse(oauthBrokerUrl)
		if err != nil {
			return nil, fmt.Errorf("oauth_broker_url must be a valid URL: %w", err)
		}
		if parsed.RawQuery != "" {
			return nil, fmt.Errorf("oauth_broker_url must not contain a query string — Port appends query parameters at runtime")
		}
		integration.OauthBrokerUrl = &oauthBrokerUrl
	}

	return integration, nil
}
