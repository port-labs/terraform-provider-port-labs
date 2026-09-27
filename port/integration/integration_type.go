package integration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// integrationTypeValue resolves the Ocean integration type.
// integration_type wins when both it and installation_app_type are set.
func integrationTypeValue(m *IntegrationModel) (string, error) {
	if m == nil {
		return "", fmt.Errorf("integration_type is required")
	}

	fromType := configuredString(m.IntegrationType)
	if fromType != "" {
		return fromType, nil
	}
	fromAppType := configuredString(m.InstallationAppType)
	if fromAppType != "" {
		return fromAppType, nil
	}
	return "", fmt.Errorf("integration_type is required")
}

func resolveIntegrationType(config, plan *IntegrationModel) (string, error) {
	if config != nil {
		fromType := configuredString(config.IntegrationType)
		if fromType != "" {
			return fromType, nil
		}
		fromAppType := configuredString(config.InstallationAppType)
		if fromAppType != "" {
			return fromAppType, nil
		}
	}
	return integrationTypeValue(plan)
}

func alignIntegrationType(config, plan *IntegrationModel) error {
	if config != nil && (config.IntegrationType.IsUnknown() || config.InstallationAppType.IsUnknown()) {
		return nil
	}
	resolved, err := resolveIntegrationType(config, plan)
	if err != nil {
		return err
	}
	plan.IntegrationType = types.StringValue(resolved)
	return nil
}

func applyIntegrationType(m *IntegrationModel, integration *cli.Integration) {
	if integration != nil && integration.IntegrationType != nil && *integration.IntegrationType != "" {
		m.IntegrationType = types.StringValue(*integration.IntegrationType)
		return
	}
	if integration != nil && integration.InstallationAppType != nil && *integration.InstallationAppType != "" {
		m.IntegrationType = types.StringValue(*integration.InstallationAppType)
		return
	}
	m.IntegrationType = types.StringNull()
}

func configuredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}
