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

func TestIntegrationToPortBodyIncludesArePortResourcesInitialized(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:              types.StringValue("my-integration"),
		ArePortResourcesInitialized: types.BoolValue(true),
	}

	body, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if body.ArePortResourcesInitialized == nil || !*body.ArePortResourcesInitialized {
		t.Fatalf("expected arePortResourcesInitialized to be true, got %v", body.ArePortResourcesInitialized)
	}
}

func TestIntegrationToPortBodyOmitsArePortResourcesInitializedWhenUnset(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	body, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if body.ArePortResourcesInitialized != nil {
		t.Fatalf("expected arePortResourcesInitialized to be omitted, got %v", body.ArePortResourcesInitialized)
	}
}
