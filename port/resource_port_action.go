package port

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
	"github.com/samber/lo"
)

var requiredGithubArguments = []string{"invocation_method.0.org", "invocation_method.0.repo", "invocation_method.0.workflow"}
var requiredWebhookArguments = []string{"invocation_method.0.url"}
var requiredAzureDevopsArguments = []string{"invocation_method.0.azure_org", "invocation_method.0.webhook"}

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
				Type:        schema.TypeString,
				Description: "The icon of the action",
				Optional:    true,
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
				Description: "The methods the action is dispatched in. Supports WEBHOOK, KAFKA, GITHUB and AZURE-DEVOPS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "One of WEBHOOK, KAFKA and GITHUB",
							ValidateFunc: validation.StringInSlice([]string{"WEBHOOK", "KAFKA", "GITHUB", "AZURE-DEVOPS"}, false),
						},
						"url": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Required when selecting type WEBHOOK. The URL to which the action is dispatched",
							ConflictsWith: []string{"invocation_method.0.org"},
						},
						"agent": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  "Relevant only when selecting type WEBHOOK. The flag that controls if the port execution agent will handle the action",
							RequiredWith: requiredWebhookArguments,
						},
						"org": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Required when selecting type GITHUB. The GitHub org that the workflow belongs to",
							RequiredWith:  requiredGithubArguments,
							ConflictsWith: []string{"invocation_method.0.url"},
						},
						"azure_org": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Required when selecting type AZURE-DEVOPS. The Azure Devops org that the webhook belongs to",
							RequiredWith: requiredAzureDevopsArguments,
						},
						"webhook": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Required when selecting type AZURE-DEVOPS. The Azure Devops webhook id",
							RequiredWith:  requiredAzureDevopsArguments,
							ConflictsWith: []string{"invocation_method.0.url"},
						},
						"repo": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Required when selecting type GITHUB. The GitHub repository that the workflow belongs to",
							RequiredWith: requiredGithubArguments,
						},
						"workflow": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Required when selecting type GITHUB. The GitHub workflow id or the workflow file name",
							RequiredWith: requiredGithubArguments,
						},
						"omit_payload": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  "Relevant only when selecting type GITHUB. The flag that controls if to omit Port's payload from workflow's dispatch input",
							RequiredWith: requiredGithubArguments,
						},
						"omit_user_inputs": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  "Relevant only when selecting type GITHUB. The flag that controls if to omit user inputs from workflow's dispatch input",
							RequiredWith: requiredGithubArguments,
						},
						"report_workflow_status": {
							Type:         schema.TypeBool,
							Optional:     true,
							Default:      true,
							Description:  "Relevant only when selecting type GITHUB. The flag that controls if to report the action status when the workflow completes",
							RequiredWith: requiredGithubArguments,
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*cli.PortClient)
	id := d.Id()
	actionIdentifier := id
	blueprintIdentifier := d.Get("blueprint_identifier").(string)
	if strings.Contains(id, ":") {
		parts := strings.SplitN(id, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return diag.FromErr(fmt.Errorf("unexpected format of ID (%s), expected blueprintId:entityId", id))
		}

		blueprintIdentifier = parts[0]
		actionIdentifier = parts[1]
	}

	action, statusCode, err := c.ReadAction(ctx, blueprintIdentifier, actionIdentifier)
	if err != nil {
		if statusCode == 404 {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}
	writeActionFieldsToResource(d, action, blueprintIdentifier)
	return diags
}

func writeActionFieldsToResource(d *schema.ResourceData, action *cli.Action, blueprintIdentifier string) {
	d.SetId(action.Identifier)
	d.Set("blueprint_identifier", blueprintIdentifier)
	d.Set("identifier", action.Identifier)
	d.Set("title", action.Title)
	d.Set("icon", action.Icon)
	d.Set("description", action.Description)

	reportWorkflowStatus := true
	if action.InvocationMethod.ReportWorkflowStatus != nil {
		reportWorkflowStatus = *action.InvocationMethod.ReportWorkflowStatus
	}

	orgKey := "org"

	if action.InvocationMethod.Type == "AZURE-DEVOPS" {
		orgKey = "azure_org"
	}

	d.Set("invocation_method", []any{map[string]any{
		"type":                   action.InvocationMethod.Type,
		"url":                    action.InvocationMethod.Url,
		"agent":                  action.InvocationMethod.Agent,
		orgKey:                   action.InvocationMethod.Org,
		"repo":                   action.InvocationMethod.Repo,
		"workflow":               action.InvocationMethod.Workflow,
		"webhook":                action.InvocationMethod.Webhook,
		"omit_payload":           action.InvocationMethod.OmitPayload,
		"omit_user_inputs":       action.InvocationMethod.OmitUserInputs,
		"report_workflow_status": reportWorkflowStatus,
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
		p["format"] = v.Format
		p["pattern"] = v.Pattern
		p["blueprint"] = v.Blueprint
		p["enum"] = v.Enum
		if lo.Contains(action.UserInputs.Required, k) {
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

		if action.InvocationMethod.Type == "GITHUB" {
			action.InvocationMethod.Org = invocationMethod.([]any)[0].(map[string]interface{})["org"].(string)
			action.InvocationMethod.Repo = invocationMethod.([]any)[0].(map[string]interface{})["repo"].(string)
			action.InvocationMethod.Workflow = invocationMethod.([]any)[0].(map[string]interface{})["workflow"].(string)
			action.InvocationMethod.OmitPayload = invocationMethod.([]any)[0].(map[string]interface{})["omit_payload"].(bool)
			action.InvocationMethod.OmitUserInputs = invocationMethod.([]any)[0].(map[string]interface{})["omit_user_inputs"].(bool)
			reportWorkflowStatus := invocationMethod.([]any)[0].(map[string]interface{})["report_workflow_status"].(bool)
			if reportWorkflowStatus {
				action.InvocationMethod.ReportWorkflowStatus = nil
			} else {
				action.InvocationMethod.ReportWorkflowStatus = new(bool)
				*action.InvocationMethod.ReportWorkflowStatus = reportWorkflowStatus
			}
		} else if action.InvocationMethod.Type == "AZURE-DEVOPS" {
			action.InvocationMethod.Org = invocationMethod.([]any)[0].(map[string]interface{})["azure_org"].(string)
			action.InvocationMethod.Webhook = invocationMethod.([]any)[0].(map[string]interface{})["webhook"].(string)
		} else if action.InvocationMethod.Type == "WEBHOOK" {
			action.InvocationMethod.Url = invocationMethod.([]any)[0].(map[string]interface{})["url"].(string)
		}
		action.InvocationMethod.Agent = invocationMethod.([]any)[0].(map[string]interface{})["agent"].(bool)

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
	id := d.Id()
	actionIdentifier := id
	blueprintIdentifier := d.Get("blueprint_identifier").(string)

	err := c.DeleteAction(ctx, blueprintIdentifier, actionIdentifier)
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

	blueprintIdentifier := d.Get("blueprint_identifier").(string)
	actionIdentifier := d.Id()
	if d.Id() != "" {
		a, err = c.UpdateAction(ctx, blueprintIdentifier, actionIdentifier, action)
	} else {
		a, err = c.CreateAction(ctx, d.Get("blueprint_identifier").(string), action)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(a.Identifier)
	return diags
}
