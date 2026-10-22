package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		"hidden_global_search_tabs": schema.SetAttribute{
			MarkdownDescription: "Spotlight tabs to hide from global search. Allowed values: `all`, `blueprint`, `page`, `action`, `workflow`, `entity`, `markdown`. At least one content tab must remain visible, and `all` cannot be the only visible tab. Which content tabs count depends on whether classic actions are disabled.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(stringvalidator.OneOf(
					"all",
					"blueprint",
					"page",
					"action",
					"workflow",
					"entity",
					"markdown",
				)),
			},
		},
	}
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Organization resource to manage organization-level settings such as name and hidden global search tabs",
		Attributes:          OrganizationSchema(),
	}
}
