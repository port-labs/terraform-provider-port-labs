package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	if plan.InstallationAppType.IsNull() || plan.InstallationAppType.IsUnknown() {
		return
	}
	if plan.InstallationId.IsNull() || plan.InstallationId.IsUnknown() {
		return
	}

	if plan.isSaasOAuth2() {
		if !specConfiguredInHCL {
			return
		}
		rawSpec := extractAppSpec(plan.Spec)
		if rawSpec.IsNull() || rawSpec.IsUnknown() || rawSpec.ValueString() == "" {
			return
		}
		spec, err := parseSpecFromConfig(rawSpec)
		if err != nil {
			diags.AddError("invalid spec", err.Error())
			return
		}
		if spec == nil || spec.AppSpec == nil {
			return
		}
		body := cli.ValidateIntegrationSpecBody{
			ValidationMode:   consts.ValidateSaasOAuth2Spec,
			InstallationId:   plan.InstallationId.ValueString(),
			InstallationType: plan.installationType(),
			Spec:             spec,
		}
		if err := r.portClient.ValidateIntegrationSpec(ctx, plan.InstallationAppType.ValueString(), body); err != nil {
			diags.AddError("invalid integration spec", err.Error())
		}
		return
	}

	if plan.Spec.IsNull() || plan.Spec.IsUnknown() || plan.Spec.ValueString() == "" {
		return
	}

	spec, err := parseSpecFromConfig(plan.Spec)
	if err != nil {
		diags.AddError("invalid spec", err.Error())
		return
	}
	if spec == nil || spec.IntegrationSpec == nil {
		return
	}

	body := cli.ValidateIntegrationSpecBody{
		ValidationMode:   consts.ValidateSaasSpec,
		InstallationId:   plan.InstallationId.ValueString(),
		InstallationType: plan.installationType(),
		Spec:             spec,
		Options: &cli.ValidateIntegrationSpecOptions{
			SkipSecretExistenceCheck: true,
		},
	}

	if err := r.portClient.ValidateIntegrationSpec(ctx, plan.InstallationAppType.ValueString(), body); err != nil {
		diags.AddError("invalid integration spec", err.Error())
	}
}
