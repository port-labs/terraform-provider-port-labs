package entity

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EntityResource{}
var _ resource.ResourceWithImportState = &EntityResource{}

func NewEntityResource() resource.Resource {
	return &EntityResource{}
}

type EntityResource struct {
	portClient *cli.PortClient
}

func (r *EntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity"
}

func (r *EntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *EntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntityModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.Blueprint.ValueString()
	e, statusCode, err := r.portClient.ReadEntity(ctx, state.Identifier.ValueString(), state.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read entity", err.Error())
		return
	}
	b, _, err := r.portClient.ReadBlueprint(ctx, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	err = refreshEntityState(ctx, state, e, b)
	if err != nil {
		resp.Diagnostics.AddError("failed writing entity fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !state.RunID.IsNull() {
		runID = state.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	writeEntityComputedFieldsToState(state, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeEntityComputedFieldsToState(state *EntityModel, e *cli.Entity) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", e.Blueprint, e.Identifier))
	state.Identifier = types.StringValue(e.Identifier)
	state.CreatedAt = types.StringValue(e.CreatedAt.String())
	state.CreatedBy = types.StringValue(e.CreatedBy)
	state.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(e.UpdatedBy)
}

func (r *EntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *EntityModel
	var previousState *EntityModel

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

	e, err := entityResourceToBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !state.RunID.IsNull() {
		runID = state.RunID.ValueString()
	}

	var en *cli.Entity

	isBlueprintChanged := !previousState.Blueprint.IsNull() && previousState.Blueprint.ValueString() != state.Blueprint.ValueString()

	if previousState.Identifier.IsNull() || isBlueprintChanged {
		en, err = r.portClient.CreateEntity(ctx, e, runID)
	} else {
		en, err = r.portClient.UpdateEntity(ctx, previousState.Identifier.ValueString(), previousState.Blueprint.ValueString(), e, runID)
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to update entity", err.Error())
		return
	}

	if isBlueprintChanged {
		// Delete the old entity
		err := r.portClient.DeleteEntity(ctx, previousState.Identifier.ValueString(), previousState.Blueprint.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("failed to delete entity", err.Error())
			return
		}
	}

	writeEntityComputedFieldsToState(state, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteEntity(ctx, state.Identifier.ValueString(), state.Blueprint.ValueString(), state.DeleteDependents.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete entity", err.Error())
		return
	}

}

func (r *EntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_id>:<entity_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}
