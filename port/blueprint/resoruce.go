package blueprint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &BlueprintResource{}
var _ resource.ResourceWithImportState = &BlueprintResource{}

func NewBlueprintResource() resource.Resource {
	return &BlueprintResource{}
}

type BlueprintResource struct {
	portClient *cli.PortClient
}

func (r *BlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *BlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *BlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *BlueprintModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	err = refreshBlueprintState(ctx, state, b)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func refreshBlueprintState(ctx context.Context, bm *BlueprintModel, b *cli.Blueprint) error {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.ID = types.StringValue(b.Identifier)
	bm.CreatedAt = types.StringValue(b.CreatedAt.String())
	bm.CreatedBy = types.StringValue(b.CreatedBy)
	bm.UpdatedAt = types.StringValue(b.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(b.UpdatedBy)

	if b.Title != nil {
		bm.Title = types.StringValue(*b.Title)
	}

	if b.Icon != nil {
		bm.Icon = types.StringValue(*b.Icon)
	}

	if b.Description != nil {
		bm.Description = types.StringValue(*b.Description)
	}

	if b.ChangelogDestination != nil {
		bm.ChangelogDestination = &ChangelogDestinationModel{
			Type:  types.StringValue(b.ChangelogDestination.Type),
			Url:   types.StringValue(b.ChangelogDestination.Url),
			Agent: types.BoolValue(b.ChangelogDestination.Agent),
		}
	}
	if b.TeamInheritance != nil {
		bm.TeamInheritance = &TeamInheritanceModel{
			Path: types.StringValue(b.TeamInheritance.Path),
		}
	}

	if len(b.Schema.Properties) > 0 {
		err := updatePropertiesToState(ctx, b, bm)
		if err != nil {
			return err
		}
	}

	if len(b.Relations) > 0 {
		addRelationsToState(b, bm)
	}

	if len(b.MirrorProperties) > 0 {
		addMirrorPropertiesToState(b, bm)
	}

	if len(b.CalculationProperties) > 0 {
		addCalculationPropertiesToState(b, bm)
	}

	return nil
}

func (r *BlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := blueprintResourceToPortRequest(ctx, state)

	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	bp, err := r.portClient.CreateBlueprint(ctx, b)
	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	state.ID = types.StringValue(bp.Identifier)

	writeBlueprintComputedFieldsToResource(state, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeBlueprintComputedFieldsToResource(state *BlueprintModel, bp *cli.Blueprint) {
	state.Identifier = types.StringValue(bp.Identifier)
	state.CreatedAt = types.StringValue(bp.CreatedAt.String())
	state.CreatedBy = types.StringValue(bp.CreatedBy)
	state.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(bp.UpdatedBy)
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	b, err := blueprintResourceToPortRequest(ctx, state)

	if err != nil {
		resp.Diagnostics.AddError("failed to transform blueprint", err.Error())
		return
	}

	var bp *cli.Blueprint

	if state.Identifier.IsNull() {
		bp, err = r.portClient.CreateBlueprint(ctx, b)
	} else {
		bp, err = r.portClient.UpdateBlueprint(ctx, b, state.Identifier.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToResource(state, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *BlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *BlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Identifier.IsNull() {
		resp.Diagnostics.AddError("failed to extract blueprint identifier", "identifier is required")
		return
	}

	err := r.portClient.DeleteBlueprint(ctx, state.Identifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete blueprint", err.Error())
		return
	}
}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}

func blueprintResourceToPortRequest(ctx context.Context, state *BlueprintModel) (*cli.Blueprint, error) {
	b := &cli.Blueprint{
		Identifier: state.Identifier.ValueString(),
	}

	if !state.Title.IsNull() {
		titleValue := state.Title.ValueString()
		b.Title = &titleValue
	}

	if !state.Icon.IsNull() {
		iconValue := state.Icon.ValueString()
		b.Icon = &iconValue
	}

	if !state.Description.IsNull() {
		descriptionTest := state.Description.ValueString()
		b.Description = &descriptionTest
	}

	if state.ChangelogDestination != nil {
		if state.ChangelogDestination.Type.ValueString() == "KAFKA" && !state.ChangelogDestination.Agent.IsNull() {
			return nil, fmt.Errorf("agent is not supported for Kafka changelog destination")
		}
		b.ChangelogDestination = &cli.ChangelogDestination{}
		b.ChangelogDestination.Type = state.ChangelogDestination.Type.ValueString()
		b.ChangelogDestination.Url = state.ChangelogDestination.Url.ValueString()
		b.ChangelogDestination.Agent = state.ChangelogDestination.Agent.ValueBool()
	} else {
		b.ChangelogDestination = nil
	}

	if state.TeamInheritance != nil {
		b.TeamInheritance = &cli.TeamInheritance{
			Path: state.TeamInheritance.Path.ValueString(),
		}
	} else {
		b.TeamInheritance = nil
	}

	required := []string{}
	props := map[string]cli.BlueprintProperty{}
	var err error
	if state.Properties != nil {
		props, required, err = readStateToPortBody(ctx, state)
		if err != nil {
			return nil, err
		}
	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relationsResourceToBody(state)
	b.MirrorProperties = mirrorPropertiesToBody(state)
	b.CalculationProperties = calculationPropertiesToBody(state)
	return b, nil
}
