package blueprint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/samber/lo"
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
		err := addPropertiesToResource(ctx, b, bm)
		if err != nil {
			return err
		}
	}

	if len(b.Relations) > 0 {
		addRelationsToResource(b, bm)
	}

	if len(b.MirrorProperties) > 0 {
		addMirrorPropertiesToResource(b, bm)
	}

	if len(b.CalculationProperties) > 0 {
		addCalculationPropertiesToResource(b, bm)
	}

	return nil
}

func addRelationsToResource(b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.Relations {
		if bm.Relations == nil {
			bm.Relations = make(map[string]RelationModel)
		}

		relationModel := &RelationModel{
			Target: types.StringValue(*v.Target),
		}

		if v.Title != nil {
			relationModel.Title = types.StringValue(*v.Title)
		}

		if v.Many != nil {
			relationModel.Many = types.BoolValue(*v.Many)
		}

		if v.Required != nil {
			relationModel.Required = types.BoolValue(*v.Required)
		}

		bm.Relations[k] = *relationModel

	}
}

func addMirrorPropertiesToResource(b *cli.Blueprint, bm *BlueprintModel) {
	if b.MirrorProperties != nil {
		for k, v := range b.MirrorProperties {
			if bm.MirrorProperties == nil {
				bm.MirrorProperties = make(map[string]MirrorPropertyModel)
			}

			mirrorPropertyModel := &MirrorPropertyModel{
				Path: types.StringValue(v.Path),
			}
			if v.Title != "" {
				mirrorPropertyModel.Title = types.StringValue(v.Title)
			}

			bm.MirrorProperties[k] = *mirrorPropertyModel

		}
	}
}

func addCalculationPropertiesToResource(b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.CalculationProperties {
		if bm.CalculationProperties == nil {
			bm.CalculationProperties = make(map[string]CalculationPropertyModel)
		}

		calculationPropertyModel := &CalculationPropertyModel{
			Calculation: types.StringValue(v.Calculation),
			Type:        types.StringValue(v.Type),
		}
		if v.Title != "" && !bm.CalculationProperties[k].Title.IsNull() {
			calculationPropertyModel.Title = types.StringValue(v.Title)
		}

		if v.Description != "" && !bm.CalculationProperties[k].Description.IsNull() {
			calculationPropertyModel.Description = types.StringValue(v.Description)
		}

		if v.Format != "" && !bm.CalculationProperties[k].Format.IsNull() {
			calculationPropertyModel.Format = types.StringValue(v.Format)
		}

		bm.CalculationProperties[k] = *calculationPropertyModel

	}
}

func addStringPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewStringValue(value.(string)))
		}

		stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
	}

	if v.EnumColors != nil {
		stringProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		stringProp.EnumColors = types.MapNull(types.StringType)
	}

	if v.Format != nil {
		stringProp.Format = types.StringValue(*v.Format)
	}

	if v.Spec != nil {
		stringProp.Spec = types.StringValue(*v.Spec)
	}

	if v.MinLength != 0 {
		stringProp.MinLength = types.Int64Value(int64(v.MinLength))
	}

	if v.MaxLength != 0 {
		stringProp.MaxLength = types.Int64Value(int64(v.MaxLength))
	}

	if v.Pattern != "" {
		stringProp.Pattern = types.StringValue(v.Pattern)
	}

	if v.SpecAuthentication != nil {
		stringProp.SpecAuthentication = &SpecAuthenticationModel{
			AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
			TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
			ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
		}
	}

	return stringProp
}

func addNumberPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *NumberPropModel {
	numberProp := &NumberPropModel{}
	if v.Minimum != nil {
		numberProp.Minimum = types.Float64Value(*v.Minimum)
	}

	if v.Maximum != nil {
		numberProp.Maximum = types.Float64Value(*v.Maximum)
	}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
		}

		numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)
	} else {
		numberProp.Enum = types.ListNull(types.Float64Type)
	}

	if v.EnumColors != nil {
		numberProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		numberProp.EnumColors = types.MapNull(types.StringType)
	}

	return numberProp
}

func addObjectPropertiesToResource(v *cli.BlueprintProperty) *ObjectPropModel {
	objectProp := &ObjectPropModel{}

	if v.Spec != nil {
		objectProp.Spec = types.StringValue(*v.Spec)
	}

	return objectProp
}

func addArrayPropertiesToResource(v *cli.BlueprintProperty) *ArrayPropModel {
	arrayProp := &ArrayPropModel{}
	if v.MinItems != nil {
		arrayProp.MinItems = types.Int64Value(int64(*v.MinItems))
	}
	if v.MaxItems != nil {
		arrayProp.MaxItems = types.Int64Value(int64(*v.MaxItems))
	}
	if v.Items != nil {
		if v.Items["type"] != "" {
			switch v.Items["type"] {
			case "string":
				arrayProp.StringItems = &StringItems{}
				if v.Default != nil {
					stringArray := make([]string, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						stringArray[i] = v.(string)
					}
					attrs := make([]attr.Value, 0, len(stringArray))
					for _, value := range stringArray {
						attrs = append(attrs, basetypes.NewStringValue(value))
					}
					arrayProp.StringItems.Default, _ = types.ListValue(types.StringType, attrs)
				} else {
					arrayProp.StringItems.Default = types.ListNull(types.StringType)
				}
				if value, ok := v.Items["format"]; ok && value != nil {
					arrayProp.StringItems.Format = types.StringValue(v.Items["format"].(string))
				}
			case "number":
				arrayProp.NumberItems = &NumberItems{}
				if v.Default != nil {
					numberArray := make([]float64, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(numberArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
					}
					arrayProp.NumberItems.Default, _ = types.ListValue(types.Float64Type, attrs)
				}

			case "boolean":
				arrayProp.BooleanItems = &BooleanItems{}
				if v.Default != nil {
					booleanArray := make([]bool, len(v.Default.([]interface{})))
					attrs := make([]attr.Value, 0, len(booleanArray))
					for _, value := range v.Default.([]interface{}) {
						attrs = append(attrs, basetypes.NewBoolValue(value.(bool)))
					}
					arrayProp.BooleanItems.Default, _ = types.ListValue(types.BoolType, attrs)
				}

			case "object":
				arrayProp.ObjectItems = &ObjectItems{}
				if v.Default != nil {
					objectArray := make([]map[string]interface{}, len(v.Default.([]interface{})))
					for i, v := range v.Default.([]interface{}) {
						objectArray[i] = v.(map[string]interface{})
					}
					attrs := make([]attr.Value, 0, len(objectArray))
					for _, value := range objectArray {
						js, _ := json.Marshal(&value)
						stringValue := string(js)
						attrs = append(attrs, basetypes.NewStringValue(stringValue))
					}
					arrayProp.ObjectItems.Default, _ = types.ListValue(types.StringType, attrs)
				}
			}
		}
	}

	return arrayProp
}

func addPropertiesToResource(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) error {
	properties := &PropertiesModel{}

	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]StringPropModel)
			}
			stringProp := addStringPropertiesToResource(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, stringProp)

			properties.StringProp[k] = *stringProp

		case "number":
			if properties.NumberProp == nil {
				properties.NumberProp = make(map[string]NumberPropModel)
			}

			numberProp := addNumberPropertiesToResource(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, numberProp)

			properties.NumberProp[k] = *numberProp

		case "array":
			if properties.ArrayProp == nil {
				properties.ArrayProp = make(map[string]ArrayPropModel)
			}

			arrayProp := addArrayPropertiesToResource(&v)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, arrayProp)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, booleanProp)

			if lo.Contains(b.Schema.Required, k) {
				booleanProp.Required = types.BoolValue(true)
			} else {
				booleanProp.Required = types.BoolValue(false)
			}

			properties.BooleanProp[k] = *booleanProp

		case "object":
			if properties.ObjectProp == nil {
				properties.ObjectProp = make(map[string]ObjectPropModel)
			}

			objectProp := addObjectPropertiesToResource(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, objectProp)

			properties.ObjectProp[k] = *objectProp

		}

	}

	bm.Properties = properties

	return nil
}

func setCommonProperties(v cli.BlueprintProperty, prop interface{}) {
	properties := []string{"Description", "Icon", "Default", "Title"}
	for _, property := range properties {
		switch property {
		case "Description":
			if v.Description != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Description = types.StringValue(*v.Description)
				case *NumberPropModel:
					p.Description = types.StringValue(*v.Description)
				case *BooleanPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ArrayPropModel:
					p.Description = types.StringValue(*v.Description)
				case *ObjectPropModel:
					p.Description = types.StringValue(*v.Description)
				}
			}
		case "Icon":
			if v.Icon != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *NumberPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *BooleanPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ArrayPropModel:
					p.Icon = types.StringValue(*v.Icon)
				case *ObjectPropModel:
					p.Icon = types.StringValue(*v.Icon)
				}
			}
		case "Title":
			if v.Title != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Title = types.StringValue(*v.Title)
				case *NumberPropModel:
					p.Title = types.StringValue(*v.Title)
				case *BooleanPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ArrayPropModel:
					p.Title = types.StringValue(*v.Title)
				case *ObjectPropModel:
					p.Title = types.StringValue(*v.Title)
				}
			}
		case "Default":
			if v.Default != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Default = types.StringValue(v.Default.(string))
				case *NumberPropModel:
					p.Default = types.Float64Value(v.Default.(float64))
				case *BooleanPropModel:
					p.Default = types.BoolValue(v.Default.(bool))
				case *ObjectPropModel:
					js, _ := json.Marshal(v.Default)
					value := string(js)
					p.Default = types.StringValue(value)
				}
			}
		}
	}
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
		props, required, err = readPropertiesFromState(ctx, state)
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

func relationsResourceToBody(state *BlueprintModel) map[string]cli.Relation {
	relations := map[string]cli.Relation{}

	for identifier, prop := range state.Relations {
		target := prop.Target.ValueString()
		relationProp := cli.Relation{
			Target: &target,
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			relationProp.Title = &title
		}
		if !prop.Many.IsNull() {
			many := prop.Many.ValueBool()
			relationProp.Many = &many
		}

		if !prop.Required.IsNull() {
			required := prop.Required.ValueBool()
			relationProp.Required = &required
		}

		relations[identifier] = relationProp
	}

	return relations
}

func mirrorPropertiesToBody(state *BlueprintModel) map[string]cli.BlueprintMirrorProperty {
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}

	for identifier, prop := range state.MirrorProperties {
		mirrorProp := cli.BlueprintMirrorProperty{
			Path: prop.Path.ValueString(),
		}

		if !prop.Title.IsNull() {
			mirrorProp.Title = prop.Title.ValueString()
		}

		mirrorProperties[identifier] = mirrorProp
	}

	return mirrorProperties
}

func calculationPropertiesToBody(state *BlueprintModel) map[string]cli.BlueprintCalculationProperty {
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	for identifier, prop := range state.CalculationProperties {
		calculationProp := cli.BlueprintCalculationProperty{
			Calculation: prop.Calculation.ValueString(),
			Type:        prop.Type.ValueString(),
		}

		if !prop.Title.IsNull() {
			calculationProp.Title = prop.Title.ValueString()
		}

		if !prop.Description.IsNull() {
			calculationProp.Description = prop.Description.ValueString()
		}

		if !prop.Format.IsNull() {
			calculationProp.Format = prop.Format.ValueString()
		}

		calculationProperties[identifier] = calculationProp
	}

	return calculationProperties
}
