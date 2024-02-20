package page_permissions

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &PagePermissionsResource{}
var _ resource.ResourceWithImportState = &PagePermissionsResource{}

func NewPagePermissionsResource() resource.Resource {
	return &PagePermissionsResource{}
}

type PagePermissionsResource struct {
	portClient *cli.PortClient
}

func (r *PagePermissionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page_permissions"
}

func (r *PagePermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *PagePermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("page_identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func (r *PagePermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *PagePermissionsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pageIdentifier := state.PageIdentifier.ValueString()

	a, statusCode, err := r.portClient.GetPagePermissions(ctx, pageIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to read page permissions", err.Error())
		return
	}

	if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = refreshPagePermissionsState(state, a, pageIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh page permissions state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *PagePermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *PagePermissionsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pageIdentifier := state.PageIdentifier.ValueString()

	pagePermissions, err := pagePermissionsToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert page permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdatePagePermissions(ctx, pageIdentifier, pagePermissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update page permissions", err.Error())
		return
	}

	state.ID = types.StringValue(pageIdentifier)
	state.PageIdentifier = types.StringValue(pageIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PagePermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *PagePermissionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// pagePermissions is not deletable resource by itself as it is tied to an page and is created by default when an page is created
	resp.State.RemoveResource(ctx)
}

func (r *PagePermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *PagePermissionsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pageIdentifier := state.PageIdentifier.ValueString()

	pagePermissions, err := pagePermissionsToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert page permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdatePagePermissions(ctx, pageIdentifier, pagePermissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update page permissions", err.Error())
		return
	}

	state.ID = types.StringValue(pageIdentifier)
	state.PageIdentifier = types.StringValue(pageIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
