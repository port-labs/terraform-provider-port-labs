package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
	oauthBrokerUrl := "https://oauth.example.com"
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		OauthBrokerUrl: types.StringValue(oauthBrokerUrl),
	}

	integration, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if integration.OauthBrokerUrl == nil || *integration.OauthBrokerUrl != oauthBrokerUrl {
		t.Fatalf("expected oauthBrokerUrl %q, got %v", oauthBrokerUrl, integration.OauthBrokerUrl)
	}
}

func TestIntegrationToPortBodyOmitsOauthBrokerUrlWhenUnset(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	integration, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if integration.OauthBrokerUrl != nil {
		t.Fatalf("expected oauthBrokerUrl to be omitted, got %v", integration.OauthBrokerUrl)
	}
}
