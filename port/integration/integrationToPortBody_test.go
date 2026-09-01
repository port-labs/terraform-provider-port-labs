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

func TestIntegrationToPortBodyIncludesVersionWhenSet(t *testing.T) {
	version := "1.33.7"
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		Version:        types.StringValue(version),
	}

	body, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("integrationToPortBody returned error: %v", err)
	}
	if body.Version == nil {
		t.Fatal("expected version to be included in integration body")
	}
	if *body.Version != version {
		t.Fatalf("expected version %q, got %q", version, *body.Version)
	}
}

func TestIntegrationToPortBodyOmitsVersionWhenUnset(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
	}

	body, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("integrationToPortBody returned error: %v", err)
	}
	if body.Version != nil {
		t.Fatalf("expected version to be omitted, got %q", *body.Version)
	}
}
