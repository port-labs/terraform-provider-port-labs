package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
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

func (r *BlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group resource",

		Attributes: map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the blueprint",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The display name of the blueprint",
				Optional:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "The icon of the blueprint",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the blueprint",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The creation date of the blueprint",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The creator of the blueprint",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The last update date of the blueprint",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "The last updater of the blueprint",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"changelog_destination": schema.SingleNestedAttribute{
				MarkdownDescription: "The changelog destination of the blueprint",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the changelog destination",
						Required:            true,
						// Validators:          []validator.String{stringvalidator.OneOf("WEBHOOK", "KAFKA")},
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "The url of the changelog destination",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.Expressions{
								path.MatchRelative(),
							}...),
						},
					},
					"agent": schema.BoolAttribute{
						MarkdownDescription: "The agent of the changelog destination",
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"properties": schema.SingleNestedAttribute{
				MarkdownDescription: "The properties of the blueprint",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"string_prop": schema.MapNestedAttribute{
						MarkdownDescription: "The string property of the blueprint",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									MarkdownDescription: "The display name of the string property",
									Optional:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "The description of the string property",
									Optional:            true,
								},
								"default": schema.StringAttribute{
									MarkdownDescription: "The default of the string property",
									Optional:            true,
								},
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the string property",
									Optional:            true,
								},
								"format": schema.StringAttribute{
									MarkdownDescription: "The format of the string property",
									Optional:            true,
								},
								"required": schema.BoolAttribute{
									MarkdownDescription: "The required of the string property",
									Optional:            true,
								},
								"min_length": schema.Int64Attribute{
									MarkdownDescription: "The min length of the string property",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"max_length": schema.Int64Attribute{
									MarkdownDescription: "The max length of the string property",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"pattern": schema.StringAttribute{
									MarkdownDescription: "The pattern of the string property",
									Optional:            true,
								},
								"enum": schema.ListAttribute{
									MarkdownDescription: "The enum of the string property",
									Optional:            true,
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.UniqueValues(),
										listvalidator.SizeAtLeast(1),
									},
								},
							},
						},
					},
					"number_prop": schema.MapNestedAttribute{
						MarkdownDescription: "The number property of the blueprint",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									MarkdownDescription: "The display name of the number property",
									Optional:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "The description of the number property",
									Optional:            true,
								},
								"default": schema.Float64Attribute{
									MarkdownDescription: "The default of the number property",
									Optional:            true,
								},
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the number property",
									Optional:            true,
								},
								"required": schema.BoolAttribute{
									MarkdownDescription: "The required of the number property",
									Optional:            true,
								},
								"maximum": schema.Float64Attribute{
									MarkdownDescription: "The min of the number property",
									Optional:            true,
								},
								"minimum": schema.Float64Attribute{
									MarkdownDescription: "The max of the number property",
									Optional:            true,
								},
								"enum": schema.ListAttribute{
									MarkdownDescription: "The enum of the number property",
									Optional:            true,
									ElementType:         types.Float64Type,
									Validators: []validator.List{
										listvalidator.UniqueValues(),
										listvalidator.SizeAtLeast(1),
									},
								},
							},
						},
					},
					"boolean_prop": schema.MapNestedAttribute{
						MarkdownDescription: "The boolean property of the blueprint",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									MarkdownDescription: "The display name of the boolean property",
									Optional:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "The description of the boolean property",
									Optional:            true,
								},
								"default": schema.BoolAttribute{
									MarkdownDescription: "The default of the boolean property",
									Optional:            true,
								},
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the boolean property",
									Optional:            true,
								},
								"required": schema.BoolAttribute{
									MarkdownDescription: "The required of the boolean property",
									Optional:            true,
								},
							},
						},
					},
					"array_prop": schema.MapNestedAttribute{
						MarkdownDescription: "The array property of the blueprint",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									MarkdownDescription: "The display name of the array property",
									Optional:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "The description of the array property",
									Optional:            true,
								},
								// "default": schema.ListAttribute{
								// 	MarkdownDescription: "The default of the array property",
								// 	Optional:            true,
								// },
								// "default": schema.ListAttribute{
								// 	MarkdownDescription: "The default of the array property",
								// 	Optional:            true,
								// 	ElementType:         types.ListUnknown(),
								// },
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the array property",
									Optional:            true,
								},
								"required": schema.BoolAttribute{
									MarkdownDescription: "The required of the array property",
									Optional:            true,
									// Default:             booldefault.StaticBool(false),
								},
								"min_items": schema.Int64Attribute{
									MarkdownDescription: "The min items of the array property",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"max_items": schema.Int64Attribute{
									MarkdownDescription: "The max items of the array property",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"items": schema.SingleNestedAttribute{
									MarkdownDescription: "The items of the array property",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "The type of the items",
											Required:            true,
											// Validators:          []validator.String{stringvalidator.OneOf("STRING", "BOOLEAN", "INTEGER", "FLOAT", "OBJECT", "ARRAY")},
										},
										"format": schema.StringAttribute{
											MarkdownDescription: "The format of the items",
											Optional:            true,
										},
										"default": schema.ListAttribute{
											ElementType:         types.StringType,
											MarkdownDescription: "The default of the items",
											Optional:            true,
										},
									},
								},
							},
						},
					},
					"object_prop": schema.MapNestedAttribute{
						MarkdownDescription: "The object property of the blueprint",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									MarkdownDescription: "The display name of the object property",
									Optional:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "The description of the object property",
									Optional:            true,
								},
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the object property",
									Optional:            true,
								},
								"default": schema.MapAttribute{
									Optional:            true,
									MarkdownDescription: "The default of the object property",
									ElementType:         types.StringType,
								},
								"required": schema.BoolAttribute{
									MarkdownDescription: "The required of the object property",
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// func getArrayDefaultAttribute(arrayType interface{}) schema.Attribute {
// 	switch arrayType {
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

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read data from the API
	b, statusCode, err := r.portClient.ReadBlueprint(ctx, data.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to read blueprint: %s", err))
		return
	}

	writeBlueprintFieldsToResource(data, b)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintFieldsToResource(bm *BlueprintModel, b *cli.Blueprint) {
	bm.Identifier = types.StringValue(b.Identifier)
	bm.Identifier = types.StringValue(b.Identifier)
	bm.Title = types.StringValue(b.Title)
	bm.Icon = types.StringValue(b.Icon)
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

	addPropertiesToResource(b, bm, properties)

	bm.Properties = properties

}

func addPropertiesToResource(b *cli.Blueprint, bm *BlueprintModel, properties *PropertiesModel) {
	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]StringPropModel)
			}

			stringProp := &StringPropModel{}

			if v.Enum != nil && !bm.Properties.StringProp[k].Enum.IsNull() {
				attrs := make([]attr.Value, 0, len(v.Enum))
				for _, value := range v.Enum {
					attrs = append(attrs, basetypes.NewStringValue(value.(string)))
				}

				stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
			} else {
				stringProp.Enum = types.ListNull(types.StringType)

			}

			setCommonProperties(v, bm.Properties.StringProp[k], stringProp)

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
			}

			setCommonProperties(v, bm.Properties.NumberProp[k], numberProp)

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
				arrayProp.Items = &ItemsModal{}
				if itemType, ok := v.Items["type"].(string); ok {
					arrayProp.Items.Type = types.StringValue(itemType)
				}
				if itemFormat, ok := v.Items["format"].(string); ok {
					arrayProp.Items.Format = types.StringValue(itemFormat)
				}
			}

			setCommonProperties(v, bm.Properties.ArrayProp[k], arrayProp)

			properties.ArrayProp[k] = *arrayProp

		case "boolean":
			if properties.BooleanProp == nil {
				properties.BooleanProp = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, bm.Properties.BooleanProp[k], booleanProp)

			properties.BooleanProp[k] = *booleanProp

		case "object":
			if properties.ObjectProp == nil {
				properties.ObjectProp = make(map[string]ObjectPropModel)
			}

			objectProp := &ObjectPropModel{}

			setCommonProperties(v, bm.Properties.ObjectProp[k], objectProp)

			properties.ObjectProp[k] = *objectProp

		}

	}
}

func setCommonProperties(v cli.BlueprintProperty, bm interface{}, prop interface{}) {
	properties := []string{"description", "icon", "default", "title"}
	for _, property := range properties {
		switch property {
		case "description":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Description == "" && bmString.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Description == "" && bmNumber.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Description == "" && bmBoolean.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Description == "" && bmArray.Description.IsNull() {
					continue
				}

				p.Description = types.StringValue(v.Description)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Description == "" && bmObject.Description.IsNull() {
					continue
				}
				p.Description = types.StringValue(v.Description)
			}
		case "icon":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Icon == "" && bmString.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Icon == "" && bmNumber.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Icon == "" && bmBoolean.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Icon == "" && bmArray.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Icon == "" && bmObject.Icon.IsNull() {
					continue
				}
				p.Icon = types.StringValue(v.Icon)
			}
		case "title":

			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Title == "" && bmString.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Title == "" && bmNumber.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Title == "" && bmBoolean.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)
			case *ArrayPropModel:
				bmArray := bm.(ArrayPropModel)
				if v.Title == "" && bmArray.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)

			case *ObjectPropModel:
				bmObject := bm.(ObjectPropModel)
				if v.Title == "" && bmObject.Title.IsNull() {
					continue
				}
				p.Title = types.StringValue(v.Title)

			}

		case "default":
			switch p := prop.(type) {
			case *StringPropModel:
				bmString := bm.(StringPropModel)
				if v.Default == nil && bmString.Default.IsNull() {
					continue
				}
				p.Default = types.StringValue(v.Default.(string))
			case *NumberPropModel:
				bmNumber := bm.(NumberPropModel)
				if v.Default == nil && bmNumber.Default.IsNull() {
					continue
				}
				p.Default = types.Float64Value(v.Default.(float64))
			case *BooleanPropModel:
				bmBoolean := bm.(BooleanPropModel)
				if v.Default == nil && bmBoolean.Default.IsNull() {
					continue
				}
				p.Default = types.BoolValue(v.Default.(bool))
				// case *cli.ObjectPropModel:
				// 	bmObject := bm.(cli.ObjectPropModel)
				// 	if bmObject.Default != nil {
				// 		p.Default = v.Default.(map[string]interface{})
				// 	}
				// }

			}
		}
	}
}

func defaultResourceToBody(value string, propFields *cli.BlueprintProperty) error {
	switch propFields.Type {
	case "string":
		propFields.Default = value
	case "number":
		defaultNum, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return err
		}
		propFields.Default = defaultNum

	case "boolean":

		defaultBool, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		propFields.Default = defaultBool

	case "object":
		defaultObj := make(map[string]interface{})
		err := json.Unmarshal([]byte(value), &defaultObj)
		if err != nil {
			return err
		}
		propFields.Default = defaultObj

	}
	return nil
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
	fmt.Printf("Creating Blueprint %+v\n", b)
	bp, err := r.portClient.CreateBlueprint(ctx, b)
	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

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

}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}

func stringPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.StringProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "string",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
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

			if !prop.Pattern.IsNull() {
				property.Pattern = prop.Pattern.ValueString()
			}

			if !prop.Description.IsNull() {
				property.Description = prop.Description.ValueString()
			}

			if !prop.Enum.IsNull() {
				property.Enum = []interface{}{}
				for _, e := range prop.Enum.Elements() {
					v, _ := e.ToTerraformValue(ctx)
					var keyValue string
					v.As(&keyValue)
					property.Enum = append(property.Enum, keyValue)
				}
			}
			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func numberPropResourceToBody(ctx context.Context, d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.NumberProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "number",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
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
					var keyValue float64
					v.As(&keyValue)
					property.Enum = append(property.Enum, keyValue)
				}
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
		}
	}
}

func booleanPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.BooleanProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "boolean",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
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
			required = append(required, propIdentifier)
		}
	}
}

func objectPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
	for propIdentifier, prop := range d.Properties.ObjectProp {
		props[propIdentifier] = cli.BlueprintProperty{
			Type:  "object",
			Title: prop.Title.ValueString(),
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				property.Default = prop.Default
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
			required = append(required, propIdentifier)
		}
	}
}

func arrayPropResourceToBody(d *BlueprintModel, props map[string]cli.BlueprintProperty, required []string) {
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
			if !prop.Items.Type.IsNull() {
				items := map[string]interface{}{}
				if !prop.Items.Type.IsNull() {
					items["type"] = prop.Items.Type.ValueString()
				}
				if !prop.Items.Format.IsNull() {
					items["format"] = prop.Items.Format.ValueString()
				}
				if !prop.Items.Default.IsNull() {
					items["default"] = prop.Items.Default
				}

				property.Items = items
			}
			props[propIdentifier] = property
		}

		if prop.Required.ValueBool() {
			required = append(required, propIdentifier)
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
	mirrorProperties := map[string]cli.BlueprintMirrorProperty{}
	calculationProperties := map[string]cli.BlueprintCalculationProperty{}
	relations := map[string]cli.Relation{}

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
			stringPropResourceToBody(ctx, d, props, required)
		}
		if d.Properties.ArrayProp != nil {
			arrayPropResourceToBody(d, props, required)
		}
		if d.Properties.NumberProp != nil {
			numberPropResourceToBody(ctx, d, props, required)
		}
		if d.Properties.BooleanProp != nil {
			booleanPropResourceToBody(d, props, required)
		}

		if d.Properties.ObjectProp != nil {
			objectPropResourceToBody(d, props, required)
		}

	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
	return b, nil
}
