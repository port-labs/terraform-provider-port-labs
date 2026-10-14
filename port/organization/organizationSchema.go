package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		"inactivity_timeout": schema.Int64Attribute{
			MarkdownDescription: "The inactivity timeout in minutes for the organization. Must be at least 10 when set. Set to `null` to clear.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(10),
			},
		},
	}
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Organization resource to manage organization-level settings such as name and hidden blueprints",
		Attributes:          OrganizationSchema(),
	}
}
