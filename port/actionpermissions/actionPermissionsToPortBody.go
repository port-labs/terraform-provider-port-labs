package actionpermissions

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func actionPermissionsResourceToBody(ctx context.Context, state *ActionPermissionsModel, a *cli.Action) (*cli.ActionPermissions, error) {
	ap := &cli.ActionPermissions{
		Action: a.Identifier,
	}

	execute := &cli.ActionExecutePermissions{
		Users:       flex.TerraformStringListToGoArray(state.Execute.Users),
		Roles:       flex.TerraformStringListToGoArray(state.Execute.Roles),
		Teams:       flex.TerraformStringListToGoArray(state.Execute.Teams),
		OwnedByTeam: true, // state.Execute.OwnedByTeam.ToTerraformValue(ctx),
	}

	ap.Execute = *execute

	return ap, nil
}
