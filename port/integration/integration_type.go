package integration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// integrationTypeValue resolves the Ocean integration type from one or more models.
// integration_type wins when both it and installation_app_type are set.
func integrationTypeValue(models ...*IntegrationModel) (string, error) {
	for _, m := range models {
		if m == nil {
			continue
		}
		if !m.IntegrationType.IsNull() && !m.IntegrationType.IsUnknown() {
			if s := m.IntegrationType.ValueString(); s != "" {
				return s, nil
			}
		}
		if !m.InstallationAppType.IsNull() && !m.InstallationAppType.IsUnknown() {
			if s := m.InstallationAppType.ValueString(); s != "" {
				return s, nil
			}
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
