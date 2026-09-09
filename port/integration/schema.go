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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			Optional: true,
		},
		"installation_app_type": schema.StringAttribute{
			MarkdownDescription: "The integrated tool name for catalog integration types (e.g. `github-ocean`, `gitlab`, `pagerduty`). Cannot be changed after creation.",
			Optional:            true,
		},
		"installation_type": schema.StringAttribute{
			MarkdownDescription: "The installation type of the integration. Use `Saas` for Ocean SaaS integrations (requires `spec`). Defaults to `OnPrem` for self-hosted integrations. Only `OnPrem` and `Saas` are supported by this resource. Cannot be changed after creation.",
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
			MarkdownDescription: "Ocean SaaS integration spec as a JSON string (use `jsonencode`). **Only supported when `installation_type` is `Saas`** — must not be set for OnPrem integrations. Required for SaaS. Contains `integrationSpec` (credentials/settings) and optionally `appSpec` (feature toggles like `liveEventsEnabled`, `sendRawDataExamples`, etc.). Sensitive `integrationSpec` values (org secret references) are preserved from your HCL since the server strips them on read. If `appSpec` fields are omitted, the server applies its own defaults — which may differ from Port UI defaults. Declare `appSpec` explicitly to match the UI behavior.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The provisioning status of the integration (e.g. `Creating`, `Running`, `Updating`, `Error`). Relevant for SaaS integrations that provision asynchronously via Ocean.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"config": schema.StringAttribute{
			MarkdownDescription: "Integration mapping and configuration as a JSON string (use `jsonencode`). **Cannot be set on creation** — integrations receive default mappings during provisioning. Add `config` after the initial `terraform apply` to override the defaults.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
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

Manages a Port integration, including self-hosted (OnPrem) and Ocean SaaS installations.

For SaaS integrations, create organization secrets first with ` + "`port_organization_secret`" + `, then reference secret names in ` + "`spec.integrationSpec`" + `.

Docs about integrations can be found [here](https://docs.getport.io/integrations-index/).

Docs about how to use Port's Terraform provider to create and manage integrations can be found [here](https://docs.getport.io/context-lake/ingestion/ingest-data-into-port/other/iac/terraform/terraform).

## Two-step workflow

Integrations receive default mappings during provisioning (from provision-service or the Ocean integration itself). Because of this, ` + "`config`" + ` **cannot be set when creating** an integration — it would be overwritten. Instead:

1. **First apply** — create the integration without ` + "`config`" + `. Provisioning sets up default blueprints and mappings.
2. **Second apply** — add a ` + "`config`" + ` block to your HCL to override the default mappings.

SaaS integrations additionally provision asynchronously (` + "`Creating`" + ` → ` + "`Running`" + `) via Ocean/Kafka. OnPrem create/update/delete are synchronous.

## SaaS example — Step 1: Create the integration

` + "```hcl" + `
# Secret naming convention: _{INSTALLATION_ID}_{INTEGRATION_TYPE}_{PROPERTY} in SCREAMING_SNAKE_CASE.
# This matches the frontend's generateSecretName() so Ocean can resolve the secret at runtime.
#
# Each integration type has its own spec fields — check the integration's .port/spec.json
# for the list of configurations (name, type, sensitive, dependencies, etc.).
locals {
  github_installation_id = "github-prod"
  github_type            = "github-ocean"
}

resource "port_organization_secret" "github_token" {
  secret_name  = "_${upper(replace("${local.github_installation_id}_${local.github_type}_github_token", "-", "_"))}"
  secret_value = var.github_token
}

resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = local.github_installation_id
  installation_app_type = local.github_type
  installation_type     = "Saas"
  title                 = "GitHub Production"

  spec = jsonencode({
    integrationSpec = {
      authenticationMode = "Personal Access Token"
      githubToken        = port_organization_secret.github_token.secret_name
    }

    # appSpec controls feature toggles. If omitted, the server applies its own
    # defaults which may differ from the Port UI defaults. Declare explicitly
    # to match the UI behaviour.
    appSpec = {
      scheduledResyncInterval  = "12h"
      sendRawDataExamples      = true
      liveEventsEnabled        = false
      actionsProcessingEnabled = false
      incrementalSyncEnabled   = false
    }
  })

  # Do NOT set config here — provisioning will populate default mappings.
  # Add config on a subsequent apply to override them.
}
` + "```" + `

## SaaS example — Step 2: Override mappings

After the first apply completes and provisioning finishes, add ` + "`config`" + ` to the same resource:

` + "```hcl" + `
resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = local.github_installation_id
  installation_app_type = local.github_type
  installation_type     = "Saas"
  title                 = "GitHub Production"

  spec = jsonencode({
    integrationSpec = {
      authenticationMode = "Personal Access Token"
      githubToken        = port_organization_secret.github_token.secret_name
    }
  })

  config = jsonencode({
    resources = [
      {
        kind = "repository"
        selector = {
          query = "true"
        }
        port = {
          entity = {
            mappings = {
              identifier = ".name"
              blueprint  = "\"githubRepository\""
              title      = ".name"
              properties = {
                url           = ".html_url"
                defaultBranch = ".default_branch"
              }
            }
          }
        }
      }
    ]
  })
}
` + "```" + `

## Self-hosted (OnPrem) example — Step 1: Create

Self-hosted integrations run on your own infrastructure (e.g. an Ocean exporter container) and pull their mapping from Port.

` + "```hcl" + `
resource "port_integration" "my_custom_integration" {
  installation_id = "my-custom-integration-id"
  title           = "My Custom Integration"

  # config is set on the next apply, after provisioning creates default mappings.
}
` + "```" + `

## Self-hosted (OnPrem) example — Step 2: Override mappings

` + "```hcl" + `
resource "port_integration" "my_custom_integration" {
  installation_id = "my-custom-integration-id"
  title           = "My Custom Integration"

  config = jsonencode({
    createMissingRelatedEntities = true
    deleteDependentEntities      = true
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
            relations = {}
          }]
        }
      }
    }]
  })
}
` + "```\n" + `

For catalog integration types, set ` + "`installation_app_type`" + ` to the integrated tool name (e.g. ` + "`github-ocean`" + `, ` + "`gitlab`" + `) and ` + "`version`" + ` if you want to pin a specific integration version. Custom integrations can omit ` + "`installation_app_type`" + `.

### NOTICE:

The following config properties (` + "`selector.query|entity.mappings.*`" + `) are jq expressions, which means that you need to input either a valid jq expression (E.g ` + "`.title`" + `), or if you want a string value, a quoted escaped string val (E.g ` + "`'my-string'`" + `).

### NOTES:

- ` + "`config`" + ` **cannot be set on creation**. Integrations receive default mappings during provisioning. Create first, then add ` + "`config`" + ` on a subsequent apply.
- ` + "`status`" + ` reflects async Ocean provisioning (` + "`Creating`" + ` → ` + "`Running`" + `) for SaaS integrations only.
- ` + "`installation_id`" + `, ` + "`installation_app_type`" + `, and ` + "`installation_type`" + ` cannot be changed after creation.
- ` + "`spec`" + ` is only supported for SaaS integrations (` + "`installation_type = \"Saas\"`" + `). Do not set it on OnPrem integrations.
- A changelog destination (` + "`webhook_changelog_destination`" + ` / ` + "`kafka_changelog_destination`" + `) can be added or updated, but not removed — the Port API does not support clearing it. To remove it, delete and recreate the integration (e.g. taint the resource).
- Existing integrations can be brought under Terraform management with ` + "`terraform import port_integration.my_integration <installation_id>`" + `.
- ` + "`terraform destroy`" + ` deletes the real integration in Port, not just removes it from state. Use ` + "`terraform state rm`" + ` if you only want to stop managing an integration with Terraform without deleting it from Port. This is especially relevant for imported resources.
`
