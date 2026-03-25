package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &OrganizationResource{}
var _ resource.ResourceWithImportState = &OrganizationResource{}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

type OrganizationResource struct {
	portClient *cli.PortClient
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *OrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := r.portClient.ReadOrganization(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to read organization", err.Error())
		return
	}

	err = refreshOrganizationState(ctx, state, org)
	if err != nil {
		resp.Diagnostics.AddError("failed writing organization fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *OrganizationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	orgUpdate, err := organizationResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert organization resource to body", err.Error())
		return
	}

	org, _, err := r.portClient.UpdateOrganization(ctx, orgUpdate)
	if err != nil {
		resp.Diagnostics.AddError("failed to update organization", err.Error())
		return
	}

	state.ID = types.StringValue(org.Name)

	err = refreshOrganizationState(ctx, state, org)
	if err != nil {
		resp.Diagnostics.AddError("failed writing organization fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *OrganizationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	orgUpdate, err := organizationResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert organization resource to body", err.Error())
		return
	}

	org, _, err := r.portClient.UpdateOrganization(ctx, orgUpdate)
	if err != nil {
		resp.Diagnostics.AddError("failed to update organization", err.Error())
		return
	}

	err = refreshOrganizationState(ctx, state, org)
	if err != nil {
		resp.Diagnostics.AddError("failed writing organization fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Organization is a singleton and cannot be deleted.
	// Removing from state is sufficient.
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
