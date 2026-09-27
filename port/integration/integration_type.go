package integration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func integrationTypeFromModel(m *IntegrationModel) string {
	if m == nil {
		return ""
	}
	if s := configuredString(m.IntegrationType); s != "" {
		return s
	}
	return configuredString(m.InstallationAppType)
}

// integrationTypeValue resolves the Ocean integration type from one or more models.
// integration_type wins when both it and installation_app_type are set.
func integrationTypeValue(models ...*IntegrationModel) (string, error) {
	for _, m := range models {
		if s := integrationTypeFromModel(m); s != "" {
			return s, nil
		}
	}
	return "", fmt.Errorf("integration_type is required")
}

func alignIntegrationType(config, plan *IntegrationModel) error {
	if config != nil && (config.IntegrationType.IsUnknown() || config.InstallationAppType.IsUnknown()) {
		return nil
	}
	resolved, err := integrationTypeValue(config, plan)
	if err != nil {
		return err
	}
	plan.IntegrationType = types.StringValue(resolved)
	return nil
}

func applyIntegrationType(m *IntegrationModel, integration *cli.Integration) {
	resolved := integrationTypeFromAPI(integration)
	if resolved == "" {
		m.IntegrationType = types.StringNull()
		return
	}
	m.IntegrationType = types.StringValue(resolved)
}

func integrationTypeFromAPI(integration *cli.Integration) string {
	if integration == nil {
		return ""
	}
	if integration.IntegrationType != nil && *integration.IntegrationType != "" {
		return *integration.IntegrationType
	}
	if integration.InstallationAppType != nil {
		return *integration.InstallationAppType
	}
	return ""
}

func configuredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}
