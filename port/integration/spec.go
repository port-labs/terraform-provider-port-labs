package integration

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// serverManagedSpecKeys are spec sub-objects that Port sets internally
// (e.g. sizing, applier backend). They must never appear in user config
// and are stripped on read via Integration.UnmarshalJSON.
var serverManagedSpecKeys = [...]string{"systemSpec", "privateSpec"}

// parseSpecFromConfig parses the user-provided spec JSON string into
// an IntegrationClientSpec, rejecting any server-managed keys.
func parseSpecFromConfig(raw types.String) (*cli.IntegrationClientSpec, error) {
	if raw.IsNull() || raw.IsUnknown() || raw.ValueString() == "" {
		return nil, nil
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(raw.ValueString()), &parsed); err != nil {
		return nil, fmt.Errorf("invalid spec JSON: %w", err)
	}

	for _, key := range serverManagedSpecKeys {
		if _, ok := parsed[key]; ok {
			return nil, fmt.Errorf(
				"spec.%s is server-managed by Port and cannot be set in Terraform; only integrationSpec and appSpec are supported", key,
			)
		}
	}

	var spec cli.IntegrationClientSpec
	raw2, _ := json.Marshal(parsed)
	if err := json.Unmarshal(raw2, &spec); err != nil {
		return nil, fmt.Errorf("invalid spec structure: %w", err)
	}

	if spec.IntegrationSpec == nil && spec.AppSpec == nil {
		return nil, nil
	}
	return &spec, nil
}

// specToState serializes an IntegrationClientSpec back to a Terraform string,
// preserving key order from the existing state when semantically equal.
func specToState(spec *cli.IntegrationClientSpec, existing types.String, escapeHTML bool) (types.String, error) {
	if spec == nil || (spec.IntegrationSpec == nil && spec.AppSpec == nil) {
		if existing.IsNull() {
			return types.StringNull(), nil
		}
		return existing, nil
	}
	return utils.GoObjectToTerraformStringPreferExisting(existing, spec, escapeHTML)
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
