package port

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
	"github.com/samber/lo"
)

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
				Type:        schema.TypeString,
				Description: "The icon of the blueprint",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the blueprint",
				Optional:    true,
			},
			"relations": {
				Description: "The blueprints that are connected to this blueprint",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier of the relation",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
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
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the relation is many",
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
						"icon": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The icon of the property",
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
						"items": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A metadata of an array's items, in case the type is an array",
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
						"default_items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of items, in case the type of default property is a list",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The format of the Property",
						},
						"spec": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"async-api", "open-api", "embedded-url"}, false),
							Description:  "The specification of the property, one of \"async-api\", \"open-api\", \"embedded-url\"",
						},
						"enum": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of allowed values for the property",
						},
						"enum_colors": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "A map of colors for the enum values",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the property is required",
						},
					},
				},
				Required: true,
			},
			"mirror_properties": {
				Type:        schema.TypeSet,
				Description: "When two Blueprints are connected via a Relation, a new set of properties becomes available to Entities in the source Blueprint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier of the property",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of this property",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path of the realtions towards the property",
						},
					},
				},
				Optional: true,
			},
			"calculation_properties": {
				Type:        schema.TypeSet,
				Description: "A set of properties that are calculated upon entity's regular properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier of the property",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of this property",
						},
						"calculation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A jq expression that calculates the value of the property, for instance \"'https://grafana.' + .identifier\"",
						},
						"icon": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The icon of the property",
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
						"format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The format of the Property",
						},
						"colorized": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the property is colorized",
						},
						"colors": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "A map of colors for the property",
						},
					},
				},
				Optional: true,
			},
			"changelog_destination": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Description: "Blueprints changelog destination, Supports WEBHOOK and KAFKA",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Changelog's destination one of WEBHOOK or KAFKA",
							ValidateFunc: validation.StringInSlice([]string{"WEBHOOK", "KAFKA"}, false),
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required when selecting type WEBHOOK. The URL to which the changelog is dispatched",
						},
					},
				},
				Optional: true,
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
	b, statusCode, err := c.ReadBlueprint(ctx, d.Id())
	if err != nil {
		if statusCode == 404 {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	writeBlueprintFieldsToResource(d, b)

	return diags
}

func writeBlueprintFieldsToResource(d *schema.ResourceData, b *cli.Blueprint) {
	d.SetId(b.Identifier)
	d.Set("title", b.Title)
	d.Set("icon", b.Icon)
	d.Set("description", b.Description)
	d.Set("created_at", b.CreatedAt.String())
	d.Set("created_by", b.CreatedBy)
	d.Set("updated_at", b.UpdatedAt.String())
	d.Set("updated_by", b.UpdatedBy)
	if b.ChangelogDestination != nil {
		d.Set("changelog_destination", []any{map[string]any{
			"type": b.ChangelogDestination.Type,
			"url":  b.ChangelogDestination.Url,
		}})
	}
	properties := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}

	mirrorPoperties := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}

	calculationProperties := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}
	for k, v := range b.Schema.Properties {
		p := map[string]interface{}{}
		p["identifier"] = k
		p["title"] = v.Title
		p["type"] = v.Type
		p["items"] = v.Items
		p["description"] = v.Description
		p["format"] = v.Format
		p["icon"] = v.Icon
		p["spec"] = v.Spec
		p["enum"] = v.Enum
		p["enum_colors"] = v.EnumColors
		if lo.Contains(b.Schema.Required, k) {
			p["required"] = true
		} else {
			p["required"] = false
		}
		if v.Default != nil {
			switch t := v.Default.(type) {
			case map[string]interface{}:
				js, _ := json.Marshal(&t)
				p["default"] = string(js)
			case []interface{}:
				p["default_items"] = t
			case float64:
				p["default"] = strconv.FormatFloat(t, 'f', -1, 64)
			case int:
				p["default"] = strconv.Itoa(t)
			case string:
				p["default"] = t
			case bool:
				p["default"] = "false"
				if t {
					p["default"] = "true"
				}
			}
		}

		properties.Add(p)
	}

	for k, v := range b.MirrorProperties {
		p := map[string]interface{}{}
		p["identifier"] = k
		p["title"] = v.Title
		p["path"] = v.Path
		mirrorPoperties.Add(p)
	}

	for k, v := range b.CalculationProperties {
		p := map[string]interface{}{}
		p["identifier"] = k
		p["title"] = v.Title
		p["description"] = v.Description
		p["icon"] = v.Icon
		p["calculation"] = v.Calculation
		p["type"] = v.Type
		p["format"] = v.Format
		p["colorized"] = v.Colorized
		p["colors"] = v.Colors

		calculationProperties.Add(p)
	}

	d.Set("properties", &properties)
	d.Set("mirror_properties", &mirrorPoperties)
	d.Set("calculation_properties", &calculationProperties)
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
	b.Description = d.Get("description").(string)
	props := d.Get("properties").(*schema.Set)
	mirrorProps := d.Get("mirror_properties").(*schema.Set)
	calcProps := d.Get("calculation_properties").(*schema.Set)

	if changelogDestination, ok := d.GetOk("changelog_destination"); ok {
		if b.ChangelogDestination == nil {
			b.ChangelogDestination = &cli.ChangelogDestination{}
		}
		b.ChangelogDestination.Type = changelogDestination.([]any)[0].(map[string]interface{})["type"].(string)
		b.ChangelogDestination.Url = changelogDestination.([]any)[0].(map[string]interface{})["url"].(string)
	}

	properties := make(map[string]cli.BlueprintProperty, props.Len())
	var required []string
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
		switch propFields.Type {
		case "string":
			if d, ok := p["default"]; ok && d.(string) != "" {
				propFields.Default = d.(string)
			}
		case "number":
			if d, ok := p["default"]; ok && d.(string) != "" {
				defaultNum, err := strconv.ParseInt(d.(string), 10, 0)
				if err != nil {
					return nil, err
				}
				propFields.Default = defaultNum
			}
		case "boolean":
			if d, ok := p["default"]; ok && d.(string) != "" {
				defaultBool, err := strconv.ParseBool(d.(string))
				if err != nil {
					return nil, err
				}
				propFields.Default = defaultBool
			}
		case "array":
			if d, ok := p["default_items"]; ok && d != nil {
				propFields.Default = d
			}
			if i, ok := p["items"]; ok && i != nil {
				items := make(map[string]any)
				for key, value := range i.(map[string]any) {
					items[key] = value.(string)
				}
				propFields.Items = items
			}
		case "object":
			if d, ok := p["default"]; ok && d.(string) != "" {
				defaultObj := make(map[string]interface{})
				err := json.Unmarshal([]byte(d.(string)), &defaultObj)
				if err != nil {
					return nil, err
				}
				propFields.Default = defaultObj
			}
		}
		if f, ok := p["format"]; ok && f != "" {
			propFields.Format = f.(string)
		}
		if i, ok := p["icon"]; ok && i != "" {
			propFields.Icon = i.(string)
		}

		if s, ok := p["spec"]; ok && s != "" {
			propFields.Spec = s.(string)
		}

		if r, ok := p["required"]; ok && r.(bool) {
			required = append(required, p["identifier"].(string))
		}
		if e, ok := p["enum"]; ok && e != nil {
			for _, v := range e.([]interface{}) {
				propFields.Enum = append(propFields.Enum, v.(string))
			}
		}
		if e, ok := p["enum_colors"]; ok && e != nil {
			enumColors := make(map[string]string)
			for key, value := range e.(map[string]interface{}) {
				enumColors[key] = value.(string)
			}
			propFields.EnumColors = enumColors
		}
		// TODO: remove the if statement when this issues is solved, https://github.com/hashicorp/terraform-plugin-sdk/pull/1042/files
		if p["identifier"] != "" {
			properties[p["identifier"].(string)] = propFields
		}
	}

	mirrorProperties := make(map[string]cli.BlueprintMirrorProperty, mirrorProps.Len())
	for _, prop := range mirrorProps.List() {
		p := prop.(map[string]interface{})
		propFields := cli.BlueprintMirrorProperty{}
		if t, ok := p["title"]; ok && t != "" {
			propFields.Title = t.(string)
		}
		if p, ok := p["path"]; ok && p != "" {
			propFields.Path = p.(string)
		}
		mirrorProperties[p["identifier"].(string)] = propFields
	}

	calculationProperties := make(map[string]cli.BlueprintCalculationProperty, calcProps.Len())
	for _, prop := range calcProps.List() {
		p := prop.(map[string]interface{})
		calcFields := cli.BlueprintCalculationProperty{}
		if t, ok := p["type"]; ok && t != "" {
			calcFields.Type = t.(string)
		}
		if t, ok := p["title"]; ok && t != "" {
			calcFields.Title = t.(string)
		}
		if d, ok := p["description"]; ok && d != "" {
			calcFields.Description = d.(string)
		}
		if f, ok := p["format"]; ok && f != "" {
			calcFields.Format = f.(string)
		}
		if i, ok := p["icon"]; ok && i != "" {
			calcFields.Icon = i.(string)
		}
		if r, ok := p["colorized"]; ok && r.(bool) {
			calcFields.Colorized = r.(bool)
		}
		if e, ok := p["colors"]; ok && e != nil {
			colors := make(map[string]string)
			for key, value := range e.(map[string]interface{}) {
				colors[key] = value.(string)
			}
			calcFields.Colors = colors
		}
		calcFields.Calculation = p["calculation"].(string)
		// TODO: remove the if statement when this issues is solved, https://github.com/hashicorp/terraform-plugin-sdk/pull/1042/files
		if p["identifier"] != "" {
			calculationProperties[p["identifier"].(string)] = calcFields
		}
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

	b.Schema = cli.BlueprintSchema{Properties: properties, Required: required}
	b.Relations = relations
	b.MirrorProperties = mirrorProperties
	b.CalculationProperties = calculationProperties
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
