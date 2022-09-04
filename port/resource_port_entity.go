package port

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

func newEntityResource() *schema.Resource {
	return &schema.Resource{
		Description:   "Port entity",
		CreateContext: createEntity,
		UpdateContext: createEntity,
		ReadContext:   readEntity,
		DeleteContext: deleteEntity,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The identifier of the entity",
				Optional:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "The display name of the entity",
				Required:    true,
			},
			"team": {
				Type:        schema.TypeString,
				Description: "The display name of the entity",
				Required:    true,
			},
			"blueprint": {
				Type:        schema.TypeString,
				Description: "The blueprint identifier the entity relates to",
				Required:    true,
			},
			"relations": {
				Description: "The other entities that are connected",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the relation",
						},
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the connected entity",
						},
					},
				},
				Optional: true,
			},
			"properties": {
				Description: "The metadata of the entity",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of this property",
						},
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"number", "string", "boolean", "array", "object"}, false),
							Optional:     true,
							Deprecated:   "property type is not required anymore",
							Description:  "The type of the property",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value for this property",
						},
						"items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of items, in case the type of this property is a list",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func deleteEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	err := c.DeleteEntity(ctx, d.Id(), d.Get("blueprint").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func convert(prop map[string]interface{}, bp *cli.Blueprint) (interface{}, error) {
	valType := prop["type"].(string)
	if valType == "" {
		if p, ok := bp.Schema.Properties[prop["name"].(string)]; ok {
			valType = p.Type
		} else {
			return nil, fmt.Errorf("no type found for property %s", prop["name"])
		}
	}
	switch valType {
	case "string", "number", "boolean":
		return prop["value"], nil
	case "array":
		return prop["items"], nil
	case "object":
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(prop["value"].(string)), &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return "", fmt.Errorf("unsupported type %s", valType)
}

func entityResourceToBody(d *schema.ResourceData, bp *cli.Blueprint) (*cli.Entity, error) {
	e := &cli.Entity{}
	if identifier, ok := d.GetOk("identifier"); ok {
		e.Identifier = identifier.(string)
	}
	id := d.Id()
	if id != "" {
		e.Identifier = id
	}
	e.Title = d.Get("title").(string)
	e.Blueprint = d.Get("blueprint").(string)
	rels := d.Get("relations").(*schema.Set)
	relations := make(map[string]string)
	for _, rel := range rels.List() {
		r := rel.(map[string]interface{})
		relations[r["name"].(string)] = r["identifier"].(string)
	}
	e.Relations = relations
	props := d.Get("properties").(*schema.Set)
	properties := make(map[string]interface{}, props.Len())
	for _, prop := range props.List() {
		p := prop.(map[string]interface{})
		propValue, err := convert(p, bp)
		if err != nil {
			return nil, err
		}
		properties[p["name"].(string)] = propValue
	}
	e.Properties = properties
	return e, nil
}

func writeEntityComputedFieldsToResource(d *schema.ResourceData, e *cli.Entity) {
	d.SetId(e.Identifier)
	d.Set("created_at", e.CreatedAt.String())
	d.Set("created_by", e.CreatedBy)
	d.Set("updated_at", e.UpdatedAt.String())
	d.Set("updated_by", e.UpdatedBy)
}

func writeEntityFieldsToResource(d *schema.ResourceData, e *cli.Entity) {
	d.SetId(e.Identifier)
	d.Set("title", e.Title)
	d.Set("created_at", e.CreatedAt.String())
	d.Set("created_by", e.CreatedBy)
	d.Set("updated_at", e.UpdatedAt.String())
	d.Set("updated_by", e.UpdatedBy)
	properties := schema.Set{F: func(i interface{}) int {
		name := (i.(map[string]interface{}))["name"].(string)
		return schema.HashString(name)
	}}
	for k, v := range e.Properties {
		if v == nil {
			continue
		}
		p := map[string]interface{}{}
		p["name"] = k
		switch t := v.(type) {
		case map[string]interface{}:
			js, _ := json.Marshal(&t)
			p["value"] = string(js)
		case []interface{}:
			p["items"] = t
		case float64:
			p["value"] = strconv.FormatFloat(t, 'f', -1, 64)
		case int:
			p["value"] = strconv.Itoa(t)
		case string:
			p["value"] = t
		case bool:
			p["value"] = "false"
			if t {
				p["value"] = "true"
			}
		}
		properties.Add(p)
	}
	d.Set("properties", &properties)
}

func createEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	bp, err := c.ReadBlueprint(ctx, d.Get("blueprint").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	e, err := entityResourceToBody(d, bp)
	if err != nil {
		return diag.FromErr(err)
	}
	en, err := c.CreateEntity(ctx, e)
	if err != nil {
		return diag.FromErr(err)
	}
	writeEntityComputedFieldsToResource(d, en)
	return diags
}

func readEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	e, err := c.ReadEntity(ctx, d.Id(), d.Get("blueprint").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	writeEntityFieldsToResource(d, e)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
