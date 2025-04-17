package action_permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *ActionPermissionsResource) refreshActionPermissionsState(state *ActionPermissionsModel, a *cli.ActionPermissions, actionId string) error {
	state.ID = types.StringValue(actionId)
	state.ActionIdentifier = types.StringValue(actionId)
	state.BlueprintIdentifier = types.StringNull()
	state.Permissions = &PermissionsModel{}

	state.Permissions.Execute = &ExecuteModel{}

	state.Permissions.Execute.Users = make([]types.String, len(a.Execute.Users))
	for i, u := range a.Execute.Users {
		state.Permissions.Execute.Users[i] = types.StringValue(u)
	}

	state.Permissions.Execute.Roles = make([]types.String, len(a.Execute.Roles))
	for i, u := range a.Execute.Roles {
		state.Permissions.Execute.Roles[i] = types.StringValue(u)
	}

	state.Permissions.Execute.Teams = make([]types.String, len(a.Execute.Teams))
	for i, u := range a.Execute.Teams {
		state.Permissions.Execute.Teams[i] = types.StringValue(u)
	}

	state.Permissions.Execute.OwnedByTeam = flex.GoBoolToFramework(a.Execute.OwnedByTeam)

	if a.Execute.Policy != nil {
		policy, err := utils.GoObjectToTerraformString(a.Execute.Policy, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.Permissions.Execute.Policy = policy
	}

	state.Permissions.Approve = &ApproveModel{}

	state.Permissions.Approve.Users = make([]types.String, len(a.Approve.Users))
	for i, u := range a.Approve.Users {
		state.Permissions.Approve.Users[i] = types.StringValue(u)
	}

	state.Permissions.Approve.Roles = make([]types.String, len(a.Approve.Roles))
	for i, u := range a.Approve.Roles {
		state.Permissions.Approve.Roles[i] = types.StringValue(u)
	}

	state.Permissions.Approve.Teams = make([]types.String, len(a.Approve.Teams))
	for i, u := range a.Approve.Teams {
		state.Permissions.Approve.Teams[i] = types.StringValue(u)
	}

	if a.Approve.Policy != nil {
		policy, err := utils.GoObjectToTerraformString(a.Approve.Policy, r.portClient.JSONEscapeHTML)
		if err != nil {
			return err
		}

		state.Permissions.Approve.Policy = policy
	}

	return nil
}
