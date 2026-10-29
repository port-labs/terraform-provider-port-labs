package adoption_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/adoption"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestTakeOver(t *testing.T) {
	existing := &cli.ResourceExistsError{ResourceKind: consts.ResourceKindBlueprintRelation, Identifier: "gitlab_group", Blueprint: "_team"}
	conflict := &cli.AdoptionConflictError{ResourceKind: consts.ResourceKindBlueprintRelation, Identifier: "gitlab_group", Field: "target", Existing: "a", Configured: "b"}

	testCases := []struct {
		name          string
		phase         adoption.Phase
		policy        cli.AdoptionPolicy
		adoptExisting types.Bool
		conflict      *cli.AdoptionConflictError
		want          bool
		wantError     string
		wantErrorPath path.Path
		wantWarning   string
	}{
		{name: "not allowed fails on the identifier", phase: adoption.AtPlan, adoptExisting: types.BoolNull(), wantError: "relation already exists", wantErrorPath: path.Root("identifier")},
		{name: "allowed by resource warns at plan", phase: adoption.AtPlan, adoptExisting: types.BoolValue(true), want: true, wantWarning: "relation will be taken over"},
		{name: "allowed by policy warns at apply", phase: adoption.AtApply, policy: cli.NewAdoptionPolicy(consts.ResourceKindBlueprintRelation), adoptExisting: types.BoolNull(), want: true, wantWarning: "relation was taken over"},
		{name: "resource opt-out wins over policy", phase: adoption.AtApply, policy: cli.NewAdoptionPolicy(consts.ResourceKindBlueprintRelation), adoptExisting: types.BoolValue(false), wantError: "relation already exists", wantErrorPath: path.Root("identifier")},
		{name: "conflict fails on its field even when allowed", phase: adoption.AtPlan, adoptExisting: types.BoolValue(true), conflict: conflict, wantError: "cannot take over existing relation", wantErrorPath: path.Root("target")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var diags diag.Diagnostics
			got := adoption.TakeOver(tc.phase, &cli.PortClient{AdoptExistingResources: tc.policy}, existing, tc.adoptExisting, tc.conflict, &diags)

			assert.Equal(t, tc.want, got)
			if tc.wantError != "" {
				if assert.Len(t, diags.Errors(), 1) {
					assert.Equal(t, tc.wantError, diags.Errors()[0].Summary())
					assert.Equal(t, tc.wantErrorPath, diags.Errors()[0].(diag.DiagnosticWithPath).Path())
				}
				assert.Empty(t, diags.Warnings(), "a refused take-over is not also reported as taken over")
				return
			}
			assert.Empty(t, diags.Errors())
			if assert.Len(t, diags.Warnings(), 1) {
				assert.Equal(t, tc.wantWarning, diags.Warnings()[0].Summary())
				assert.Contains(t, diags.Warnings()[0].Detail(), `on blueprint "_team"`)
			}
		})
	}
}

func TestTakeOverIntegrationFailsOnInstallationID(t *testing.T) {
	var diags diag.Diagnostics
	existing := &cli.ResourceExistsError{ResourceKind: consts.ResourceKindIntegration, Identifier: "my-gitlab"}

	adoption.TakeOver(adoption.AtPlan, &cli.PortClient{}, existing, types.BoolNull(), nil, &diags)

	if assert.Len(t, diags.Errors(), 1) {
		assert.Equal(t, path.Root("installation_id"), diags.Errors()[0].(diag.DiagnosticWithPath).Path())
	}
}
