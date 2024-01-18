package page_permissions

import (
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func pagePermissionsToPortBody(state *PagePermissionsModel) (*cli.PagePermissions, error) {
	if state == nil {
		return nil, nil
	}

	pagePermissions := cli.PagePermissions{
		Read: cli.PageReadPermissions{
			Users: flex.TerraformStringListToGoArray(state.Read.Users),
			Roles: flex.TerraformStringListToGoArray(state.Read.Roles),
			Teams: flex.TerraformStringListToGoArray(state.Read.Teams),
		},
	}
	return &pagePermissions, nil
}
