package entity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EntityResource{}
var _ resource.ResourceWithImportState = &EntityResource{}

func NewEntityResource() resource.Resource {
	return &EntityResource{}
}

type EntityResource struct {
	portClient *cli.PortClient
}

func (r *EntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity"
}

func (r *EntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *EntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EntityModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := data.Blueprint.ValueString()
	e, statusCode, err := r.portClient.ReadEntity(ctx, data.Blueprint.ValueString(), data.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	err = writeEntityFieldsToResource(ctx, e, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed writing blueprint fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeEntityFieldsToResource(ctx context.Context, e *cli.Entity, blueprint string) error {
	return nil
}
func (r *EntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !data.RunID.IsNull() {
		runID = data.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	data.ID = types.StringValue(e.Identifier)

	writeEntityComputedFieldsToResource(data, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func entityResourceToBody(em *EntityModel, bp *cli.Blueprint) (*cli.Entity, error) {
	e := &cli.Entity{}
	return e, nil
}

func writeEntityComputedFieldsToResource(data *EntityModel, e *cli.Entity) {
	data.Identifier = types.StringValue(e.Identifier)
	data.CreatedAt = types.StringValue(e.CreatedAt.String())
	data.CreatedBy = types.StringValue(e.CreatedBy)
	data.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	data.UpdatedBy = types.StringValue(e.UpdatedBy)
}

func (r *EntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO
}

func (r *EntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}

func (r *EntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO
}

// func newEntityResource() *schema.Resource {
// 	return &schema.Resource{
// 		Description:   "Port entity",
// 		CreateContext: createEntity,
// 		UpdateContext: createEntity,
// 		ReadContext:   readEntity,
// 		DeleteContext: deleteEntity,
// 		Schema: map[string]*schema.Schema{
// 			"identifier": {
// 				Type:        schema.TypeString,
// 				Description: "The identifier of the entity",
// 				Optional:    true,
// 			},
// 			"title": {
// 				Type:        schema.TypeString,
// 				Description: "The display name of the entity",
// 				Required:    true,
// 			},
// 			"run_id": {
// 				Type:        schema.TypeString,
// 				Description: "The runID of the action run that created the entity",
// 				Optional:    true,
// 			},
// 			"team": {
// 				Type:          schema.TypeString,
// 				Description:   "The team related to the entity",
// 				Optional:      true,
// 				ConflictsWith: []string{"teams"},
// 			},
// 			"teams": {
// 				Type:        schema.TypeSet,
// 				Description: "The teams related to the entity",
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 				Optional:      true,
// 				ConflictsWith: []string{"team"},
// 			},
// 			"blueprint": {
// 				Type:        schema.TypeString,
// 				Description: "The blueprint identifier the entity relates to",
// 				Required:    true,
// 			},
// 			"relations": {
// 				Description: "The other entities that are connected",
// 				Type:        schema.TypeSet,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "The name of the relation",
// 						},
// 						"identifier": {
// 							Type:        schema.TypeString,
// 							Optional:    true,
// 							Description: "The id of the connected entity",
// 						},
// 						"identifiers": {
// 							Type:        schema.TypeSet,
// 							Optional:    true,
// 							Description: "The ids of the connected entities",
// 							Elem: &schema.Schema{
// 								Type: schema.TypeString,
// 							},
// 						},
// 					},
// 				},
// 				Optional: true,
// 			},
// 			"properties": {
// 				Description: "The metadata of the entity",
// 				Type:        schema.TypeSet,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "The name of this property",
// 						},
// 						"type": {
// 							Type:         schema.TypeString,
// 							ValidateFunc: validation.StringInSlice([]string{"number", "string", "boolean", "array", "object"}, false),
// 							Optional:     true,
// 							Deprecated:   "property type is not required anymore",
// 							Description:  "The type of the property",
// 						},
// 						"value": {
// 							Type:        schema.TypeString,
// 							Optional:    true,
// 							Description: "The value for this property",
// 						},
// 						"items": {
// 							Type:        schema.TypeList,
// 							Optional:    true,
// 							Description: "The list of items, in case the type of this property is a list",
// 							Elem: &schema.Schema{
// 								Type: schema.TypeString,
// 							},
// 						},
// 					},
// 				},
// 				Required: true,
// 			},
// 			"created_at": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"updated_at": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"created_by": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"updated_by": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 		},
// 		Importer: &schema.ResourceImporter{
// 			StateContext: schema.ImportStatePassthroughContext,
// 		},
// 	}
// }

// func deleteEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	id := d.Id()
// 	entityIdentifier := id
// 	blueprintIdentifier := d.Get("blueprint").(string)
// 	err := c.DeleteEntity(ctx, entityIdentifier, blueprintIdentifier)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return diags
// }

// func convert(prop map[string]interface{}, bp *cli.Blueprint) (interface{}, error) {
// 	valType := prop["type"].(string)
// 	if valType == "" {
// 		if p, ok := bp.Schema.Properties[prop["name"].(string)]; ok {
// 			valType = p.Type
// 		} else {
// 			return nil, fmt.Errorf("no type found for property %s", prop["name"])
// 		}
// 	}
// 	switch valType {
// 	case "string", "number", "boolean":
// 		return prop["value"], nil
// 	case "array":
// 		var itemsArray []interface{}
// 		if itemsType, ok := bp.Schema.Properties[prop["name"].(string)].Items["type"]; ok {

// 			for _, item := range prop["items"].([]interface{}) {
// 				convertItem := item.(string)
// 				switch itemsType {
// 				case "number":
// 					convertItem, err := strconv.ParseFloat(item.(string), 64)
// 					if err != nil {
// 						return nil, fmt.Errorf("failed to convert item to number")
// 					}
// 					itemsArray = append(itemsArray, convertItem)
// 				case "boolean":
// 					convertItem, err := strconv.ParseBool(item.(string))
// 					if err != nil {
// 						return nil, fmt.Errorf("failed to convert item to boolean")
// 					}
// 					itemsArray = append(itemsArray, convertItem)
// 				case "object":
// 					convertItem := make(map[string]interface{})
// 					err := json.Unmarshal([]byte(item.(string)), &convertItem)
// 					if err != nil {
// 						return nil, fmt.Errorf("failed to convert item to object")
// 					}
// 					itemsArray = append(itemsArray, convertItem)
// 				default:
// 					itemsArray = append(itemsArray, convertItem)
// 				}
// 			}

// 			return itemsArray, nil
// 		} else {
// 			return prop["items"], nil
// 		}

// 	case "object":
// 		obj := make(map[string]interface{})
// 		err := json.Unmarshal([]byte(prop["value"].(string)), &obj)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return obj, nil
// 	}
// 	return "", fmt.Errorf("unsupported type %s", valType)
// }

// func validateRelation(rel map[string]interface{}) error {
// 	if rel["identifier"] == "" && len(rel["identifiers"].(*schema.Set).List()) == 0 {
// 		return fmt.Errorf("either relation identifier or identifiers is required for %s", rel["name"])
// 	}

// 	if rel["identifier"] != "" && len(rel["identifiers"].(*schema.Set).List()) > 0 {
// 		return fmt.Errorf("either relation identifier or identifiers is required for %s", rel["name"])
// 	}

// 	return nil
// }

// func entityResourceToBody(d *schema.ResourceData, bp *cli.Blueprint) (*cli.Entity, error) {
// 	e := &cli.Entity{}
// 	if identifier, ok := d.GetOk("identifier"); ok {
// 		e.Identifier = identifier.(string)
// 	}
// 	id := d.Id()
// 	if id != "" {
// 		e.Identifier = id
// 	}

// 	e.Title = d.Get("title").(string)
// 	e.Blueprint = d.Get("blueprint").(string)

// 	teams := []string{}

// 	if team, ok := d.GetOk("team"); ok {
// 		teams = append(teams, team.(string))

// 	}

// 	if resourceTeams, ok := d.Get("teams").(*schema.Set); ok {
// 		for _, team := range resourceTeams.List() {
// 			teams = append(teams, team.(string))
// 		}
// 	}
// 	e.Team = teams

// 	rels := d.Get("relations").(*schema.Set)
// 	relations := make(map[string]interface{})
// 	for _, rel := range rels.List() {
// 		r := rel.(map[string]interface{})
// 		identifier := r["identifier"].(string)
// 		identifiers := r["identifiers"].(*schema.Set).List()
// 		err := validateRelation(r)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if identifier != "" {
// 			relations[r["name"].(string)] = identifier
// 		} else {
// 			relations[r["name"].(string)] = identifiers
// 		}
// 	}

// 	e.Relations = relations
// 	props := d.Get("properties").(*schema.Set)
// 	properties := make(map[string]interface{}, props.Len())
// 	for _, prop := range props.List() {
// 		p := prop.(map[string]interface{})
// 		propValue, err := convert(p, bp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		properties[p["name"].(string)] = propValue
// 	}
// 	e.Properties = properties
// 	return e, nil
// }

// func writeEntityComputedFieldsToResource(d *schema.ResourceData, e *cli.Entity) {
// 	d.SetId(e.Identifier)
// 	d.Set("created_at", e.CreatedAt.String())
// 	d.Set("created_by", e.CreatedBy)
// 	d.Set("updated_at", e.UpdatedAt.String())
// 	d.Set("updated_by", e.UpdatedBy)
// }

// func writeEntityFieldsToResource(d *schema.ResourceData, e *cli.Entity, blueprintIdentifier string) {
// 	d.SetId(e.Identifier)
// 	d.Set("title", e.Title)
// 	d.Set("blueprint", e.Blueprint)

// 	entityTeams := e.Team
// 	if len(entityTeams) > 0 {
// 		team := d.Get("team")

// 		if team != "" {
// 			d.Set("team", e.Team[0])
// 		} else {
// 			d.Set("teams", e.Team)
// 		}
// 	}

// 	d.Set("created_at", e.CreatedAt.String())
// 	d.Set("created_by", e.CreatedBy)
// 	d.Set("updated_at", e.UpdatedAt.String())
// 	d.Set("updated_by", e.UpdatedBy)
// 	properties := schema.Set{F: func(i interface{}) int {
// 		name := (i.(map[string]interface{}))["name"].(string)
// 		return schema.HashString(name)
// 	}}
// 	for k, v := range e.Properties {
// 		if v == nil {
// 			continue
// 		}
// 		p := map[string]interface{}{}
// 		p["name"] = k
// 		switch t := v.(type) {
// 		case map[string]interface{}:
// 			js, _ := json.Marshal(&t)
// 			p["value"] = string(js)
// 		case []interface{}:
// 			p["items"] = t
// 		case float64:
// 			p["value"] = strconv.FormatFloat(t, 'f', -1, 64)
// 		case int:
// 			p["value"] = strconv.Itoa(t)
// 		case string:
// 			p["value"] = t
// 		case bool:
// 			p["value"] = "false"
// 			if t {
// 				p["value"] = "true"
// 			}
// 		}
// 		properties.Add(p)
// 	}
// 	d.Set("properties", &properties)

// 	relations := schema.Set{F: func(i interface{}) int {
// 		name := (i.(map[string]interface{}))["name"].(string)
// 		return schema.HashString(name)
// 	}}

// 	for k, v := range e.Relations {
// 		if v == nil {
// 			continue
// 		}
// 		r := map[string]interface{}{}
// 		r["name"] = k
// 		switch t := v.(type) {
// 		case []interface{}:
// 			r["identifiers"] = t
// 		case string:
// 			r["identifier"] = t
// 		}
// 		relations.Add(r)

// 	}
// 	d.Set("relations", &relations)
// }

// func createEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	bp, _, err := c.ReadBlueprint(ctx, d.Get("blueprint").(string))
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	e, err := entityResourceToBody(d, bp)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	runID := ""
// 	if rid, ok := d.GetOk("run_id"); ok {
// 		runID = rid.(string)
// 	}
// 	en, err := c.CreateEntity(ctx, e, runID)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	writeEntityComputedFieldsToResource(d, en)
// 	return diags
// }

// func readEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	c := m.(*cli.PortClient)
// 	id := d.Id()
// 	blueprintIdentifier := d.Get("blueprint").(string)
// 	entityIdentifier := d.Id()

// 	if strings.Contains(id, ":") {
// 		parts := strings.SplitN(id, ":", 2)
// 		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
// 			return diag.FromErr(fmt.Errorf("unexpected format of ID (%s), expected blueprintId:entityId", id))
// 		}
// 		blueprintIdentifier = parts[0]
// 		entityIdentifier = parts[1]
// 	}

// 	e, statusCode, err := c.ReadEntity(ctx, entityIdentifier, blueprintIdentifier)
// 	if err != nil {
// 		if statusCode == 404 {
// 			d.SetId("")
// 			return diags
// 		}

// 		return diag.FromErr(err)
// 	}
// 	writeEntityFieldsToResource(d, e, blueprintIdentifier)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return diags
// }
