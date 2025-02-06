package folder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &FolderResource{}
var _ resource.ResourceWithImportState = &FolderResource{}

func NewFolderResource() resource.Resource {
	return &FolderResource{}
}

type FolderResource struct {
	portClient *cli.PortClient
}

func (r *FolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *FolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *FolderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	f, statusCode, err := r.portClient.GetFolder(ctx,
		state.SidebarIdentifier.ValueString(), state.FolderIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to get folder", err.Error())
		return
	}

	err = refreshFolderToState(state, f)

	if err != nil {
		resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *FolderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	statusCode, err := r.portClient.DeleteFolder(ctx, state.SidebarIdentifier.ValueString(), state.FolderIdentifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete folder", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *FolderModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	folder, err := FolderToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
		return
	}

	f, err := r.portClient.CreateFolder(ctx, folder)
	if err != nil {
		resp.Diagnostics.AddError("failed to create folder", err.Error())
		return
	}
	if f == nil {
		f, _, err = r.portClient.GetFolder(ctx,
			state.SidebarIdentifier.ValueString(), state.FolderIdentifier.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to get folder", err.Error())
			return
		}
	}

	state.FolderIdentifier = types.StringValue(f.FolderIdentifier)
	state.SidebarIdentifier = types.StringValue(f.SidebarIdentifier)
	state.Parent = types.StringValue(f.Parent)
	state.After = types.StringValue(f.After)
	state.Title = types.StringValue(f.Title)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *FolderModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	f, _, err := r.portClient.GetFolder(ctx,
		state.SidebarIdentifier.ValueString(), state.FolderIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get folder", err.Error())
		return
	}

	folder, err := FolderToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
		return
	}

	_, err = r.portClient.UpdateFolder(ctx, folder)

	if err != nil {
		resp.Diagnostics.AddError("failed to update folder", err.Error())
		return
	}

	state.FolderIdentifier = types.StringValue(f.FolderIdentifier)
	state.SidebarIdentifier = types.StringValue(f.SidebarIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
