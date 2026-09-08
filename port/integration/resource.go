package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &IntegrationResource{}
var _ resource.ResourceWithImportState = &IntegrationResource{}

func NewIntegrationResource() resource.Resource {
	return &IntegrationResource{}
}

type IntegrationResource struct {
	portClient *cli.PortClient
}

func (r *IntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *IntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.portClient = req.ProviderData.(*cli.PortClient)
	}
}

func (r *IntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("installation_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateSaasSpec(&plan); err != nil {
		resp.Diagnostics.AddError("invalid SaaS integration configuration", err.Error())
		return
	}

	body, err := integrationToPortBody(&plan)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	created, err := r.portClient.CreateIntegration(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("failed to create integration", err.Error())
		return
	}

	if plan.isSaas() {
		ready, pollErr := r.portClient.WaitForIntegrationReady(ctx, created.InstallationId)
		if pollErr != nil {
			// The integration was created in Port but provisioning hasn't finished.
			// Save it to state so the user can re-apply (to re-poll) or destroy.
			resp.Diagnostics.AddWarning(
				"integration created but not yet ready",
				pollErr.Error()+". The integration has been saved to state. Run 'terraform apply' again to re-check, or 'terraform destroy' to clean up.",
			)
		} else {
			created = ready
		}
	}

	if err := r.refreshIntegrationState(&plan, created); err != nil {
		resp.Diagnostics.AddError("failed to refresh state after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, statusCode, err := r.portClient.GetIntegration(ctx, state.InstallationId.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read integration", err.Error())
		return
	}

	if err := r.refreshIntegrationState(&state, integration); err != nil {
		resp.Diagnostics.AddError("failed to refresh integration state", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, err := integrationToPortBody(&plan)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, plan.InstallationId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("failed to update integration", err.Error())
		return
	}

	if plan.isSaas() {
		ready, pollErr := r.portClient.WaitForIntegrationReady(ctx, plan.InstallationId.ValueString())
		if pollErr != nil {
			resp.Diagnostics.AddWarning(
				"integration updated but not yet ready",
				pollErr.Error()+". Run 'terraform apply' again to re-check status.",
			)
		} else {
			updated = ready
		}
	}

	if err := r.refreshIntegrationState(&plan, updated); err != nil {
		resp.Diagnostics.AddError("failed to refresh state after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.portClient.DeleteIntegration(ctx, state.InstallationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete integration", err.Error())
		return
	}

	if state.isSaas() {
		if err := r.portClient.WaitForIntegrationDeleted(ctx, state.InstallationId.ValueString()); err != nil {
			resp.Diagnostics.AddError("integration deletion did not complete", err.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
