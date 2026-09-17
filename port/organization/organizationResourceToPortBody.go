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

	if !state.IdleTimeMS.IsUnknown() {
		update.IncludeIdleTimeMS = true
		if state.IdleTimeMS.IsNull() {
			update.IdleTimeMS = nil
		} else {
			update.IdleTimeMS = utils.PtrTo(int(state.IdleTimeMS.ValueInt64()))
		}
	}

	return update, nil
}
