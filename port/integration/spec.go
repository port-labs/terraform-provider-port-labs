package integration

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// parseSpecFromConfig parses the user-provided spec JSON string into
// an IntegrationClientSpec. Only integrationSpec and appSpec are read; any
// other top-level keys in the JSON are ignored.
func parseSpecFromConfig(raw types.String) (*cli.IntegrationClientSpec, error) {
	if raw.IsNull() || raw.IsUnknown() || raw.ValueString() == "" {
		return nil, nil
	}

	var spec cli.IntegrationClientSpec
	if err := json.Unmarshal([]byte(raw.ValueString()), &spec); err != nil {
		return nil, fmt.Errorf("invalid spec JSON: %w", err)
	}
	if spec.IsEmpty() {
		return nil, nil
	}
	return &spec, nil
}

func validateIntegrationModel(m *IntegrationModel) error {
	if err := validateSaasSpec(m); err != nil {
		return err
	}

	if !m.isSaas() && specIsConfigured(m.Spec) {
		return fmt.Errorf(
			"spec is only supported when installation_type is %q",
			consts.InstallationTypeSaas,
		)
	}
	
	return nil
}

func specIsConfigured(spec types.String) bool {
	return !spec.IsNull() && !spec.IsUnknown() && spec.ValueString() != ""
}

// validateSaasSpec checks that a Saas integration has the required spec fields.
func validateSaasSpec(m *IntegrationModel) error {
	if !m.isSaas() {
		return nil
	}

	if m.Spec.IsNull() || m.Spec.ValueString() == "" {
		return fmt.Errorf("spec is required when installation_type is %q", m.installationType())
	}

	spec, err := parseSpecFromConfig(m.Spec)
	if err != nil {
		return err
	}
	if spec == nil || spec.IntegrationSpec == nil {
		return fmt.Errorf("spec.integrationSpec is required when installation_type is %q", m.installationType())
	}
	return nil
}
