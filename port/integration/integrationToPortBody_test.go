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

func TestOauthBrokerUrlValidation(t *testing.T) {
	validURL := "https://oauth.example.com/broker"
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		OauthBrokerUrl: types.StringValue(validURL),
	}

	integration, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("expected valid oauth_broker_url to succeed, got: %v", err)
	}
	if integration.OauthBrokerUrl == nil || *integration.OauthBrokerUrl != validURL {
		t.Fatalf("expected oauthBrokerUrl %q, got %#v", validURL, integration.OauthBrokerUrl)
	}
}

func TestOauthBrokerUrlRejectsQueryString(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		OauthBrokerUrl: types.StringValue("https://oauth.example.com/broker?foo=bar"),
	}

	_, err := integrationToPortBody(state)
	if err == nil {
		t.Fatal("expected oauth_broker_url with query string to fail")
	}
}

func TestOauthBrokerUrlRejectsInvalidURL(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		OauthBrokerUrl: types.StringValue("not-a-url"),
	}

	_, err := integrationToPortBody(state)
	if err == nil {
		t.Fatal("expected invalid oauth_broker_url to fail")
	}
}
