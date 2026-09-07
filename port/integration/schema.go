package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func IntegrationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"installation_id": schema.StringAttribute{
			MarkdownDescription: "The installation ID of the integration. Must contain only lowercase letters, numbers, and dashes (pattern: `" + installationIdPattern + "`). Cannot be changed after creation.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(installationIdRegex, "installation_id must match the pattern "+installationIdPattern+": must contain only lowercase letters, numbers, and dashes"),
			},
		},
		"version": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"title": schema.StringAttribute{
			Optional: true,
		},
		"installation_app_type": schema.StringAttribute{
			MarkdownDescription: "Deprecated. The integrated tool name for catalog integration types (e.g. `GitHub`, `GitLab`, `K8S EXPORTER`). Custom integrations can omit this field. Cannot be changed after creation.",
			Optional:            true,
		},
		"config": schema.StringAttribute{
			MarkdownDescription: "Integration Config Raw JSON string (use `jsonencode`)",
			Optional:            true,
		},
		"webhook_changelog_destination": schema.SingleNestedAttribute{
			MarkdownDescription: "The webhook changelog destination of the integration",
			Optional:            true,
			Computed:            true,
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
			Computed:            true,
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

This resource can be used to create new self-hosted integrations directly from Terraform (no prior installation via the UI is required), as well as to manage the config and mappings of existing integrations (including those installed from Port's catalog of native integrations).

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

Docs about how to use Port's Terraform provider to create and manage integrations can be found [here](https://docs.getport.io/context-lake/ingestion/ingest-data-into-port/other/iac/terraform/terraform).

## Self-hosted (OnPrem) integration

Self-hosted integrations run on your own infrastructure (e.g. an Ocean exporter container) and pull their mapping from Port.

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

For catalog integration types, set ` + "`installation_app_type`" + ` to the integrated tool name (e.g. ` + "`GitHub`" + `, ` + "`GitLab`" + `, ` + "`K8S EXPORTER`" + `) and ` + "`version`" + ` if you want to pin a specific integration version. Custom integrations can omit ` + "`installation_app_type`" + `.

### NOTICE:

The following config properties (` + "`selector.query|entity.mappings.*`" + `) are jq expressions, which means that you need to input either a valid jq expression (E.g ` + "`.title`" + `), or if you want a string value, a qouted escaped string val (E.g ` + "`'my-string'`" + `).

### NOTES:

- ` + "`installation_id`" + ` and ` + "`installation_app_type`" + ` cannot be changed after creation.
- A changelog destination (` + "`webhook_changelog_destination`" + ` / ` + "`kafka_changelog_destination`" + `) can be added or updated, but not removed - the Port API does not support clearing it. To remove it, delete and recreate the integration (e.g. taint the resource).
- Existing integrations can be brought under Terraform management with ` + "`terraform import port_integration.my_integration <installation_id>`" + `.
- ` + "`terraform destroy`" + ` deletes the real integration in Port, not just removes it from state. Use ` + "`terraform state rm`" + ` if you only want to stop managing an integration with Terraform without deleting it from Port. This is especially relevant for imported resources.
`
