package user

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func userResourceToPortBody(ctx context.Context, state *UserModel) (*cli.UserUpdate, error) {
	update := &cli.UserUpdate{}

	if state.Roles != nil {
		update.IncludeRoles = true
		update.Roles = utils.TFStringListToStringArray(state.Roles)
	}

	if state.Teams != nil {
		update.IncludeTeams = true
		update.Teams = utils.TFStringListToStringArray(state.Teams)
	}

	if !state.InactivityTimeout.IsUnknown() {
		update.IncludeInactivityTimeout = true
		if state.InactivityTimeout.IsNull() {
			update.InactivityTimeout = nil
		} else {
			update.InactivityTimeout = utils.PtrTo(int(state.InactivityTimeout.ValueInt64()))
		}
	}

	return update, nil
}

func userInviteFromState(state *UserModel) *cli.UserInviteRequest {
	invite := &cli.UserInviteRequest{
		Invitee: cli.UserInvitee{
			Email: state.Email.ValueString(),
		},
	}

	if state.Roles != nil {
		invite.Invitee.Roles = utils.TFStringListToStringArray(state.Roles)
	}
	if state.Teams != nil {
		invite.Invitee.Teams = utils.TFStringListToStringArray(state.Teams)
	}

	return invite
}
