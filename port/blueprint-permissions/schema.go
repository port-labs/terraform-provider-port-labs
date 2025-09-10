package blueprint_permissions

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getAssigneeProps(permName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"users": schema.SetAttribute{
			MarkdownDescription: fmt.Sprintf("Users with %+v permissions", permName),
			Optional:            true,
			ElementType:         types.StringType,
		},
		"roles": schema.SetAttribute{
			MarkdownDescription: fmt.Sprintf("Roles with %+v permissions", permName),
			Optional:            true,
			ElementType:         types.StringType,
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: fmt.Sprintf("Teams with %+v permissions", permName),
			Optional:            true,
			ElementType:         types.StringType,
		},
		"owned_by_team": schema.BoolAttribute{
			MarkdownDescription: "Owned by team",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
	}

}

func BlueprintPermissionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"blueprint_identifier": schema.StringAttribute{
			Required: true,
		},
		"entities": schema.SingleNestedAttribute{
			MarkdownDescription: "Entities permissions to read the blueprint",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"register": schema.SingleNestedAttribute{
					MarkdownDescription: "Manage permissions to register entities of the blueprint",
					Required:            true,
					Attributes:          getAssigneeProps("register"),
				},
				"unregister": schema.SingleNestedAttribute{
					MarkdownDescription: "Manage permissions to unregister entities of the blueprint",
					Required:            true,
					Attributes:          getAssigneeProps("unregister"),
				},
				"update": schema.SingleNestedAttribute{
					MarkdownDescription: "Manage permissions to update entities of the blueprint",
					Required:            true,
					Attributes:          getAssigneeProps("update"),
				},
				"update_metadata_properties": schema.SingleNestedAttribute{
					Required: true,
					MarkdownDescription: `Manage permissions to the metadata properties (` + "`" + `$icon|$title|$team|$identifier` + "`)" + `
These are translated to the updateProperties in the Port Api, proxied since we can't have Terraform properties starting with ` + "`$`" + `signs.
See [here](https://docs.getport.io/build-your-software-catalog/customize-integrations/configure-data-model/setup-blueprint/properties/meta-properties/) for more details.`,
					Attributes: map[string]schema.Attribute{
						"icon": schema.SingleNestedAttribute{
							MarkdownDescription: "The entity's icon",
							Required:            true,
							Attributes:          getAssigneeProps("update `$icon` metadata"),
						},
						"title": schema.SingleNestedAttribute{
							MarkdownDescription: "A human-readable name for the entity",
							Required:            true,
							Attributes:          getAssigneeProps("update `$title` metadata"),
						},
						"team": schema.SingleNestedAttribute{
							MarkdownDescription: "The team this entity belongs to",
							Required:            true,
							Attributes:          getAssigneeProps("update `$team` metadata"),
						},
						"identifier": schema.SingleNestedAttribute{
							MarkdownDescription: "Unique Entity identifier, used for API calls, programmatic access and distinguishing between different entities",
							Required:            true,
							Attributes:          getAssigneeProps("update `$identifier` metadata"),
						},
					},
				},
				"update_properties": schema.MapNestedAttribute{
					MarkdownDescription: "Manage permissions to update the entity properties",
					Optional:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: getAssigneeProps("update specific property"),
					},
					Validators: []validator.Map{
						mapvalidator.KeysAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[^\$]`),
								"`update_properties` can not start with `$`, those are reserved to update_metadata_properties",
							),
						),
					},
				},
				"update_relations": schema.MapNestedAttribute{
					MarkdownDescription: "Manage permissions to update the entity relations",
					Optional:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: getAssigneeProps("update specific relation"),
					},
				},
			},
		},
	}
}

func (r *BlueprintPermissionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: BlueprintPermissionsResourceMarkdownDescription,
		Attributes:          BlueprintPermissionsSchema(),
	}
}

var BlueprintPermissionsResourceMarkdownDescription = `

# Blueprint Permissions resource

Docs about blueprint permissions can be found [here](https://docs.getport.io/build-your-software-catalog/set-catalog-rbac/examples/#setting-blueprint-permissions)

` + "```hcl" + `
resource "port_blueprint_permissions" "microservices_permissions" {
	blueprint_identifier = "my_blueprint_identifier"
		entities             = {
			"register" = {
				"roles" : [
					"Member",
				],
				"users" : [],
				"teams" : []
			},
		}
	}
}` + "\n```" + `

## Example Usage

### Allow access to all members:

` + "```hcl" + `
resource "port_blueprint_permissions" "microservices_permissions" {
	blueprint_identifier = "my_blueprint_identifier"
		entities             = {
			"register" = {
				"roles" : [
					"Member",
				],
				"users" : [],
				"teams" : []
			},
			"unregister" = {
				"roles" : [
					"Member",
				],
				"users" : [],
				"teams" : []
			},
			"update" = {
				"roles" : [
					"Member",
				],
				"users" : ["test-admin-user@test.com"],
				"teams" : []
			},
			"update_metadata_properties" = {
				"icon" = {
					"roles" : [
						"Member",
					],
					"users" : [],
					"teams" : []
				},
				"identifier" = {
					"roles" : [
						"Member",
					],
					"users" : [],
					"teams" : ["Team Spiderman"]
				},
				"team" = {
					"roles" : [
						"Admin",
					],
					"users" : [],
					"teams" : []
				},
				"title" = {
					"roles" : [
						"Member",
					],
					"users" : [],
					"teams" : []
				}
			}
		}
}` + "\n```" + `


### Allow update ` + "`" + `myStringProperty` + "``" + ` for admins and a specific user and team:

` + "```hcl" + `
resource "port_blueprint_permissions" "microservices_permissions" {
	blueprint_identifier = "my_blueprint_identifier"
		entities = {
			# all properties from the previous example...
			"update_properties" = {
				"myStringProperty" = {
					"roles": [
						"Admin",
					],
					"users": ["test-admin-user@test.com"],
					"teams": ["Team Spiderman"],
				}
			}
		}
	}
}` + "\n```" + `

### Allow update relations for a specific team for admins and a specific user and team:

` + "```hcl" + `
resource "port_blueprint_permissions" "microservices_permissions" {
	blueprint_identifier = "my_blueprint_identifier"
		entities = {
			# all properties from the first example...
			"update_relations" = {
				"myRelations" = {
					"roles": [
						"Admin",
					],
					"users": ["test-admin-user@test.com"],
					"teams": ["Team Spiderman"],
				}
			}
		}
}` + "\n```" + `

## Disclaimer

- Blueprint permissions are created by default when blueprint is first created, this means that you should use this resource when you want to change the default permissions of a blueprint.
- When deleting a blueprint permissions resource using terraform, the blueprint permissions will not be deleted from Port, as they are required for the action to work, instead, the blueprint permissions will be removed from the terraform state.
- You always need to explicity set ` + "`" + `register|unregister|update|update_metadata_properties` + "`" + ` properties.
- All the permission lists (roles, users, teams) are managed by Port in a sorted manner, this means that if your ` + "`" + `.tf` + "`" + ` has for example roles defined out of order, your state will be invalid
    E.g:

    ` + "```hcl" + `
	resource "port_blueprint_permissions" "microservices_permissions" {
		blueprint_identifier = "my_blueprint_identifier"
			entities             = {
				# invalid:
				"register" = {
					"roles" : [
						"Member",
					"Admin",
					],
					"users" : [],
					"teams" : []
				},
				# valid
				"register" = {
					"roles" : [
						"Admin",
					"Member",
					],
					"users" : [],
					"teams" : []
				},
				...
			},
		},
	}` + "\n```"
