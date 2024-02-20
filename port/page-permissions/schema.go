package page_permissions

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PagePermissionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"page_identifier": schema.StringAttribute{
			Required: true,
		},
		"read": schema.SingleNestedAttribute{
			MarkdownDescription: "The permission to read the page",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"users": schema.ListAttribute{
					MarkdownDescription: "The users with read permission",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"roles": schema.ListAttribute{
					MarkdownDescription: "The roles with read permission",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"teams": schema.ListAttribute{
					MarkdownDescription: "The teams with read permission",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		}}
}

func (r *PagePermissionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: PagePermissionsResourceMarkdownDescription,
		Attributes:          PagePermissionsSchema(),
	}
}

var PagePermissionsResourceMarkdownDescription = `

# Page Permissions resource

Docs about page permissions can be found [here](https://docs.getport.io/customize-pages-dashboards-and-plugins/page/page-permissions?view-permissions=api).

## Example Usage

### Allow read access to all members:

` + "```hcl" + `
resource "port_page_permissions" "microservices_permissions" {
  page_identifier = "microservices"
  read = {
    "roles": ["Member"],
    "users": [],
    "teams": [],
  }
}` + "\n```" + `

### Allow read access to all admins and a specific user and team:

` + "```hcl" + `
resource "port_page_permissions" "microservices_permissions" {
  page_identifier = "microservices"
  read = {
    "roles": [
      "Admin",
    ],
    "users": ["test-admin-user@test.com"],
    "teams": ["Team Spiderman"],
  }
}` + "\n```" + `

### Allow read access to specific users and teams:

` + "```hcl" + `
resource "port_page_permissions" "microservices_permissions" {
  page_identifier = "microservices"
  read = {
    "roles": [],
    "users": ["test-admin-user@test.com"],
    "teams": ["Team Spiderman"],
  }
}` + "\n```" + `

## Disclaimer 

- Page permissions are created by default when page is first created, this means that you should use this resource when you want to change the default permissions of a page.
- When deleting a page permissions resource using terraform, the page permissions will not be deleted from Port, as they are required for the action to work, instead, the page permissions will be removed from the terraform state.
`
