package integration

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/adoption"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// adoptionConflict reports the first field that cannot be changed after creation and differs between the
// configuration and the existing integration, which makes taking it over impossible
func adoptionConflict(plan *IntegrationModel, existing *cli.Integration) *cli.AdoptionConflictError {
	existingInstallationType := consts.InstallationTypeOnPrem
	if existing.InstallationType != nil && *existing.InstallationType != "" {
		existingInstallationType = *existing.InstallationType
	}
	if plan.installationType() != existingInstallationType {
		return &cli.AdoptionConflictError{
			ResourceKind: consts.ResourceKindIntegration,
			Identifier:   plan.InstallationId.ValueString(),
			Field:        "installation_type",
			Existing:     existingInstallationType,
			Configured:   plan.installationType(),
		}
	}

	existingAppType := ""
	if existing.InstallationAppType != nil {
		existingAppType = *existing.InstallationAppType
	}
	if isKnown(plan.InstallationAppType) && plan.InstallationAppType.ValueString() != existingAppType {
		return &cli.AdoptionConflictError{
			ResourceKind: consts.ResourceKindIntegration,
			Identifier:   plan.InstallationId.ValueString(),
			Field:        "installation_app_type",
			Existing:     existingAppType,
			Configured:   plan.InstallationAppType.ValueString(),
		}
	}
	return nil
}

// checkExistingAtPlan checks, before apply, whether an integration about to be created already exists in Port,
// so that taking it over is never silent: it fails unless adoption is allowed and possible, and warns when it is.
func (r *IntegrationResource) checkExistingAtPlan(ctx context.Context, plan *IntegrationModel, diags *diag.Diagnostics) {
	if r.portClient == nil || !isKnown(plan.InstallationId) {
		return
	}

	existing, statusCode, err := r.portClient.GetIntegration(ctx, plan.InstallationId.ValueString())
	if err != nil {
		if statusCode != http.StatusNotFound {
			diags.AddWarning("failed to check whether the integration already exists", err.Error())
		}
		return
	}

	if !adoption.TakeOver(adoption.AtPlan, r.portClient, existingIntegration(plan), plan.AdoptExisting, adoptionConflict(plan, existing), diags) {
		return
	}
	if isKnown(plan.CreatePortResourcesOrigin) {
		diags.AddAttributeWarning(path.Root("create_port_resources_origin"), "create_port_resources_origin is ignored",
			"`create_port_resources_origin` only applies when an integration is created. The existing integration keeps the resources it already created.")
	}
}

func existingIntegration(plan *IntegrationModel) *cli.ResourceExistsError {
	return &cli.ResourceExistsError{ResourceKind: consts.ResourceKindIntegration, Identifier: plan.InstallationId.ValueString()}
}

// adoptIntegration takes ownership of an integration that already exists in Port by updating it to match the plan
func (r *IntegrationResource) adoptIntegration(ctx context.Context, plan *IntegrationModel, existing *cli.Integration, resp *resource.CreateResponse) {
	installationId := plan.InstallationId.ValueString()
	if !adoption.TakeOver(adoption.AtApply, r.portClient, existingIntegration(plan), plan.AdoptExisting, adoptionConflict(plan, existing), &resp.Diagnostics) {
		return
	}

	body, err := integrationToPortBody(plan, false)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, installationId, body)
	if err != nil {
		resp.Diagnostics.AddError("failed to take over integration", err.Error())
		return
	}

	applyWriteResult(plan, updated, installationId)

	if plan.isHosted() {
		r.awaitInfra(ctx, plan, installationId, "taken over", true, false, &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
