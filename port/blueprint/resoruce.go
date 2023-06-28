package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
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
	var data *BlueprintModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, data.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	err = writeBlueprintFieldsToResource(ctx, data, b)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintFieldsToResource(ctx context.Context, bm *BlueprintModel, b *cli.Blueprint) error {
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

	if len(b.Schema.Properties) != 0 {
		err := addPropertiesToResource(ctx, b, bm)
		if err != nil {
			return err
		}
	}

	if len(b.Relations) != 0 {
		addRelationsToResource(b, bm)
	}

	if len(b.MirrorProperties) != 0 {
		addMirrorPropertiesToResource(b, bm)
	}

	if len(b.CalculationProperties) != 0 {
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
	if b.CalculationProperties != nil {
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
}

func addStingPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
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
			stringProp := addStingPropertiesToResource(ctx, &v)

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

			if !bm.Properties.ArrayProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					arrayProp.Required = types.BoolValue(true)
				} else {
					arrayProp.Required = types.BoolValue(false)
				}
			}

			setCommonProperties(v, arrayProp)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, booleanProp)

			if !bm.Properties.BooleanProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					booleanProp.Required = types.BoolValue(true)
				} else {
					booleanProp.Required = types.BoolValue(false)
				}
			}

			properties.BooleanProp[k] = *booleanProp

		case "object":
			if properties.ObjectProp == nil {
				properties.ObjectProp = make(map[string]ObjectPropModel)
			}

			objectProp := addObjectPropertiesToResource(&v)

			if !bm.Properties.ObjectProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					objectProp.Required = types.BoolValue(true)
				} else {
					objectProp.Required = types.BoolValue(false)
				}
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
	var data *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := blueprintResourceToBody(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	bp, err := r.portClient.CreateBlueprint(ctx, b)
	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	data.ID = types.StringValue(bp.Identifier)

	writeBlueprintComputedFieldsToResource(data, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintComputedFieldsToResource(bm *BlueprintModel, bp *cli.Blueprint) {
	bm.Identifier = types.StringValue(bp.Identifier)
	bm.CreatedAt = types.StringValue(bp.CreatedAt.String())
	bm.CreatedBy = types.StringValue(bp.CreatedBy)
	bm.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(bp.UpdatedBy)
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	b, err := blueprintResourceToBody(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError("failed to transform blueprint", err.Error())
		return
	}

	var bp *cli.Blueprint

	if data.Identifier.IsNull() {
		bp, err = r.portClient.CreateBlueprint(ctx, b)
	} else {
		bp, err = r.portClient.UpdateBlueprint(ctx, b, data.Identifier.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to update blueprint", err.Error())
		return
	}

	writeBlueprintComputedFieldsToResource(data, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *BlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *BlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Identifier.IsNull() {
		resp.Diagnostics.AddError("failed to extract blueprint identifier", "identifier is required")
		return
	}

	err := r.portClient.DeleteBlueprint(ctx, data.Identifier.ValueString())

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

func stringPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.Properties.StringProp {
		property := cli.BlueprintProperty{
			Type: "string",
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			property.Title = &title
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
		}

		if !prop.Format.IsNull() {
			format := prop.Format.ValueString()
			property.Format = &format
		}

		if !prop.Icon.IsNull() {
			icon := prop.Icon.ValueString()
			property.Icon = &icon
		}

		if !prop.MinLength.IsNull() {
			property.MinLength = int(prop.MinLength.ValueInt64())
		}

		if !prop.MaxLength.IsNull() {
			property.MaxLength = int(prop.MaxLength.ValueInt64())
		}

		if !prop.Spec.IsNull() {
			spec := prop.Spec.ValueString()
			property.Spec = &spec
		}

		if prop.SpecAuthentication != nil {
			specAuth := &cli.SpecAuthentication{
				AuthorizationUrl: prop.SpecAuthentication.AuthorizationUrl.ValueString(),
				TokenUrl:         prop.SpecAuthentication.TokenUrl.ValueString(),
				ClientId:         prop.SpecAuthentication.ClientId.ValueString(),
			}
			property.SpecAuthentication = specAuth
		}

		if !prop.Pattern.IsNull() {
			property.Pattern = prop.Pattern.ValueString()
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			property.Description = &description
		}

		if !prop.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "string")
			if err != nil {
				return err
			}
			property.Enum = enumList
		}

		if !prop.EnumColors.IsNull() {
			enumColor := map[string]string{}
			for k, v := range prop.EnumColors.Elements() {
				value, _ := v.ToTerraformValue(ctx)
				var keyValue string
				value.As(&keyValue)
				enumColor[k] = keyValue
			}

			property.EnumColors = enumColor
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func numberPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.Properties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "number",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueFloat64()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Minimum.IsNull() {
				minimum := prop.Minimum.ValueFloat64()
				property.Minimum = &minimum
			}

			if !prop.Maximum.IsNull() {
				maximum := prop.Maximum.ValueFloat64()
				property.Maximum = &maximum
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Enum.IsNull() {
				enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "float64")
				if err != nil {
					return err
				}
				property.Enum = enumList
			}

			if !prop.EnumColors.IsNull() {
				property.EnumColors = map[string]string{}
				for k, v := range prop.EnumColors.Elements() {
					value, _ := v.ToTerraformValue(ctx)
					var keyValue string
					value.As(&keyValue)
					property.EnumColors[k] = keyValue
				}
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func booleanPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "boolean",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueBool()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func objectPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.ObjectProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "object",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				defaultAsString := prop.Default.ValueString()
				defaultObj := make(map[string]interface{})
				err := json.Unmarshal([]byte(defaultAsString), &defaultObj)
				if err != nil {
					log.Fatal(err)
				} else {
					property.Default = defaultObj
				}
			}

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.Spec.IsNull() {
				spec := prop.Spec.ValueString()
				property.Spec = &spec
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func arrayPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.Properties.ArrayProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "array",
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}
			if !prop.MinItems.IsNull() {
				minItems := int(prop.MinItems.ValueInt64())
				property.MinItems = &minItems
			}

			if !prop.MaxItems.IsNull() {
				maxItems := int(prop.MaxItems.ValueInt64())
				property.MaxItems = &maxItems
			}

			if prop.StringItems != nil {
				items := map[string]interface{}{}
				items["type"] = "string"
				if !prop.StringItems.Format.IsNull() {
					items["format"] = prop.StringItems.Format.ValueString()
				}
				if !prop.StringItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.StringItems.Default, "string")
					if err != nil {
						return err
					}

					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.NumberItems != nil {
				items := map[string]interface{}{}
				items["type"] = "number"
				if !prop.NumberItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.NumberItems.Default, "float64")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.BooleanItems != nil {
				items := map[string]interface{}{}
				items["type"] = "boolean"
				if !prop.BooleanItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.BooleanItems.Default, "bool")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.ObjectItems != nil {
				items := map[string]interface{}{}
				items["type"] = "object"
				if !prop.ObjectItems.Default.IsNull() {
					defaultList, err := utils.TerraformListToGoArray(ctx, prop.ObjectItems.Default, "object")
					if err != nil {
						return err
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}

	return nil
}

func blueprintResourceToBody(ctx context.Context, d *BlueprintModel) (*cli.Blueprint, error) {
	b := &cli.Blueprint{
		Identifier: d.Identifier.ValueString(),
	}

	if !d.Title.IsNull() {
		titleValue := d.Title.ValueString()
		b.Title = &titleValue
	}

	if !d.Icon.IsNull() {
		iconValue := d.Icon.ValueString()
		b.Icon = &iconValue
	}

	if !d.Description.IsNull() {
		descriptionTest := d.Description.ValueString()
		b.Description = &descriptionTest
	}
	props := map[string]cli.BlueprintProperty{}

	if d.ChangelogDestination != nil {
		if d.ChangelogDestination.Type.ValueString() == "KAFKA" && !d.ChangelogDestination.Agent.IsNull() {
			return nil, fmt.Errorf("agent is not supported for Kafka changelog destination")
		}
		b.ChangelogDestination = &cli.ChangelogDestination{}
		b.ChangelogDestination.Type = d.ChangelogDestination.Type.ValueString()
		b.ChangelogDestination.Url = d.ChangelogDestination.Url.ValueString()
		b.ChangelogDestination.Agent = d.ChangelogDestination.Agent.ValueBool()
	} else {
		b.ChangelogDestination = nil
	}

	if d.TeamInheritance != nil {
		b.TeamInheritance = &cli.TeamInheritance{
			Path: d.TeamInheritance.Path.ValueString(),
		}
	} else {
		b.TeamInheritance = nil
	}

	required := []string{}

	if d.Properties != nil {
		if d.Properties.StringProp != nil {
			err := stringPropResourceToBody(ctx, d, props, &required)
			if err != nil {
				return nil, err
			}
		}
		if d.Properties.ArrayProp != nil {
			err := arrayPropResourceToBody(ctx, d, props, &required)
			if err != nil {
				return nil, err
			}
		}
		if d.Properties.NumberProp != nil {
			err := numberPropResourceToBody(ctx, d, props, &required)
			if err != nil {
				return nil, err
			}
		}
		if d.Properties.BooleanProp != nil {
			booleanPropResourceToBody(d, props, &required)
		}

		if d.Properties.ObjectProp != nil {
			objectPropResourceToBody(d, props, &required)
		}

	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relationsResourceToBody(d)
	b.MirrorProperties = mirrorPropertiesToBody(d)
	b.CalculationProperties = calculationPropertiesToBody(d)
	return b, nil
}

func relationsResourceToBody(d *BlueprintModel) map[string]cli.Relation {
	relations := map[string]cli.Relation{}

	for identifier, prop := range d.Relations {
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

func mirrorPropertiesToBody(d *BlueprintModel) map[string]cli.BlueprintMirrorProperty {
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}

	for identifier, prop := range d.MirrorProperties {
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

func calculationPropertiesToBody(d *BlueprintModel) map[string]cli.BlueprintCalculationProperty {
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	for identifier, prop := range d.CalculationProperties {
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
