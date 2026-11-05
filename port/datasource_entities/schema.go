package datasource_entities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func Schema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"datasource_prefix": schema.StringAttribute{
			MarkdownDescription: "The datasource prefix to match entities against.",
			Required:            true,
		},
		"datasource_suffix": schema.StringAttribute{
			MarkdownDescription: "The datasource suffix to match entities against.",
			Required:            true,
		},
		"limit": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of entities to return per request. When omitted, Port uses its default page size.",
			Optional:            true,
		},
		"before": schema.StringAttribute{
			MarkdownDescription: "Return only entities updated before this RFC3339 timestamp.",
			Optional:            true,
		},
		"entities": schema.ListNestedAttribute{
			MarkdownDescription: "Entities whose datasource matches the prefix and suffix.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"identifier": schema.StringAttribute{
						MarkdownDescription: "The identifier of the entity.",
						Computed:            true,
					},
					"blueprint": schema.StringAttribute{
						MarkdownDescription: "The blueprint identifier of the entity.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *DatasourceEntitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: DatasourceEntitiesDataSourceMarkdownDescription,
		Attributes:          Schema(),
	}
}

var DatasourceEntitiesDataSourceMarkdownDescription = `

# Datasource Entities Data Source

The datasource entities data source looks up entity identifiers by matching their datasource path against a prefix and suffix.
This mirrors the Port API endpoint updated in [port-labs/Port#18099](https://github.com/port-labs/Port/pull/18099), which no longer accepts an ` + "`extra_conditions.not_modified_since_creation`" + ` filter.

## Example Usage

` + "```hcl" + `

data "port_datasource_entities" "integration_entities" {
  datasource_prefix = "port-ocean/github/"
  datasource_suffix = "/my-github-integration/resync"
}

output "integration_entity_identifiers" {
  value = [for entity in data.port_datasource_entities.integration_entities.entities : entity.identifier]
}

` + "\n```" + `

`
