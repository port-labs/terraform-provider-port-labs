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

func TestIntegrationToPortBodyPassesEnableMergeEntityInConfig(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		Config:         types.StringValue(`{"enableMergeEntity":false,"resources":[]}`),
	}

	body, err := integrationToPortBody(state)
	if err != nil {
		t.Fatalf("integrationToPortBody: %v", err)
	}
	if body.Config == nil {
		t.Fatal("expected config to be set")
	}
	if (*body.Config)["enableMergeEntity"] != false {
		t.Fatalf("enableMergeEntity = %v, want false", (*body.Config)["enableMergeEntity"])
	}
}
