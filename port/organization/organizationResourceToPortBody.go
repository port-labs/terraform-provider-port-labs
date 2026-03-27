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
		name := state.Name.ValueString()
		update.Name = &name
	}

	// Settings fields
	if !state.HiddenBlueprints.IsNull() && !state.HiddenBlueprints.IsUnknown() {
		var blueprints []string
		diags := state.HiddenBlueprints.ElementsAs(ctx, &blueprints, false)
		if diags.HasError() {
			return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
		}
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		update.Settings.HiddenBlueprints = blueprints
	}
	if !state.FederatedLogout.IsNull() && !state.FederatedLogout.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.FederatedLogout.ValueBool()
		update.Settings.FederatedLogout = &v
	}
	if !state.PortalIcon.IsNull() && !state.PortalIcon.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.PortalIcon.ValueString()
		update.Settings.PortalIcon = &v
	}
	if !state.PortalTitle.IsNull() && !state.PortalTitle.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.PortalTitle.ValueString()
		update.Settings.PortalTitle = &v
	}
	if !state.SupportUserPermission.IsNull() && !state.SupportUserPermission.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.SupportUserPermission.ValueString()
		update.Settings.SupportUserPermission = &v
	}
	if !state.SupportUserTTL.IsNull() && !state.SupportUserTTL.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.SupportUserTTL.ValueString()
		update.Settings.SupportUserTTL = &v
	}
	if !state.SupportUserExpiresAt.IsNull() && !state.SupportUserExpiresAt.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.SupportUserExpiresAt.ValueString()
		update.Settings.SupportUserExpiresAt = &v
	}
	if !state.PortAgentStreamerName.IsNull() && !state.PortAgentStreamerName.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.PortAgentStreamerName.ValueString()
		update.Settings.PortAgentStreamerName = &v
	}
	if !state.IncludeBlueprintsInGlobalSearchByDefault.IsNull() && !state.IncludeBlueprintsInGlobalSearchByDefault.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettings{}
		}
		v := state.IncludeBlueprintsInGlobalSearchByDefault.ValueBool()
		update.Settings.IncludeBlueprintsInGlobalSearchByDefault = &v
	}

	// Top-level fields
	if !state.IsOnboarded.IsNull() && !state.IsOnboarded.IsUnknown() {
		v := state.IsOnboarded.ValueBool()
		update.IsOnboarded = &v
	}

	// Tool selection provisioning
	if !state.ToolSelectionProvisioningStatus.IsNull() && !state.ToolSelectionProvisioningStatus.IsUnknown() {
		v := state.ToolSelectionProvisioningStatus.ValueString()
		update.ToolSelectionProvisioning = &cli.OrganizationToolSelectionProvisioning{Status: &v}
	}

	// Announcement
	hasAnnouncement := (!state.AnnouncementEnabled.IsNull() && !state.AnnouncementEnabled.IsUnknown()) ||
		(!state.AnnouncementContent.IsNull() && !state.AnnouncementContent.IsUnknown()) ||
		(!state.AnnouncementLink.IsNull() && !state.AnnouncementLink.IsUnknown()) ||
		(!state.AnnouncementColor.IsNull() && !state.AnnouncementColor.IsUnknown())
	if hasAnnouncement {
		update.Announcement = &cli.OrganizationAnnouncement{}
		if !state.AnnouncementEnabled.IsNull() && !state.AnnouncementEnabled.IsUnknown() {
			v := state.AnnouncementEnabled.ValueBool()
			update.Announcement.Enabled = &v
		}
		if !state.AnnouncementContent.IsNull() && !state.AnnouncementContent.IsUnknown() {
			v := state.AnnouncementContent.ValueString()
			update.Announcement.Content = &v
		}
		if !state.AnnouncementLink.IsNull() && !state.AnnouncementLink.IsUnknown() {
			v := state.AnnouncementLink.ValueString()
			update.Announcement.Link = &v
		}
		if !state.AnnouncementColor.IsNull() && !state.AnnouncementColor.IsUnknown() {
			v := state.AnnouncementColor.ValueString()
			update.Announcement.Color = &v
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
