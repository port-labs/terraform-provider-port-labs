package organization

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func organizationResourceToPortBody(ctx context.Context, state *OrganizationModel) (*cli.OrganizationUpdate, error) {
	update := &cli.OrganizationUpdate{}

	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name := state.Name.ValueString()
		update.Name = &name
	}

	if !state.InactivityTimeout.IsUnknown() {
		update.IncludeInactivityTimeout = true
		if state.InactivityTimeout.IsNull() {
			update.InactivityTimeout = nil
		} else {
			update.InactivityTimeout = utils.PtrTo(int(state.InactivityTimeout.ValueInt64()))
		}
	}

	return update, nil
}
