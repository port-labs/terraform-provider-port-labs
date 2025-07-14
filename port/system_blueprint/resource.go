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

func (r *Resource) refreshBlueprintState(ctx context.Context, bm *SystemBlueprintModel, b *cli.Blueprint, systemBp *cli.Blueprint) error {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.ID = types.StringValue(b.Identifier)

	if len(b.Schema.Properties)-len(systemBp.Schema.Properties) > 0 {
		err := r.updatePropertiesToState(ctx, b, systemBp, bm)
		if err != nil {
			return err
		}
	}

	if len(b.Relations)-len(systemBp.Relations) > 0 {
		addRelationsToState(b, systemBp, bm)
	}

	if len(b.MirrorProperties)-len(systemBp.MirrorProperties) > 0 {
		addMirrorPropertiesToState(b, systemBp, bm)
	}

	if len(b.CalculationProperties)-len(systemBp.CalculationProperties) > 0 {
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

	b, statusCode, err := r.client.ReadBlueprint(ctx, state.Identifier.ValueString())
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

	systemBp, statusCode, err := r.client.ReadSystemBlueprintStructure(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("System blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading system blueprint", err.Error())
		return
	}

	err = r.refreshBlueprintState(ctx, state, b, systemBp)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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

	err = r.refreshBlueprintState(ctx, state, b, systemBp)
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

	structure, statusCode, err := r.client.ReadSystemBlueprintStructure(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	prevB, err := r.mergeSystemBlueprint(ctx, previousState, existingBp, structure)
	if err != nil {
		resp.Diagnostics.AddError("Failed to merge in remote to previous state", err.Error())
		return
	}

	b, err := r.mergeSystemBlueprint(ctx, state, existingBp, structure)
	if err != nil {
		resp.Diagnostics.AddError("Failed to merge in remote to current state", err.Error())
		return
	}

	propsWithChangedTypes := make(map[string]string, 0)
	for propKey, prop := range b.Schema.Properties {
		if prevProp, prevHasProp := prevB.Schema.Properties[propKey]; prevHasProp && prop.Type != prevProp.Type {
			propsWithChangedTypes[propKey] = prevProp.Type
			delete(prevB.Schema.Properties, propKey)
		}
	}
	if len(propsWithChangedTypes) > 0 && r.client.BlueprintPropertyTypeChangeProtection {
		for propKey, prevPropType := range propsWithChangedTypes {
			currentPropType := b.Schema.Properties[propKey].Type
			resp.Diagnostics.AddAttributeError(
				path.Root("properties").AtName(fmt.Sprintf("%s_props", currentPropType)).
					AtName(propKey).AtName("type"),
				"Property type changed while protection is enabled",
				fmt.Sprintf("The type of property %q changed from %q to %q. Applying this change will cause "+
					"you to lose the data for that property. If you wish to continue disable the protection in the "+
					"provider configuration by setting %q to false", propKey, prevPropType, currentPropType,
					"blueprint_property_type_change_protection"),
			)
		}
		return
	}
	if len(propsWithChangedTypes) > 0 {
		_, err = r.client.UpdateBlueprint(ctx, prevB, previousState.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to pre-delete properties that changed their type", err.Error())
			return
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

func (r *Resource) mergeSystemBlueprint(ctx context.Context, state *SystemBlueprintModel, existingBp, structure *cli.Blueprint) (*cli.Blueprint, error) {

	// For system blueprints, we merge properties, relations, mirror properties and calculation properties
	// Everything else should be preserved exactly as is
	props, _, err := MergeProperties(ctx, structure.Schema.Properties, state.Properties)
	if err != nil {
		return nil, fmt.Errorf("error merging properties: %w", err)
	}

	relations := MergeRelations(structure.Relations, state.Relations)
	mirrorProps := MergeMirrorProperties(structure.MirrorProperties, state.MirrorProperties)
	calcProps := MergeCalculationProperties(ctx, structure.CalculationProperties, state.CalculationProperties)

	b := &cli.Blueprint{
		Identifier:           existingBp.Identifier,
		Title:                existingBp.Title,
		Icon:                 existingBp.Icon,
		Description:          existingBp.Description,
		TeamInheritance:      existingBp.TeamInheritance,
		ChangelogDestination: existingBp.ChangelogDestination,
		Schema: cli.BlueprintSchema{
			Properties: props,
			Required:   structure.Schema.Required,
		},
		Relations:             relations,
		MirrorProperties:      mirrorProps,
		CalculationProperties: calcProps,
		AggregationProperties: existingBp.AggregationProperties,
	}

	if existingBp.Ownership != nil {
		b.Ownership = &cli.Ownership{
			Type:  existingBp.Ownership.Type,
			Path:  existingBp.Ownership.Path,
			Title: existingBp.Ownership.Title,
		}
	}

	return b, err
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
