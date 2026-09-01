package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshIntegrationStateOauthBrokerUrl(t *testing.T) {
	oauthBrokerUrl := "https://oauth.example.com"
	resource := &IntegrationResource{}
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	err := resource.refreshIntegrationState(state, &cli.Integration{
		InstallationId: "my-integration",
		OauthBrokerUrl: &oauthBrokerUrl,
	}, "my-integration")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.OauthBrokerUrl.ValueString() != oauthBrokerUrl {
		t.Fatalf("expected oauth_broker_url %q, got %q", oauthBrokerUrl, state.OauthBrokerUrl.ValueString())
	}
}
