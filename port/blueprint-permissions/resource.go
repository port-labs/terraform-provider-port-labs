package blueprint_permissions

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &BlueprintPermissionsResource{}
var _ resource.ResourceWithImportState = &BlueprintPermissionsResource{}

func NewBlueprintPermissionsResource() resource.Resource {
	return &BlueprintPermissionsResource{}
}

type BlueprintPermissionsResource struct {
	portClient *cli.PortClient
}

func (r *BlueprintPermissionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_permissions"
}

func (r *BlueprintPermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *BlueprintPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("blueprint_identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func (r *BlueprintPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *BlueprintPermissionsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()

	a, statusCode, err := r.portClient.GetBlueprintPermissions(ctx, blueprintIdentifier)

	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint permissions", err.Error())
		return
	}

	if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = refreshBlueprintPermissionsState(state, a, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh blueprint permissions state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *BlueprintPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *BlueprintPermissionsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()

	blueprintPermissions, err := blueprintPermissionsToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert blueprint permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateBlueprintPermissions(ctx, blueprintIdentifier, blueprintPermissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint permissions", err.Error())
		return
	}

	state.ID = types.StringValue(blueprintIdentifier)
	state.BlueprintIdentifier = types.StringValue(blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *BlueprintPermissionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// blueprintPermissions is not deletable resource by itself, as it is tied to a blueprint and is created by default when a blueprint is created
	resp.State.RemoveResource(ctx)
}

func (r *BlueprintPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *BlueprintPermissionsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()

	blueprintPermissions, err := blueprintPermissionsToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert blueprint permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateBlueprintPermissions(ctx, blueprintIdentifier, blueprintPermissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint permissions", err.Error())
		return
	}

	state.ID = types.StringValue(blueprintIdentifier)
	state.BlueprintIdentifier = types.StringValue(blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
