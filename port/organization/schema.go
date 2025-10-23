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
		MarkdownDescription: "Organization secret resource",
		Attributes:          OrganizationSecretSchema(),
	}
}
