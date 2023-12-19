package action_permissions

import (
	"encoding/json"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func policyToPortBody(policy *string) (*map[string]any, error) {
	// if policy is empty, set it to nil, so it will override the existing policy on server,
	// as opposed to merging it, due to only having a PATCH endpoint

	if policy == nil || *policy == "" {
		return nil, nil
	}

	policyMap := make(map[string]any)
	if err := json.Unmarshal([]byte(*policy), &policyMap); err != nil {
		return nil, err
	}

	return &policyMap, nil
}

func actionPermissionsToPortBody(state *PermissionsModel) (*cli.ActionPermissions, error) {
	if state == nil {
		return nil, nil
	}

	actionPermissions := cli.ActionPermissions{
		Execute: cli.ActionExecutePermissions{
			Users:       flex.TerraformStringListToGoArray(state.Execute.Users),
			Roles:       flex.TerraformStringListToGoArray(state.Execute.Roles),
			Teams:       flex.TerraformStringListToGoArray(state.Execute.Teams),
			OwnedByTeam: state.Execute.OwnedByTeam.ValueBoolPointer(),
		},
		Approve: cli.ActionApprovePermissions{
			Users: flex.TerraformStringListToGoArray(state.Approve.Users),
			Roles: flex.TerraformStringListToGoArray(state.Approve.Roles),
			Teams: flex.TerraformStringListToGoArray(state.Approve.Teams),
		},
	}

	approvePolicyMap, err := policyToPortBody(state.Approve.Policy.ValueStringPointer())
	if err != nil {
		return nil, err
	}

	executePolicyMap, err := policyToPortBody(state.Execute.Policy.ValueStringPointer())
	if err != nil {
		return nil, err
	}

	actionPermissions.Approve.Policy = approvePolicyMap
	actionPermissions.Execute.Policy = executePolicyMap

	return &actionPermissions, nil
}
