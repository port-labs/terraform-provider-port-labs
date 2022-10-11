package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

func newActionResource() *schema.Resource {
	return &schema.Resource{
		Description:   "Port action",
		CreateContext: createAction,
		UpdateContext: createAction,
		ReadContext:   readAction,
		DeleteContext: deleteAction,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The identifier of the action",
				Required:    true,
			},
			"blueprint_identifier": {
				Type:        schema.TypeString,
				Description: "The identifier of the blueprint",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "The display name of the action",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the action",
				Optional:    true,
			},
			"icon": {
				Type:         schema.TypeString,
				Description:  "The icon of the action",
				ValidateFunc: validation.StringInSlice(ICONS, false),
				Optional:     true,
			},
			"user_properties": {
				Description: "The input properties of the action",
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
							Description: "A nicely written name for the property",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the property",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description of the property. This value is visible to users when hovering on the info icon in the UI. It provides detailed information about the use of a specific property.",
						},
						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A default value for this property in case an entity is created without explicitly providing a value.",
						},
						"format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A specific data format to pair with some of the available types",
						},
						"blueprint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "When selecting format 'entity', the identifier of the target blueprint",
						},
						"pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A regular expression (regex) pattern to specify the set of allowed values for the property",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the property is required or not",
						},
						"enum": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of allowed values for the property",
						},
					},
				},
				Optional: true,
			},
			"invocation_method": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Description: "The methods the action is dispatched in, Supports WEBHOOK and KAFKA",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "How to invoke the action using WEBHOOK or KAFKA",
							ValidateFunc: validation.StringInSlice([]string{"WEBHOOK", "KAFKA"}, false),
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required when selecting type WEBHOOK. The URL to which the action is dispatched",
						},
					},
				},
				Required: true,
			},
			"trigger": {
				Type:         schema.TypeString,
				Description:  "The type of the action, one of CREATE, DAY-2, DELETE",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATE", "DAY-2", "DELETE"}, false),
			},
		},
	}
}

func readAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	action, err := c.ReadAction(ctx, d.Get("blueprint_identifier").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeActionFieldsToResource(d, action)
	return diags
}

func writeActionFieldsToResource(d *schema.ResourceData, action *cli.Action) {
	d.SetId(action.Identifier)
	d.Set("title", action.Title)
	d.Set("icon", action.Icon)
	d.Set("description", action.Description)
	d.Set("invocation_method", []any{map[string]any{
		"type": action.InvocationMethod.Type,
		"url":  action.InvocationMethod.Url,
	}})

	d.Set("trigger", action.Trigger)
	properties := schema.Set{F: func(i interface{}) int {
		id := (i.(map[string]interface{}))["identifier"].(string)
		return schema.HashString(id)
	}}
	for k, v := range action.UserInputs.Properties {
		p := map[string]interface{}{}
		p["identifier"] = k
		p["title"] = v.Title
		p["type"] = v.Type
		p["description"] = v.Description
		p["default"] = v.Default
		p["format"] = v.Format
		p["pattern"] = v.Pattern
		p["blueprint"] = v.Blueprint
		p["enum"] = v.Enum
		if contains(action.UserInputs.Required, k) {
			p["required"] = true
		} else {
			p["required"] = false
		}
		properties.Add(p)
	}
	d.Set("user_properties", &properties)
}

func actionResourceToBody(d *schema.ResourceData) (*cli.Action, error) {
	action := &cli.Action{}
	if identifier, ok := d.GetOk("identifier"); ok {
		action.Identifier = identifier.(string)
	}

	action.Title = d.Get("title").(string)
	action.Icon = d.Get("icon").(string)
	action.Description = d.Get("description").(string)
	if invocationMethod, ok := d.GetOk("invocation_method"); ok {
		if action.InvocationMethod == nil {
			action.InvocationMethod = &cli.InvocationMethod{}
		}

		action.InvocationMethod.Type = invocationMethod.([]any)[0].(map[string]interface{})["type"].(string)
		action.InvocationMethod.Url = invocationMethod.([]any)[0].(map[string]interface{})["url"].(string)
	}
	action.Trigger = d.Get("trigger").(string)

	props := d.Get("user_properties").(*schema.Set)
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
		if d, ok := p["default"]; ok && d != "" {
			propFields.Default = d.(string)
		}
		if f, ok := p["format"]; ok && f != "" {
			propFields.Format = f.(string)
		}
		if b, ok := p["blueprint"]; ok && b != "" {
			propFields.Blueprint = b.(string)
		}
		if p, ok := p["pattern"]; ok && p != "" {
			propFields.Pattern = p.(string)
		}
		if r, ok := p["required"]; ok && r.(bool) {
			required = append(required, p["identifier"].(string))
		}
		if e, ok := p["enum"]; ok && e != nil {
			for _, v := range e.([]interface{}) {
				propFields.Enum = append(propFields.Enum, v.(string))
			}
		}
		properties[p["identifier"].(string)] = propFields
	}

	action.UserInputs = cli.ActionUserInputs{Properties: properties, Required: required}
	return action, nil
}

func deleteAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	err := c.DeleteAction(ctx, d.Get("blueprint_identifier").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	action, err := actionResourceToBody(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var a *cli.Action
	if d.Id() != "" {
		a, err = c.UpdateAction(ctx, d.Get("blueprint_identifier").(string), d.Id(), action)
	} else {
		a, err = c.CreateAction(ctx, d.Get("blueprint_identifier").(string), action)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(a.Identifier)
	return diags
}
