package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func OrganizationSecretSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"secret_name": schema.StringAttribute{
			MarkdownDescription: "The name of the organization secret",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"secret_value": schema.StringAttribute{
			MarkdownDescription: "The value of the organization secret",
			Required:            true,
			Sensitive:           true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the organization secret",
			Optional:            true,
		},
	}
}

func (r *OrganizationSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: organizationSecretMarkdownDescription,
		Attributes:          OrganizationSecretSchema(),
	}
}

var organizationSecretMarkdownDescription = `

# Organization Secret Resource

Manages a Port organization secret. Organization secrets store sensitive values (API tokens, credentials) that can be referenced by integrations without exposing the actual values.

Docs about organization secrets can be found [here](https://docs.getport.io/build-your-software-catalog/sync-data-to-catalog/organization-secrets/).

~> **Important:** The ` + "`secret_value`" + ` is stored in plaintext in Terraform state. Use appropriate state encryption and access controls.

~> **Note:** The API does not return secret values, so Terraform cannot detect drift on the value. If you update the secret outside Terraform, run ` + "`terraform apply`" + ` to ensure the value matches your configuration.

## Example Usage

### Basic secret

` + "```hcl" + `
variable "github_token" {
  type      = string
  sensitive = true
}

resource "port_organization_secret" "github_token" {
  secret_name  = "github_token"
  secret_value = var.github_token
  description  = "GitHub personal access token"
}
` + "```" + `

### Use secret in an integration

` + "```hcl" + `
variable "github_token" {
  type      = string
  sensitive = true
}

resource "port_organization_secret" "github_token" {
  secret_name  = "github_token"
  secret_value = var.github_token
}

resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = "github-prod"
  installation_app_type = "github-ocean"
  installation_type     = "Saas"

  spec = jsonencode({
    integrationSpec = {
      githubToken = port_organization_secret.github_token.secret_name
    }
  })
}
` + "```" + `

## Import

Organization secrets can be imported using the secret name:

` + "```shell" + `
terraform import port_organization_secret.github_token github_token
` + "```" + `

Note that the ` + "`secret_value`" + ` cannot be imported from the API since it is write-only. After import, you must provide the ` + "`secret_value`" + ` in your configuration and run ` + "`terraform apply`" + ` to update the state.
`
