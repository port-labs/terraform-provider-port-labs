package organization

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func organizationResourceToPortBody(ctx context.Context, state *OrganizationModel) (*cli.OrganizationUpdate, error) {
	update := &cli.OrganizationUpdate{}

	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name := state.Name.ValueString()
		update.Name = &name
	}

	if !state.PortalIcon.IsNull() && !state.PortalIcon.IsUnknown() {
		if update.Settings == nil {
			update.Settings = &cli.OrganizationSettingsUpdate{}
		}
		portalIcon := state.PortalIcon.ValueString()
		update.Settings.PortalIcon = &portalIcon
	}

	return update, nil
}
