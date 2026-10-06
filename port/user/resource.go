package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResource struct {
	portClient *cli.PortClient
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *UserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, statusCode, err := r.portClient.ReadUser(ctx, state.Email.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read user", err.Error())
		return
	}

	err = refreshUserState(ctx, state, user)
	if err != nil {
		resp.Diagnostics.AddError("failed writing user fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	email := state.Email.ValueString()
	_, statusCode, err := r.portClient.ReadUser(ctx, email)
	if err != nil && statusCode != 404 {
		resp.Diagnostics.AddError("failed to check if user exists", err.Error())
		return
	}

	if statusCode == 404 {
		err = r.portClient.InviteUser(ctx, userInviteFromState(state), false)
		if err != nil {
			resp.Diagnostics.AddError("failed to invite user", err.Error())
			return
		}
	}

	update, err := userResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert user resource to body", err.Error())
		return
	}

	if update.IncludeRoles || update.IncludeTeams || update.IncludeInactivityTimeout {
		_, err = r.portClient.UpdateUser(ctx, email, update)
		if err != nil {
			resp.Diagnostics.AddError("failed to update user", err.Error())
			return
		}
	}

	user, _, err := r.portClient.ReadUser(ctx, email)
	if err != nil {
		resp.Diagnostics.AddError("failed to read user after create", err.Error())
		return
	}

	err = refreshUserState(ctx, state, user)
	if err != nil {
		resp.Diagnostics.AddError("failed writing user fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	update, err := userResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert user resource to body", err.Error())
		return
	}

	email := state.Email.ValueString()
	var user *cli.User
	if update.IncludeRoles || update.IncludeTeams || update.IncludeInactivityTimeout {
		user, err = r.portClient.UpdateUser(ctx, email, update)
		if err != nil {
			resp.Diagnostics.AddError("failed to update user", err.Error())
			return
		}
	} else {
		user, _, err = r.portClient.ReadUser(ctx, email)
		if err != nil {
			resp.Diagnostics.AddError("failed to read user", err.Error())
			return
		}
	}

	err = refreshUserState(ctx, state, user)
	if err != nil {
		resp.Diagnostics.AddError("failed writing user fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteUser(ctx, state.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete user", err.Error())
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("email"), req.ID,
	)...)
}
