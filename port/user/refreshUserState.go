package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshUserState(ctx context.Context, state *UserModel, user *cli.User) error {
	state.ID = types.StringValue(user.Email)
	state.Email = types.StringValue(user.Email)
	state.Status = flex.GoStringToFramework(&user.Status)
	state.ManagedByScim = flex.GoBoolToFramework(user.ManagedByScim)
	state.InactivityTimeout = flex.GoInt64ToFramework(user.InactivityTimeout)

	if len(user.Roles) != 0 {
		roleNames := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roleNames[i] = role.Name
		}
		if len(state.Roles) != 0 {
			roleNames = utils.SortStringSliceByOther(roleNames, utils.TFStringListToStringArray(state.Roles))
		}
		state.Roles = utils.Map(roleNames, types.StringValue)
	} else {
		state.Roles = nil
	}

	if len(user.Teams) != 0 {
		teams := user.Teams
		if len(state.Teams) != 0 {
			teams = utils.SortStringSliceByOther(teams, utils.TFStringListToStringArray(state.Teams))
		}
		state.Teams = utils.Map(teams, types.StringValue)
	} else {
		state.Teams = nil
	}

	return nil
}
