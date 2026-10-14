package action_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateExecuteOwnedByTeamRejectsMemberRole(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	validateExecuteOwnedByTeamWithMemberRole(resp, &PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(true),
			Roles:       []types.String{types.StringValue("Member"), types.StringValue("Admin")},
		},
	})
	assert.True(t, resp.Diagnostics.HasError())
}

func TestValidateExecuteOwnedByTeamAllowsMemberWhenDisabled(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	validateExecuteOwnedByTeamWithMemberRole(resp, &PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(false),
			Roles:       []types.String{types.StringValue("Member")},
		},
	})
	assert.False(t, resp.Diagnostics.HasError())
}

func TestValidateExecuteOwnedByTeamAllowsAdminWithOwnedByTeam(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	validateExecuteOwnedByTeamWithMemberRole(resp, &PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(true),
			Roles:       []types.String{types.StringValue("Admin")},
		},
	})
	assert.False(t, resp.Diagnostics.HasError())
}

func TestWithoutMemberRoleIsCaseInsensitive(t *testing.T) {
	assert.Equal(t, []string{"Admin"}, withoutMemberRole([]string{"member", "Admin"}))
}

func TestValidateExecuteOwnedByTeamRejectsLowercaseMemberRole(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	validateExecuteOwnedByTeamWithMemberRole(resp, &PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(true),
			Roles:       []types.String{types.StringValue("member")},
		},
	})
	assert.True(t, resp.Diagnostics.HasError())
}
