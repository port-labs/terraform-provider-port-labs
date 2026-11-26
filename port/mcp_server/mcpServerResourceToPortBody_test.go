package mcp_server

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestMcpServerResourceToPortBodyWithClientCredentialsGrantType(t *testing.T) {
	state := &McpServerModel{
		Title:      types.StringValue("GitHub MCP"),
		URL:        types.StringValue("https://api.githubcopilot.com/mcp/"),
		AuthMethod: types.StringValue("oauth"),
		OAuthConfig: &OAuthConfigModel{
			ClientID:     types.StringValue("client-id"),
			ClientSecret: types.StringValue("client-secret"),
			Scope:        types.StringValue("repo"),
			GrantType:    types.StringValue(GrantTypeClientCredentials),
		},
	}

	entity, err := mcpServerResourceToPortBody(state)
	require.NoError(t, err)
	require.Equal(t, McpServerBlueprintIdentifier, entity.Blueprint)
	require.Equal(t, "GitHub MCP", entity.Title)

	oauthConfig, ok := entity.Properties["oauth_config"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "client-id", oauthConfig["client_id"])
	require.Equal(t, "client-secret", oauthConfig["client_secret"])
	require.Equal(t, "repo", oauthConfig["scope"])
	require.Equal(t, GrantTypeClientCredentials, oauthConfig["grant_type"])
}

func TestMcpServerResourceToPortBodyDefaultsAuthorizationCodeGrantType(t *testing.T) {
	state := &McpServerModel{
		Title: types.StringValue("Example MCP"),
		URL:   types.StringValue("https://mcp.example.com/v1"),
		OAuthConfig: &OAuthConfigModel{
			ClientID: types.StringValue("client-id"),
		},
	}

	entity, err := mcpServerResourceToPortBody(state)
	require.NoError(t, err)

	oauthConfig, ok := entity.Properties["oauth_config"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, GrantTypeAuthorizationCode, oauthConfig["grant_type"])
}

func TestRefreshMcpServerStateReadsGrantType(t *testing.T) {
	now := time.Now()
	entity := &cli.Entity{
		Identifier: "my-mcp-server",
		Title:      "My MCP Server",
		Blueprint:  McpServerBlueprintIdentifier,
		Properties: map[string]interface{}{
			"url":         "https://mcp.example.com/v1",
			"auth_method": "oauth",
			"oauth_config": map[string]interface{}{
				"client_id":     "client-id",
				"scope":         "read",
				"grant_type":    GrantTypeClientCredentials,
			},
		},
	}
	entity.CreatedAt = &now
	entity.UpdatedAt = &now

	state := &McpServerModel{}
	err := refreshMcpServerState(context.Background(), state, entity)
	require.NoError(t, err)
	require.Equal(t, "my-mcp-server", state.Identifier.ValueString())
	require.Equal(t, "https://mcp.example.com/v1", state.URL.ValueString())
	require.NotNil(t, state.OAuthConfig)
	require.Equal(t, GrantTypeClientCredentials, state.OAuthConfig.GrantType.ValueString())
}

func TestRefreshMcpServerStateDefaultsMissingGrantType(t *testing.T) {
	entity := &cli.Entity{
		Identifier: "my-mcp-server",
		Title:      "My MCP Server",
		Properties: map[string]interface{}{
			"url": "https://mcp.example.com/v1",
			"oauth_config": map[string]interface{}{
				"client_id": "client-id",
			},
		},
	}

	state := &McpServerModel{}
	err := refreshMcpServerState(context.Background(), state, entity)
	require.NoError(t, err)
	require.NotNil(t, state.OAuthConfig)
	require.Equal(t, GrantTypeAuthorizationCode, state.OAuthConfig.GrantType.ValueString())
}

func TestRefreshMcpServerStateRetainsClientSecretOnRead(t *testing.T) {
	entity := &cli.Entity{
		Identifier: "my-mcp-server",
		Title:      "My MCP Server",
		Properties: map[string]interface{}{
			"url": "https://mcp.example.com/v1",
			"oauth_config": map[string]interface{}{
				"client_id": "client-id",
				"grant_type": GrantTypeAuthorizationCode,
			},
		},
	}

	state := &McpServerModel{
		OAuthConfig: &OAuthConfigModel{
			ClientSecret: types.StringValue("retained-secret"),
		},
	}

	err := refreshMcpServerState(context.Background(), state, entity)
	require.NoError(t, err)
	require.Equal(t, "retained-secret", state.OAuthConfig.ClientSecret.ValueString())
}
