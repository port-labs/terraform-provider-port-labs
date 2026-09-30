package action_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestActionPermissionsToPortBodyStripsMemberWhenOwnedByTeam(t *testing.T) {
	body, err := actionPermissionsToPortBody(&PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(true),
			Roles:       []types.String{types.StringValue("Member"), types.StringValue("Admin")},
		},
		Approve: &ApproveModel{},
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"Admin"}, body.Execute.Roles)
}

func TestActionPermissionsToPortBodyKeepsMemberWhenOwnedByTeamDisabled(t *testing.T) {
	body, err := actionPermissionsToPortBody(&PermissionsModel{
		Execute: &ExecuteModel{
			OwnedByTeam: types.BoolValue(false),
			Roles:       []types.String{types.StringValue("Member")},
		},
		Approve: &ApproveModel{},
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"Member"}, body.Execute.Roles)
}
