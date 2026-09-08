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

func (r *IntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *IntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *IntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("installation_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *IntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateSaasSpec(plan); err != nil {
		resp.Diagnostics.AddError("invalid SaaS integration configuration", err.Error())
		return
	}

	body, err := integrationToPortBody(plan)
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
			resp.Diagnostics.AddWarning(
				"integration created but not yet ready",
				pollErr.Error()+". The integration has been saved to state. Run 'terraform apply' again to re-check, or 'terraform destroy' to clean up.",
			)
		} else {
			created = ready
		}
	}

	applyWriteResult(plan, created, created.InstallationId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationIdentifier := state.InstallationId.ValueString()

	a, statusCode, err := r.portClient.GetIntegration(ctx, integrationIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading integration", err.Error())
		return
	}

	r.refreshIntegrationState(state, a, integrationIdentifier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *IntegrationModel
	var state *IntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationIdentifier := state.InstallationId.ValueString()

	hadDestination := !state.KafkaChangelogDestination.IsNull() || state.WebhookChangelogDestination != nil
	lostDestination := plan.KafkaChangelogDestination.IsNull() && plan.WebhookChangelogDestination == nil
	if hadDestination && lostDestination {
		resp.Diagnostics.AddError(
			"cannot remove changelog destination",
			"The Port API does not support removing a changelog destination from an existing integration. To remove it, delete and recreate the integration (e.g. taint the resource).",
		)
		return
	}

	body, err := integrationToPortBody(plan)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, integrationIdentifier, body)
	if err != nil {
		resp.Diagnostics.AddError("failed to update integration", err.Error())
		return
	}

	if plan.isSaas() {
		ready, pollErr := r.portClient.WaitForIntegrationReady(ctx, integrationIdentifier)
		if pollErr != nil {
			resp.Diagnostics.AddWarning(
				"integration updated but not yet ready",
				pollErr.Error()+". Run 'terraform apply' again to re-check status.",
			)
		} else {
			updated = ready
		}
	}

	applyWriteResult(plan, updated, integrationIdentifier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *IntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationIdentifier := state.InstallationId.ValueString()

	_, err := r.portClient.DeleteIntegration(ctx, integrationIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete integration", err.Error())
		return
	}

	if state.isSaas() {
		if err := r.portClient.WaitForIntegrationDeleted(ctx, integrationIdentifier); err != nil {
			resp.Diagnostics.AddError("integration deletion did not complete", err.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
