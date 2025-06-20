package blueprint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
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

	err = r.refreshBlueprintState(ctx, state, b)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintResource) refreshBlueprintState(ctx context.Context, bm *BlueprintModel, b *cli.Blueprint) error {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.ID = types.StringValue(b.Identifier)
	bm.CreatedAt = types.StringValue(b.CreatedAt.String())
	bm.CreatedBy = types.StringValue(b.CreatedBy)
	bm.UpdatedAt = types.StringValue(b.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(b.UpdatedBy)

	if bm.CreateCatalogPage.IsNull() {
		// backwards compatibility, if the field is not set, we assume that the user wants to create a catalog page
		bm.CreateCatalogPage = types.BoolValue(true)
	}

	bm.Title = types.StringValue(b.Title)
	bm.Icon = flex.GoStringToFramework(b.Icon)
	bm.Description = flex.GoStringToFramework(b.Description)

	if bm.ForceDeleteEntities.IsNull() {
		bm.ForceDeleteEntities = types.BoolValue(false)
	}

	if b.ChangelogDestination != nil {
		if b.ChangelogDestination.Type == consts.Kafka {
			bm.KafkaChangelogDestination, _ = types.ObjectValue(nil, nil)
		} else {
			bm.WebhookChangelogDestination = &WebhookChangelogDestinationModel{
				Url: types.StringValue(b.ChangelogDestination.Url),
			}
			if b.ChangelogDestination.Agent != nil {
				bm.WebhookChangelogDestination.Agent = types.BoolValue(*b.ChangelogDestination.Agent)
			}
		}
	}
	if b.TeamInheritance != nil {
		bm.TeamInheritance = &TeamInheritanceModel{
			Path: types.StringValue(b.TeamInheritance.Path),
		}
	}

	if b.Ownership != nil {
		bm.Ownership = &OwnershipModel{
			Type: types.StringValue(b.Ownership.Type),
		}
		if b.Ownership.Path != nil {
			bm.Ownership.Path = types.StringValue(*b.Ownership.Path)
		}
		if b.Ownership.Title != nil {
			bm.Ownership.Title = types.StringValue(*b.Ownership.Title)
		}
	}
	if b.Ownership == nil && bm.Ownership != nil {
		bm.Ownership = nil
	}

	if len(b.Schema.Properties) > 0 {
		err := r.updatePropertiesToState(ctx, b, bm)
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
		addCalculationPropertiesToState(ctx, b, bm)
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

	createCatalogPage := state.CreateCatalogPage.ValueBoolPointer()

	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	bp, err := r.portClient.CreateBlueprint(ctx, b, createCatalogPage)
	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToState(state, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeBlueprintComputedFieldsToState(state *BlueprintModel, bp *cli.Blueprint) {
	state.ID = types.StringValue(bp.Identifier)
	state.Identifier = types.StringValue(bp.Identifier)
	state.CreatedAt = types.StringValue(bp.CreatedAt.String())
	state.CreatedBy = types.StringValue(bp.CreatedBy)
	state.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(bp.UpdatedBy)

	if state.ForceDeleteEntities.IsNull() {
		state.ForceDeleteEntities = types.BoolValue(false)
	}
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *BlueprintModel
	var previousState *BlueprintModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prevB, err := blueprintResourceToPortRequest(ctx, previousState)
	if err != nil {
		resp.Diagnostics.AddError("failed to transform previous state into a blueprint", err.Error())
		return
	}

	b, err := blueprintResourceToPortRequest(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to transform blueprint", err.Error())
		return
	}

	var bp *cli.Blueprint
	if previousState.Identifier.IsNull() {
		bp, err = r.portClient.CreateBlueprint(ctx, b, state.CreateCatalogPage.ValueBoolPointer())
		if err != nil {
			resp.Diagnostics.AddError("failed to create blueprint", err.Error())
			return
		}
	} else {
		var existingBp *cli.Blueprint
		var statusCode int
		existingBp, statusCode, err = r.portClient.ReadBlueprint(ctx, previousState.Identifier.ValueString())
		if err != nil {
			if statusCode == 404 {
				resp.Diagnostics.AddError("Blueprint doesn't exists, it is required to update the blueprint", err.Error())
				return
			}
			resp.Diagnostics.AddError("failed reading blueprint", err.Error())
			return
		}
		// aggregation properties are managed in a different resource, so we need to keep them in the update
		// to avoid losing them
		b.AggregationProperties = existingBp.AggregationProperties
		prevB.AggregationProperties = existingBp.AggregationProperties

		propsWithChangedTypes := make(map[string]string, 0)
		for propKey, prop := range b.Schema.Properties {
			if prevProp, prevHasProp := prevB.Schema.Properties[propKey]; prevHasProp && prop.Type != prevProp.Type {
				propsWithChangedTypes[propKey] = prevProp.Type
				delete(prevB.Schema.Properties, propKey)
			}
		}
		if len(propsWithChangedTypes) > 0 && r.portClient.BlueprintPropertyTypeChangeProtection {
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
			_, err = r.portClient.UpdateBlueprint(ctx, prevB, previousState.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("failed to pre-delete properties that changed their type", err.Error())
				return
			}
		}

		bp, err = r.portClient.UpdateBlueprint(ctx, b, previousState.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to update blueprint", err.Error())
			return
		}
	}

	writeBlueprintComputedFieldsToState(state, bp)

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

	// if deletion protection is not set, this means that the user destroyed the resource, right after upgrading to a version that supports deletion protection
	// therefor we want to be backwards compatible and assume that the user want to have deletion protection
	forceDeleteEntities := state.ForceDeleteEntities.ValueBool()

	if !forceDeleteEntities {
		err := r.portClient.DeleteBlueprint(ctx, state.Identifier.ValueString())
		if err != nil {
			if strings.Contains(err.Error(), "has_dependents") {
				resp.Diagnostics.AddError("failed to delete blueprint", fmt.Sprintf(`Blueprint %s has dependant entities that aren't managed by terraform, if you still wish to destroy the blueprint and delete all entities, set the force_delete_entities argument to true`, state.Identifier.ValueString()))
				return
			}
			resp.Diagnostics.AddError("failed to delete blueprint", err.Error())
			return
		}
	} else {
		forceDeleteBlueprint(ctx, r.portClient, state, resp)
	}

}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func forceDeleteBlueprint(ctx context.Context, portClient *cli.PortClient, state *BlueprintModel, resp *resource.DeleteResponse) {
	migrationId, err := portClient.DeleteBlueprintWithAllEntities(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete blueprint", err.Error())
		return
	}
	// query migration status until status is SUCCESS or FAILED
	for {
		migration, err := portClient.GetMigration(ctx, *migrationId)
		if err != nil {
			resp.Diagnostics.AddError("failed to get migration status", err.Error())
			return
		}
		if migration.Status == consts.Failure {
			resp.Diagnostics.AddError("failed to delete blueprint", "migration failed")
			return
		}
		if migration.Status == consts.Cancelled {
			resp.Diagnostics.AddError("failed to delete blueprint", "migration was cancelled")
			return
		}
		if migration.Status == consts.Completed {
			tflog.Info(ctx, "Migration completed successfully", map[string]interface{}{
				"migration_id": migration.Id,
			})
			break
		}
		if err != nil {
			resp.Diagnostics.AddError("failed to get migration status", err.Error())
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func blueprintResourceToPortRequest(ctx context.Context, state *BlueprintModel) (*cli.Blueprint, error) {
	b := &cli.Blueprint{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	if !state.Icon.IsNull() {
		iconValue := state.Icon.ValueString()
		b.Icon = &iconValue
	}

	if !state.Description.IsNull() {
		descriptionTest := state.Description.ValueString()
		b.Description = &descriptionTest
	}

	if !state.KafkaChangelogDestination.IsNull() {
		b.ChangelogDestination = &cli.ChangelogDestination{
			Type: consts.Kafka,
		}
	}

	if state.WebhookChangelogDestination != nil {
		b.ChangelogDestination = &cli.ChangelogDestination{
			Type: consts.Webhook,
			Url:  state.WebhookChangelogDestination.Url.ValueString(),
		}
		if !state.WebhookChangelogDestination.Agent.IsNull() {
			agent := state.WebhookChangelogDestination.Agent.ValueBool()
			b.ChangelogDestination.Agent = &agent
		}
	}

	if state.TeamInheritance != nil {
		b.TeamInheritance = &cli.TeamInheritance{
			Path: state.TeamInheritance.Path.ValueString(),
		}
	}

	if state.Ownership != nil && !state.Ownership.Type.IsNull() {
		ownershipType := state.Ownership.Type.ValueString()
		if ownershipType == "Inherited" && state.Ownership.Path.IsNull() {
			return nil, fmt.Errorf("path is required when ownership type is Inherited")
		}

		ownership := &cli.Ownership{
			Type: ownershipType,
		}
		if !state.Ownership.Path.IsNull() {
			path := state.Ownership.Path.ValueString()
			ownership.Path = &path
		}
		if !state.Ownership.Title.IsNull() {
			title := state.Ownership.Title.ValueString()
			ownership.Title = &title
		}
		b.Ownership = ownership
	}

	required := []string{}
	props := map[string]cli.BlueprintProperty{}
	var err error
	if state.Properties != nil {
		props, required, err = PropsResourceToBody(ctx, state.Properties)
		if err != nil {
			return nil, err
		}
	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = RelationsResourceToBody(state.Relations)
	b.MirrorProperties = MirrorPropertiesToBody(state.MirrorProperties)
	b.CalculationProperties = CalculationPropertiesToBody(ctx, state.CalculationProperties)
	if err != nil {
		return nil, err
	}
	return b, nil
}
