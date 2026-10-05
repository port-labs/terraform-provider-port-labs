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

// validateIntegrationConfig checks static configuration rules that apply even
// when the resource is being destroyed. Keep this limited to issues that are
// always invalid in HCL; defer SaaS spec requirements to ModifyPlan/CRUD.
func validateIntegrationConfig(m *IntegrationModel) error {
	if !m.isHosted() && specIsConfigured(m.Spec) {
		return fmt.Errorf(
			"spec is only supported when installation_type is %q or %q",
			consts.InstallationTypeSaas,
			consts.InstallationTypeSaasOAuth2,
		)
	}

	if m.isSaasOAuth2() {
		if err := validateSaasOAuth2Spec(m); err != nil {
			return err
		}
	}

	if specIsConfigured(m.Spec) {
		if _, err := parseSpecFromConfig(m.Spec); err != nil {
			return err
		}
	}

	return nil
}

func validateIntegrationModel(m *IntegrationModel) error {
	if err := validateIntegrationConfig(m); err != nil {
		return err
	}
	if err := validateSaasSpec(m); err != nil {
		return err
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

	if m.Spec.IsNull() || (!m.Spec.IsUnknown() && m.Spec.ValueString() == "") {
		return fmt.Errorf("spec is required when installation_type is %q", m.installationType())
	}
	if m.Spec.IsUnknown() {
		return nil
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

func validateSaasOAuth2Spec(m *IntegrationModel) error {
	if !specIsConfigured(m.Spec) {
		return nil
	}
	if m.Spec.IsUnknown() {
		return nil
	}

	var sections map[string]json.RawMessage
	if err := json.Unmarshal([]byte(m.Spec.ValueString()), &sections); err != nil {
		return fmt.Errorf("invalid spec JSON: %w", err)
	}
	if raw, ok := sections["integrationSpec"]; ok && len(raw) > 0 && string(raw) != "null" {
		var integrationSpec map[string]any
		if err := json.Unmarshal(raw, &integrationSpec); err != nil {
			return fmt.Errorf("invalid spec.integrationSpec JSON: %w", err)
		}
		if len(integrationSpec) > 0 {
			return fmt.Errorf(
				"spec.integrationSpec cannot be set when installation_type is %q; "+
					"OAuth-managed credentials are not Terraform-managed. Use spec.appSpec and config only",
				consts.InstallationTypeSaasOAuth2,
			)
		}
	}
	return nil
}
