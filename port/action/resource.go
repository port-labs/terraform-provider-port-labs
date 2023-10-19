package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithImportState = &ActionResource{}

func NewActionResource() resource.Resource {
	return &ActionResource{}
}

type ActionResource struct {
	portClient *cli.PortClient
}

func (r *ActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action"
}

func (r *ActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_id>:<action_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)

}

func (r *ActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ActionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.Blueprint.ValueString()
	a, statusCode, err := r.portClient.ReadAction(ctx, blueprintIdentifier, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading action", err.Error())
		return
	}

	err = refreshActionState(ctx, state, a, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed writing action fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteAction(ctx, state.Blueprint.ValueString(), state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete action", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *ActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	action, err := actionStateToPortBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert action resource to body", err.Error())
		return
	}

	a, err := r.portClient.CreateAction(ctx, bp.Identifier, action)
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", bp.Identifier, a.Identifier))
	state.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *ActionModel
	var previousState *ActionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	action, err := actionStateToPortBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	var a *cli.Action
	if previousState.ID.IsNull() {
		a, err = r.portClient.CreateAction(ctx, bp.Identifier, action)
	} else {
		a, err = r.portClient.UpdateAction(ctx, bp.Identifier, previousState.Identifier.ValueString(), action)
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", bp.Identifier, a.Identifier))
	state.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
