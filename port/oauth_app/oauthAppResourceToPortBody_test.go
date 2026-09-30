package oauth_app

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestValidateRedirectURI(t *testing.T) {
	t.Run("valid redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("https://api.port.io/v1/mcp/oauth2/callback"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("wildcard redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("https://*.example.com/callback"); err == nil {
			t.Fatal("expected error for wildcard redirect uri")
		}
	})

	t.Run("relative redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("/callback"); err == nil {
			t.Fatal("expected error for relative redirect uri")
		}
	})

	t.Run("http localhost redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("http://localhost:3000/callback"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("http non-localhost redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("http://example.com/callback"); err == nil {
			t.Fatal("expected error for non-https non-localhost redirect uri")
		}
	})

	t.Run("fragment redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("https://example.com/callback#token"); err == nil {
			t.Fatal("expected error for redirect uri with fragment")
		}
	})

	t.Run("disallowed characters redirect uri", func(t *testing.T) {
		if err := validateRedirectURI("https://example.com/callback?foo=bar|baz"); err == nil {
			t.Fatal("expected error for redirect uri with disallowed characters")
		}
	})

	t.Run("redirect uri exceeds max length", func(t *testing.T) {
		if err := validateRedirectURI("https://example.com/" + strings.Repeat("a", maxRedirectURILength)); err == nil {
			t.Fatal("expected error for redirect uri exceeding max length")
		}
	})
}

func TestValidateRedirectURIs(t *testing.T) {
	t.Run("valid redirect uris", func(t *testing.T) {
		if err := validateRedirectURIs([]string{
			"https://example.com/callback",
			"http://localhost:3000/callback",
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty redirect uris", func(t *testing.T) {
		if err := validateRedirectURIs([]string{}); err == nil {
			t.Fatal("expected error for empty redirect uris")
		}
	})

	t.Run("too many redirect uris", func(t *testing.T) {
		if err := validateRedirectURIs([]string{
			"https://example.com/callback-1",
			"https://example.com/callback-2",
			"https://example.com/callback-3",
			"https://example.com/callback-4",
			"https://example.com/callback-5",
			"https://example.com/callback-6",
		}); err == nil {
			t.Fatal("expected error for too many redirect uris")
		}
	})

	t.Run("duplicate redirect uris", func(t *testing.T) {
		if err := validateRedirectURIs([]string{
			"https://example.com/callback",
			"https://example.com/callback",
		}); err == nil {
			t.Fatal("expected error for duplicate redirect uris")
		}
	})
}

func TestRedirectURIsEqual(t *testing.T) {
	if !redirectURIsEqual(
		[]string{"https://b.example.com", "https://a.example.com"},
		[]string{"https://a.example.com", "https://b.example.com"},
	) {
		t.Fatal("expected redirect uri sets with different order to be equal")
	}

	if redirectURIsEqual(
		[]string{"https://a.example.com"},
		[]string{"https://b.example.com"},
	) {
		t.Fatal("expected different redirect uri sets to be unequal")
	}
}

func TestOAuthAppResourceToPortBodyCreate(t *testing.T) {
	ctx := context.Background()
	redirectURIs, _ := types.ListValueFrom(ctx, types.StringType, []string{"https://example.com/callback"})

	state := &OAuthAppModel{
		Name:         types.StringValue("My App"),
		RedirectURIs: redirectURIs,
	}

	body, err := oauthAppResourceToPortBodyCreate(ctx, state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Name != "My App" || len(body.RedirectURIs) != 1 || body.RedirectURIs[0] != "https://example.com/callback" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestOAuthAppResourceToPortBodyUpdate(t *testing.T) {
	ctx := context.Background()
	redirectURIs, _ := types.ListValueFrom(ctx, types.StringType, []string{"https://example.com/new-callback"})

	state := &OAuthAppModel{
		Name:         types.StringValue("Updated App"),
		RedirectURIs: redirectURIs,
	}

	body, err := oauthAppResourceToPortBodyUpdate(ctx, state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Name == nil || *body.Name != "Updated App" {
		t.Fatalf("expected name to be set, got %+v", body.Name)
	}
	if len(body.RedirectURIs) != 1 || body.RedirectURIs[0] != "https://example.com/new-callback" {
		t.Fatalf("expected redirect uris to be set, got %+v", body.RedirectURIs)
	}
}

func TestRefreshOAuthAppStateLastLoginAt(t *testing.T) {
	ctx := context.Background()
	lastLogin := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)

	t.Run("sets last_login_at when present", func(t *testing.T) {
		state := &OAuthAppModel{}
		app := &cli.OAuthApp{
			ID:           "app-id",
			Name:         "My App",
			RedirectURIs: []string{"https://example.com/callback"},
			ClientID:     "client-id",
			LastLoginAt:  &lastLogin,
		}

		if err := refreshOAuthAppState(ctx, state, app); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if state.LastLoginAt.IsNull() || state.LastLoginAt.ValueString() != lastLogin.String() {
			t.Fatalf("expected last_login_at %q, got %v", lastLogin.String(), state.LastLoginAt)
		}
	})

	t.Run("null last_login_at when absent", func(t *testing.T) {
		state := &OAuthAppModel{}
		app := &cli.OAuthApp{
			ID:           "app-id",
			Name:         "My App",
			RedirectURIs: []string{"https://example.com/callback"},
			ClientID:     "client-id",
		}

		if err := refreshOAuthAppState(ctx, state, app); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !state.LastLoginAt.IsNull() {
			t.Fatalf("expected null last_login_at, got %v", state.LastLoginAt)
		}
	})
}
