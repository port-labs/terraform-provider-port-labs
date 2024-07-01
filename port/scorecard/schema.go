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
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint of the scorecard",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the scorecard",
			Required:            true,
		},
		"levels": schema.ListNestedAttribute{
			MarkdownDescription: "The Levels of the scorecard",
			Optional:            true,
			Computed:            true,
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

## Example Usage

Create a parent blueprint with a child blueprint and an aggregation property to count the parent kids:

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
            property = "author",
            operator : "=",
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

` + "```"
