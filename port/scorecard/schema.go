package scorecard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func LevelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"color": schema.StringAttribute{
			MarkdownDescription: "The color of the level",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the level",
			Required:            true,
		},
	}
}

func RuleSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the rule",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the rule",
			Required:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the rule",
			Optional:            true,
		},
		"level": schema.StringAttribute{
			MarkdownDescription: "The level of the rule",
			Required:            true,
		},
		"query": schema.SingleNestedAttribute{
			MarkdownDescription: "The query of the rule",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"combinator": schema.StringAttribute{
					MarkdownDescription: "The combinator of the query",
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
					Required: true,
				},
				"conditions": schema.ListAttribute{
					MarkdownDescription: "The conditions of the query. Each condition object should be encoded to a string",
					Required:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
			},
		},
	}
}
func ScorecardSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the scorecard",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint of the scorecard",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the scorecard",
			Required:            true,
		},
		"filter": schema.SingleNestedAttribute{
			MarkdownDescription: "The filter to apply on the entities before calculating the scorecard",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"combinator": schema.StringAttribute{
					MarkdownDescription: "The combinator of the filter",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
				},
				"conditions": schema.ListAttribute{
					MarkdownDescription: "The conditions of the filter. Each condition object should be encoded to a string",
					Required:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
			},
		},
		"levels": schema.ListNestedAttribute{
			MarkdownDescription: "The levels of the scorecard. This overrides the default levels (Basic, Bronze, Silver, Gold) if provided",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: LevelSchema(),
			},
		},
		"rules": schema.ListNestedAttribute{
			MarkdownDescription: "The rules of the scorecard",
			Required:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: RuleSchema(),
			},
		},
		"properties": schema.StringAttribute{
			MarkdownDescription: "Additional `_scorecard` blueprint properties applied to the scorecard entity, as a JSON encoded string. Property keys must match custom properties you added to the `_scorecard` blueprint.",
			Optional:            true,
		},
		"relations": schema.StringAttribute{
			MarkdownDescription: "Additional `_scorecard` blueprint relations applied to the scorecard entity, as a JSON encoded string. Relation values can be a string, an array of strings, or `null` to clear a relation. The `group` relation is managed by Port and cannot be set here.",
			Optional:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the scorecard",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the scorecard",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the scorecard",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the scorecard",
			Computed:            true,
		},
	}
}

func (r *ScorecardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          ScorecardSchema(),
	}
}

var ResourceMarkdownDescription = `

# Scorecard

This resource allows you to manage a scorecard.

See the [Port documentation](https://docs.getport.io/promote-scorecards/) for more information about scorecards.

` + "`properties`" + ` and ` + "`relations`" + ` set additional ` + "`_scorecard`" + ` blueprint property and relation values on the scorecard entity. Define the schema on the system blueprint first (for example with ` + "`port_system_blueprint`" + `), then reference those keys in ` + "`jsonencode({...})`" + `. Use ` + "`depends_on`" + ` so the scorecard is created only after the blueprint schema exists.

## Example Usage

This will create a blueprint with a Scorecard measuring the readiness of a microservice.

` + "```hcl" + `

resource "port_blueprint" "microservice" {
  title      = "microservice"
  icon       = "Terraform"
  identifier = "microservice"
  properties = {
    string_props = {
      "author" = {
        title = "Author"
      }
      "url" = {
        title = "URL"
      }
    }
    boolean_props = {
      "required" = {
        type = "boolean"
      }
    }
    number_props = {
      "sum" = {
        type = "number"
      }
    }
  }
}

resource "port_scorecard" "readiness" {
  identifier = "Readiness"
  title      = "Readiness"
  blueprint  = port_blueprint.microservice.identifier
  rules      = [
    {
      identifier = "hasOwner"
      title      = "Has Owner"
      level      = "Gold"
      query      = {
        combinator = "and"
        conditions = [
          jsonencode({
            property = "$team"
            operator = "isNotEmpty"
          }),
          jsonencode({
            property = "author"
            operator : "="
            value : "myValue"
          })
        ]
      }
    },
    {
      identifier = "hasUrl"
      title      = "Has URL"
      level      = "Silver"
      query      = {
        combinator = "and"
        conditions = [
          jsonencode({
            property = "url"
            operator = "isNotEmpty"
          })
        ]
      }
    },
    {
      identifier = "checkSumIfRequired"
      title      = "Check Sum If Required"
      level      = "Bronze"
      query      = {
        combinator = "or"
        conditions = [
          jsonencode({
            property = "required"
            operator : "="
            value : false
          }),
          jsonencode({
            property = "sum"
            operator : ">"
            value : 2
          })
        ]
      }
    }
  ]
  depends_on = [
    port_blueprint.microservice
  ]
}

` + "```" + `

## Example Usage with Properties and Relations

This will set custom ` + "`_scorecard`" + ` blueprint properties and relations on the scorecard entity.

` + "```hcl" + `

resource "port_system_blueprint" "scorecard" {
  identifier = "_scorecard"
  properties = {
    string_props = {
      owner = {
        type  = "string"
        title = "Owner"
      }
    }
    number_props = {
      priority = {
        type  = "number"
        title = "Priority"
      }
    }
  }
  relations = {
    owner_team = {
      title  = "Owner Team"
      target = "_team"
    }
  }
}

resource "port_scorecard" "readiness" {
  identifier = "Readiness"
  title      = "Readiness"
  blueprint  = port_blueprint.microservice.identifier
  properties = jsonencode({
    owner    = "platform-team"
    priority = 1
  })
  relations = jsonencode({
    owner_team = "platform-team"
  })
  rules = [{
    identifier = "hasOwner"
    title      = "Has Owner"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [jsonencode({
        property = "$team"
        operator = "isNotEmpty"
      })]
    }
  }]
  depends_on = [
    port_blueprint.microservice,
    port_system_blueprint.scorecard,
  ]
}

` + "```" + `

## Example Usage with Levels and Filter

This will override the default levels (Basic, Bronze, Silver, Gold) with the provided levels: Not Ready, Partially Ready, Ready.


` + "```hcl" + `

resource "port_blueprint" "microservice" {
  title      = "microservice"
  icon       = "Terraform"
  identifier = "microservice"
  properties = {
    string_props = {
      "author" = {
        title = "Author"
      }
      "url" = {
        title = "URL"
      }
    }
    boolean_props = {
      "required" = {
        type = "boolean"
      }
    }
    number_props = {
      "sum" = {
        type = "number"
      }
    }
  }
}

resource "port_scorecard" "readiness" {
  identifier = "Readiness"
  title      = "Readiness"
  blueprint  = port_blueprint.microservice.identifier
  filter = {
    combinator = "and"
    conditions = [
      jsonencode({
        property = "sum"
        operator = ">"
        value = 0
      })
    ]
  }
  levels = [
    {
      color = "red"
      title = "No Ready"
    },
    {
      color = "yellow"
      title = "Partially Ready"
    },
    {
      color = "green"
      title = "Ready"
    }
  ]
  rules = [
    {
      identifier = "hasOwner"
      title      = "Has Owner"
      level      = "Ready"
      query = {
        combinator = "and"
        conditions = [
          jsonencode({
            property = "$team"
            operator = "isNotEmpty"
          }),
          jsonencode({
            property = "author"
            operator : "="
            value : "myValue"
          })
        ]
      }
    },
    {
      identifier = "hasUrl"
      title      = "Has URL"
      level      = "Partially Ready"
      query = {
        combinator = "and"
        conditions = [
          jsonencode({
            property = "url"
            operator = "isNotEmpty"
          })
        ]
      }
    },
    {
      identifier = "checkSumIfRequired"
      title      = "Check Sum If Required"
      level      = "Partially Ready"
      query = {
        combinator = "or"
        conditions = [
          jsonencode({
            property = "required"
            operator : "="
            value : false
          }),
          jsonencode({
            property = "sum"
            operator : ">"
            value : 2
          })
        ]
      }
    }
  ]
  depends_on = [
    port_blueprint.microservice
  ]
}

` + "```"
