package integration

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// SaaS-provisioned integrations (installation_type = Saas) run through an async pipeline in
// Port (the integration record moves Creating/Updating/Deleting -> a terminal status while the
// underlying Ocean workload is rolled out/torn down). We poll until that settles so `apply`/
// `destroy` only report success once the integration is actually usable/gone, mirroring the
// polling Ocean itself does for provisioned defaults.
const (
	saasStatusPollInterval = 5 * time.Second
	saasStatusPollMaxTries = 60 // ~5 minutes
)

func isSaasInstallation(state *IntegrationModel) bool {
	return state.InstallationType.ValueString() == consts.InstallationTypeSaas
}

// waitForSaasIntegrationSettled polls the integration until its status leaves the given
// in-progress state, returning the final integration or an error if it lands on Error or the
// poll times out.
func (r *IntegrationResource) waitForSaasIntegrationSettled(ctx context.Context, installationId string, inProgressStatus string) (*cli.Integration, error) {
	for i := 0; i < saasStatusPollMaxTries; i++ {
		integration, statusCode, err := r.portClient.GetIntegration(ctx, installationId)
		if err != nil {
			if statusCode == http.StatusNotFound && inProgressStatus == consts.IntegrationStatusDeleting {
				return nil, nil
			}
			return nil, err
		}

		if integration.StatusInfo == nil || integration.StatusInfo.IntegrationStatus.Status != inProgressStatus {
			if integration.StatusInfo != nil && integration.StatusInfo.IntegrationStatus.Status == consts.IntegrationStatusError {
				message := "no additional details"
				if integration.StatusInfo.IntegrationStatus.Message != nil {
					message = *integration.StatusInfo.IntegrationStatus.Message
				}
				return nil, fmt.Errorf("integration %q entered Error status: %s", installationId, message)
			}
			return integration, nil
		}

		time.Sleep(saasStatusPollInterval)
	}

	return nil, fmt.Errorf("timed out waiting for integration %q to leave %q status after %s", installationId, inProgressStatus, saasStatusPollMaxTries*saasStatusPollInterval)
}

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
		if statusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read integration", err.Error())
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
	var state *IntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationIdentifier := state.InstallationId.ValueString()

	integration, err := integrationToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert integration to port body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, integrationIdentifier, integration)

	if err != nil {
		resp.Diagnostics.AddError("failed to update integration", err.Error())
		return
	}

	if isSaasInstallation(state) {
		updated, err = r.waitForSaasIntegrationSettled(ctx, integrationIdentifier, consts.IntegrationStatusUpdating)
		if err != nil {
			resp.Diagnostics.AddError("failed to update integration", err.Error())
			return
		}
	}

	err = r.refreshIntegrationState(state, updated, integrationIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh integration state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	if isSaasInstallation(state) {
		if _, err := r.waitForSaasIntegrationSettled(ctx, integrationIdentifier, consts.IntegrationStatusDeleting); err != nil {
			resp.Diagnostics.AddError("failed to delete integration", err.Error())
			return
		}
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

	if isSaasInstallation(state) {
		created, err = r.waitForSaasIntegrationSettled(ctx, created.InstallationId, consts.IntegrationStatusCreating)
		if err != nil {
			resp.Diagnostics.AddError("failed to create integration", err.Error())
			return
		}
	}

	err = r.refreshIntegrationState(state, created, created.InstallationId)

	if err != nil {
		resp.Diagnostics.AddError("failed to create integration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
