package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func validatePlanRules(plan, state *IntegrationModel, isCreate bool, diags *diag.Diagnostics) {
	if isCreate && plan.isSaasOAuth2() {
		diags.AddError(
			"SaasOAuth2 integrations cannot be created via Terraform",
			"Authorize the integration in the Port UI, then import it with "+
				"`terraform import port_integration.<name> <installation_id>`.",
		)
	}

	if isCreate && !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		diags.AddError(
			"config cannot be set on creation",
			"Integrations receive default mappings during provisioning. "+
				"Create the integration first (without config), then add config "+
				"to your HCL and run 'terraform apply' again to override the defaults.",
		)
	}

	if !isCreate && !plan.CreatePortResourcesOrigin.Equal(state.CreatePortResourcesOrigin) {
		diags.AddError(
			"cannot change create_port_resources_origin",
			"`create_port_resources_origin` can only be set when creating an integration. To use a different value, destroy this resource and create a new `port_integration`.",
		)
	}

	if !isCreate {
		hadDestination := isConfigured(state.KafkaChangelogDestination) || isConfigured(state.WebhookChangelogDestination)
		lostDestination := !isConfigured(plan.KafkaChangelogDestination) && !isConfigured(plan.WebhookChangelogDestination)
		if hadDestination && lostDestination {
			diags.AddError(
				"cannot remove changelog destination",
				"The Port API does not support removing a changelog destination from an existing integration. "+
					"To remove it, delete and recreate the integration (e.g. taint the resource).",
			)
		}
	}

	if err := validateIntegrationModel(plan); err != nil {
		diags.AddError("invalid integration configuration", err.Error())
	}
}

func (r *IntegrationResource) validateSpecAtPlan(ctx context.Context, plan *IntegrationModel, specConfiguredInHCL bool, diags *diag.Diagnostics) {
	if r.portClient == nil || !plan.isHosted() {
		return
	}
	if !isKnown(plan.InstallationAppType) || !isKnown(plan.InstallationId) {
		return
	}

	body := cli.ValidateIntegrationSpecBody{
		InstallationId:   plan.InstallationId.ValueString(),
		InstallationType: plan.installationType(),
	}

	rawSpec := plan.Spec
	if plan.isSaasOAuth2() {
		if !specConfiguredInHCL {
			return
		}
		rawSpec = extractAppSpec(rawSpec)
		body.ValidationMode = consts.ValidateSaasOAuth2Spec
	} else {
		body.ValidationMode = consts.ValidateSaasSpec
		body.Options = &cli.ValidateIntegrationSpecOptions{SkipSecretExistenceCheck: true}
	}

	spec, err := parseSpecFromConfig(rawSpec)
	if err != nil {
		diags.AddError("invalid spec", err.Error())
		return
	}
	if spec == nil || (plan.isSaas() && spec.IntegrationSpec == nil) {
		return
	}
	body.Spec = spec

	if err := r.portClient.ValidateIntegrationSpec(ctx, plan.InstallationAppType.ValueString(), body); err != nil {
		diags.AddError("invalid integration spec", err.Error())
	}
}

func isKnown(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown()
}
