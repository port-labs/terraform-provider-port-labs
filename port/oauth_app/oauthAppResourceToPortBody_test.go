package oauth_app

import (
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

func TestOAuthAppResourceToPortBodyCreate(t *testing.T) {
	state := &OAuthAppModel{
		Name:        types.StringValue("My App"),
		RedirectURI: types.StringValue("https://example.com/callback"),
	}

	body, err := oauthAppResourceToPortBodyCreate(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Name != "My App" || body.RedirectURI != "https://example.com/callback" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestOAuthAppResourceToPortBodyUpdate(t *testing.T) {
	state := &OAuthAppModel{
		Name:        types.StringValue("Updated App"),
		RedirectURI: types.StringValue("https://example.com/new-callback"),
	}

	body, err := oauthAppResourceToPortBodyUpdate(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Name == nil || *body.Name != "Updated App" {
		t.Fatalf("expected name to be set, got %+v", body.Name)
	}
	if body.RedirectURI == nil || *body.RedirectURI != "https://example.com/new-callback" {
		t.Fatalf("expected redirect uri to be set, got %+v", body.RedirectURI)
	}
}
