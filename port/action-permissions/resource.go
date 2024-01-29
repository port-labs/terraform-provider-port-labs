package action_permissions

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"strings"
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
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_identifier>:<action_identifier>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint_identifier"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("action_identifier"), idParts[1])...)
}

func (r *ActionPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ActionPermissionsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()
	actionIdentifier := state.ActionIdentifier.ValueString()

	a, statusCode, err := r.portClient.GetActionPermissions(ctx, blueprintIdentifier, actionIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to read action permissions", err.Error())
		return
	}

	if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = refreshActionPermissionsState(ctx, state, a, blueprintIdentifier, actionIdentifier)
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

	permissions, err := actionPermissionsToPortBody(state.Permissions)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert action permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateActionPermissions(ctx, blueprintIdentifier, actionIdentifier, permissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update action permissions", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintIdentifier, actionIdentifier))
	state.ActionIdentifier = types.StringValue(actionIdentifier)
	state.BlueprintIdentifier = types.StringValue(blueprintIdentifier)

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

	permissions, err := actionPermissionsToPortBody(state.Permissions)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert action permissions to port body", err.Error())
		return
	}

	_, err = r.portClient.UpdateActionPermissions(ctx, blueprintIdentifier, actionIdentifier, permissions)

	if err != nil {
		resp.Diagnostics.AddError("failed to update action permissions", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintIdentifier, actionIdentifier))
	state.ActionIdentifier = types.StringValue(actionIdentifier)
	state.BlueprintIdentifier = types.StringValue(blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
