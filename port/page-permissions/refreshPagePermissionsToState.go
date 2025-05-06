package page_permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshPagePermissionsState(state *PagePermissionsModel, a *cli.PagePermissions, pageId string) error {
	oldPermissions := state.Read
	state.ID = types.StringValue(pageId)
	state.PageIdentifier = types.StringValue(pageId)
	state.Read = &ReadPagePermissionsModel{}

	if oldPermissions == nil {
		state.Read.Users = utils.Map(a.Read.Users, types.StringValue)
		state.Read.Roles = utils.Map(a.Read.Roles, types.StringValue)
		state.Read.Teams = utils.Map(a.Read.Teams, types.StringValue)
	} else {
		state.Read.Users = utils.Map(utils.SortStringSliceByOther(a.Read.Users, utils.TFStringListToStringArray(oldPermissions.Users)), types.StringValue)
		state.Read.Roles = utils.Map(utils.SortStringSliceByOther(a.Read.Roles, utils.TFStringListToStringArray(oldPermissions.Roles)), types.StringValue)
		state.Read.Teams = utils.Map(utils.SortStringSliceByOther(a.Read.Teams, utils.TFStringListToStringArray(oldPermissions.Teams)), types.StringValue)
	}

	return nil
}
