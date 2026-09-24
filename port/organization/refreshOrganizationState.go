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
	state.IdleTimeMS = flex.GoInt64ToFramework(org.IdleTimeMS)

	return nil
}
