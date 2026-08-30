package blueprint_relation

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
)

var _ resource.Resource = &BlueprintRelationResource{}
var _ resource.ResourceWithImportState = &BlueprintRelationResource{}

func NewBlueprintRelationResource() resource.Resource {
	return &BlueprintRelationResource{}
}

type BlueprintRelationResource struct {
	portClient *cli.PortClient
}

func (r *BlueprintRelationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_relation"
}

func (r *BlueprintRelationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func resourceId(blueprintIdentifier string, relationIdentifier string) string {
	return fmt.Sprintf("%s:%s", blueprintIdentifier, relationIdentifier)
}

func relationToPortBody(state *BlueprintRelationModel) *cli.Relation {
	target := state.Target.ValueString()
	relation := &cli.Relation{
		Target:   &target,
		Many:     state.Many.ValueBoolPointer(),
		Required: state.Required.ValueBoolPointer(),
	}

	if !state.Title.IsNull() {
		title := state.Title.ValueString()
		relation.Title = &title
	}

	if !state.Description.IsNull() {
		description := state.Description.ValueString()
		relation.Description = &description
	}

	return relation
}

func refreshRelationState(state *BlueprintRelationModel, relation *cli.Relation) {
	state.ID = types.StringValue(resourceId(state.Blueprint.ValueString(), state.Identifier.ValueString()))
	state.Target = flex.GoStringToFramework(relation.Target)
	state.Title = flex.GoStringToFramework(relation.Title)
	state.Description = flex.GoStringToFramework(relation.Description)
	state.Many = flex.GoBoolToFramework(relation.Many)
	state.Required = flex.GoBoolToFramework(relation.Required)

	// many and required are always sent, and Port always returns them, but keep the
	// schema defaults rather than writing nulls if the API ever omits them.
	if state.Many.IsNull() {
		state.Many = types.BoolValue(false)
	}
	if state.Required.IsNull() {
		state.Required = types.BoolValue(false)
	}
}

func (r *BlueprintRelationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *BlueprintRelationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	relation, ok := b.Relations[state.Identifier.ValueString()]
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	refreshRelationState(state, &relation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintRelationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *BlueprintRelationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.Blueprint.ValueString()
	relationIdentifier := state.Identifier.ValueString()

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, blueprintIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exist, it is required to create a relation", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	// A relation must be owned by exactly one place. Refusing here surfaces an overlap with an
	// inline relations block on port_blueprint, instead of the two resources silently fighting.
	if _, ok := b.Relations[relationIdentifier]; ok {
		resp.Diagnostics.AddError(
			"relation already exists",
			fmt.Sprintf("Relation %q already exists on blueprint %q. Remove it from the blueprint's relations block, "+
				"or import it with `terraform import <address> %q`.",
				relationIdentifier, blueprintIdentifier, resourceId(blueprintIdentifier, relationIdentifier)),
		)
		return
	}

	if _, err := r.portClient.PatchBlueprintRelation(ctx, blueprintIdentifier, relationIdentifier, relationToPortBody(state)); err != nil {
		resp.Diagnostics.AddError("failed to create relation", err.Error())
		return
	}

	state.ID = types.StringValue(resourceId(blueprintIdentifier, relationIdentifier))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintRelationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *BlueprintRelationModel
	var previousState *BlueprintRelationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.Blueprint.ValueString()
	relationIdentifier := state.Identifier.ValueString()
	relation := relationToPortBody(state)

	// PATCH deep-merges and the API rejects null for title and description, so a field that was
	// set and is now unset can only be removed by rewriting the whole blueprint.
	clearsField := (state.Title.IsNull() && !previousState.Title.IsNull()) ||
		(state.Description.IsNull() && !previousState.Description.IsNull())

	var err error
	if clearsField {
		_, err = r.portClient.PutBlueprintRelation(ctx, blueprintIdentifier, relationIdentifier, relation)
	} else {
		_, err = r.portClient.PatchBlueprintRelation(ctx, blueprintIdentifier, relationIdentifier, relation)
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to update relation", err.Error())
		return
	}

	state.ID = types.StringValue(resourceId(blueprintIdentifier, relationIdentifier))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintRelationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *BlueprintRelationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	statusCode, err := r.portClient.DeleteBlueprintRelation(ctx, state.Blueprint.ValueString(), state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			// the blueprint is already gone, and the relation with it
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete relation", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *BlueprintRelationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	blueprintIdentifier, relationIdentifier, found := strings.Cut(req.ID, ":")
	if !found || blueprintIdentifier == "" || relationIdentifier == "" {
		resp.Diagnostics.AddError(
			"unexpected import identifier",
			fmt.Sprintf("Expected an identifier of the form `<blueprint>:<relation>`, got %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), blueprintIdentifier)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), relationIdentifier)...)
}
