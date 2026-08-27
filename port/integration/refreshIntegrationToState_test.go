package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshIntegrationStateSetsVersionFromAPI(t *testing.T) {
	version := "1.33.7"
	resource := &IntegrationResource{}
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	err := resource.refreshIntegrationState(state, &cli.Integration{
		Title:   stringPtr("My Integration"),
		Version: &version,
	}, "my-integration")
	if err != nil {
		t.Fatalf("refreshIntegrationState returned error: %v", err)
	}
	if state.Version.ValueString() != version {
		t.Fatalf("expected version %q, got %q", version, state.Version.ValueString())
	}
}

func stringPtr(value string) *string {
	return &value
}
