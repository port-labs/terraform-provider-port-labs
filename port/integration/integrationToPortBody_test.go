package integration

import (
	"encoding/json"
	"strings"
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

func TestIntegrationToPortBodyOmitsNullInstallationAppType(t *testing.T) {
	body, err := integrationToPortBody(&IntegrationModel{
		InstallationId:      types.StringValue("my-integration"),
		InstallationAppType: types.StringNull(),
	})
	if err != nil {
		t.Fatalf("integrationToPortBody: %v", err)
	}
	if body.InstallationAppType != nil {
		t.Fatalf("installation app type = %v, want nil", body.InstallationAppType)
	}

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if strings.Contains(string(payload), "installationAppType") {
		t.Fatalf("expected null installation_app_type to be omitted from payload, got %s", payload)
	}
}

func TestIntegrationToPortBodyIncludesInstallationAppTypeWhenSet(t *testing.T) {
	body, err := integrationToPortBody(&IntegrationModel{
		InstallationId:      types.StringValue("my-integration"),
		InstallationAppType: types.StringValue("kafka"),
	})
	if err != nil {
		t.Fatalf("integrationToPortBody: %v", err)
	}
	if body.InstallationAppType == nil || *body.InstallationAppType != "kafka" {
		t.Fatalf("installation app type = %v, want kafka", body.InstallationAppType)
	}
}
