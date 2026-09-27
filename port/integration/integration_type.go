package integration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// integrationTypeValue resolves the Ocean integration type.
// integration_type is required. A deprecated installation_app_type is copied
// into it when integration_type is omitted. integ-service stores that value as
// integrationType and copies it onto installationAppType.
func integrationTypeValue(m *IntegrationModel) (string, error) {
	if m == nil {
		return "", fmt.Errorf("integration_type is required")
	}

	fromType := configuredString(m.IntegrationType)
	fromAppType := configuredString(m.InstallationAppType)
	resolved := fromType
	if resolved == "" {
		resolved = fromAppType
	}
	if resolved == "" {
		return "", fmt.Errorf("integration_type is required")
	}
	if fromType != "" && fromAppType != "" && fromType != fromAppType {
		return "", fmt.Errorf("installation_app_type (%q) must match integration_type (%q)", fromAppType, fromType)
	}
	return resolved, nil
}

// resolveIntegrationType prefers an explicit integration_type. A deprecated
// installation_app_type is used only when integration_type is absent from
// configuration. Values already stored on the plan fill in when the
// configuration omits both.
func resolveIntegrationType(config, plan *IntegrationModel) (string, error) {
	if config != nil {
		fromType := configuredString(config.IntegrationType)
		fromAppType := configuredString(config.InstallationAppType)
		if fromType != "" && fromAppType != "" && fromType != fromAppType {
			return "", fmt.Errorf("installation_app_type (%q) must match integration_type (%q)", fromAppType, fromType)
		}
		if fromType != "" {
			return fromType, nil
		}
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
	plan.InstallationAppType = types.StringValue(resolved)
	return nil
}

func applyIntegrationType(m *IntegrationModel, integration *cli.Integration) {
	resolved := ""
	if integration != nil && integration.IntegrationType != nil && *integration.IntegrationType != "" {
		resolved = *integration.IntegrationType
	} else if integration != nil && integration.InstallationAppType != nil {
		resolved = *integration.InstallationAppType
	}
	if resolved == "" {
		m.IntegrationType = types.StringNull()
		m.InstallationAppType = types.StringNull()
		return
	}
	m.IntegrationType = types.StringValue(resolved)
	m.InstallationAppType = types.StringValue(resolved)
}

func configuredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}
