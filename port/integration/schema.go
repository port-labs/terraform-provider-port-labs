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
			MarkdownDescription: "The installation ID of the integration. Must contain only lowercase letters, numbers, and dashes (pattern: `" + installationIdPattern + "`). Changing this forces replacement.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(installationIdRegex, "must contain only lowercase letters, numbers, and dashes"),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
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
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"installation_type": schema.StringAttribute{
			MarkdownDescription: "How the integration is hosted: `OnPrem`, `Saas`, `SaasOAuth2`, `CustomGithubApp`, or `EnterpriseGithubApp`. Defaults to `OnPrem` when omitted on create.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"OnPrem",
					"Saas",
					"SaasOAuth2",
					"CustomGithubApp",
					"EnterpriseGithubApp",
				),
			},
		},
		"config": schema.StringAttribute{
			MarkdownDescription: "Integration Config Raw JSON string (use `jsonencode`)",
			Optional:            true,
		},
		"spec": schema.StringAttribute{
			MarkdownDescription: "SaaS integration spec JSON (`integrationSpec` / `appSpec`). Use `jsonencode`. Required for most SaaS integrations.",
			Optional:            true,
			Sensitive:           true,
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

This resource can be used to create new integrations directly from Terraform (no prior installation via the UI is required), as well as to manage the config and mappings of existing integrations (including those installed from Port's catalog of native integrations).

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

Docs about how to use Port's Terraform provider to create and manage integrations can be found [here](https://docs.getport.io/context-lake/ingestion/ingest-data-into-port/other/iac/terraform/terraform).

## Self-hosted (OnPrem) integration

Self-hosted integrations run on your own infrastructure (e.g. an Ocean exporter container) and pull their mapping from Port. This is the default hosting method.

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

## Port-hosted (SaaS) integration

Port-hosted integrations run on Port's infrastructure. Set ` + "`installation_type = \"Saas\"`" + ` and provide the integration ` + "`spec`" + `, which holds the hosting configuration (` + "`integrationSpec`" + `) and runtime preferences (` + "`appSpec`" + `).

Sensitive ` + "`integrationSpec`" + ` fields (tokens, keys, etc.) must reference a [Port organization secret](https://docs.getport.io/platform-administration/secrets-management/port-secrets/) **by name** - plain values are rejected. Use the ` + "`port_organization_secret`" + ` resource to manage them:

` + "```hcl" + `
resource "port_organization_secret" "gitlab_token_mapping" {
	secret_name  = "my-gitlab-token-mapping"
	secret_value = jsonencode({ "glpat-example" : ["**"] })
	description  = "GitLab tokens and the group scopes they may ingest"
}

resource "port_integration" "my_gitlab_saas" {
	installation_id       = "my-gitlab-saas"
	installation_type     = "Saas"
	installation_app_type = "gitlab"
	title                 = "My GitLab (Port-hosted)"

	spec = jsonencode({
		integrationSpec = {
			gitlabHost   = "https://gitlab.com"
			tokenMapping = port_organization_secret.gitlab_token_mapping.secret_name
		}
		appSpec = {
			scheduledResyncInterval = "2h"
		}
	})

	config = jsonencode({
		createMissingRelatedEntities = true
		deleteDependentEntities      = true
		resources = [{
			kind = "project"
			selector = {
				query = "true"
			}
			port = {
				entity = {
					mappings = [{
						identifier = ".path_with_namespace"
						title      = ".name"
						blueprint  = "'gitlabProject'"
					}]
				}
			}
		}]
	})
}
` + "```\n" + `

The required ` + "`integrationSpec`" + ` fields depend on the integration type - see each integration's documentation page for its configuration reference.

### NOTICE:

The following config properties (` + "`selector.query|entity.mappings.*`" + `) are jq expressions, which means that you need to input either a valid jq expression (E.g ` + "`.title`" + `), or if you want a string value, a qouted escaped string val (E.g ` + "`'my-string'`" + `).

### NOTES:

- ` + "`installation_id`" + ` and ` + "`installation_app_type`" + ` cannot be changed after creation - changing them recreates the integration.
- For Port-hosted (SaaS) integrations, ` + "`version`" + ` is managed by Port once the integration is provisioned. Omit it from your configuration to avoid a perpetual diff.
- A changelog destination (` + "`webhook_changelog_destination`" + ` / ` + "`kafka_changelog_destination`" + `) can be added or updated, but not removed - the Port API does not support clearing it. To remove it, recreate the integration.
- Existing integrations can be brought under Terraform management with ` + "`terraform import port_integration.my_integration <installation_id>`" + `.
`
