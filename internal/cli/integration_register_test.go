package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterIntegration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/organization":
			_ = json.NewEncoder(w).Encode(PortBody{
				OK: true,
				Organization: &Organization{
					Id:   "org_123",
					Name: "test-org",
				},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/integration/register":
			var body RegisterIntegrationRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode register body: %v", err)
			}

			if body.OrgId != "org_123" {
				t.Fatalf("expected orgId org_123, got %q", body.OrgId)
			}
			if body.InstallationId != "github-prod" {
				t.Fatalf("expected installationId github-prod, got %q", body.InstallationId)
			}
			if body.ShouldUpdate {
				t.Fatalf("expected shouldUpdate false on create")
			}

			_ = json.NewEncoder(w).Encode(RegisterIntegrationResponse{
				OK:      true,
				Created: true,
				Integration: Integration{
					InstallationId:      "github-prod",
					InstallationAppType: strPtr("github-ocean"),
					Version:             strPtr("1.0.0"),
					Title:               strPtr("GitHub Production"),
					Config: &map[string]any{
						"resources": []any{
							map[string]any{"kind": "repository"},
						},
					},
				},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	title := "GitHub Production"
	appType := "github-ocean"
	version := "1.0.0"
	created, err := client.RegisterIntegration(context.Background(), &Integration{
		InstallationId:      "github-prod",
		InstallationAppType: &appType,
		Version:             &version,
		Title:               &title,
		Config: &map[string]any{
			"resources": []any{
				map[string]any{"kind": "repository"},
			},
		},
	}, false)
	if err != nil {
		t.Fatalf("RegisterIntegration returned error: %v", err)
	}
	if created.InstallationId != "github-prod" {
		t.Fatalf("expected installation id github-prod, got %q", created.InstallationId)
	}
}

func TestValidateRegisterIntegrationRequiresResources(t *testing.T) {
	t.Parallel()

	appType := "github-ocean"
	version := "1.0.0"
	err := validateRegisterIntegration(&Integration{
		InstallationId:      "github-prod",
		InstallationAppType: &appType,
		Version:             &version,
		Config:              &map[string]any{},
	})
	if err == nil {
		t.Fatal("expected validation error for missing config.resources")
	}
}

func strPtr(value string) *string {
	return &value
}
