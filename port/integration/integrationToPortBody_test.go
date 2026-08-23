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
	t.Run("valid oauth broker url", func(t *testing.T) {
		state := &IntegrationModel{
			InstallationId: types.StringValue("my-integration"),
			OauthBrokerUrl: types.StringValue("https://oauth.example.com"),
		}

		integration, err := integrationToPortBody(state)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if integration.OauthBrokerUrl == nil || *integration.OauthBrokerUrl != "https://oauth.example.com" {
			t.Fatalf("expected oauth broker url to be set, got %#v", integration.OauthBrokerUrl)
		}
	})

	t.Run("rejects oauth broker url with query string", func(t *testing.T) {
		state := &IntegrationModel{
			InstallationId: types.StringValue("my-integration"),
			OauthBrokerUrl: types.StringValue("https://oauth.example.com?foo=bar"),
		}

		_, err := integrationToPortBody(state)
		if err == nil {
			t.Fatal("expected error for oauth broker url with query string")
		}
	})
}
