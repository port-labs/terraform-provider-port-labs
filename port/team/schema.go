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
		MarkdownDescription: "Team resource",
		Attributes:          TeamSchema(),
	}
}
