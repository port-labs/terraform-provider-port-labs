package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func IntegrationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"installation_id": schema.StringAttribute{
			MarkdownDescription: "The installation ID of the integration. Must contain only lowercase letters, numbers, dashes and underscores (pattern: `" + installationIdPattern + "`).",
			Required:            true,
		},
		"version": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"title": schema.StringAttribute{
			Optional: true,
		},
		"installation_app_type": schema.StringAttribute{
			Optional: true,
		},
		"config": schema.StringAttribute{
			MarkdownDescription: "Integration Config Raw JSON string (use `jsonencode`)",
			Optional:            true,
		},
		"webhook_changelog_destination": schema.SingleNestedAttribute{
			MarkdownDescription: "The webhook changelog destination of the integration",
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
			MarkdownDescription: "The changelog destination of the blueprint (just an empty `{}`)",
			Optional:            true,
			AttributeTypes:      map[string]attr.Type{},
		},
	}
}

func (r *IntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: IntegrationResourceMarkdownDescription,
		Attributes:          IntegrationSchema(),
	}
}

var IntegrationResourceMarkdownDescription = `

# Integration resource

**NOTE:** This resource manages existing integration and integration mappings, not for creating new integrations.

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

Docs about how to import existing integrations and manage their mappings can be found [here](https://docs.getport.io/guides/all/import-and-manage-integration).


` + "```hcl" + `
resource "port_integration" "my_custom_integration" {
	installation_id       = "my-custom-integration-id"
	title                 = "My Custom Integration"
	config = jsonencode({
		createMissingRelatedEntitiesboolean = true
		deleteDependentEntities = true,
		resources = [{
			kind = "my-custom-kind"
			selector = {
				query = ".title"
			}
			port = {
				entity = {
					mappings = [{
						identifier = "'my-identifier'"
						title      = ".title"
						blueprint  = "'my-blueprint'"
						properties = {
							my_property = 123
						}
						relations  = {}
					}]
				}
			}
		}]
	})
}


` + "```\n" + `
### NOTICE:

The following config properties (` + "`selector.query|entity.mappings.*`" + `) are jq expressions, which means that you need to input either a valid jq expression (E.g ` + "`.title`" + `), or if you want a string value, a qouted escaped string val (E.g ` + "`'my-string'`" + `).
`
