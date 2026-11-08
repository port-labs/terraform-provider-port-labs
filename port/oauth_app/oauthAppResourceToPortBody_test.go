package oauth_app

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
