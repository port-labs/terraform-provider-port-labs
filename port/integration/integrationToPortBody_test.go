package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func TestIntegrationToPortBodyOauthBrokerUrl(t *testing.T) {
	oauthBrokerUrl := "https://oauth.example.com/broker"
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		OauthBrokerUrl: types.StringValue(oauthBrokerUrl),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body)
	assert.Equal(t, oauthBrokerUrl, *body.OauthBrokerUrl)
}

func TestIntegrationToPortBodyOauthBrokerUrlOmitted(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body)
	assert.Nil(t, body.OauthBrokerUrl)
}
