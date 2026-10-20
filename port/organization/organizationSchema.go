package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OrganizationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the organization",
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the organization",
			Optional:            true,
			Computed:            true,
		},
	}
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: organizationMarkdownDescription,
		Attributes:          OrganizationSchema(),
	}
}

var organizationMarkdownDescription = `

# Organization Resource

Manages Port organization settings. Currently supports setting the organization name.

~> **Note:** This resource manages organization-wide settings. Changes affect all users in the organization.

## Example Usage

` + "```hcl" + `
resource "port_organization" "my_org" {
  name = "My Organization"
}
` + "```" + `
`
