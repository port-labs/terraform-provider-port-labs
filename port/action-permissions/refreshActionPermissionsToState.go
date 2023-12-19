package action_permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func refreshActionPermissionsState(ctx context.Context, state *ActionPermissionsModel, a *cli.ActionPermissions, blueprintId string, actionId string) error {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprintId, actionId))
	state.ActionIdentifier = types.StringValue(actionId)
	state.BlueprintIdentifier = types.StringValue(blueprintId)
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
		policy, err := json.Marshal(a.Execute.Policy)
		if err != nil {
			return err
		}

		state.Permissions.Execute.Policy = types.StringValue(string(policy))
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
		policy, err := json.Marshal(a.Approve.Policy)
		if err != nil {
			return err
		}

		state.Permissions.Approve.Policy = types.StringValue(string(policy))
	}

	return nil
}
