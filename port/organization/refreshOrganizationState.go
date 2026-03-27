package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshOrganizationState(ctx context.Context, state *OrganizationModel, org *cli.Organization) error {
	state.ID = types.StringValue(org.Name)
	state.Name = types.StringValue(org.Name)

	// Settings
	if org.Settings != nil {
		list, err := hiddenBlueprintsToList(ctx, org.Settings.HiddenBlueprints)
		if err != nil {
			return err
		}
		state.HiddenBlueprints = list
		state.FederatedLogout = optionalBoolValue(org.Settings.FederatedLogout)
		state.PortalIcon = optionalStringValue(org.Settings.PortalIcon)
		state.PortalTitle = optionalStringValue(org.Settings.PortalTitle)
		state.SupportUserPermission = optionalStringValue(org.Settings.SupportUserPermission)
		state.SupportUserTTL = optionalStringValue(org.Settings.SupportUserTTL)
		state.SupportUserExpiresAt = optionalStringValue(org.Settings.SupportUserExpiresAt)
		state.PortAgentStreamerName = optionalStringValue(org.Settings.PortAgentStreamerName)
		state.IncludeBlueprintsInGlobalSearchByDefault = optionalBoolValue(org.Settings.IncludeBlueprintsInGlobalSearchByDefault)
	} else {
		state.HiddenBlueprints = types.ListNull(types.StringType)
		state.FederatedLogout = types.BoolNull()
		state.PortalIcon = types.StringNull()
		state.PortalTitle = types.StringNull()
		state.SupportUserPermission = types.StringNull()
		state.SupportUserTTL = types.StringNull()
		state.SupportUserExpiresAt = types.StringNull()
		state.PortAgentStreamerName = types.StringNull()
		state.IncludeBlueprintsInGlobalSearchByDefault = types.BoolNull()
	}

	// IsOnboarded
	state.IsOnboarded = optionalBoolValue(org.IsOnboarded)

	// Tool selection provisioning
	if org.ToolSelectionProvisioning != nil {
		state.ToolSelectionProvisioningStatus = optionalStringValue(org.ToolSelectionProvisioning.Status)
	} else {
		state.ToolSelectionProvisioningStatus = types.StringNull()
	}

	// Announcement
	if org.Announcement != nil {
		state.AnnouncementEnabled = optionalBoolValue(org.Announcement.Enabled)
		state.AnnouncementContent = optionalStringValue(org.Announcement.Content)
		state.AnnouncementLink = optionalStringValue(org.Announcement.Link)
		state.AnnouncementColor = optionalStringValue(org.Announcement.Color)
	} else {
		state.AnnouncementEnabled = types.BoolNull()
		state.AnnouncementContent = types.StringNull()
		state.AnnouncementLink = types.StringNull()
		state.AnnouncementColor = types.StringNull()
	}

	return nil
}

func optionalStringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func optionalBoolValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}
