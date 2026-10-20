package entity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EntitySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the entity",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the entity",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the entity",
			Optional:            true,
		},
		"run_id": schema.StringAttribute{
			MarkdownDescription: "The runID of the action run that created the entity",
			Optional:            true,
		},
		"create_missing_related_entities": schema.BoolAttribute{
			MarkdownDescription: "Whether to create missing related entities",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: "The teams the entity belongs to",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier the entity relates to",
			Required:            true,
		},
		"properties": schema.SingleNestedAttribute{
			MarkdownDescription: "The properties of the entity",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_props": schema.MapAttribute{
					MarkdownDescription: "The string properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"number_props": schema.MapAttribute{
					MarkdownDescription: "The number properties of the entity",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"boolean_props": schema.MapAttribute{
					MarkdownDescription: "The bool properties of the entity",
					Optional:            true,
					ElementType:         types.BoolType,
				},
				"object_props": schema.MapAttribute{
					MarkdownDescription: "The object properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"array_props": schema.SingleNestedAttribute{
					MarkdownDescription: "The array properties of the entity",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"string_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Optional:    true,
						},
						"number_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.Float64Type},
							Optional:    true,
						},
						"boolean_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.BoolType},
							Optional:    true,
						},
						"object_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Optional:    true,
						},
					},
				},
			},
		},
		"relations": schema.SingleNestedAttribute{
			MarkdownDescription: "The relations of the entity",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"single_relations": schema.MapAttribute{
					MarkdownDescription: "The single relation of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"many_relations": schema.MapAttribute{
					MarkdownDescription: "The many relation of the entity",
					Optional:            true,
					ElementType:         types.ListType{ElemType: types.StringType},
				},
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the entity",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the entity",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the entity",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the entity",
			Computed:            true,
		},
	}
}

func (r *EntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: entityResourceMarkdownDescription,
		Attributes:          EntitySchema(),
	}
}

var entityResourceMarkdownDescription = `

# Entity Resource

Manages a Port entity. Entities are instances of blueprints representing your software catalog items (services, environments, packages, etc.).

Docs about entities can be found [here](https://docs.getport.io/build-your-software-catalog/sync-data-to-catalog/).

## Example Usage

### Basic entity with properties

` + "```hcl" + `
resource "port_blueprint" "microservice" {
  identifier = "microservice"
  title      = "Microservice"
  properties = {
    string_props = {
      "name" = {
        title = "Name"
      }
      "language" = {
        title = "Language"
        enum  = ["go", "python", "javascript"]
      }
    }
    boolean_props = {
      "is_active" = {
        title   = "Is Active"
        default = true
      }
    }
  }
}

resource "port_entity" "my_service" {
  identifier = "my-service"
  title      = "My Service"
  blueprint  = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      "name"     = "my-service"
      "language" = "go"
    }
    boolean_props = {
      "is_active" = true
    }
  }
}
` + "```" + `

### Entity with relations

` + "```hcl" + `
resource "port_blueprint" "environment" {
  identifier = "environment"
  title      = "Environment"
}

resource "port_blueprint" "microservice" {
  identifier = "microservice"
  title      = "Microservice"
  relations = {
    "environment" = {
      title  = "Environment"
      target = port_blueprint.environment.identifier
    }
  }
}

resource "port_entity" "production" {
  identifier = "production"
  title      = "Production"
  blueprint  = port_blueprint.environment.identifier
}

resource "port_entity" "my_service" {
  identifier = "my-service"
  title      = "My Service"
  blueprint  = port_blueprint.microservice.identifier
  relations = {
    single_relations = {
      "environment" = port_entity.production.identifier
    }
  }
}
` + "```" + `

### Entity with teams

` + "```hcl" + `
resource "port_team" "backend" {
  name = "Backend Team"
}

resource "port_blueprint" "microservice" {
  identifier = "microservice"
  title      = "Microservice"
  ownership = {
    type = "Direct"
  }
}

resource "port_entity" "my_service" {
  identifier = "my-service"
  title      = "My Service"
  blueprint  = port_blueprint.microservice.identifier
  teams      = [port_team.backend.identifier]
}
` + "```" + `

### Entity with create_missing_related_entities

When ` + "`create_missing_related_entities`" + ` is true, Port will create stub entities for any relations that reference non-existent entities.

` + "```hcl" + `
resource "port_entity" "my_service" {
  identifier = "my-service"
  title      = "My Service"
  blueprint  = "microservice"
  relations = {
    single_relations = {
      "environment" = "staging"  # Will be created if it doesn't exist
    }
  }
  create_missing_related_entities = true
}
` + "```" + `

## Import

Entities can be imported using the format ` + "`blueprint_identifier:entity_identifier`" + `:

` + "```shell" + `
terraform import port_entity.my_service microservice:my-service
` + "```" + `
`
