package system_blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

type Resource struct {
	client *cli.PortClient
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*cli.PortClient)
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	blueprintSchemas := blueprint.BlueprintSchema()

	resp.Schema = schema.Schema{
		MarkdownDescription: systemBlueprintMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the system blueprint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Identifier of the system blueprint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"properties":               blueprintSchemas["properties"],
			"relations":                blueprintSchemas["relations"],
			"mirror_properties":        blueprintSchemas["mirror_properties"],
			"calculation_properties":   blueprintSchemas["calculation_properties"],
			"include_in_global_search": blueprintSchemas["include_in_global_search"],
		},
	}
}

var systemBlueprintMarkdownDescription = `

# System Blueprint Resource

Extends Port system blueprints (` + "`_user`" + `, ` + "`_team`" + `, etc.) with additional properties, relations, and calculation properties. System blueprints are built-in blueprints that Port provides for common concepts.

~> **Note:** System blueprints cannot be created or deleted. This resource only manages extensions (properties, relations) added to them. Changes affect all entities of that system blueprint across the organization.

## Example Usage

### Add custom properties to the _user blueprint

` + "```hcl" + `
resource "port_system_blueprint" "user" {
  identifier = "_user"
  properties = {
    string_props = {
      "department" = {
        title = "Department"
      }
      "github_username" = {
        title = "GitHub Username"
      }
    }
  }
}
` + "```" + `

### Add relation to the _team blueprint

` + "```hcl" + `
resource "port_blueprint" "cost_center" {
  identifier = "cost_center"
  title      = "Cost Center"
}

resource "port_system_blueprint" "team" {
  identifier = "_team"
  relations = {
    "cost_center" = {
      title  = "Cost Center"
      target = port_blueprint.cost_center.identifier
    }
  }
}
` + "```" + `

## Available System Blueprints

- ` + "`_user`" + ` - Port users
- ` + "`_team`" + ` - Port teams
- ` + "`_rule`" + ` - Scorecard rules (when scorecard groups are enabled)

## Import

System blueprints can be imported using their identifier:

` + "```shell" + `
terraform import port_system_blueprint.user _user
` + "```" + `
`
