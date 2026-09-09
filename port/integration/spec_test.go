package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationToPortBody_Saas(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:      types.StringValue("pagerduty-prod"),
		InstallationAppType: types.StringValue("pagerduty"),
		InstallationType:    types.StringValue(consts.InstallationTypeSaas),
		Version:             types.StringValue("0.1.0"),
		Spec:                types.StringValue(`{"integrationSpec":{"token":"pagerduty-api-token"},"appSpec":{"scheduledResyncInterval":"12h"}}`),
		Config:              types.StringValue(`{"resources":[]}`),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	assert.Equal(t, "pagerduty-prod", body.InstallationId)
	require.NotNil(t, body.InstallationType)
	assert.Equal(t, consts.InstallationTypeSaas, *body.InstallationType)
	require.NotNil(t, body.Spec)
	assert.Equal(t, "pagerduty-api-token", body.Spec.IntegrationSpec["token"])
	assert.Equal(t, "12h", body.Spec.AppSpec["scheduledResyncInterval"])
}

func TestIntegrationToPortBody_OnPremOmitsSpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:      types.StringValue("my-kafka"),
		InstallationAppType: types.StringValue("kafka"),
		InstallationType:    types.StringValue(consts.InstallationTypeOnPrem),
		Spec:                types.StringValue(`{"integrationSpec":{"token":"should-not-send"}}`),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	assert.Equal(t, consts.InstallationTypeOnPrem, *body.InstallationType)
	assert.Nil(t, body.Spec, "OnPrem integrations should not send spec")
}

func TestValidateSaasSpec_RequiresSpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:   types.StringValue("pagerduty-prod"),
		InstallationType: types.StringValue(consts.InstallationTypeSaas),
	}
	assert.ErrorContains(t, validateSaasSpec(state), "spec is required")
}

func TestValidateSaasSpec_RequiresIntegrationSpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:   types.StringValue("pagerduty-prod"),
		InstallationType: types.StringValue(consts.InstallationTypeSaas),
		Spec:             types.StringValue(`{"appSpec":{"scheduledResyncInterval":"12h"}}`),
	}
	assert.ErrorContains(t, validateSaasSpec(state), "integrationSpec is required")
}

func TestValidateSaasSpec_OnPremPassesWithoutSpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:   types.StringValue("my-kafka"),
		InstallationType: types.StringValue(consts.InstallationTypeOnPrem),
	}
	assert.NoError(t, validateSaasSpec(state))
}

func TestValidateIntegrationModel_OnPremRejectsSpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:   types.StringValue("my-kafka"),
		InstallationType: types.StringValue(consts.InstallationTypeOnPrem),
		Spec:             types.StringValue(`{"integrationSpec":{"token":"x"}}`),
	}
	assert.ErrorContains(t, validateIntegrationModel(state), "spec is only supported")
}

func TestValidateIntegrationModel_OnPremAllowsEmptySpec(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:   types.StringValue("my-kafka"),
		InstallationType: types.StringValue(consts.InstallationTypeOnPrem),
		Spec:             types.StringNull(),
	}
	assert.NoError(t, validateIntegrationModel(state))
}

func TestParseSpecFromConfig_OmitsUnknownSpecSections(t *testing.T) {
	spec, err := parseSpecFromConfig(types.StringValue(
		`{"integrationSpec":{"token":"x"},"appSpec":{"liveEventsEnabled":true},"systemSpec":{"size":"M"}}`,
	))
	require.NoError(t, err)
	require.NotNil(t, spec)
	assert.Equal(t, "x", spec.IntegrationSpec["token"])
	assert.Equal(t, true, spec.AppSpec["liveEventsEnabled"])
}
