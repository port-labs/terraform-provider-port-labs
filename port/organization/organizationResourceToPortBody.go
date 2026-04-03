package organization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func organizationResourceToPortBody(ctx context.Context, state *OrganizationModel) (*cli.OrganizationUpdate, error) {
	update := &cli.OrganizationUpdate{}

	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		update.Name = state.Name.ValueStringPointer()
	}

	// Settings: fill a local struct and assign once. If we set update.Settings to a non-nil
	// pointer up front, encoding/json still emits "settings":{} (omitempty only drops nil
	// pointers), which can differ from omitting the key entirely.
	settings := &cli.OrganizationSettings{}
	settingsSet := false

	if !state.HiddenBlueprints.IsNull() && !state.HiddenBlueprints.IsUnknown() {
		var blueprints []string
		diags := state.HiddenBlueprints.ElementsAs(ctx, &blueprints, false)
		if diags.HasError() {
			return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
		}
		settings.HiddenBlueprints = blueprints
		settingsSet = true
	}
	if !state.FederatedLogout.IsNull() && !state.FederatedLogout.IsUnknown() {
		settings.FederatedLogout = state.FederatedLogout.ValueBoolPointer()
		settingsSet = true
	}
	if !state.PortalIcon.IsNull() && !state.PortalIcon.IsUnknown() {
		settings.PortalIcon = state.PortalIcon.ValueStringPointer()
		settingsSet = true
	}
	if !state.PortalTitle.IsNull() && !state.PortalTitle.IsUnknown() {
		settings.PortalTitle = state.PortalTitle.ValueStringPointer()
		settingsSet = true
	}
	if !state.SupportUserPermission.IsNull() && !state.SupportUserPermission.IsUnknown() {
		settings.SupportUserPermission = state.SupportUserPermission.ValueStringPointer()
		settingsSet = true
	}
	if !state.SupportUserTTL.IsNull() && !state.SupportUserTTL.IsUnknown() {
		settings.SupportUserTTL = state.SupportUserTTL.ValueStringPointer()
		settingsSet = true
	}
	if !state.SupportUserExpiresAt.IsNull() && !state.SupportUserExpiresAt.IsUnknown() {
		settings.SupportUserExpiresAt = state.SupportUserExpiresAt.ValueStringPointer()
		settingsSet = true
	}
	if !state.PortAgentStreamerName.IsNull() && !state.PortAgentStreamerName.IsUnknown() {
		settings.PortAgentStreamerName = state.PortAgentStreamerName.ValueStringPointer()
		settingsSet = true
	}
	if !state.IncludeBlueprintsInGlobalSearchByDefault.IsNull() && !state.IncludeBlueprintsInGlobalSearchByDefault.IsUnknown() {
		settings.IncludeBlueprintsInGlobalSearchByDefault = state.IncludeBlueprintsInGlobalSearchByDefault.ValueBoolPointer()
		settingsSet = true
	}
	if settingsSet {
		update.Settings = settings
	}

	// Top-level fields
	if !state.IsOnboarded.IsNull() && !state.IsOnboarded.IsUnknown() {
		update.IsOnboarded = state.IsOnboarded.ValueBoolPointer()
	}

	// Tool selection provisioning
	if !state.ToolSelectionProvisioningStatus.IsNull() && !state.ToolSelectionProvisioningStatus.IsUnknown() {
		update.ToolSelectionProvisioning = &cli.OrganizationToolSelectionProvisioning{
			Status: state.ToolSelectionProvisioningStatus.ValueStringPointer(),
		}
	}

	// Announcement: default enabled to false when any announcement attribute is set, then apply
	// known values from state (so content/link/color without an explicit enabled still disables by default).
	if state.AnnouncementEnabled.ValueBoolPointer() != nil ||
		state.AnnouncementContent.ValueStringPointer() != nil ||
		state.AnnouncementLink.ValueStringPointer() != nil ||
		state.AnnouncementColor.ValueStringPointer() != nil {
		enabledFalse := false
		update.Announcement = &cli.OrganizationAnnouncement{Enabled: &enabledFalse}
		if !state.AnnouncementEnabled.IsNull() && !state.AnnouncementEnabled.IsUnknown() {
			update.Announcement.Enabled = state.AnnouncementEnabled.ValueBoolPointer()
		}
		if !state.AnnouncementContent.IsNull() && !state.AnnouncementContent.IsUnknown() {
			update.Announcement.Content = state.AnnouncementContent.ValueStringPointer()
		}
		if !state.AnnouncementLink.IsNull() && !state.AnnouncementLink.IsUnknown() {
			update.Announcement.Link = state.AnnouncementLink.ValueStringPointer()
		}
		if !state.AnnouncementColor.IsNull() && !state.AnnouncementColor.IsUnknown() {
			update.Announcement.Color = state.AnnouncementColor.ValueStringPointer()
		}
	}

	return update, nil
}

func hiddenBlueprintsToList(ctx context.Context, blueprints []string) (types.List, error) {
	if len(blueprints) == 0 {
		return types.ListNull(types.StringType), nil
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, blueprints)
	if diags.HasError() {
		return types.ListNull(types.StringType), fmt.Errorf("%s", diags.Errors()[0].Detail())
	}
	return list, nil
}
