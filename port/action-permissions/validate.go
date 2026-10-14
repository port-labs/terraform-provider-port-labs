package action_permissions

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const memberRoleName = "Member"

var _ resource.ResourceWithValidateConfig = &ActionPermissionsResource{}

func (r *ActionPermissionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state ActionPermissionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validateExecuteOwnedByTeamWithMemberRole(resp, state.Permissions)
}

func validateExecuteOwnedByTeamWithMemberRole(resp *resource.ValidateConfigResponse, permissions *PermissionsModel) {
	if permissions == nil || permissions.Execute == nil {
		return
	}

	execute := permissions.Execute
	if execute.OwnedByTeam.IsNull() || !execute.OwnedByTeam.ValueBool() {
		return
	}

	for i, role := range execute.Roles {
		if role.IsNull() || role.IsUnknown() {
			continue
		}
		if isMemberRole(role.ValueString()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("permissions").AtName("execute").AtName("roles").AtListIndex(i),
				"Invalid execute permission configuration",
				"When permissions.execute.owned_by_team is true, the Member role cannot be granted org-wide execute permission. "+
					"Remove Member from permissions.execute.roles or set owned_by_team to false. "+
					"This matches the Port action permissions UI and the API behavior introduced in owned-by-team action permissions migration.",
			)
			return
		}
	}
}

func isMemberRole(role string) bool {
	return strings.EqualFold(role, memberRoleName)
}

func executeOwnedByTeamEnabled(ownedByTeam types.Bool) bool {
	return !ownedByTeam.IsNull() && ownedByTeam.ValueBool()
}

func withoutMemberRole(roles []string) []string {
	filtered := make([]string, 0, len(roles))
	for _, role := range roles {
		if !isMemberRole(role) {
			filtered = append(filtered, role)
		}
	}
	return filtered
}
