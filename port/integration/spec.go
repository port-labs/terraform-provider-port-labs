package integration

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// serverManagedSpecKeys are spec sub-objects that Port sets internally
// (e.g. sizing, applier backend). They must never appear in user config, and
// IntegrationClientSpec drops them from responses.
var serverManagedSpecKeys = []string{"systemSpec", "privateSpec"}

// parseSpecFromConfig parses the user-provided spec JSON string into
// an IntegrationClientSpec, rejecting any server-managed keys.
func parseSpecFromConfig(raw types.String) (*cli.IntegrationClientSpec, error) {
	if raw.IsNull() || raw.IsUnknown() || raw.ValueString() == "" {
		return nil, nil
	}
	data := []byte(raw.ValueString())

	var sections map[string]json.RawMessage
	if err := json.Unmarshal(data, &sections); err != nil {
		return nil, fmt.Errorf("invalid spec JSON: %w", err)
	}
	for _, key := range serverManagedSpecKeys {
		if _, ok := sections[key]; ok {
			return nil, fmt.Errorf(
				"spec.%s is server-managed by Port and cannot be set in Terraform; only integrationSpec and appSpec are supported", key,
			)
		}
	}

	var spec cli.IntegrationClientSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("invalid spec structure: %w", err)
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
