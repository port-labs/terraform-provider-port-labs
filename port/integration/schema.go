package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func IntegrationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"installation_id": schema.StringAttribute{
			MarkdownDescription: "The installation ID of the integration. Must contain only lowercase letters, numbers, and dashes (pattern: `" + installationIdPattern + "`). Changing this forces the integration to be destroyed and recreated under the new ID.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"installation_type": schema.StringAttribute{
			MarkdownDescription: "Where the integration runs: `OnPrem` (default) registers a self-hosted integration you run yourself (e.g. via Helm/Docker) - Terraform only manages the Port-side record and mapping. `Saas` has Port provision and run a hosted Ocean integration on your behalf; `terraform apply` waits for it to finish provisioning before returning. Other installation types (OAuth-based GitHub App installs, etc.) require a manual browser consent step and are not supported by this resource. Immutable after creation - changing it forces recreation.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(consts.InstallationTypeOnPrem),
			Validators: []validator.String{
				stringvalidator.OneOf(consts.InstallationTypeOnPrem, consts.InstallationTypeSaas),
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
			MarkdownDescription: "The name of the integrated tool/platform (e.g. `kubernetes`, `pagerduty`). Required when creating a new integration.",
			Optional:            true,
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

This resource manages the full lifecycle (create/read/update/delete) of a Port integration and its mapping - including creating a brand-new integration, not just adopting one installed by hand.

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

If you have an integration that was installed before you started managing it with Terraform, you can bring it under Terraform management by importing it instead of creating a new one - see [Import and manage an integration](https://docs.getport.io/guides/all/import-and-manage-integration).

## Creating a self-hosted (OnPrem) integration

` + "```hcl" + `
resource "port_integration" "my_custom_integration" {
	installation_id       = "my-custom-integration-id"
	installation_app_type = "my-custom-kind"
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

With ` + "`installation_type = \"OnPrem\"`" + ` (the default), Terraform only creates/manages the Port-side registration and mapping - you still need to run the integration's own agent/container (e.g. via Helm or Docker) separately, pointed at the same ` + "`installation_id`" + `.

## Creating a Port-hosted (SaaS) integration

` + "```hcl" + `
resource "port_integration" "my_hosted_integration" {
	installation_id       = "my-hosted-integration-id"
	installation_app_type = "githubOcean"
	installation_type     = "Saas"
	title                 = "My Hosted Integration"
}
` + "```\n" + `

With ` + "`installation_type = \"Saas\"`" + `, Port provisions and runs the integration for you; ` + "`terraform apply`" + ` waits for provisioning to finish (and ` + "`terraform destroy`" + ` waits for teardown to finish) before returning.

**NOTE:** OAuth-based installation types (e.g. GitHub App installs) require a manual browser consent step and cannot be created through this resource.

### NOTICE:

The following config properties (` + "`selector.query|entity.mappings.*`" + `) are jq expressions, which means that you need to input either a valid jq expression (E.g ` + "`.title`" + `), or if you want a string value, a qouted escaped string val (E.g ` + "`'my-string'`" + `).
`
