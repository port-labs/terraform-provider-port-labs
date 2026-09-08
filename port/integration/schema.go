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
			MarkdownDescription: "The installation ID of the integration. Must contain only lowercase letters, numbers, and dashes (pattern: `" + installationIdPattern + "`).",
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"installation_type": schema.StringAttribute{
			MarkdownDescription: "The installation type of the integration. Use `Saas` for Ocean SaaS integrations (requires `spec`). Defaults to `OnPrem` for self-hosted integrations. Only `OnPrem` and `Saas` are supported by this resource.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(consts.InstallationTypeOnPrem),
			Validators: []validator.String{
				stringvalidator.OneOf(
					consts.InstallationTypeOnPrem,
					consts.InstallationTypeSaas,
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"spec": schema.StringAttribute{
			MarkdownDescription: "Ocean SaaS integration spec as a JSON string (use `jsonencode`). Required when `installation_type` is `Saas`. Only `integrationSpec` and `appSpec` are supported. `systemSpec` and `privateSpec` are server-managed by Port (e.g. sizing, applier backend) and are ignored on read — changes made internally by Port will not cause Terraform drift.",
			Optional:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The provisioning status of the integration. Populated for SaaS installations.",
			Computed:            true,
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

Manages a Port integration, including self-hosted (OnPrem) and Ocean SaaS installations.

For SaaS integrations, create organization secrets first with ` + "`port_organization_secret`" + `, then reference secret names in ` + "`spec.integrationSpec`" + `.

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

Docs about how to import existing integrations and manage their mappings can be found [here](https://docs.getport.io/guides/all/import-and-manage-integration).


` + "```hcl" + `
resource "port_organization_secret" "pagerduty_token" {
  secret_name  = "pagerduty-api-token"
  secret_value = var.pagerduty_token
}

resource "port_integration" "pagerduty" {
  depends_on = [port_organization_secret.pagerduty_token]

  installation_id       = "pagerduty-prod"
  installation_app_type = "pagerduty"
  installation_type     = "Saas"
  version               = "0.1.0"
  title                 = "PagerDuty Production"

  spec = jsonencode({
    integrationSpec = {
      token = port_organization_secret.pagerduty_token.secret_name
    }
    appSpec = {
      scheduledResyncInterval = "12h"
    }
  })

  config = jsonencode({
    resources = []
  })
}


` + "```" + `

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
