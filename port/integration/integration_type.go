package integration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// integrationTypeValue resolves the Ocean integration type from Terraform models.
// integration_type wins when both attributes are set; installation_app_type is a
// silent legacy fallback. Used when we need the type string without mutating the
// model: API body construction, SaaS spec validation, and immutability checks.
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

// alignIntegrationType runs during ModifyPlan before validation. It resolves the
// type from config (and plan when needed), copies it onto plan.integration_type,
// and errors if neither attribute is set. This keeps legacy configs working without
// exposing the migration in docs. An unset installation_app_type is planned as
// null because Read never populates it, so leaving it unknown fails the apply.
func alignIntegrationType(config, plan *IntegrationModel) error {
	if config != nil && config.InstallationAppType.IsNull() && plan.InstallationAppType.IsUnknown() {
		plan.InstallationAppType = types.StringNull()
	}
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

// applyIntegrationType runs on read after GET/create/update responses. Port may
// return integrationType and/or installationAppType; this writes the resolved
// value into state.integration_type only.
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
