package action_permissions

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ActionPermissionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"blueprint_identifier": schema.StringAttribute{
			Description: "The ID of the blueprint",
			Required:    true,
		},
		"action_identifier": schema.StringAttribute{
			Description: "The ID of the action",
			Required:    true,
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
							MarkdownDescription: "The policy to use for approval",
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
		MarkdownDescription: "Action Permissions resource",
		Attributes:          ActionPermissionsSchema(),
	}
}
