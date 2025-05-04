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

	state.Read.Users = utils.Map(utils.SortStringSliceByOther(a.Read.Users, utils.TFStringListToStringArray(oldPermissions.Users)), types.StringValue)
	state.Read.Roles = utils.Map(utils.SortStringSliceByOther(a.Read.Roles, utils.TFStringListToStringArray(oldPermissions.Roles)), types.StringValue)
	state.Read.Teams = utils.Map(utils.SortStringSliceByOther(a.Read.Teams, utils.TFStringListToStringArray(oldPermissions.Teams)), types.StringValue)

	return nil
}
