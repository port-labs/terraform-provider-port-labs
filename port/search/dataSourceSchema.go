package search

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EntitySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the entity",
			Computed:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the entity",
			Computed:            true,
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the entity",
			Computed:            true,
			Optional:            true,
		},
		"run_id": schema.StringAttribute{
			MarkdownDescription: "The runID of the action run that created the entity",
			Computed:            true,
			Optional:            true,
		},
		"teams": schema.ListAttribute{
			MarkdownDescription: "The teams the entity belongs to",
			Computed:            true,
			Optional:            true,
			ElementType:         types.StringType,
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier the entity relates to",
			Computed:            true,
		},
		"properties": schema.SingleNestedAttribute{
			MarkdownDescription: "The properties of the entity",
			Computed:            true,
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_props": schema.MapAttribute{
					MarkdownDescription: "The string properties of the entity",
					Computed:            true,
					Optional:            true,
					ElementType:         types.StringType,
				},
				"number_props": schema.MapAttribute{
					MarkdownDescription: "The number properties of the entity",
					Computed:            true,
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"boolean_props": schema.MapAttribute{
					MarkdownDescription: "The bool properties of the entity",
					Computed:            true,
					Optional:            true,
					ElementType:         types.BoolType,
				},
				"object_props": schema.MapAttribute{
					MarkdownDescription: "The object properties of the entity",
					Computed:            true,
					Optional:            true,
					ElementType:         types.StringType,
				},
				"array_props": schema.SingleNestedAttribute{
					MarkdownDescription: "The array properties of the entity",
					Computed:            true,
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"string_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Computed:    true,
							Optional:    true,
						},
						"number_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.Float64Type},
							Computed:    true,
							Optional:    true,
						},
						"boolean_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.BoolType},
							Computed:    true,
							Optional:    true,
						},
						"object_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
		},
		"scorecards": schema.MapAttribute{
			MarkdownDescription: "The scorecards of the entity",
			Computed:            true,
			Optional:            true,
			ElementType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"rules": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"identifier": types.StringType,
								"status":     types.StringType,
								"level":      types.StringType,
							},
						},
					},
					"level": types.StringType,
				},
			},
		},
		"relations": schema.SingleNestedAttribute{
			MarkdownDescription: "The relations of the entity",
			Computed:            true,
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"single_relations": schema.MapAttribute{
					MarkdownDescription: "The single relation of the entity",
					Computed:            true,
					Optional:            true,
					ElementType:         types.StringType,
				},
				"many_relations": schema.MapAttribute{
					MarkdownDescription: "The many relation of the entity",
					Computed:            true,
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

func Schema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"query": schema.StringAttribute{
			MarkdownDescription: "The search query",
			Required:            true,
		},
		"exclude_calculated_properties": schema.BoolAttribute{
			MarkdownDescription: "Exclude calculated properties",
			Optional:            true,
		},
		"include": schema.ListAttribute{
			MarkdownDescription: "Properties to include in the results",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"exclude": schema.ListAttribute{
			MarkdownDescription: "Properties to exclude from the results",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"attach_title_to_relation": schema.BoolAttribute{
			MarkdownDescription: "Attach title to relation",
			Optional:            true,
		},
		"matching_blueprints": schema.ListAttribute{
			MarkdownDescription: "The matching blueprints for the search query",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"entities": schema.ListNestedAttribute{
			MarkdownDescription: "A list of entities matching the search query",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: EntitySchema(),
			},
		},
	}
}

func (d *SearchDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: SearchDataSourceMarkdownDescription,
		Attributes:          Schema(),
	}
}

var SearchDataSourceMarkdownDescription = `

# Search Data Source

The search data source allows you to search for entities in Port.

See the [Port documentation](https://docs.getport.io/search-and-query/) for more information about the search capabilities in Port.

## Example Usage

### Search for all entities in a specific blueprint:

` + "```hcl" + `

data "port_search" "all_service" {
  query = jsonencode({
    "combinator" : "and", "rules" : [
      { "operator" : "=", "property" : "$blueprint", "value" : "Service" },
    ]
  })
}

` + "\n```" + `

### Search for entity with specific identifier in a specific blueprint to create another resource based on the values of the entity:


` + "```hcl" + `

data "port_search" "ads_service" {
  query = jsonencode({
    "combinator" : "and", "rules" : [
      { "operator" : "=", "property" : "$blueprint", "value" : "Service" },
      { "operator" : "=", "property" : "$identifier", "value" : "Ads" },
    ]
  })
}

` + "\n```" + `

### Scorecards automation example
In this example we are creating a jira task for each service that its Ownership Scorecard hasn't reached Gold level : 

` + "```hcl" + `

data "port_search" "all_services" {
  query = jsonencode({
    "combinator" : "and", "rules" : [
      { "operator" : "=", "property" : "$blueprint", "value" : "microservice" },
    ]
  })
}

locals {
  // Count the number of services that are not owned by a team with a Gold level
  microservice_ownership_without_gold_level = length([
    for entity in data.port_search.all_services.entities : entity.scorecards["ownership"].level
    if entity.scorecards["ownership"].level != "Gold"
  ])
}

// create jira issue per service that is not owned by a team with a Gold level
resource "jira_issue" "microservice_ownership_without_gold_level" {
  count      = local.microservice_ownership_without_gold_level
  issue_type = "Task"

  project_key = "PORT"

  summary     = "Service ${data.port_search.backend_services.entities[count.index].title} hasn't reached Gold level in Ownership Scorecard"
  description = "[Service](https://app.getport.io/${port_blueprint.microservice.identifier}Entity/${data.port_search.backend_services.entities[count.index].identifier}) is not owned by a team with a Gold level, please assign a team with a Gold level to the service"
}

` + "\n```" + ``
