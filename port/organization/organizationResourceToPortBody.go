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

	return update, nil
}
