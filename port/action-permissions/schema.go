package action_permissions

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ActionPermissionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"action_identifier": schema.StringAttribute{
			Description: "The ID of the action",
			Required:    true,
		},
		"blueprint_identifier": schema.StringAttribute{
			Description:        "The ID of the blueprint",
			Optional:           true,
			DeprecationMessage: "Action is not attached to blueprint anymore. This value is ignored",
			Validators:         []validator.String{stringvalidator.OneOf("")},
		},
		"permissions": schema.SingleNestedAttribute{
			MarkdownDescription: "The permissions for the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"execute": schema.SingleNestedAttribute{
					MarkdownDescription: "The permission to execute the action",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"users": schema.ListAttribute{
							MarkdownDescription: "The users with execution permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"roles": schema.ListAttribute{
							MarkdownDescription: "The roles with execution permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"teams": schema.ListAttribute{
							MarkdownDescription: "The teams with execution permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"owned_by_team": schema.BoolAttribute{
							MarkdownDescription: "Give execution permission to the teams who own the entity",
							Optional:            true,
						},
						"policy": schema.StringAttribute{
							MarkdownDescription: "The policy to use for execution",
							Optional:            true,
						},
					},
				},
				"approve": schema.SingleNestedAttribute{
					MarkdownDescription: "The permission to approve the action's runs",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"users": schema.ListAttribute{
							MarkdownDescription: "The users with approval permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"roles": schema.ListAttribute{
							MarkdownDescription: "The roles with approval permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"teams": schema.ListAttribute{
							MarkdownDescription: "The teams with approval permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"policy": schema.StringAttribute{
							MarkdownDescription: "The policy to use for approval",
							Optional:            true,
						},
					},
				},
			},
		}}
}

func (r *ActionPermissionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ActionPermissionsResourceMarkdownDescription,
		Attributes:          ActionPermissionsSchema(),
	}
}

var ActionPermissionsResourceMarkdownDescription = `

# Action Permissions resource

Docs for the Action Permissions resource can be found [here](https://docs.getport.io/create-self-service-experiences/set-self-service-actions-rbac/examples).

## Example Usage

` + "```hcl" + `
resource "port_action_permissions" "restart_microservice_permissions" {
  action_identifier = port_action.restart_microservice.identifier
  permissions = {
    "execute": {
      "roles": [
        "Admin"
      ],
      "users": [],
      "teams": [],
      "owned_by_team": true
    },
    "approve": {
      "roles": ["Member", "Admin"],
      "users": [],
      "teams": []
    }
  }
}` + "\n```" + `

## Example Usage with Policy

Port allows setting dynamic permissions for executing and/or approving execution of self-service actions, based on any properties/relations of an action's corresponding blueprint.

Docs about the Policy language can be found [here](https://docs.getport.io/create-self-service-experiences/set-self-service-actions-rbac/dynamic-permissions#configuring-permissions).

Policy is expected to be passed as a JSON string and not as an object, this means that the evaluation of the policy will be done by Port and not by Terraform.
To pass a JSON string to Terraform, you can use the [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function.

` + "```hcl" + `
resource "port_action_permissions" "restart_microservice_permissions" {
  action_identifier = port_action.restart_microservice.identifier
  permissions = {
    "execute": {
      "roles": [
        "Admin"
      ],
      "users": [],
      "teams": [],
      "owned_by_team": true
    },
    "approve": {
      "roles": ["Member", "Admin"],
      "users": [],
      "teams": []
      # Terraform's "jsonencode" function converts a
      # Terraform expression result to valid JSON syntax.
      "policy": jsonencode(
        {
          queries: {
            executingUser: {
              rules: [
                {
                  value: "user",
                  operator: "=",
                  property: "$blueprint"
                },
                {
                    value: "true",
                    operator: "=",
                    property: "$owned_by_team"

                }
              ],
              combinator: "and"
            }
          },
          conditions: [
          "true"]
        }
      )
    }
  }
}` + "\n```" + `

## Disclaimer 

- Action permissions are created by default when creating a new action, this means that you should use this resource when you want to change the default permissions of an action.
- When deleting an action permissions resource using terraform, the action permissions will not be deleted from Port, as they are required for the action to work, instead, the action permissions will be removed from the terraform state.
`
