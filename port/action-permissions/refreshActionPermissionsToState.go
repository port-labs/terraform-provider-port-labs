package action_permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *ActionPermissionsResource) refreshActionPermissionsState(state *ActionPermissionsModel, a *cli.ActionPermissions, actionId string) error {
	oldPermissions := state.Permissions

	state.ID = types.StringValue(actionId)
	state.ActionIdentifier = types.StringValue(actionId)
	state.BlueprintIdentifier = types.StringNull()
	state.Permissions = &PermissionsModel{}

	state.Permissions.Execute = &ExecuteModel{}

	if oldPermissions == nil || oldPermissions.Execute == nil {
		state.Permissions.Execute.Users = utils.Map(a.Execute.Users, types.StringValue)
		state.Permissions.Execute.Roles = utils.Map(a.Execute.Roles, types.StringValue)
		state.Permissions.Execute.Teams = utils.Map(a.Execute.Teams, types.StringValue)
	} else {
		state.Permissions.Execute.Users = utils.Map(utils.SortStringSliceByOther(a.Execute.Users, utils.TFStringListToStringArray(oldPermissions.Execute.Users)), types.StringValue)
		state.Permissions.Execute.Roles = utils.Map(utils.SortStringSliceByOther(a.Execute.Roles, utils.TFStringListToStringArray(oldPermissions.Execute.Roles)), types.StringValue)
		state.Permissions.Execute.Teams = utils.Map(utils.SortStringSliceByOther(a.Execute.Teams, utils.TFStringListToStringArray(oldPermissions.Execute.Teams)), types.StringValue)
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

	if oldPermissions == nil || oldPermissions.Approve == nil {
		state.Permissions.Approve.Users = utils.Map(a.Execute.Users, types.StringValue)
		state.Permissions.Approve.Roles = utils.Map(a.Execute.Roles, types.StringValue)
		state.Permissions.Approve.Teams = utils.Map(a.Execute.Teams, types.StringValue)
	} else {
		state.Permissions.Approve.Users = utils.Map(utils.SortStringSliceByOther(a.Approve.Users, utils.TFStringListToStringArray(oldPermissions.Approve.Users)), types.StringValue)
		state.Permissions.Approve.Roles = utils.Map(utils.SortStringSliceByOther(a.Approve.Roles, utils.TFStringListToStringArray(oldPermissions.Approve.Roles)), types.StringValue)
		state.Permissions.Approve.Teams = utils.Map(utils.SortStringSliceByOther(a.Approve.Teams, utils.TFStringListToStringArray(oldPermissions.Approve.Teams)), types.StringValue)
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
