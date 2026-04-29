package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationModel struct {
	ID                                       types.String `tfsdk:"id"`
	Name                                     types.String `tfsdk:"name"`
	HiddenBlueprints                         types.List   `tfsdk:"hidden_blueprints"`
	FederatedLogout                          types.Bool   `tfsdk:"federated_logout"`
	PortalIcon                               types.String `tfsdk:"portal_icon"`
	PortalTitle                              types.String `tfsdk:"portal_title"`
	SupportUserPermission                    types.String `tfsdk:"support_user_permission"`
	SupportUserTTL                           types.String `tfsdk:"support_user_ttl"`
	SupportUserExpiresAt                     types.String `tfsdk:"support_user_expires_at"`
	PortAgentStreamerName                    types.String `tfsdk:"port_agent_streamer_name"`
	IncludeBlueprintsInGlobalSearchByDefault types.Bool   `tfsdk:"include_blueprints_in_global_search_by_default"`
	IsOnboarded                              types.Bool   `tfsdk:"is_onboarded"`
	ToolSelectionProvisioningStatus          types.String `tfsdk:"tool_selection_provisioning_status"`
	AnnouncementEnabled                      types.Bool   `tfsdk:"announcement_enabled"`
	AnnouncementContent                      types.String `tfsdk:"announcement_content"`
	AnnouncementLink                         types.String `tfsdk:"announcement_link"`
	AnnouncementColor                        types.String `tfsdk:"announcement_color"`
}
