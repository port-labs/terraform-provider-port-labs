package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIntegrationToPortBodyOauthBrokerUrl(t *testing.T) {
	t.Run("valid oauth broker url", func(t *testing.T) {
		state := &IntegrationModel{
			InstallationId: types.StringValue("my-integration"),
			OauthBrokerUrl: types.StringValue("https://oauth.example.com/broker"),
		}

		integration, err := integrationToPortBody(state)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if integration.OauthBrokerUrl == nil || *integration.OauthBrokerUrl != "https://oauth.example.com/broker" {
			t.Fatalf("expected oauth broker url to be set, got %v", integration.OauthBrokerUrl)
		}
	})

	t.Run("oauth broker url with query string", func(t *testing.T) {
		state := &IntegrationModel{
			InstallationId: types.StringValue("my-integration"),
			OauthBrokerUrl: types.StringValue("https://oauth.example.com/broker?foo=bar"),
		}

		_, err := integrationToPortBody(state)
		if err == nil {
			t.Fatal("expected error for oauth broker url with query string")
		}
		if err.Error() != "oauth_broker_url must not contain a query string — Port appends query parameters at runtime" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("null oauth broker url omitted", func(t *testing.T) {
		state := &IntegrationModel{
			InstallationId: types.StringValue("my-integration"),
			OauthBrokerUrl: types.StringNull(),
		}

		integration, err := integrationToPortBody(state)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if integration.OauthBrokerUrl != nil {
			t.Fatalf("expected oauth broker url to be omitted, got %v", integration.OauthBrokerUrl)
		}
	})
}
