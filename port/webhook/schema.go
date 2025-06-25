package webhook

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func WebhookSecuritySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"secret": schema.StringAttribute{
			MarkdownDescription: "The secret of the webhook",
			Optional:            true,
		},
		"signature_header_name": schema.StringAttribute{
			MarkdownDescription: "The signature header name of the webhook",
			Optional:            true,
		},
		"signature_algorithm": schema.StringAttribute{
			MarkdownDescription: "The signature algorithm of the webhook",
			Optional:            true,
		},
		"signature_prefix": schema.StringAttribute{
			MarkdownDescription: "The signature prefix of the webhook",
			Optional:            true,
		},
		"request_identifier_path": schema.StringAttribute{
			MarkdownDescription: "The request identifier path of the webhook",
			Optional:            true,
		},
	}
}

func WebhookMappingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint of the mapping",
			Required:            true,
		},
		"filter": schema.StringAttribute{
			MarkdownDescription: "The filter of the mapping",
			Optional:            true,
		},
		"items_to_parse": schema.StringAttribute{
			MarkdownDescription: "The items to parser of the mapping",
			Optional:            true,
		},
		"operation": schema.SingleNestedAttribute{
			MarkdownDescription: "The operation of the mapping",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					MarkdownDescription: "The type of the operation",
					Validators: []validator.String{
						stringvalidator.OneOf("create", "delete"),
					},
					Required: true,
				},
				"delete_dependents": schema.BoolAttribute{
					MarkdownDescription: "Whether to delete dependents entities, only relevant for delete operations",
					Optional:            true,
				},
			},
		},
		"entity": schema.SingleNestedAttribute{
			MarkdownDescription: "The entity of the mapping",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier of the entity",
					Required:            true,
				},
				"title": schema.StringAttribute{
					MarkdownDescription: "The title of the entity",
					Optional:            true,
				},
				"icon": schema.StringAttribute{
					MarkdownDescription: "The icon of the entity",
					Optional:            true,
				},
				"team": schema.StringAttribute{
					MarkdownDescription: "The team of the entity",
					Optional:            true,
				},
				"properties": schema.MapAttribute{
					MarkdownDescription: "The properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"relations": schema.MapAttribute{
					MarkdownDescription: "The relations of the entity. Relations can be defined as either simple JQ expressions (strings) or search query objects. When using objects, the rules array must be encoded with jsonencode().",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		},
	}
}

func WebhookSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the webhook",
			Optional:            true,
			Computed:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the webhook",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the webhook",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the webhook",
			Optional:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether the webhook is enabled",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"url": schema.StringAttribute{
			MarkdownDescription: "The url of the webhook",
			Computed:            true,
		},
		"webhook_key": schema.StringAttribute{
			MarkdownDescription: "The webhook key of the webhook",
			Computed:            true,
		},
		"security": schema.SingleNestedAttribute{
			MarkdownDescription: "The security of the webhook",
			Optional:            true,
			Attributes:          WebhookSecuritySchema(),
		},

		"mappings": schema.ListNestedAttribute{
			MarkdownDescription: "The mappings of the webhook",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: WebhookMappingSchema(),
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the webhook",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the webhook",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the webhook",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the webhook",
			Computed:            true,
		},
	}
}

func (r *WebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          WebhookSchema(),
	}
}

var ResourceMarkdownDescription = `


# Webhook resource
Webhook resource can be used to create webhooks integrations in Port.

## Example Usage

` + "```hcl" + `
resource "port_blueprint" "author" {
    title = "Author"
    icon = "User"
    identifier = "author"
    properties = {
      string_props = {
        "name" = {
          type = "string"
          title = "Name"
        }
      }
    }
  }

  resource "port_blueprint" "team" {
    title = "Team"
    icon = "Team"
    identifier = "team"
    properties = {
      string_props = {
        "name" = {
          type = "string"
          title = "Team Name"
        }
      }
    }
  }

  resource "port_blueprint" "microservice" {
    title = "TF test microservice"
    icon = "Terraform"
    identifier = "microservice"
    properties = {
      string_props = {
        "url" = {
          type = "string"
          title = "URL"
        }
      }
    }
    relations = {
      "author" = {
        title = "Author"
        target = port_blueprint.author.identifier
      }
      "team" = {
        title = "Team"
        target = port_blueprint.team.identifier
      }
    }
  }

  resource "port_webhook" "create_pr" {
    identifier = "pr_webhook"
    title      = "Webhook with mixed relations"
    icon       = "Terraform"
    enabled    = true
    
    mappings = [
      {
        blueprint = port_blueprint.microservice.identifier
        operation = { "type" = "create" }
        filter    = ".headers.\"x-github-event\" == \"pull_request\""
        entity = {
          identifier = ".body.pull_request.id | tostring"
          title      = ".body.pull_request.title"
          properties = {
            url = ".body.pull_request.html_url"
          }
          relations = {
            # Complex object relation with search query
            author = jsonencode({
              combinator = "'and'",
              rules = [
                {
                  property = "'$identifier'"
                  operator = "'='"
                  value    = ".body.pull_request.user.login | tostring"
                }
              ]
            })
            
            # Simple string relation
            team = ".body.repository.owner.login | tostring"
          }
        }
      }
    ]
    
    depends_on = [
      port_blueprint.microservice,
      port_blueprint.author,
      port_blueprint.team
    ]
  }
` + "\n```" + `

## Notes

- When using object format for relations, ` + "`combinator`, `property`" + ` and ` + "`operator`" + ` fields should be enclosed in single quotes, while ` + "`value`" + ` should not have quotes as it's a JQ expression. The single quotes are required because these fields contain literal string values that must be passed as-is to the Port API, whereas ` + "`value`" + ` contains a JQ expression that should be evaluated dynamically.
- For all available operators, see the [Port comparison operators documentation](https://docs.port.io/search-and-query/comparison-operators).

`
