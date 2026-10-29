package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UserSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "The email address of the user",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"roles": schema.SetAttribute{
			MarkdownDescription: "The roles assigned to the user",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"teams": schema.SetAttribute{
			MarkdownDescription: "The teams assigned to the user",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"inactivity_timeout": schema.Int64Attribute{
			MarkdownDescription: "The inactivity timeout in minutes for the user. Must be at least 10 when set. Set to `null` to clear.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(10),
			},
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The status of the user",
			Computed:            true,
		},
		"managed_by_scim": schema.BoolAttribute{
			MarkdownDescription: "Whether the user is managed by SCIM",
			Computed:            true,
		},
	}
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "User resource to manage Port users, including roles, teams, and inactivity timeout",
		Attributes:          UserSchema(),
	}
}
