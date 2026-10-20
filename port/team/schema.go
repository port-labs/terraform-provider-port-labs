package team

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TeamSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the team",
			Required:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The machine-readable identifier of the _team entity",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the team",
			Optional:            true,
		},
		"users": schema.SetAttribute{
			MarkdownDescription: "The users of the team",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"provider_name": schema.StringAttribute{
			MarkdownDescription: "The provider of the team",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the team",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the team",
			Computed:            true,
		},
	}
}
func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: teamResourceMarkdownDescription,
		Attributes:          TeamSchema(),
	}
}

var teamResourceMarkdownDescription = `

# Team Resource

Manages a Port team. Teams can be assigned ownership of entities and used in permissions.

Docs about teams can be found [here](https://docs.getport.io/sso-rbac/rbac-overview/).

## Example Usage

### Basic team

` + "```hcl" + `
resource "port_team" "backend" {
  name        = "Backend Team"
  description = "The backend engineering team"
  users = [
    "developer1@example.com",
    "developer2@example.com",
  ]
}
` + "```" + `

### Team with entity ownership

Use the team's ` + "`identifier`" + ` attribute to assign entities to teams. The identifier is automatically generated from the team name.

` + "```hcl" + `
resource "port_team" "platform" {
  name = "Platform Team"
  users = [
    "platform-lead@example.com",
  ]
}

resource "port_blueprint" "service" {
  identifier = "service"
  title      = "Service"
  ownership = {
    type = "Direct"
  }
}

resource "port_entity" "api_gateway" {
  identifier = "api-gateway"
  title      = "API Gateway"
  blueprint  = port_blueprint.service.identifier
  teams      = [port_team.platform.identifier]
}
` + "```" + `

## Import

Teams can be imported using the team name:

` + "```shell" + `
terraform import port_team.backend "Backend Team"
` + "```" + `
`
