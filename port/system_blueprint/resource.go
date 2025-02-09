package system_blueprint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

func writeBlueprintComputedFieldsToState(b *cli.Blueprint, state *Model) {
	state.ID = types.StringValue(b.Identifier)
	state.Identifier = types.StringValue(b.Identifier)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to read the existing system blueprint
	b, statusCode, err := r.client.ReadBlueprint(ctx, plan.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError(
				"Unsupported Operation",
				"The system_blueprint resource does not support creation. Please import an existing system blueprint instead.",
			)
			return
		}
		resp.Diagnostics.AddError("Error Reading Blueprint", err.Error())
		return
	}

	// Write computed fields to state
	writeBlueprintComputedFieldsToState(b, &plan)

	if err := blueprint.UpdatePropertiesToState(ctx, b, plan.Properties); err != nil {
		resp.Diagnostics.AddError("Error updating properties", err.Error())
		return
	}

	if err := blueprint.UpdateRelationsToState(b, plan.Relations); err != nil {
		resp.Diagnostics.AddError("Error updating relations", err.Error())
		return
	}

	if err := blueprint.UpdateMirrorPropertiesToState(b, plan.MirrorProperties); err != nil {
		resp.Diagnostics.AddError("Error updating mirror properties", err.Error())
		return
	}

	if err := blueprint.UpdateCalculationPropertiesToState(ctx, b, plan.CalculationProperties); err != nil {
		resp.Diagnostics.AddError("Error updating calculation properties", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *Model
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

	// Write computed fields to state
	writeBlueprintComputedFieldsToState(b, state)

	if err := blueprint.UpdatePropertiesToState(ctx, b, state.Properties); err != nil {
		resp.Diagnostics.AddError("Error updating properties", err.Error())
		return
	}

	if err := blueprint.UpdateRelationsToState(b, state.Relations); err != nil {
		resp.Diagnostics.AddError("Error updating relations", err.Error())
		return
	}

	if err := blueprint.UpdateMirrorPropertiesToState(b, state.MirrorProperties); err != nil {
		resp.Diagnostics.AddError("Error updating mirror properties", err.Error())
		return
	}

	if err := blueprint.UpdateCalculationPropertiesToState(ctx, b, state.CalculationProperties); err != nil {
		resp.Diagnostics.AddError("Error updating calculation properties", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingBp, statusCode, err := r.client.ReadBlueprint(ctx, plan.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	fmt.Printf("Existing blueprint from Port:\n%+v\n", existingBp)
	fmt.Printf("Existing blueprint schema:\n%+v\n", existingBp.Schema)
	fmt.Printf("Existing blueprint ownership:\n%+v\n", existingBp.Ownership)

	// For system blueprints, we merge properties, relations, mirror properties and calculation properties
	// Everything else should be preserved exactly as is
	props, _, err := blueprint.MergeProperties(ctx, existingBp.Schema.Properties, plan.Properties)
	if err != nil {
		resp.Diagnostics.AddError("Error merging properties", err.Error())
		return
	}

	relations := blueprint.MergeRelations(existingBp.Relations, plan.Relations)
	mirrorProps := blueprint.MergeMirrorProperties(existingBp.MirrorProperties, plan.MirrorProperties)
	calcProps := blueprint.MergeCalculationProperties(ctx, existingBp.CalculationProperties, plan.CalculationProperties)

	// Create update request with ALL fields from the existing blueprint
	b := &cli.Blueprint{
		Identifier:            existingBp.Identifier,
		Title:                existingBp.Title,
		Icon:                 existingBp.Icon,
		Description:          existingBp.Description,
		TeamInheritance:      existingBp.TeamInheritance,
		ChangelogDestination: existingBp.ChangelogDestination,
		Schema: cli.BlueprintSchema{
			Properties: props,
			Required:  existingBp.Schema.Required,
		},
		Relations:             relations,
		MirrorProperties:     mirrorProps,
		CalculationProperties: calcProps,
		AggregationProperties: existingBp.AggregationProperties,
	}

	// Handle ownership - preserve the existing ownership
	if existingBp.Ownership != nil {
		b.Ownership = &cli.Ownership{
			Type: existingBp.Ownership.Type,
			Path: existingBp.Ownership.Path,
			Title: existingBp.Ownership.Title,
		}
	}

	fmt.Printf("Update blueprint request:\n%+v\n", b)
	fmt.Printf("Update blueprint schema:\n%+v\n", b.Schema)
	fmt.Printf("Update blueprint ownership:\n%+v\n", b.Ownership)

	bp, err := r.client.UpdateBlueprint(ctx, b, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint", err.Error())
		return
	}

	// Write computed fields to state
	writeBlueprintComputedFieldsToState(bp, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// System blueprints cannot be deleted, so this is a no-op
	// We just remove it from the state without trying to delete from Port
	return
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_blueprint"
} 