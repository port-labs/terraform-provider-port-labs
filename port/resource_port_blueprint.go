package port

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// type BlueprintModel struct {
// 	Identifier types.String `tfsdk:"identifier"`
// 	Title      types.String `tfsdk:"title"`
// 	Icon       types.String `tfsdk:"icon"`
// 	Description       types.String `tfsdk:"description"`
// 	CreatedAt         types.String `tfsdk:"created_at"`
// 	CreatedBy         types.String `tfsdk:"created_by"`
// 	UpdatedAt         types.String `tfsdk:"updated_at"`
// 	UpdatedBy         types.String `tfsdk:"updated_by"`
// 	ChangelogDestination *ChangelogDestinationModel `tfsdk:"changelog_destination"`

// }

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
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The creator of the blueprint",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The last update date of the blueprint",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "The last updater of the blueprint",
				Computed:            true,
			},
			"changelog_destination": schema.MapNestedAttribute{
				MarkdownDescription: "The changelog destination of the blueprint",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the changelog destination",
							Required:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The url of the changelog destination",
							Computed:            true,
						},
						"agent": schema.BoolAttribute{
							MarkdownDescription: "The agent of the changelog destination",
							Computed:            true,
						},
					},
				},
			},
			"string_prop": schema.MapNestedAttribute{
				MarkdownDescription: "The string property of the blueprint",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The identifier of the string property",
							Required:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The display name of the string property",
							Optional:            true,
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
	bm.Title = types.StringValue(b.Title)
	bm.Icon = types.StringValue(b.Icon)
	bm.Description = types.StringValue(b.Description)
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

	b, err := blueprintResourceToBody(data)

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

}

// 		p["default"] = value
// 	} else {
// 		mapDefault := make(map[string]string)
// 		mapDefault["value"] = value
// 		p["default_value"] = mapDefault
// 	}
// }

func blueprintResourceToBody(d *cli.BlueprintModel) (*cli.Blueprint, error) {
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
		b.ChangelogDestination.Type = d.ChangelogDestination.Type.ValueString()
		b.ChangelogDestination.Url = d.ChangelogDestination.Url.ValueString()
		b.ChangelogDestination.Agent = d.ChangelogDestination.Agent.ValueBool()
	} else {
		b.ChangelogDestination = nil
	}

	properties := props
	// var required []string
	// for _, prop := range props.List() {
	// 	p := prop.(map[string]interface{})
	// 	propFields := cli.BlueprintProperty{}
	// 	if t, ok := p["type"]; ok && t != "" {
	// 		propFields.Type = t.(string)
	// 	}
	// 	if t, ok := p["title"]; ok && t != "" {
	// 		propFields.Title = t.(string)
	// 	}
	// 	if d, ok := p["description"]; ok && d != "" {
	// 		propFields.Description = d.(string)
	// 	}

	// 	df, defaultOk := p["default"]
	// 	dv, defaultValueOk := p["default_value"].(map[string]interface{})
	// 	di, defaultItemsOk := p["default_items"].([]interface{})

	// 	if propFields.Type == "array" {
	// 		if (defaultValueOk && len(dv) != 0) || (defaultOk && df != "") {
	// 			return nil, fmt.Errorf("default or default_value can't be used when type is array for property %s", p["identifier"].(string))
	// 		}

	// 	} else {
	// 		if defaultItemsOk && len(di) != 0 {
	// 			return nil, fmt.Errorf("default_items can't be used when type is not array for property %s", p["identifier"].(string))
	// 		}
	// 		if (defaultValueOk && defaultOk) && (len(dv) != 0 && df != "") {
	// 			return nil, fmt.Errorf("default and default_value can't be used together for property %s", p["identifier"].(string))
	// 		}

	// 		if _, ok := dv["value"]; !ok && defaultValueOk && len(dv) != 0 {
	// 			return nil, fmt.Errorf("value key is missing in default_value for property %s", p["identifier"].(string))
	// 		}
	// 	}

	// 	if defaultItemsOk && len(di) != 0 && propFields.Type == "array" {
	// 		propFields.Default = di
	// 	}

	// 	if propFields.Type == "array" {
	// 		if i, ok := p["items"]; ok && i != nil {
	// 			items := make(map[string]any)
	// 			for key, value := range i.(map[string]any) {
	// 				items[key] = value.(string)
	// 			}
	// 			propFields.Items = items
	// 		}
	// 	}

	// 	if defaultOk && df != "" {
	// 		err := defaultResourceToBody(df.(string), &propFields)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	} else {
	// 		if defaultValueOk && len(dv) != 0 {
	// 			err := defaultResourceToBody(dv["value"].(string), &propFields)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 		}
	// 	}

	// 	if f, ok := p["format"]; ok && f != "" {
	// 		propFields.Format = f.(string)
	// 	}

	// 	if i, ok := p["max_length"]; ok && i != 0 {
	// 		propFields.MaxLength = i.(int)
	// 	}

	// 	if i, ok := p["min_length"]; ok && i != 0 {
	// 		propFields.MinLength = i.(int)
	// 	}

	// 	if i, ok := p["min_items"]; ok && i != 0 {
	// 		if propFields.Type != "array" {
	// 			return nil, fmt.Errorf("min_items can only be used when type is array for property %s", p["identifier"].(string))
	// 		}

	// 		propFields.MinItems = i.(int)
	// 	}

	// 	if i, ok := p["max_items"]; ok && i != 0 {
	// 		if propFields.Type != "array" {
	// 			return nil, fmt.Errorf("max_items can only be used when type is array for property %s", p["identifier"].(string))
	// 		}

	// 		propFields.MaxItems = i.(int)
	// 	}

	// 	if i, ok := p["icon"]; ok && i != "" {
	// 		propFields.Icon = i.(string)
	// 	}

	// 	if s, ok := p["spec"]; ok && s != "" {
	// 		propFields.Spec = s.(string)
	// 	}

	// 	if s, ok := p["spec_authentication"]; ok && len(s.([]any)) > 0 {
	// 		if propFields.Spec != "embedded-url" {
	// 			return nil, fmt.Errorf("spec_authentication can only be used when spec is embedded-url for property %s", p["identifier"].(string))
	// 		}

	// 		if propFields.SpecAuthentication == nil {
	// 			propFields.SpecAuthentication = &cli.SpecAuthentication{}
	// 		}

	// 		propFields.SpecAuthentication.TokenUrl = s.([]any)[0].(map[string]interface{})["token_url"].(string)
	// 		propFields.SpecAuthentication.AuthorizationUrl = s.([]any)[0].(map[string]interface{})["authorization_url"].(string)
	// 		propFields.SpecAuthentication.ClientId = s.([]any)[0].(map[string]interface{})["client_id"].(string)

	// 	}

	// 	if r, ok := p["required"]; ok && r.(bool) {
	// 		required = append(required, p["identifier"].(string))
	// 	}
	// 	if e, ok := p["enum"]; ok && len(e.([]interface{})) > 0 {
	// 		if propFields.Type != "number" && propFields.Type != "string" {
	// 			return nil, fmt.Errorf("enum can only be used when type is number or string for property %s", p["identifier"].(string))
	// 		}
	// 		for _, v := range e.([]interface{}) {
	// 			if propFields.Type == "number" {
	// 				enumValue, err := strconv.ParseInt(v.(string), 10, 0)
	// 				if err != nil {
	// 					return nil, fmt.Errorf("enum value %s is not a valid number for property %s", v.(string), p["identifier"].(string))
	// 				}
	// 				propFields.Enum = append(propFields.Enum, enumValue)
	// 			}
	// 			if propFields.Type == "string" {
	// 				enumValue := v.(string)
	// 				propFields.Enum = append(propFields.Enum, enumValue)
	// 			}

	// 		}

	// 	}
	// 	if e, ok := p["enum_colors"]; ok && e != nil {
	// 		enumColors := make(map[string]string)
	// 		for key, value := range e.(map[string]interface{}) {
	// 			enumColors[key] = value.(string)
	// 		}
	// 		propFields.EnumColors = enumColors
	// 	}
	// 	// TODO: remove the if statement when this issues is solved, https://github.com/hashicorp/terraform-plugin-sdk/pull/1042/files
	// 	if p["identifier"] != "" {
	// 		properties[p["identifier"].(string)] = propFields
	// 	}
	// }

	// mirrorProperties := mirrorProps
	// for _, prop := range mirrorProps.List() {
	// 	p := prop.(map[string]interface{})
	// 	propFields := cli.BlueprintMirrorProperty{}
	// 	if t, ok := p["title"]; ok && t != "" {
	// 		propFields.Title = t.(string)
	// 	}
	// 	if p, ok := p["path"]; ok && p != "" {
	// 		propFields.Path = p.(string)
	// 	}
	// 	mirrorProperties[p["identifier"].(string)] = propFields
	// }

	// calculationProperties := calcProps
	// for _, prop := range calcProps.List() {
	// 	p := prop.(map[string]interface{})
	// 	calcFields := cli.BlueprintCalculationProperty{}
	// 	if t, ok := p["type"]; ok && t != "" {
	// 		calcFields.Type = t.(string)
	// 	}
	// 	if t, ok := p["title"]; ok && t != "" {
	// 		calcFields.Title = t.(string)
	// 	}
	// 	if d, ok := p["description"]; ok && d != "" {
	// 		calcFields.Description = d.(string)
	// 	}
	// 	if f, ok := p["format"]; ok && f != "" {
	// 		calcFields.Format = f.(string)
	// 	}
	// 	if i, ok := p["icon"]; ok && i != "" {
	// 		calcFields.Icon = i.(string)
	// 	}
	// 	if r, ok := p["colorized"]; ok && r.(bool) {
	// 		calcFields.Colorized = r.(bool)
	// 	}
	// 	if e, ok := p["colors"]; ok && e != nil {
	// 		colors := make(map[string]string)
	// 		for key, value := range e.(map[string]interface{}) {
	// 			colors[key] = value.(string)
	// 		}
	// 		calcFields.Colors = colors
	// 	}
	// 	calcFields.Calculation = p["calculation"].(string)
	// 	// TODO: remove the if statement when this issues is solved, https://github.com/hashicorp/terraform-plugin-sdk/pull/1042/files
	// 	if p["identifier"] != "" {
	// 		calculationProperties[p["identifier"].(string)] = calcFields
	// 	}
	// }

	// rels := d.Get("relations").(*schema.Set)
	// relations := make(map[string]cli.Relation, unsafe.Sizeof(d.StringProp))
	// for _, rel := range rels.List() {
	// 	p := rel.(map[string]interface{})
	// 	relationFields := cli.Relation{}
	// 	if t, ok := p["required"]; ok && t != "" {
	// 		relationFields.Required = t.(bool)
	// 	}
	// 	if t, ok := p["many"]; ok && t != "" {
	// 		relationFields.Many = t.(bool)
	// 	}
	// 	if d, ok := p["title"]; ok && d != "" {
	// 		relationFields.Title = d.(string)
	// 	}
	// 	if d, ok := p["target"]; ok && d != "" {
	// 		relationFields.Target = d.(string)
	// 	}

	// 	relations[p["identifier"].(string)] = relationFields
	// }

	var arr []string = []string{}

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: arr}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
	return b, nil
}

// func deleteBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	err := c.DeleteBlueprint(ctx, d.Id())
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return diags
// }

// func createBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	b, err := blueprintResourceToBody(d)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	var bp *cli.Blueprint
// 	if d.Id() != "" {
// 		bp, err = c.UpdateBlueprint(ctx, b, d.Id())
// 	} else {
// 		bp, err = c.CreateBlueprint(ctx, b)
// 	}
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	writeBlueprintComputedFieldsToResource(d, bp)

// 	return diags
// }

// func updateBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	b, err := blueprintResourceToBody(d)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	var bp *cli.Blueprint
// 	if d.Id() != "" {
// 		bp, err = c.UpdateBlueprint(ctx, b, d.Id())
// 	} else {
// 		bp, err = c.CreateBlueprint(ctx, b)
// 	}
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	writeBlueprintComputedFieldsToResource(d, bp)
// 	return diags
// }
