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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
									Computed:            true,
									Default:             booldefault.StaticBool(false),
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
				},
			},
		},
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
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func writeBlueprintComputedFieldsToResource(bm *cli.BlueprintModel, bp *cli.Blueprint) {
	bm.Identifier = types.StringValue(bp.Identifier)
	bm.CreatedAt = types.StringValue(bp.CreatedAt.String())
	bm.CreatedBy = types.StringValue(bp.CreatedBy)
	bm.UpdatedAt = types.StringValue(bp.UpdatedAt.String())
	bm.UpdatedBy = types.StringValue(bp.UpdatedBy)
}

func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
		fmt.Printf("Properties %+v\n", d.Properties)
		if d.Properties.StringProp != nil {
			for propIdentifier, prop := range d.Properties.StringProp {
				// 		fmt.Printf("Prop %+v\n", prop)
				props[propIdentifier] = cli.BlueprintProperty{
					Type:      "string",
					Title:     prop.Title.ValueString(),
					Format:    prop.Format.ValueString(),
					Default:   prop.Default.ValueString(),
					Icon:      prop.Icon.ValueString(),
					MinLength: int(prop.MinLength.ValueInt64()),
					MaxLength: int(prop.MaxLength.ValueInt64()),
					Pattern:   prop.Pattern.ValueString(),
				}
				if prop.Required.ValueBool() {
					required = append(required, propIdentifier)
				}
			}
		}

	}

	// for _, prop := range d.Properties {
	// 	if prop == "string_prop" {
	// 		continue
	// 	}
	// 	props[prop.Identifier.ValueString()] = cli.BlueprintProperty{
	// 		Type:      prop.Type.ValueString(),
	// 		Title:     prop.Title.ValueString(),
	// 		Format:    prop.Format.ValueString(),
	// 		Default:   prop.Default.ValueString(),
	// 		Icon:      prop.Icon.ValueString(),
	// 		MinLength: int(prop.MinLength.ValueInt64()),
	// 		MaxLength: int(prop.MaxLength.ValueInt64()),
	// 		Pattern:   prop.Pattern.ValueString(),
	// 	}
	// 	if prop.Required.ValueBool() {
	// 		required = append(required, prop.Identifier.ValueString())
	// 	}
	// }

	properties := props

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
	return b, nil
}
