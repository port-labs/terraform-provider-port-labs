package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveIntegrationTypeCopiesDeprecatedValue(t *testing.T) {
	config := &IntegrationModel{
		InstallationAppType: types.StringValue("github-ocean"),
	}
	plan := &IntegrationModel{}

	resolved, err := integrationTypeValue(config, plan)
	require.NoError(t, err)
	assert.Equal(t, "github-ocean", resolved)
}

func TestResolveIntegrationTypePrefersConfiguredIntegrationType(t *testing.T) {
	config := &IntegrationModel{
		IntegrationType: types.StringValue("github-ocean"),
	}
	plan := &IntegrationModel{
		InstallationAppType: types.StringValue("gitlab"),
	}

	resolved, err := integrationTypeValue(config, plan)
	require.NoError(t, err)
	assert.Equal(t, "github-ocean", resolved)
}

func TestResolveIntegrationTypePrefersIntegrationTypeOnConflict(t *testing.T) {
	config := &IntegrationModel{
		IntegrationType:     types.StringValue("github-ocean"),
		InstallationAppType: types.StringValue("gitlab"),
	}

	resolved, err := integrationTypeValue(config)
	require.NoError(t, err)
	assert.Equal(t, "github-ocean", resolved)
}

func TestResolveIntegrationTypeFallsBackToPlan(t *testing.T) {
	plan := &IntegrationModel{
		IntegrationType: types.StringValue("pagerduty"),
	}

	resolved, err := integrationTypeValue(plan)
	require.NoError(t, err)
	assert.Equal(t, "pagerduty", resolved)
}

func TestAlignIntegrationTypeDefersUnknownConfig(t *testing.T) {
	config := &IntegrationModel{IntegrationType: types.StringUnknown()}
	plan := &IntegrationModel{IntegrationType: types.StringUnknown()}

	err := alignIntegrationType(config, plan)
	require.NoError(t, err)
	assert.True(t, plan.IntegrationType.IsUnknown())
}

func TestAlignIntegrationTypePlansUnsetInstallationAppTypeAsNull(t *testing.T) {
	config := &IntegrationModel{
		IntegrationType:     types.StringValue("github-ocean"),
		InstallationAppType: types.StringNull(),
	}
	plan := &IntegrationModel{
		IntegrationType:     types.StringValue("github-ocean"),
		InstallationAppType: types.StringUnknown(),
	}

	err := alignIntegrationType(config, plan)
	require.NoError(t, err)
	assert.True(t, plan.InstallationAppType.IsNull())
}

func TestApplyIntegrationTypePrefersServerIntegrationType(t *testing.T) {
	integrationType := "github-ocean"
	installationAppType := "gitlab"
	model := &IntegrationModel{}

	applyIntegrationType(model, &cli.Integration{
		IntegrationType:     &integrationType,
		InstallationAppType: &installationAppType,
	})

	assert.Equal(t, "github-ocean", model.IntegrationType.ValueString())
}

func TestApplyIntegrationTypeFallsBackToInstallationAppType(t *testing.T) {
	installationAppType := "github-ocean"
	model := &IntegrationModel{}

	applyIntegrationType(model, &cli.Integration{
		InstallationAppType: &installationAppType,
	})

	assert.Equal(t, "github-ocean", model.IntegrationType.ValueString())
}
