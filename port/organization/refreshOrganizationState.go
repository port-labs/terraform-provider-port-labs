package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
)

func refreshOrganizationState(ctx context.Context, state *OrganizationModel, org *cli.Organization) error {
	state.ID = types.StringValue(org.Name)
	state.Name = types.StringValue(org.Name)

	// Settings
	orgSettings := org.Settings
	if orgSettings == nil {
		orgSettings = &cli.OrganizationSettings{}
	}
	list, err := hiddenBlueprintsToList(ctx, orgSettings.HiddenBlueprints)
	if err != nil {
		return err
	}
	state.HiddenBlueprints = list
	state.FederatedLogout = flex.GoBoolToFramework(orgSettings.FederatedLogout)
	state.PortalIcon = flex.GoStringToFramework(orgSettings.PortalIcon)
	state.PortalTitle = flex.GoStringToFramework(orgSettings.PortalTitle)
	state.SupportUserPermission = flex.GoStringToFramework(orgSettings.SupportUserPermission)
	state.SupportUserTTL = flex.GoStringToFramework(orgSettings.SupportUserTTL)
	state.SupportUserExpiresAt = flex.GoStringToFramework(orgSettings.SupportUserExpiresAt)
	state.PortAgentStreamerName = flex.GoStringToFramework(orgSettings.PortAgentStreamerName)
	state.IncludeBlueprintsInGlobalSearchByDefault = flex.GoBoolToFramework(orgSettings.IncludeBlueprintsInGlobalSearchByDefault)

	// IsOnboarded
	state.IsOnboarded = flex.GoBoolToFramework(org.IsOnboarded)

	// Tool selection provisioning
	if org.ToolSelectionProvisioning != nil {
		state.ToolSelectionProvisioningStatus = flex.GoStringToFramework(org.ToolSelectionProvisioning.Status)
	} else {
		state.ToolSelectionProvisioningStatus = types.StringNull()
	}

	// Announcement
	orgAnnouncement := org.Announcement
	if orgAnnouncement == nil {
		orgAnnouncement = &cli.OrganizationAnnouncement{}
	}
	state.AnnouncementEnabled = flex.GoBoolToFramework(orgAnnouncement.Enabled)
	state.AnnouncementContent = flex.GoStringToFramework(orgAnnouncement.Content)
	state.AnnouncementLink = flex.GoStringToFramework(orgAnnouncement.Link)
	state.AnnouncementColor = flex.GoStringToFramework(orgAnnouncement.Color)

	return nil
}
