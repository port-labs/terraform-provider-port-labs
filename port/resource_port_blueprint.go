package port

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BlueprintResource{}
var _ resource.ResourceWithImportState = &BlueprintResource{}

func newBlueprintResource() resource.Resource {
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
				},
			},
		},
	}
}

func getArrayDefaultAttribute(arrayType interface{}) schema.Attribute {
	switch arrayType {
	case "string":
		return schema.ListAttribute{
			MarkdownDescription: "The default of the array property",
			Optional:            true,
			ElementType:         types.StringType,
		}
	case "boolean":
		return schema.ListAttribute{
			MarkdownDescription: "The default of the array property",
			Optional:            true,
			ElementType:         types.BoolType,
		}
	}
	return schema.ListAttribute{
		MarkdownDescription: "The default of the array property",
		Optional:            true,
		ElementType:         types.StringType,
	}
}

func (r *BlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *cli.BlueprintModel

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

func writeBlueprintFieldsToResource(bm *cli.BlueprintModel, b *cli.Blueprint) {
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
		bm.ChangelogDestination = &cli.ChangelogDestinationModel{
			Type:  types.StringValue(b.ChangelogDestination.Type),
			Url:   types.StringValue(b.ChangelogDestination.Url),
			Agent: types.BoolValue(b.ChangelogDestination.Agent),
		}
	}

	properties := &cli.PropertiesModel{}

	addPropertiesToResource(b, bm, properties)

	bm.Properties = properties

}

func addPropertiesToResource(b *cli.Blueprint, bm *cli.BlueprintModel, properties *cli.PropertiesModel) {
	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProp == nil {
				properties.StringProp = make(map[string]cli.StringPropModel)
			}

			stringProp := &cli.StringPropModel{}

			setCommonProperties(v, bm.Properties.StringProp[k], stringProp)

			properties.StringProp[k] = *stringProp

		case "number":
			if properties.NumberProp == nil {
				properties.NumberProp = make(map[string]cli.NumberPropModel)
			}

			numberProp := &cli.NumberPropModel{}

			if v.Minimum != 0 && !bm.Properties.NumberProp[k].Minimum.IsNull() {
				numberProp.Minimum = types.Float64Value(v.Minimum)
			}

			if v.Maximum != 0 && !bm.Properties.NumberProp[k].Maximum.IsNull() {
				numberProp.Maximum = types.Float64Value(v.Maximum)
			}

			setCommonProperties(v, bm.Properties.NumberProp[k], numberProp)

			properties.NumberProp[k] = *numberProp

		case "array":
			if properties.ArrayProp == nil {
				properties.ArrayProp = make(map[string]cli.ArrayPropModel)
			}

			arrayProp := &cli.ArrayPropModel{}

			if v.MinItems != 0 && !bm.Properties.ArrayProp[k].MinItems.IsNull() {
				arrayProp.MinItems = types.Int64Value(int64(v.MinItems))
			}
			if v.MaxItems != 0 && !bm.Properties.ArrayProp[k].MaxItems.IsNull() {
				arrayProp.MaxItems = types.Int64Value(int64(v.MaxItems))
			}

			if v.Items != nil {
				arrayProp.Items = &cli.ItemsModal{}
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
				properties.BooleanProp = make(map[string]cli.BooleanPropModel)
			}

			booleanProp := &cli.BooleanPropModel{}

			setCommonProperties(v, bm.Properties.BooleanProp[k], booleanProp)

			properties.BooleanProp[k] = *booleanProp
		}

	}
}

func setCommonProperties(v cli.BlueprintProperty, bm interface{}, prop interface{}) {
	properties := []string{"description", "icon", "default", "title"}
	for _, property := range properties {
		switch property {
		case "description":
			if v.Description != "" {
				switch p := prop.(type) {
				case *cli.StringPropModel:
					bmString := bm.(cli.StringPropModel)
					if !bmString.Description.IsNull() {
						p.Description = types.StringValue(v.Description)
					}
				case *cli.NumberPropModel:
					bmNumber := bm.(cli.NumberPropModel)
					if !bmNumber.Description.IsNull() {
						p.Description = types.StringValue(v.Description)
					}
				case *cli.BooleanPropModel:
					bmBoolean := bm.(cli.BooleanPropModel)
					if !bmBoolean.Description.IsNull() {
						p.Description = types.StringValue(v.Description)
					}

				case *cli.ArrayPropModel:
					bmArray := bm.(cli.ArrayPropModel)
					if !bmArray.Description.IsNull() {
						p.Description = types.StringValue(v.Description)
					}
				}
			}
		case "icon":
			if v.Icon != "" {
				switch p := prop.(type) {
				case *cli.StringPropModel:
					bmString := bm.(cli.StringPropModel)
					if !bmString.Icon.IsNull() {
						p.Icon = types.StringValue(v.Icon)
					}
				case *cli.NumberPropModel:
					bmNumber := bm.(cli.NumberPropModel)
					if !bmNumber.Icon.IsNull() {
						p.Icon = types.StringValue(v.Icon)
					}
				case *cli.BooleanPropModel:
					bmBoolean := bm.(cli.BooleanPropModel)
					if !bmBoolean.Icon.IsNull() {
						p.Icon = types.StringValue(v.Icon)
					}
				case *cli.ArrayPropModel:
					bmArray := bm.(cli.ArrayPropModel)
					if !bmArray.Icon.IsNull() {
						p.Icon = types.StringValue(v.Icon)
					}
				}
			}
		case "required":
			// Handle "required" property (add your logic here)
		case "title":
			if v.Title != "" {
				switch p := prop.(type) {
				case *cli.StringPropModel:
					bmString := bm.(cli.StringPropModel)
					if !bmString.Title.IsNull() {
						p.Title = types.StringValue(v.Title)
					}
				case *cli.NumberPropModel:
					bmNumber := bm.(cli.NumberPropModel)
					if !bmNumber.Title.IsNull() {
						p.Title = types.StringValue(v.Title)
					}
				case *cli.BooleanPropModel:
					bmBoolean := bm.(cli.BooleanPropModel)
					if !bmBoolean.Title.IsNull() {
						p.Title = types.StringValue(v.Title)
					}
				case *cli.ArrayPropModel:
					bmArray := bm.(cli.ArrayPropModel)
					if !bmArray.Title.IsNull() {
						p.Title = types.StringValue(v.Title)
					}
				}
			}

		case "default":
			if v.Default != "" {
				switch p := prop.(type) {
				case *cli.StringPropModel:
					bmString := bm.(cli.StringPropModel)
					if !bmString.Default.IsNull() {
						p.Default = types.StringValue(v.Default.(string))
					}
				case *cli.NumberPropModel:
					bmNumber := bm.(cli.NumberPropModel)
					if !bmNumber.Default.IsNull() {
						p.Default = types.Float64Value(v.Default.(float64))
					}

				case *cli.BooleanPropModel:
					bmBoolean := bm.(cli.BooleanPropModel)
					if !bmBoolean.Default.IsNull() {
						p.Default = types.BoolValue(v.Default.(bool))
					}
				}
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
	var data *cli.BlueprintModel
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
		resp.Diagnostics.AddError("failed to create R2 bucket", err.Error())
		return
	}

	writeBlueprintComputedFieldsToResource(data, bp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeBlueprintComputedFieldsToResource(bm *cli.BlueprintModel, bp *cli.Blueprint) {
	bm.Identifier = types.StringValue(bp.Identifier)
	bm.CreatedAt = types.StringValue(bp.CreatedAt.String())
	bm.CreatedBy = types.StringValue(bp.CreatedBy)
	bm.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(bp.UpdatedBy)
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *cli.BlueprintModel
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
}

func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}

func blueprintResourceToBody(ctx context.Context, d *cli.BlueprintModel) (*cli.Blueprint, error) {
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
			for propIdentifier, prop := range d.Properties.StringProp {
				props[propIdentifier] = cli.BlueprintProperty{
					Type:        "string",
					Title:       prop.Title.ValueString(),
					Format:      prop.Format.ValueString(),
					Default:     prop.Default.ValueString(),
					Icon:        prop.Icon.ValueString(),
					MinLength:   int(prop.MinLength.ValueInt64()),
					MaxLength:   int(prop.MaxLength.ValueInt64()),
					Pattern:     prop.Pattern.ValueString(),
					Description: prop.Description.ValueString(),
				}
				if prop.Required.ValueBool() {
					required = append(required, propIdentifier)
				}
			}
		}
		if d.Properties.ArrayProp != nil {
			for propIdentifier, prop := range d.Properties.ArrayProp {
				items := map[string]interface{}{}
				if prop.Items != nil {
					if !prop.Items.Type.IsNull() {
						items["type"] = prop.Items.Type.ValueString()
					}
					if !prop.Items.Format.IsNull() {
						items["format"] = prop.Items.Format.ValueString()
					}
					if !prop.Items.Default.IsNull() {
						items["default"] = prop.Items.Default
					}

				}

				props[propIdentifier] = cli.BlueprintProperty{
					Type:     "array",
					Title:    prop.Title.ValueString(),
					Icon:     prop.Icon.ValueString(),
					MaxItems: int(prop.MaxItems.ValueInt64()),
					MinItems: int(prop.MinItems.ValueInt64()),
					Items:    items,
				}
			}
		}
		if d.Properties.NumberProp != nil {
			for propIdentifier, prop := range d.Properties.NumberProp {
				props[propIdentifier] = cli.BlueprintProperty{
					Type:        "number",
					Title:       prop.Title.ValueString(),
					Default:     prop.Default.ValueFloat64(),
					Icon:        prop.Icon.ValueString(),
					Maximum:     prop.Maximum.ValueFloat64(),
					Minimum:     prop.Minimum.ValueFloat64(),
					Description: prop.Description.ValueString(),
				}
				if prop.Required.ValueBool() {
					required = append(required, propIdentifier)
				}
			}
		}
		if d.Properties.BooleanProp != nil {
			for propIdentifier, prop := range d.Properties.BooleanProp {
				props[propIdentifier] = cli.BlueprintProperty{
					Type:        "boolean",
					Title:       prop.Title.ValueString(),
					Default:     prop.Default.ValueBool(),
					Icon:        prop.Icon.ValueString(),
					Description: prop.Description.ValueString(),
				}
				if prop.Required.ValueBool() {
					required = append(required, propIdentifier)
				}
			}
		}

	}

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
	return b, nil
}
