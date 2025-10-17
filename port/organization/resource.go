package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &OrganizationSecretResource{}
var _ resource.ResourceWithImportState = &OrganizationSecretResource{}

func NewOrganizationSecretResource() resource.Resource {
	return &OrganizationSecretResource{}
}

type OrganizationSecretResource struct {
	portClient *cli.PortClient
}

func (r *OrganizationSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_secret"
}

func (r *OrganizationSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *OrganizationSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *OrganizationSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secretName := state.SecretName.ValueString()
	secret, statusCode, err := r.portClient.ReadOrganizationSecret(ctx, secretName)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read organization secret", err.Error())
		return
	}

	err = refreshOrganizationSecretState(ctx, state, secret)
	if err != nil {
		resp.Diagnostics.AddError("failed writing organization secret fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *OrganizationSecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := organizationSecretResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert organization secret resource to body", err.Error())
		return
	}

	s, err := r.portClient.CreateOrganizationSecret(ctx, secret)
	if err != nil {
		resp.Diagnostics.AddError("failed to create organization secret", err.Error())
		return
	}

	state.ID = types.StringValue(s.SecretName)
	state.SecretName = types.StringValue(s.SecretName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *OrganizationSecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := organizationSecretResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert organization secret resource to body", err.Error())
		return
	}

	secretName := state.SecretName.ValueString()
	s, err := r.portClient.UpdateOrganizationSecret(ctx, secretName, secret)
	if err != nil {
		resp.Diagnostics.AddError("failed to update organization secret", err.Error())
		return
	}

	state.ID = types.StringValue(s.SecretName)
	state.SecretName = types.StringValue(s.SecretName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *OrganizationSecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteOrganizationSecret(ctx, state.SecretName.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete organization secret", err.Error())
		return
	}
}

func (r *OrganizationSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("secret_name"), req.ID,
	)...)
}
