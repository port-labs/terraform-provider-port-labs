package port

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

var ICONS = []string{"Actions", "Airflow", "Ansible", "Argo", "AuditLog", "Aws", "Azure", "Blueprint", "Bucket", "Cloud", "Cluster", "CPU", "Customer", "Datadog", "Day2Operation", "DefaultEntity", "DefaultProperty", "DeployedAt", "Deployment", "DevopsTool", "Docs", "Environment", "Git", "Github", "GitVersion", "GoogleCloud", "GPU", "Grafana", "Infinity", "Jenkins", "Lambda", "Link", "Lock", "Microservice", "Moon", "Node", "Okta", "Package", "Permission", "Relic", "Server", "Service", "Team", "Terraform", "User"}

func newBlueprintResource() *schema.Resource {
	return &schema.Resource{
		Description:   "Port blueprint",
		CreateContext: createBlueprint,
		UpdateContext: updateBlueprint,
		ReadContext:   readBlueprint,
		DeleteContext: deleteBlueprint,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The identifier of the blueprint",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "The display name of the blueprint",
				Required:    true,
			},
			"data_source": {
				Type:        schema.TypeString,
				Description: "The data source for entities of this blueprint",
				Default:     "Port",
				Optional:    true,
				Deprecated:  "Data source is ignored",
			},
			"icon": {
				Type:         schema.TypeString,
				Description:  "The icon of the blueprint",
				ValidateFunc: validation.StringInSlice(ICONS, false),
				Required:     true,
			},
			"relations": {
				Description: "The blueprints that are connected to this blueprint",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The identifier of the relation",
						},
						"title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The display name of the relation",
						},
						"target": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the connected blueprint",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the relation is required",
						},
						"many": {
							Type:     schema.TypeBool,
							Optional: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								if i.(bool) {
									return diag.Errorf("Many relations are not supported")
								}
								return nil
							},
							Description: "Unsupported ATM.\nWhether or not the relation is many",
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
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier of the property",
						},
						"title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of this property",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the property",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the property",
						},
						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The default value of the property",
						},
						"format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The format of the Property",
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

func readBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	b, err := c.ReadBlueprint(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeBlueprintFieldsToResource(d, b)

	return diags
}

func writeBlueprintRelationsToResource(d *schema.ResourceData, relations []*cli.Relation) {
	rels := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}
	for _, v := range relations {
		r := map[string]interface{}{
			"identifier": v.Identifier,
			"title":      v.Title,
			"target":     v.Target,
			"required":   v.Required,
			"many":       v.Many,
		}
		rels.Add(r)
	}
	d.Set("relations", &rels)
}

func writeBlueprintFieldsToResource(d *schema.ResourceData, b *cli.Blueprint) {
	d.SetId(b.Identifier)
	d.Set("title", b.Title)
	d.Set("icon", b.Icon)
	d.Set("created_at", b.CreatedAt.String())
	d.Set("created_by", b.CreatedBy)
	d.Set("updated_at", b.UpdatedAt.String())
	d.Set("updated_by", b.UpdatedBy)
	properties := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}
	for k, v := range b.Schema.Properties {
		p := map[string]interface{}{}
		p["identifier"] = k
		p["title"] = v.Title
		p["type"] = v.Type
		p["description"] = v.Description
		p["default"] = v.Default
		p["format"] = v.Format
		properties.Add(p)
	}
	d.Set("properties", &properties)
}

func blueprintResourceToBody(d *schema.ResourceData) (*cli.Blueprint, error) {
	b := &cli.Blueprint{}
	if identifier, ok := d.GetOk("identifier"); ok {
		b.Identifier = identifier.(string)
	}
	id := d.Id()
	if id != "" {
		b.Identifier = id
	}

	b.Title = d.Get("title").(string)
	b.Icon = d.Get("icon").(string)

	props := d.Get("properties").(*schema.Set)
	properties := make(map[string]cli.BlueprintProperty, props.Len())
	for _, prop := range props.List() {
		p := prop.(map[string]interface{})
		propFields := cli.BlueprintProperty{}
		if t, ok := p["type"]; ok && t != "" {
			propFields.Type = t.(string)
		}
		if t, ok := p["title"]; ok && t != "" {
			propFields.Title = t.(string)
		}
		if d, ok := p["description"]; ok && d != "" {
			propFields.Description = d.(string)
		}
		if d, ok := p["default"]; ok && d != "" {
			propFields.Default = d.(string)
		}
		if f, ok := p["format"]; ok && f != "" {
			propFields.Format = f.(string)
		}
		properties[p["identifier"].(string)] = propFields
	}

	rels := d.Get("relations").(*schema.Set)
	relations := make(map[string]cli.Relation, props.Len())
	for _, rel := range rels.List() {
		p := rel.(map[string]interface{})
		relationFields := cli.Relation{}
		if t, ok := p["required"]; ok && t != "" {
			relationFields.Required = t.(bool)
		}
		if t, ok := p["many"]; ok && t != "" {
			relationFields.Many = t.(bool)
		}
		if d, ok := p["title"]; ok && d != "" {
			relationFields.Title = d.(string)
		}
		if d, ok := p["target"]; ok && d != "" {
			relationFields.Target = d.(string)
		}

		relations[p["identifier"].(string)] = relationFields
	}

	b.Schema = cli.BlueprintSchema{Properties: properties}
	b.Relations = relations
	return b, nil
}

func deleteBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	err := c.DeleteBlueprint(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func getRelations(d *schema.ResourceData) (rel []*cli.Relation) {
	relations, ok := d.GetOk("relations")
	if !ok {
		return nil
	}
	for _, relation := range relations.(*schema.Set).List() {
		relation := relation.(map[string]interface{})
		r := &cli.Relation{}
		if t, ok := relation["title"]; ok {
			r.Title = t.(string)
		}
		if t, ok := relation["target"]; ok {
			r.Target = t.(string)
		}
		if i, ok := relation["identifier"]; ok {
			r.Identifier = i.(string)
		}
		if req, ok := relation["required"]; ok {
			r.Required = req.(bool)
		}
		if m, ok := relation["many"]; ok {
			r.Many = m.(bool)
		}
		rel = append(rel, r)
	}
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	b, err := blueprintResourceToBody(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var bp *cli.Blueprint
	if d.Id() != "" {
		bp, err = c.UpdateBlueprint(ctx, b, d.Id())
	} else {
		bp, err = c.CreateBlueprint(ctx, b)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	writeBlueprintComputedFieldsToResource(d, bp)

	return diags
}

func updateBlueprint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	b, err := blueprintResourceToBody(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var bp *cli.Blueprint
	if d.Id() != "" {
		bp, err = c.UpdateBlueprint(ctx, b, d.Id())
	} else {
		bp, err = c.CreateBlueprint(ctx, b)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	writeBlueprintComputedFieldsToResource(d, bp)
	return diags
}

func writeBlueprintComputedFieldsToResource(d *schema.ResourceData, b *cli.Blueprint) {
	d.SetId(b.Identifier)
	d.Set("created_at", b.CreatedAt.String())
	d.Set("created_by", b.CreatedBy)
	d.Set("updated_at", b.UpdatedAt.String())
	d.Set("updated_by", b.UpdatedBy)
}
