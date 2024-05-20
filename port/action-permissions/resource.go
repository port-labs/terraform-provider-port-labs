package action_permissions

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &ActionPermissionsResource{}
var _ resource.ResourceWithImportState = &ActionPermissionsResource{}

func NewActionPermissionsResource() resource.Resource {
	return &ActionPermissionsResource{}
}

type ActionPermissionsResource struct {
	portClient *cli.PortClient
}

func (r *ActionPermissionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action_permissions"
}

func (r *ActionPermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ActionPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("action_identifier"), req.ID)...)
}

func (r *ActionPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ActionPermissionsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()
	actionIdentifier := state.ActionIdentifier.ValueString()
	// For the first time a user is migrating from action v1 to v2
	if blueprintIdentifier != "" {
		actionIdentifier = fmt.Sprintf("%s_%s", blueprintIdentifier, actionIdentifier)
	}

	a, statusCode, err := r.portClient.GetActionPermissions(ctx, actionIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read action permissions", err.Error())
		return
	}

	err = refreshActionPermissionsState(state, a, actionIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to refresh action permissions state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *ActionPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *ActionPermissionsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()
	actionIdentifier := state.ActionIdentifier.ValueString()
	if blueprintIdentifier != "" {
		actionIdentifier = fmt.Sprintf("%s_%s", blueprintIdentifier, actionIdentifier)
	}

	permissions, err := actionPermissionsToPortBody(state.Permissions)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert action permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateActionPermissions(ctx, actionIdentifier, permissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update action permissions", err.Error())
		return
	}

	state.ID = types.StringValue(actionIdentifier)
	state.ActionIdentifier = types.StringValue(actionIdentifier)
	state.BlueprintIdentifier = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ActionPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ActionPermissionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// actionPermissions is not deletable resource by itself as it is tied to an action and is created by default when an action is created
	resp.State.RemoveResource(ctx)
}

func (r *ActionPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *ActionPermissionsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()
	actionIdentifier := state.ActionIdentifier.ValueString()
	if blueprintIdentifier != "" {
		actionIdentifier = fmt.Sprintf("%s_%s", blueprintIdentifier, actionIdentifier)
	}

	permissions, err := actionPermissionsToPortBody(state.Permissions)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert action permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateActionPermissions(ctx, actionIdentifier, permissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update action permissions", err.Error())
		return
	}

	state.ID = types.StringValue(actionIdentifier)
	state.ActionIdentifier = types.StringValue(actionIdentifier)
	state.BlueprintIdentifier = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
