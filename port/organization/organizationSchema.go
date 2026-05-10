package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		"hidden_blueprints": schema.ListAttribute{
			MarkdownDescription: "A list of blueprint identifiers to hide from the portal",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"federated_logout": schema.BoolAttribute{
			MarkdownDescription: "Enable or disable federated logout",
			Optional:            true,
			Computed:            true,
		},
		"portal_icon": schema.StringAttribute{
			MarkdownDescription: "The icon URI for the portal (must be a valid URI)",
			Optional:            true,
			Computed:            true,
		},
		"portal_title": schema.StringAttribute{
			MarkdownDescription: "The title for the portal",
			Optional:            true,
			Computed:            true,
		},
		"support_user_permission": schema.StringAttribute{
			MarkdownDescription: "The support user permission level (e.g. OPT_OUT)",
			Optional:            true,
			Computed:            true,
		},
		"support_user_ttl": schema.StringAttribute{
			MarkdownDescription: "The TTL for support user access (e.g. ONE_DAY)",
			Optional:            true,
			Computed:            true,
		},
		"support_user_expires_at": schema.StringAttribute{
			MarkdownDescription: "The expiration date for support user access (ISO 8601 format)",
			Optional:            true,
			Computed:            true,
		},
		"port_agent_streamer_name": schema.StringAttribute{
			MarkdownDescription: "The name of the Port agent streamer (e.g. KAFKA)",
			Optional:            true,
			Computed:            true,
		},
		"include_blueprints_in_global_search_by_default": schema.BoolAttribute{
			MarkdownDescription: "Whether to include blueprints in global search by default",
			Optional:            true,
			Computed:            true,
		},
		"is_onboarded": schema.BoolAttribute{
			MarkdownDescription: "Whether the organization has completed onboarding",
			Optional:            true,
			Computed:            true,
		},
		"tool_selection_provisioning_status": schema.StringAttribute{
			MarkdownDescription: "The status of tool selection provisioning (e.g. IN_PROGRESS)",
			Optional:            true,
			Computed:            true,
		},
		"announcement_enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether the organization announcement is enabled",
			Optional:            true,
			Computed:            true,
		},
		"announcement_content": schema.StringAttribute{
			MarkdownDescription: "The content of the organization announcement",
			Optional:            true,
			Computed:            true,
		},
		"announcement_link": schema.StringAttribute{
			MarkdownDescription: "A link for the organization announcement",
			Optional:            true,
			Computed:            true,
		},
		"announcement_color": schema.StringAttribute{
			MarkdownDescription: "The color of the organization announcement (e.g. blue)",
			Optional:            true,
			Computed:            true,
		},
	}
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Organization resource to manage organization-level settings such as name and hidden blueprints",
		Attributes:          OrganizationSchema(),
	}
}
