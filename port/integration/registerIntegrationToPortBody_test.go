package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIntegrationToRegisterRequest(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:      types.StringValue("my-integration"),
		InstallationAppType: types.StringValue("kafka"),
		Version:             types.StringValue("1.0.0"),
		Title:               types.StringValue("My Integration"),
		Config: types.StringValue(`{
			"resources": [{
				"kind": "ZOMG"
			}]
		}`),
	}

	request, err := integrationToRegisterRequest(state, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if request.InstallationId != "my-integration" {
		t.Fatalf("expected installationId my-integration, got %q", request.InstallationId)
	}
	if request.InstallationAppType != "kafka" {
		t.Fatalf("expected installationAppType kafka, got %q", request.InstallationAppType)
	}
	if request.Version != "1.0.0" {
		t.Fatalf("expected version 1.0.0, got %q", request.Version)
	}
	if request.Title == nil || *request.Title != "My Integration" {
		t.Fatalf("expected title My Integration, got %v", request.Title)
	}
	if request.ShouldUpdate {
		t.Fatal("expected shouldUpdate to be false")
	}
	if len(request.Config["resources"].([]any)) != 1 {
		t.Fatal("expected one resource mapping in config")
	}
}

func TestIntegrationToRegisterRequestRequiresInstallationAppType(t *testing.T) {
	state := &IntegrationModel{
		InstallationId: types.StringValue("my-integration"),
		Version:        types.StringValue("1.0.0"),
		Config: types.StringValue(`{
			"resources": [{
				"kind": "ZOMG"
			}]
		}`),
	}

	_, err := integrationToRegisterRequest(state, false)
	if err == nil {
		t.Fatal("expected error when installation_app_type is missing")
	}
}

func TestIntegrationToRegisterRequestRequiresVersion(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:      types.StringValue("my-integration"),
		InstallationAppType: types.StringValue("kafka"),
		Config: types.StringValue(`{
			"resources": [{
				"kind": "ZOMG"
			}]
		}`),
	}

	_, err := integrationToRegisterRequest(state, false)
	if err == nil {
		t.Fatal("expected error when version is missing")
	}
}

func TestIntegrationToRegisterRequestRequiresResources(t *testing.T) {
	state := &IntegrationModel{
		InstallationId:      types.StringValue("my-integration"),
		InstallationAppType: types.StringValue("kafka"),
		Version:             types.StringValue("1.0.0"),
		Config:              types.StringValue(`{"deleteDependentEntities": true}`),
	}

	_, err := integrationToRegisterRequest(state, false)
	if err == nil {
		t.Fatal("expected error when config.resources is missing")
	}
}

func TestValidateIntegrationConfigResources(t *testing.T) {
	err := validateIntegrationConfigResources(map[string]any{"resources": []any{}})
	if err == nil {
		t.Fatal("expected error for empty resources array")
	}

	err = validateIntegrationConfigResources(map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing resources")
	}

	err = validateIntegrationConfigResources(map[string]any{"resources": []any{map[string]any{"kind": "test"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
