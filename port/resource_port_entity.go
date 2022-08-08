package port

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
							Required:     true,
							Description:  "The type of the properrty",
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
	client := m.(*resty.Client)
	url := "v0.1/entities/{identifier}"
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetPathParam("identifier", d.Id()).
		Delete(url)
	if err != nil {
		return diag.FromErr(err)
	}
	responseBody := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	if !(responseBody["ok"].(bool)) {
		return diag.FromErr(fmt.Errorf("failed to delete entity. got:\n%s", string(resp.Body())))
	}
	return diags
}

func convert(prop map[string]interface{}) (interface{}, error) {
	valType := prop["type"].(string)
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

func entityResourceToBody(d *schema.ResourceData) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	if identifier, ok := d.GetOk("identifier"); ok {
		body["identifier"] = identifier
	}
	id := d.Id()
	if id != "" {
		body["identifier"] = id
	}
	body["title"] = d.Get("title").(string)
	body["blueprint"] = d.Get("blueprint").(string)
	body["blueprintIdentifier"] = d.Get("blueprint").(string)
	rels := d.Get("relations").(*schema.Set)
	relations := make(map[string]string)
	for _, rel := range rels.List() {
		r := rel.(map[string]interface{})
		bpName := r["name"].(string)
		relID := r["identifier"].(string)
		relations[bpName] = relID
	}
	body["relations"] = relations
	props := d.Get("properties").(*schema.Set)
	properties := map[string]interface{}{}
	for _, prop := range props.List() {
		p := prop.(map[string]interface{})
		var propValue interface{}
		var err error
		propValue, err = convert(p)
		if err != nil {
			return nil, err
		}
		properties[p["name"].(string)] = propValue
	}
	body["properties"] = properties
	return body, nil
}

func writeEntityComputedFieldsToResource(d *schema.ResourceData, e Entity) {
	d.SetId(e.Identifier)
	d.Set("created_at", e.CreatedAt.String())
	d.Set("created_by", e.CreatedBy)
	d.Set("updated_at", e.UpdatedAt.String())
	d.Set("updated_by", e.UpdatedBy)
}

func writeEntityFieldsToResource(d *schema.ResourceData, e Entity) {
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
			p["type"] = "object"
			js, _ := json.Marshal(&t)
			p["value"] = string(js)
		case []interface{}:
			p["type"] = "array"
			p["items"] = t
		case float64:
			p["type"] = "number"
			p["value"] = strconv.FormatFloat(t, 'f', -1, 64)
		case int:
			p["type"] = "number"
			p["value"] = strconv.Itoa(t)
		case string:
			p["type"] = "string"
			p["value"] = t
		case bool:
			p["type"] = "boolean"
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
	client := m.(*resty.Client)
	url := "v0.1/entities"
	body, err := entityResourceToBody(d)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.R().
		SetBody(body).
		SetQueryParam("upsert", "true").
		Post(url)
	if err != nil {
		return diag.FromErr(err)
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return diag.FromErr(err)
	}
	if !pb.OK {
		return diag.FromErr(fmt.Errorf("failed to create entity, got: %s", resp.Body()))
	}
	writeEntityComputedFieldsToResource(d, pb.Entity)
	return diags
}

func readEntity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*resty.Client)
	url := "v0.1/entities/{identifier}"
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("exclude_mirror_properties", "true").
		SetPathParam("identifier", d.Id()).
		Get(url)
	if err != nil {
		return diag.FromErr(err)
	}
	var pb PortBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return diag.FromErr(err)
	}
	e := pb.Entity
	writeEntityFieldsToResource(d, e)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
