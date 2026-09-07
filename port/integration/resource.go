package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &IntegrationResource{}
var _ resource.ResourceWithImportState = &IntegrationResource{}
var _ resource.ResourceWithModifyPlan = &IntegrationResource{}

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

func (r *IntegrationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state IntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.InstallationId.Equal(state.InstallationId) {
		resp.Diagnostics.AddError(
			"cannot change installation_id",
			"The Port API does not support changing `installation_id` on an existing integration. To use a different ID, destroy this resource (which deletes the integration from Port) and create a new `port_integration`.",
		)
	}

	if !plan.InstallationAppType.Equal(state.InstallationAppType) {
		resp.Diagnostics.AddError(
			"cannot change installation_app_type",
			"The Port API does not support changing `installation_app_type` on an existing integration. To use a different app type, destroy this resource (which deletes the integration from Port) and create a new `port_integration`.",
		)
	}
}

func (r *IntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("installation_id"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
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

	err = r.refreshIntegrationState(state, a, integrationIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh integration state", err.Error())
		return
	}

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

	integration, err := integrationToPortBody(plan)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert integration to port body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, integrationIdentifier, integration)

	if err != nil {
		resp.Diagnostics.AddError("failed to update integration", err.Error())
		return
	}

	err = r.refreshIntegrationState(plan, updated, integrationIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh integration state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *IntegrationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	integrationIdentifier := state.InstallationId.ValueString()

	_, err := r.portClient.DeleteIntegration(ctx, integrationIdentifier)

	if err != nil {
		resp.Diagnostics.AddError("failed to delete integration", err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *IntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := integrationToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert integration to port body", err.Error())
		return
	}

	created, err := r.portClient.CreateIntegration(ctx, integration)

	if err != nil {
		resp.Diagnostics.AddError("failed to create integration", err.Error())
		return
	}

	err = r.refreshIntegrationState(state, created, created.InstallationId)

	if err != nil {
		resp.Diagnostics.AddError("failed to create integration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
