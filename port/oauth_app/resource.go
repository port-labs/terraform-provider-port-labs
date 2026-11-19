package oauth_app

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &OAuthAppResource{}
var _ resource.ResourceWithImportState = &OAuthAppResource{}

func NewOAuthAppResource() resource.Resource {
	return &OAuthAppResource{}
}

type OAuthAppResource struct {
	portClient *cli.PortClient
}

func (r *OAuthAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_app"
}

func (r *OAuthAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *OAuthAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *OAuthAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	oauthApp, err := oauthAppResourceToPortBodyCreate(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert oauth app resource to body", err.Error())
		return
	}

	createdApp, err := r.portClient.CreateOAuthApp(ctx, oauthApp)
	if err != nil {
		resp.Diagnostics.AddError("failed to create oauth app", err.Error())
		return
	}

	err = refreshOAuthAppState(ctx, state, createdApp)
	if err != nil {
		resp.Diagnostics.AddError("failed writing oauth app fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OAuthAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *OAuthAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	app, statusCode, err := r.portClient.ReadOAuthApp(ctx, state.ID.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read oauth app", err.Error())
		return
	}

	existingClientSecret := state.ClientSecret
	err = refreshOAuthAppState(ctx, state, app)
	if err != nil {
		resp.Diagnostics.AddError("failed writing oauth app fields to resource", err.Error())
		return
	}

	if !existingClientSecret.IsNull() && existingClientSecret.ValueString() != "" {
		state.ClientSecret = existingClientSecret
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OAuthAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *OAuthAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	oauthAppUpdate, err := oauthAppResourceToPortBodyUpdate(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert oauth app resource to body", err.Error())
		return
	}

	updatedApp, err := r.portClient.UpdateOAuthApp(ctx, state.ID.ValueString(), oauthAppUpdate)
	if err != nil {
		resp.Diagnostics.AddError("failed to update oauth app", err.Error())
		return
	}

	existingClientSecret := state.ClientSecret
	err = refreshOAuthAppState(ctx, state, updatedApp)
	if err != nil {
		resp.Diagnostics.AddError("failed writing oauth app fields to resource", err.Error())
		return
	}

	if !existingClientSecret.IsNull() && existingClientSecret.ValueString() != "" {
		state.ClientSecret = existingClientSecret
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OAuthAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *OAuthAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteOAuthApp(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete oauth app", err.Error())
		return
	}
}

func (r *OAuthAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
