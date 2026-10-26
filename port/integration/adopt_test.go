package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestAdoptionConflict(t *testing.T) {
	saas, onPrem, gitlab, github := consts.InstallationTypeSaas, consts.InstallationTypeOnPrem, "gitlab", "github-ocean"

	testCases := []struct {
		name          string
		plan          *IntegrationModel
		existing      *cli.Integration
		conflictField string
	}{
		{
			name:     "matching integration can be taken over",
			plan:     &IntegrationModel{InstallationType: types.StringValue(saas), InstallationAppType: types.StringValue(gitlab)},
			existing: &cli.Integration{InstallationType: &saas, InstallationAppType: &gitlab},
		},
		{
			name:     "unset installation type matches an on-prem integration",
			plan:     &IntegrationModel{InstallationType: types.StringNull(), InstallationAppType: types.StringNull()},
			existing: &cli.Integration{InstallationType: &onPrem, InstallationAppType: &gitlab},
		},
		{
			name:          "different installation type",
			plan:          &IntegrationModel{InstallationType: types.StringValue(saas), InstallationAppType: types.StringValue(gitlab)},
			existing:      &cli.Integration{InstallationType: &onPrem, InstallationAppType: &gitlab},
			conflictField: "installation_type",
		},
		{
			name:          "different installation app type",
			plan:          &IntegrationModel{InstallationType: types.StringValue(saas), InstallationAppType: types.StringValue(github)},
			existing:      &cli.Integration{InstallationType: &saas, InstallationAppType: &gitlab},
			conflictField: "installation_app_type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.plan.InstallationId = types.StringValue("my-gitlab")
			conflictErr := adoptionConflict(tc.plan, tc.existing)
			if tc.conflictField == "" {
				assert.Nil(t, conflictErr)
				return
			}
			if assert.NotNil(t, conflictErr) {
				assert.Equal(t, tc.conflictField, conflictErr.Field)
				assert.Equal(t, consts.ResourceKindIntegration, conflictErr.ResourceKind)
				assert.Equal(t, "my-gitlab", conflictErr.Identifier)
			}
		})
	}
}
