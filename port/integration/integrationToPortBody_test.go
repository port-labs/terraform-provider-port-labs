package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallationIdPattern(t *testing.T) {
	valid := []string{
		"myintegration",
		"my-integration",
		"my-integration-123",
		"integration1",
		"123-integration",
		"my-integration-",
		"-my-integration",
		"12345",
		"---",
	}

	for _, id := range valid {
		if !installationIdRegex.MatchString(id) {
			t.Errorf("expected %q to match installation ID pattern", id)
		}
	}

	invalid := []string{
		"my integration with spaces",
		"MyIntegration",
		"my_integration",
		"my-integration!",
		"my@integration",
	}

	for _, id := range invalid {
		if installationIdRegex.MatchString(id) {
			t.Errorf("expected %q not to match installation ID pattern", id)
		}
	}
}

func TestIntegrationToPortBody_CreatePortResourcesOrigin(t *testing.T) {
	model := &IntegrationModel{
		InstallationId:            types.StringValue("github-prod"),
		InstallationType:          types.StringValue(consts.InstallationTypeSaas),
		CreatePortResourcesOrigin: types.StringValue(consts.CreatePortResourcesOriginEmpty),
	}

	createBody, err := integrationToPortBody(model, true)
	require.NoError(t, err)
	require.NotNil(t, createBody.CreatePortResourcesOrigin)
	assert.Equal(t, consts.CreatePortResourcesOriginEmpty, *createBody.CreatePortResourcesOrigin)

	updateBody, err := integrationToPortBody(model, false)
	require.NoError(t, err)
	assert.Nil(t, updateBody.CreatePortResourcesOrigin)
}

func TestIntegrationToPortBody_CreatePortResourcesOriginPort(t *testing.T) {
	model := &IntegrationModel{
		InstallationId:            types.StringValue("github-prod"),
		CreatePortResourcesOrigin: types.StringValue(consts.CreatePortResourcesOriginPort),
	}

	body, err := integrationToPortBody(model, true)
	require.NoError(t, err)
	require.NotNil(t, body.CreatePortResourcesOrigin)
	assert.Equal(t, consts.CreatePortResourcesOriginPort, *body.CreatePortResourcesOrigin)
}

func TestIntegrationToPortBody_CreatePortResourcesOriginOmitted(t *testing.T) {
	model := &IntegrationModel{
		InstallationId:   types.StringValue("github-prod"),
		InstallationType: types.StringValue(consts.InstallationTypeSaas),
	}

	body, err := integrationToPortBody(model, true)
	require.NoError(t, err)
	assert.Nil(t, body.CreatePortResourcesOrigin)
}
