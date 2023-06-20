package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
	"github.com/samber/lo"
)

// Ensure provider defined types fully satisfy framework interfaces
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

// 	case "string":
// 		return schema.ListAttribute{
// 			MarkdownDescription: "The default of the array property",
// 			Optional:            true,
// 			ElementType:         types.StringType,
// 		}
// 	case "boolean":
// 		return schema.ListAttribute{
// 			MarkdownDescription: "The default of the array property",
// 			Optional:            true,
// 			ElementType:         types.BoolType,
// 		}
// 	}
// 	return schema.ListAttribute{
// 		MarkdownDescription: "The default of the array property",
// 		Optional:            true,
// 		ElementType:         types.StringType,
// 	}
// }

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to read blueprint: %s", err))
		return
	}

	writeBlueprintFieldsToResource(ctx, data, b)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintFieldsToResource(ctx context.Context, bm *BlueprintModel, b *cli.Blueprint) {
	bm.Identifier = types.StringValue(b.Identifier)
	if !bm.Title.IsNull() {
		bm.Title = types.StringValue(b.Title)
	}

	if !bm.Icon.IsNull() {
		bm.Icon = types.StringValue(b.Icon)
	}

	if !bm.Description.IsNull() {
		bm.Description = types.StringValue(b.Description)
	}
	bm.CreatedAt = types.StringValue(b.CreatedAt.String())
	bm.CreatedBy = types.StringValue(b.CreatedBy)
	bm.UpdatedAt = types.StringValue(b.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(b.UpdatedBy)
	if b.ChangelogDestination != nil {
		bm.ChangelogDestination = &ChangelogDestinationModel{
			Type:  types.StringValue(b.ChangelogDestination.Type),
			Url:   types.StringValue(b.ChangelogDestination.Url),
			Agent: types.BoolValue(b.ChangelogDestination.Agent),
		}
	}

	properties := &PropertiesModel{}

	if bm.Properties == nil && len(b.Schema.Properties) == 0 {
		bm.Properties = nil
	} else {
		if bm.Properties == nil {
			bm.Properties = &PropertiesModel{}
		}
		addPropertiesToResource(ctx, b, bm, properties)
		bm.Properties = properties
	}

	if bm.Relations == nil && len(b.Relations) == 0 {
		bm.Relations = nil
	} else {
		addRelationsToResource(b, bm)
	}

	if bm.MirrorProperties == nil && len(b.MirrorProperties) == 0 {
		bm.MirrorProperties = nil
	} else {
		addMirrorPropertiesToResource(b, bm)
	}

}

func addRelationsToResource(b *cli.Blueprint, bm *BlueprintModel) {
	if b.Relations != nil {
		for k, v := range b.Relations {
			if bm.Relations == nil {
				bm.Relations = make(map[string]RelationModel)
			}

			relationModel := &RelationModel{
				Target: types.StringValue(v.Target),
			}
			if v.Title != "" {
				relationModel.Title = types.StringValue(v.Title)
			}
			if !bm.Relations[k].Many.IsNull() {
				relationModel.Many = types.BoolValue(v.Many)
			}

			if !bm.Relations[k].Required.IsNull() {
				relationModel.Required = types.BoolValue(v.Required)
			}

			bm.Relations[k] = *relationModel

		}
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

func addPropertiesToResource(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel, properties *PropertiesModel) {
	for k, v := range b.Schema.Properties {
		isImportActive := false
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]StringPropModel)
			}

			stringProp := &StringPropModel{}

			if bm.Properties.StringProp == nil {
				isImportActive = true
				bm.Properties.StringProp = make(map[string]StringPropModel)
				bm.Properties.StringProp[k] = *stringProp

			}
			if v.Enum != nil && !bm.Properties.StringProp[k].Enum.IsNull() {
				attrs := make([]attr.Value, 0, len(v.Enum))
				for _, value := range v.Enum {
					attrs = append(attrs, basetypes.NewStringValue(value.(string)))
				}

				stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
			} else {
				stringProp.Enum = types.ListNull(types.StringType)
			}

			if v.EnumColors != nil && !bm.Properties.StringProp[k].EnumColors.IsNull() {
				stringProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
			} else {
				stringProp.EnumColors = types.MapNull(types.StringType)
			}

			if v.Format != "" && !bm.Properties.StringProp[k].Format.IsNull() {
				stringProp.Format = types.StringValue(v.Format)
			}

			if v.Spec != "" && !bm.Properties.StringProp[k].Spec.IsNull() {
				stringProp.Spec = types.StringValue(v.Spec)
			}

			if v.MinLength != 0 && !bm.Properties.StringProp[k].MinLength.IsNull() {
				stringProp.MinLength = types.Int64Value(int64(v.MinLength))
			}

			if v.MaxLength != 0 && !bm.Properties.StringProp[k].MaxLength.IsNull() {
				stringProp.MaxLength = types.Int64Value(int64(v.MaxLength))
			}

			if v.Pattern != "" && !bm.Properties.StringProp[k].Pattern.IsNull() {
				stringProp.Pattern = types.StringValue(v.Pattern)
			}

			if !bm.Properties.StringProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					stringProp.Required = types.BoolValue(true)
				} else {
					stringProp.Required = types.BoolValue(false)
				}
			}

			if v.SpecAuthentication != nil && bm.Properties.StringProp[k].SpecAuthentication != nil {
				stringProp.SpecAuthentication = &SpecAuthenticationModel{
					AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
					TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
					ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
				}
			}

			setCommonProperties(v, bm.Properties.StringProp[k], stringProp, isImportActive)

			properties.StringProp[k] = *stringProp

		case "number":
			if properties.NumberProp == nil {
				properties.NumberProp = make(map[string]NumberPropModel)
			}

			numberProp := &NumberPropModel{}

			if v.Minimum != 0 && !bm.Properties.NumberProp[k].Minimum.IsNull() {
				numberProp.Minimum = types.Float64Value(v.Minimum)
			}

			if v.Maximum != 0 && !bm.Properties.NumberProp[k].Maximum.IsNull() {
				numberProp.Maximum = types.Float64Value(v.Maximum)
			}

			if v.Enum != nil && !bm.Properties.NumberProp[k].Enum.IsNull() {
				attrs := make([]attr.Value, 0, len(v.Enum))
				for _, value := range v.Enum {
					attrs = append(attrs, basetypes.NewFloat64Value(value.(float64)))
				}

				numberProp.Enum, _ = types.ListValue(types.Float64Type, attrs)
			} else {
				numberProp.Enum = types.ListNull(types.Float64Type)
			}

			if v.EnumColors != nil && !bm.Properties.NumberProp[k].EnumColors.IsNull() {
				numberProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
			}

			if !bm.Properties.NumberProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					numberProp.Required = types.BoolValue(true)
				} else {
					numberProp.Required = types.BoolValue(false)
				}
			}

			setCommonProperties(v, bm.Properties.NumberProp[k], numberProp, isImportActive)

			properties.NumberProp[k] = *numberProp

		case "array":
			if properties.ArrayProp == nil {
				properties.ArrayProp = make(map[string]ArrayPropModel)
			}

			arrayProp := &ArrayPropModel{}

			if v.MinItems != 0 && !bm.Properties.ArrayProp[k].MinItems.IsNull() {
				arrayProp.MinItems = types.Int64Value(int64(v.MinItems))
			}
			if v.MaxItems != 0 && !bm.Properties.ArrayProp[k].MaxItems.IsNull() {
				arrayProp.MaxItems = types.Int64Value(int64(v.MaxItems))
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
					}
				}
			}

			if !bm.Properties.ArrayProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					arrayProp.Required = types.BoolValue(true)
				} else {
					arrayProp.Required = types.BoolValue(false)
				}
			}

			setCommonProperties(v, bm.Properties.ArrayProp[k], arrayProp, isImportActive)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, bm.Properties.BooleanProp[k], booleanProp, isImportActive)

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

			objectProp := &ObjectPropModel{}

			if v.Spec != "" && !bm.Properties.ObjectProp[k].Spec.IsNull() {
				objectProp.Spec = types.StringValue(v.Spec)
			}

			if !bm.Properties.ObjectProp[k].Required.IsNull() {
				if lo.Contains(b.Schema.Required, k) {
					objectProp.Required = types.BoolValue(true)
				} else {
					objectProp.Required = types.BoolValue(false)
				}
			}

			setCommonProperties(v, bm.Properties.ObjectProp[k], objectProp, isImportActive)

			properties.ObjectProp[k] = *objectProp

		}

	}

}

func setCommonProperties(v cli.BlueprintProperty, bm interface{}, prop interface{}, isImportActive bool) {
	properties := []string{"Description", "Icon", "Default", "Title", "Required"}
	for _, property := range properties {
		switch property {
		case "Description":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Description == "" && bmString.Description.IsNull() && !isImportActive {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Description == "" && bmNumber.Description.IsNull() && !isImportActive {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Description == "" && bmBoolean.Description.IsNull() && !isImportActive {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Description == "" && bmArray.Description.IsNull() && !isImportActive {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Description == "" && bmObject.Description.IsNull() && !isImportActive {
					continue
				}
				p.Description = types.StringValue(v.Description)
			}
		case "Icon":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Icon == "" && bmString.Icon.IsNull() && !isImportActive {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Icon == "" && bmNumber.Icon.IsNull() && !isImportActive {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Icon == "" && bmBoolean.Icon.IsNull() && !isImportActive {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Icon == "" && bmArray.Icon.IsNull() && !isImportActive {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Icon == "" && bmObject.Icon.IsNull() && !isImportActive {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			}
		case "Title":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Title == "" && bmString.Title.IsNull() && !isImportActive {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Title == "" && bmNumber.Title.IsNull() && !isImportActive {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Title == "" && bmBoolean.Title.IsNull() && !isImportActive {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Title == "" && bmArray.Title.IsNull() && !isImportActive {
					continue
				}
				p.Title = types.StringValue(v.Title)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Title == "" && bmObject.Title.IsNull() && !isImportActive {
					continue
				}
				p.Title = types.StringValue(v.Title)

			}

		case "Default":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Default == nil && bmString.Default.IsNull() && !isImportActive {
					continue
				}
				p.Default = types.StringValue(v.Default.(string))
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Default == nil && bmNumber.Default.IsNull() && !isImportActive {
					continue
				}
				p.Default = types.Float64Value(v.Default.(float64))
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Default == nil && bmBoolean.Default.IsNull() && !isImportActive {
					continue
				}
				p.Default = types.BoolValue(v.Default.(bool))
			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Default == nil && bmObject.Default.IsNull() && !isImportActive {
					continue
				}
				js, _ := json.Marshal(v.Default)
				value := string(js)
				p.Default = types.StringValue(value)
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

func stringPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.StringProp {
		property := cli.BlueprintProperty{
			Type:  "string",
			Title: prop.Title.ValueString(),
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
		}

		if !prop.Format.IsNull() {
			property.Format = prop.Format.ValueString()
		}

		if !prop.Icon.IsNull() {
			property.Icon = prop.Icon.ValueString()
		}

		if !prop.MinLength.IsNull() {
			property.MinLength = int(prop.MinLength.ValueInt64())
		}

		if !prop.MaxLength.IsNull() {
			property.MaxLength = int(prop.MaxLength.ValueInt64())
		}

		if !prop.Spec.IsNull() {
			property.Spec = prop.Spec.ValueString()
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
			property.Description = prop.Description.ValueString()
		}

		if !prop.Enum.IsNull() {
			enumList := []interface{}{}
			for _, e := range prop.Enum.Elements() {
				v, _ := e.ToTerraformValue(ctx)
				var keyValue string
				v.As(&keyValue)
				enumList = append(enumList, keyValue)
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
}

func numberPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "number",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueFloat64()
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Minimum.IsNull() {
				property.Minimum = prop.Minimum.ValueFloat64()
			}

			if !prop.Maximum.IsNull() {
				property.Maximum = prop.Maximum.ValueFloat64()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Enum.IsNull() {
				property.Enum = []interface{}{}
				for _, e := range prop.Enum.Elements() {
					v, _ := e.ToTerraformValue(ctx)
					var keyValue big.Float
					v.As(&keyValue)
					floatValue, _ := keyValue.Float64()
					property.Enum = append(property.Enum, floatValue)
				}
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
}

func booleanPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "boolean",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueBool()
			}

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
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
			Type:  "object",
			Title: prop.Title.ValueString(),
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

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Spec.IsNull() {
				property.Spec = prop.Spec.ValueString()
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func arrayPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.Properties.ArrayProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "array",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {

			if !prop.Icon.IsNull() {
				property.Icon = prop.Icon.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}
			if !prop.MinItems.IsNull() {
				property.MinItems = int(prop.MinItems.ValueInt64())
			}

			if !prop.MaxItems.IsNull() {
				property.MaxItems = int(prop.MaxItems.ValueInt64())
			}

			if prop.StringItems != nil {
				items := map[string]interface{}{}
				items["type"] = "string"
				if !prop.StringItems.Format.IsNull() {
					items["format"] = prop.StringItems.Format.ValueString()
				}
				if !prop.StringItems.Default.IsNull() {
					defaultList := []interface{}{}
					for _, e := range prop.StringItems.Default.Elements() {
						v, _ := e.ToTerraformValue(ctx)
						var keyValue string
						v.As(&keyValue)
						defaultList = append(defaultList, keyValue)
					}
					property.Default = defaultList
				}
				property.Items = items
			}

			if prop.NumberItems != nil {
				items := map[string]interface{}{}
				items["type"] = "number"
				if !prop.NumberItems.Default.IsNull() {
					items["default"] = prop.NumberItems.Default
				}
				property.Items = items
			}

			if prop.BooleanItems != nil {
				items := map[string]interface{}{}
				items["type"] = "boolean"
				if !prop.BooleanItems.Default.IsNull() {
					items["default"] = prop.BooleanItems.Default
				}
				property.Items = items
			}

			if prop.ObjectItems != nil {
				items := map[string]interface{}{}
				items["type"] = "object"
				if !prop.ObjectItems.Default.IsNull() {
					items["default"] = prop.ObjectItems.Default
				}
				property.Items = items
			}

			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}

func blueprintResourceToBody(ctx context.Context, d *BlueprintModel) (*cli.Blueprint, error) {
	b := &cli.Blueprint{}
	b.Identifier = d.Identifier.ValueString()

	b.Title = d.Title.ValueString()
	b.Icon = d.Icon.ValueString()
	b.Description = d.Description.ValueString()
	props := map[string]cli.BlueprintProperty{}
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}

	if d.ChangelogDestination != nil {
		b.ChangelogDestination = &cli.ChangelogDestination{}
		b.ChangelogDestination.Type = d.ChangelogDestination.Type.ValueString()
		b.ChangelogDestination.Url = d.ChangelogDestination.Url.ValueString()
		b.ChangelogDestination.Agent = d.ChangelogDestination.Agent.ValueBool()
	} else {
		b.ChangelogDestination = nil
	}

	required := []string{}

	if d.Properties != nil {
		if d.Properties.StringProp != nil {
			stringPropResourceToBody(ctx, d, props, &required)
		}
		if d.Properties.ArrayProp != nil {
			arrayPropResourceToBody(ctx, d, props, &required)
		}
		if d.Properties.NumberProp != nil {
			numberPropResourceToBody(ctx, d, props, &required)
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
	b.CalculationProperties = calculationProperties
	return b, nil
}

func relationsResourceToBody(d *BlueprintModel) map[string]cli.Relation {
	relations := map[string]cli.Relation{}

	for identifier, prop := range d.Relations {
		relationProp := cli.Relation{
			Target: prop.Target.ValueString(),
		}

		if !prop.Title.IsNull() {
			relationProp.Title = prop.Title.ValueString()
		}
		if !prop.Many.IsNull() {
			relationProp.Many = prop.Many.ValueBool()
		}

		if !prop.Required.IsNull() {
			relationProp.Required = prop.Required.ValueBool()
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
