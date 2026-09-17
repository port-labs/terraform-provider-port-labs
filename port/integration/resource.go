package integration

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		resp.Diagnostics.AddError(
			"config cannot be set on creation",
			"Integrations receive default mappings during provisioning. "+
				"Create the integration first (without config), then add config "+
				"to your HCL and run 'terraform apply' again to override the defaults.",
		)
		return
	}

	if err := validateIntegrationModel(plan); err != nil {
		resp.Diagnostics.AddError("invalid integration configuration", err.Error())
		return
	}

	body, err := integrationToPortBody(plan, true)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	created, err := r.portClient.CreateIntegration(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("failed to create integration", err.Error())
		return
	}

	applyWriteResult(plan, created, created.InstallationId)

	waitReady, waitProvisioned := createAwaitParams(plan)
	if waitReady || waitProvisioned {
		r.awaitInfra(ctx, plan, created.InstallationId, "created", waitReady, waitProvisioned, &resp.Diagnostics)
	}

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

	if err := validateIntegrationModel(plan); err != nil {
		resp.Diagnostics.AddError("invalid integration configuration", err.Error())
		return
	}

	body, err := integrationToPortBody(plan, false)
	if err != nil {
		resp.Diagnostics.AddError("failed to build request body", err.Error())
		return
	}

	updated, err := r.portClient.UpdateIntegration(ctx, integrationIdentifier, body)
	if err != nil {
		resp.Diagnostics.AddError("failed to update integration", err.Error())
		return
	}

	applyWriteResult(plan, updated, integrationIdentifier)

	if plan.isSaas() {
		r.awaitInfra(ctx, plan, integrationIdentifier, "updated", true, false, &resp.Diagnostics)
	}

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

func createAwaitParams(plan *IntegrationModel) (waitReady, waitProvisioned bool) {
	return plan.isSaas(), true
}

// awaitInfra waits for integration operation and/or default resource provisioning to
// finish, then syncs status and version into state when ready polling succeeds.
func (r *IntegrationResource) awaitInfra(ctx context.Context, model *IntegrationModel, installationId, verb string, waitReady, waitProvisioned bool, diags *diag.Diagnostics) {
	var ready *cli.Integration
	var readyErr error
	var provisionedErr error

	var wg sync.WaitGroup
	if waitReady {
		wg.Go(func() {
			ready, readyErr = r.portClient.WaitForIntegrationReady(ctx, installationId)
		})
	}
	if waitProvisioned {
		wg.Go(func() {
			_, provisionedErr = r.portClient.WaitForIntegrationProvisioned(ctx, installationId)
		})
	}
	wg.Wait()

	if waitReady {
		switch {
		case readyErr != nil:
			diags.AddWarning(
				fmt.Sprintf("integration %s but operation did not complete", verb),
				readyErr.Error()+". The integration has been saved to state. Check the Port UI or run 'terraform plan' to inspect.",
			)
		case ready != nil:
			if ready.StatusInfo != nil {
				model.Status = types.StringValue(ready.StatusInfo.IntegrationStatus.Status)
			}
			if ready.Version != nil {
				model.Version = types.StringPointerValue(ready.Version)
			}
		}
	}

	if waitProvisioned && provisionedErr != nil {
		diags.AddWarning(
			fmt.Sprintf("integration %s but resource provisioning did not complete", verb),
			provisionedErr.Error()+". Default mappings may not be available yet. Run 'terraform apply' again or check the Port UI.",
		)
	}
	// ready == nil && readyErr == nil → operation timeout, keep the planned values.
}
