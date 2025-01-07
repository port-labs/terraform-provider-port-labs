package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func MetadataProperties() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the property",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the property",
			Optional:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "Whether the property is required",
			Computed:            true,
			Optional:            true,
			Default:             booldefault.StaticBool(false),
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the property",
			Optional:            true,
		}}

}

func StringPropertySchema() schema.Attribute {
	stringPropertySchema := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			MarkdownDescription: "The default of the string property",
			Optional:            true,
		},
		"format": schema.StringAttribute{
			MarkdownDescription: "The format of the string property",
			Optional:            true,
		},
		"min_length": schema.Int64Attribute{
			MarkdownDescription: "The min length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "The max length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"pattern": schema.StringAttribute{
			MarkdownDescription: "The pattern of the string property",
			Optional:            true,
		},
		"spec": schema.StringAttribute{
			MarkdownDescription: "The spec of the string property",
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("open-api", "async-api", "embedded-url")},
		},
		"spec_authentication": schema.SingleNestedAttribute{
			MarkdownDescription: "The spec authentication of the string property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"client_id": schema.StringAttribute{
					MarkdownDescription: "The clientId of the spec authentication",
					Required:            true,
				},
				"token_url": schema.StringAttribute{
					MarkdownDescription: "The tokenUrl of the spec authentication",
					Required:            true,
				},
				"authorization_url": schema.StringAttribute{
					MarkdownDescription: "The authorizationUrl of the spec authentication",
					Required:            true,
				},
			},
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the string property",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_colors": schema.MapAttribute{
			MarkdownDescription: "The enum colors of the string property",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}

	utils.CopyMaps(stringPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The string property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: stringPropertySchema,
		},
	}
}

func NumberPropertySchema() schema.Attribute {
	numberPropertySchema := map[string]schema.Attribute{
		"default": schema.Float64Attribute{
			MarkdownDescription: "The default of the number property",
			Optional:            true,
		},
		"maximum": schema.Float64Attribute{
			MarkdownDescription: "The min of the number property",
			Optional:            true,
		},
		"minimum": schema.Float64Attribute{
			MarkdownDescription: "The max of the number property",
			Optional:            true,
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the number property",
			Optional:            true,
			ElementType:         types.Float64Type,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_colors": schema.MapAttribute{
			MarkdownDescription: "The enum colors of the number property",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}

	utils.CopyMaps(numberPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The number property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: numberPropertySchema,
		},
	}
}

func BooleanPropertySchema() schema.Attribute {
	booleanPropertySchema := map[string]schema.Attribute{
		"default": schema.BoolAttribute{
			MarkdownDescription: "The default of the boolean property",
			Optional:            true,
		},
	}

	utils.CopyMaps(booleanPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The boolean property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: booleanPropertySchema,
		},
	}
}

func ArrayPropertySchema() schema.MapNestedAttribute {
	arrayPropertySchema := map[string]schema.Attribute{
		"min_items": schema.Int64Attribute{
			MarkdownDescription: "The min items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_items": schema.Int64Attribute{
			MarkdownDescription: "The max items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"string_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of the items",
					Optional:            true,
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		},
		"number_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
			},
		},
		"boolean_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.BoolType,
				},
			},
		},
		"object_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		},
	}

	utils.CopyMaps(arrayPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The array property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: arrayPropertySchema,
		},
	}
}

func OwnershipSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Optional ownership field for Blueprint. 'type' can be Inherited or Direct. If 'Inherited', then 'path' is required and must be a valid relation identifier.",
		Optional:            true,
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(path.MatchRoot("type")),
		},
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "Ownership type: either 'Inherited' or 'Direct'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Inherited", "Direct"),
				},
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path for the Inherited ownership type. Required when type is 'Inherited'. Must be a valid relation identifier (e.g. '$relations.parent').",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						`^\$relations\.[a-zA-Z0-9_-]+$`,
						"path must be a valid relation identifier starting with '$relations.' followed by the relation name",
					),
				},
			},
		},
	}
}

func ObjectPropertySchema() schema.MapNestedAttribute {

	objectPropertySchema := map[string]schema.Attribute{
		"spec": schema.StringAttribute{
			MarkdownDescription: "The spec of the object property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("async-api", "open-api"),
			},
		},
		"default": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The default of the object property",
		},
	}

	utils.CopyMaps(objectPropertySchema, MetadataProperties())

	return schema.MapNestedAttribute{
		MarkdownDescription: "The object property of the blueprint",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: objectPropertySchema,
		},
	}
}

func BlueprintSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display name of the blueprint",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the blueprint",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the blueprint",
			Optional:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the blueprint",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the blueprint",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the blueprint",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the blueprint",
			Computed:            true,
		},
		"team_inheritance": schema.SingleNestedAttribute{
			MarkdownDescription: "The team inheritance of the blueprint",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"path": schema.StringAttribute{
					MarkdownDescription: "The path of the team inheritance",
					Required:            true,
				},
			},
		},
		"webhook_changelog_destination": schema.SingleNestedAttribute{
			MarkdownDescription: "The webhook changelog destination of the blueprint",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "The url of the webhook changelog destination",
					Required:            true,
				},
				"agent": schema.BoolAttribute{
					MarkdownDescription: "The agent of the webhook changelog destination",
					Optional:            true,
				},
			},
		},
		"kafka_changelog_destination": schema.ObjectAttribute{
			MarkdownDescription: "The changelog destination of the blueprint",
			Optional:            true,
			AttributeTypes:      map[string]attr.Type{},
		},
		"properties": schema.SingleNestedAttribute{
			MarkdownDescription: "The properties of the blueprint",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_props":  StringPropertySchema(),
				"number_props":  NumberPropertySchema(),
				"boolean_props": BooleanPropertySchema(),
				"array_props":   ArrayPropertySchema(),
				"object_props":  ObjectPropertySchema(),
			},
		},
		"relations": schema.MapNestedAttribute{
			MarkdownDescription: "The relations of the blueprint",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"title": schema.StringAttribute{
						MarkdownDescription: "The title of the relation",
						Optional:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "The description of the relation",
						Optional:            true,
					},
					"target": schema.StringAttribute{
						MarkdownDescription: "The target of the relation",
						Required:            true,
					},
					"many": schema.BoolAttribute{
						MarkdownDescription: "The many of the relation",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"required": schema.BoolAttribute{
						MarkdownDescription: "The required of the relation",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
		},
		"mirror_properties": schema.MapNestedAttribute{
			MarkdownDescription: "The mirror properties of the blueprint",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						MarkdownDescription: "The path of the mirror property",
						Required:            true,
					},
					"title": schema.StringAttribute{
						MarkdownDescription: "The title of the mirror property",
						Optional:            true,
					},
				},
			},
		},
		"calculation_properties": schema.MapNestedAttribute{
			MarkdownDescription: "The calculation properties of the blueprint",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"calculation": schema.StringAttribute{
						MarkdownDescription: "The calculation of the calculation property",
						Required:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the calculation property",
						Required:            true,
					},
					"title": schema.StringAttribute{
						MarkdownDescription: "The title of the calculation property",
						Optional:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "The description of the calculation property",
						Optional:            true,
					},
					"icon": schema.StringAttribute{
						MarkdownDescription: "The icon of the calculation property",
						Optional:            true,
					},
					"format": schema.StringAttribute{
						MarkdownDescription: "The format of the calculation property",
						Optional:            true,
					},
					"colorized": schema.BoolAttribute{
						MarkdownDescription: "The colorized of the calculation property",
						Optional:            true,
					},
					"colors": schema.MapAttribute{
						MarkdownDescription: "The colors of the calculation property",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
		"force_delete_entities": schema.BoolAttribute{
			MarkdownDescription: "If set to true, the blueprint will be deleted with all its entities, even if they are not managed by Terraform",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"create_catalog_page": schema.BoolAttribute{
			MarkdownDescription: "This flag is only relevant for blueprint creation, by default if not set, a catalog page will be created for the blueprint",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"ownership": OwnershipSchema(),
	}
}

func (r *BlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: blueprintMarkdownDescription,
		Attributes:          BlueprintSchema(),
	}
}

var blueprintMarkdownDescription = `

# Blueprint Resource

Docs about the blueprint resource in Port can be found [here](https://docs.getport.io/build-your-software-catalog/define-your-data-model/setup-blueprint/).


## Example Usage

` + "```hcl" + `

resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "environment"
  properties = {
    string_props = {
      "aws-region" = {
        title = "AWS Region"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}

` + "```" + `

## Example Usage with Relations

` + "```hcl" + `

resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "environment"
  properties = {
    string_props = {
      "aws-region" = {
        title = "AWS Region"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}

resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "microservice"
  properties = {
    string_props = {
      "domain" = {
        title = "Domain"
      }
      "slack-channel" = {
        title  = "Slack Channel"
        format = "url"
      }
    }
  }
  relations = {
    "environment" = {
      target   = port_blueprint.environment.identifier
      required = true
      many     = false
    }
  }
}

` + "```" + `


## Example Usage with Mirror Properties

` + "```hcl" + `

resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "microservice"
  properties = {
    string_props = {
      "domain" = {
        title = "Domain"
      }
      "slack-channel" = {
        title  = "Slack Channel"
        format = "url"
      }
    }
  }
  mirror_properties = {
    "aws-region" = {
      path = "environment.aws-region"
    }
  }
  relations = {
    "environment" = {
      target   = port_blueprint.environment.identifier
      required = true
      many     = false
    }
  }
}

` + "```" + `

## Force Deleting a Blueprint

There could be cases where a blueprint will be managed by Terraform, but entities will get created from other sources (e.g. Port UI, API or other supported integrations).
In this case, when trying to delete the blueprint, Terraform will fail because it will try to delete the blueprint without deleting the entities first as they are not managed by Terraform.

To overcome this behavior, you can set the argument ` + "`force_delete_entities=true`" + `.
On the blueprint destroy it will trigger a migration that will delete all the entities in the blueprint and then delete the blueprint itself.

` + "```hcl" + `
resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "microservice"
  properties = {
    string_props = {
      "domain" = {
        title = "Domain"
      }
      "slack-channel" = {
        title  = "Slack Channel"
        format = "url"
      }
    }
  }
  force_delete_entities = false
}

` + "```" + `

`
