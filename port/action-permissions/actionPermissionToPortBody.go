package action_permissions

import (
	"encoding/json"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
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

	approvePolicyMap := make(map[string]interface{})
	if state.Approve.Policy.ValueString() != "" {
		if err := json.Unmarshal([]byte(state.Approve.Policy.ValueString()), &approvePolicyMap); err != nil {
			return nil, err
		}
	}
	if approvePolicyMap != nil {
		actionPermissions.Approve.Policy = approvePolicyMap
	}

	executePolicyMap := make(map[string]interface{})
	if state.Execute.Policy.ValueString() != "" {
		if err := json.Unmarshal([]byte(state.Execute.Policy.ValueString()), &executePolicyMap); err != nil {
			return nil, err
		}
	}
	return &actionPermissions, nil
}
