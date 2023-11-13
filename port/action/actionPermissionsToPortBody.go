package action

import (
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func policyToPortBody(data *PolicyModel) *cli.Policy {
	if data == nil {
		return nil
	}

	return &cli.Policy{
		Queries:    utils.TerraformMapToGoMap(data.Queries),
		Conditions: flex.TerraformStringListToGoArray(data.Conditions),
	}
}

func actionPermissionsToPortBody(data *PermissionsModel) *cli.ActionPermissions {
	if data == nil {
		return nil
	}

	ap := &cli.ActionPermissions{}

	if data.Execute != nil {
		ap.Execute = &cli.ActionExecutePermissions{
			Users:       flex.TerraformStringListToGoArray(data.Execute.Users),
			Roles:       flex.TerraformStringListToGoArray(data.Execute.Roles),
			Teams:       flex.TerraformStringListToGoArray(data.Execute.Teams),
			OwnedByTeam: data.Execute.OwnedByTeam.ValueBoolPointer(),
			Policy:      policyToPortBody(data.Execute.Policy),
		}
	}

	if data.Approve != nil {
		ap.Approve = &cli.ActionApprovePermissions{
			Users:  flex.TerraformStringListToGoArray(data.Approve.Users),
			Roles:  flex.TerraformStringListToGoArray(data.Approve.Roles),
			Teams:  flex.TerraformStringListToGoArray(data.Approve.Teams),
			Policy: policyToPortBody(data.Approve.Policy),
		}
	}

	return ap
}
