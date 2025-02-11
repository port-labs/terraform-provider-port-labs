package system_blueprint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func writeBlueprintComputedFieldsToState(b *cli.Blueprint, state *SystemBlueprintModel) {
	state.ID = types.StringValue(b.Identifier)
	state.Identifier = types.StringValue(b.Identifier)
}

func refreshBlueprintState(ctx context.Context, bm *SystemBlueprintModel, b *cli.Blueprint, systemBp *cli.Blueprint) error {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.ID = types.StringValue(b.Identifier)

	if len(b.Schema.Properties) - len(systemBp.Schema.Properties) > 0 {
		err := updatePropertiesToState(ctx, b, systemBp, bm)
		if err != nil {
			return err
		}
	}

	if len(b.Relations) - len(systemBp.Relations) > 0 {
		addRelationsToState(b, systemBp, bm)
	}

	if len(b.MirrorProperties) - len(systemBp.MirrorProperties) > 0 {
		addMirrorPropertiesToState(b, systemBp, bm)
	}

	if len(b.CalculationProperties) - len(systemBp.CalculationProperties) > 0 {
		addCalculationPropertiesToState(ctx, b, systemBp, bm)
	}

	return nil
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *SystemBlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError(
		"System Blueprint Creation Not Supported",
		fmt.Sprintf("System blueprints cannot be created. To manage the system blueprint '%s', please import it using:\n\nterraform import port_system_blueprint.%s %s", 
			state.Identifier.ValueString(),
			state.Identifier.ValueString(),
			state.Identifier.ValueString(),
		),
	)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *SystemBlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.client.ReadBlueprint(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Blueprint", err.Error())
		return
	}

	systemBp, statusCode, err := r.client.ReadSystemBlueprintStructure(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("System blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading system blueprint", err.Error())
		return
	}

	err = refreshBlueprintState(ctx, state, b, systemBp)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *SystemBlueprintModel
	var previousState *SystemBlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	existingBp, statusCode, err := r.client.ReadBlueprint(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	systemBp, statusCode, err := r.client.ReadSystemBlueprintStructure(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	// For system blueprints, we merge properties, relations, mirror properties and calculation properties
	// Everything else should be preserved exactly as is
	props, _, err := MergeProperties(ctx, systemBp.Schema.Properties, state.Properties)
	if err != nil {
		resp.Diagnostics.AddError("Error merging properties", err.Error())
		return
	}

	relations := MergeRelations(systemBp.Relations, state.Relations)
	mirrorProps := MergeMirrorProperties(systemBp.MirrorProperties, state.MirrorProperties)
	calcProps := MergeCalculationProperties(ctx, systemBp.CalculationProperties, state.CalculationProperties)

	b := &cli.Blueprint{
		Identifier:           existingBp.Identifier,
		Title:                existingBp.Title,
		Icon:                 existingBp.Icon,
		Description:          existingBp.Description,
		TeamInheritance:      existingBp.TeamInheritance,
		ChangelogDestination: existingBp.ChangelogDestination,
		Schema: cli.BlueprintSchema{
			Properties: props,
			Required:  systemBp.Schema.Required,
		},
		Relations:             relations,
		MirrorProperties:     mirrorProps,
		CalculationProperties: calcProps,
		AggregationProperties: existingBp.AggregationProperties,
	}
	
	if existingBp.Ownership != nil {
		b.Ownership = &cli.Ownership{
			Type: existingBp.Ownership.Type,
			Path: existingBp.Ownership.Path,
			Title: existingBp.Ownership.Title,
		}
	}

	bp, err := r.client.UpdateBlueprint(ctx, b, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToState(bp, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// System blueprints cannot be deleted, so this is a no-op
	// We just remove it from the state without trying to delete from Port
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_blueprint"
} 