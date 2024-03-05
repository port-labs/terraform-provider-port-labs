package action_permissions

import (
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func actionPermissionsToPortBody(state *PermissionsModel) (*cli.ActionPermissions, error) {
	if state == nil {
		return nil, nil
	}

	actionPermissions := cli.ActionPermissions{
		Execute: cli.ActionExecutePermissions{
			Users:       flex.TerraformStringListToGoArray(state.Execute.Users),
			Roles:       flex.TerraformStringListToGoArray(state.Execute.Roles),
			Teams:       flex.TerraformStringListToGoArray(state.Execute.Teams),
			OwnedByTeam: state.Execute.OwnedByTeam.ValueBoolPointer(),
		},
		Approve: cli.ActionApprovePermissions{
			Users: flex.TerraformStringListToGoArray(state.Approve.Users),
			Roles: flex.TerraformStringListToGoArray(state.Approve.Roles),
			Teams: flex.TerraformStringListToGoArray(state.Approve.Teams),
		},
	}

	approvePolicyMap, err := utils.TerraformJsonStringToGoObject(state.Approve.Policy.ValueStringPointer())
	if err != nil {
		return nil, err
	}

	executePolicyMap, err := utils.TerraformJsonStringToGoObject(state.Execute.Policy.ValueStringPointer())
	if err != nil {
		return nil, err
	}

	actionPermissions.Approve.Policy = approvePolicyMap
	actionPermissions.Execute.Policy = executePolicyMap

	return &actionPermissions, nil
}
