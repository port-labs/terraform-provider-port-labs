package page_permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshPagePermissionsState(state *PagePermissionsModel, a *cli.PagePermissions, pageId string) error {
	state.ID = types.StringValue(pageId)
	state.PageIdentifier = types.StringValue(pageId)
	state.Read = &ReadPagePermissionsModel{}

	state.Read.Users = make([]types.String, len(a.Read.Users))
	for i, u := range a.Read.Users {
		state.Read.Users[i] = types.StringValue(u)
	}

	state.Read.Roles = make([]types.String, len(a.Read.Roles))
	for i, u := range a.Read.Roles {
		state.Read.Roles[i] = types.StringValue(u)
	}

	state.Read.Teams = make([]types.String, len(a.Read.Teams))
	for i, u := range a.Read.Teams {
		state.Read.Teams[i] = types.StringValue(u)
	}

	return nil
}
