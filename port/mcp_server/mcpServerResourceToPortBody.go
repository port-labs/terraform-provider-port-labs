package mcp_server

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func mcpServerResourceToPortBody(state *McpServerModel) (*cli.Entity, error) {
	e := &cli.Entity{
		Blueprint: McpServerBlueprintIdentifier,
		Title:     state.Title.ValueString(),
	}

	if !state.Identifier.IsUnknown() && !state.Identifier.IsNull() {
		e.Identifier = state.Identifier.ValueString()
	}

	if !state.Icon.IsNull() {
		e.Icon = state.Icon.ValueString()
	}

	if state.Teams != nil {
		e.Team = make([]string, len(state.Teams))
		for i, team := range state.Teams {
			e.Team[i] = team.ValueString()
		}
	}

	properties := make(map[string]interface{})

	if !state.URL.IsNull() {
		properties["url"] = state.URL.ValueString()
	}

	if !state.AuthMethod.IsNull() {
		properties["auth_method"] = state.AuthMethod.ValueString()
	}

	if state.OAuthConfig != nil {
		oauthConfig := make(map[string]interface{})

		if !state.OAuthConfig.ClientID.IsNull() {
			oauthConfig["client_id"] = state.OAuthConfig.ClientID.ValueString()
		}

		if !state.OAuthConfig.ClientSecret.IsNull() {
			oauthConfig["client_secret"] = state.OAuthConfig.ClientSecret.ValueString()
		}

		if !state.OAuthConfig.Scope.IsNull() {
			oauthConfig["scope"] = state.OAuthConfig.Scope.ValueString()
		}

		grantType := GrantTypeAuthorizationCode
		if !state.OAuthConfig.GrantType.IsNull() {
			grantType = state.OAuthConfig.GrantType.ValueString()
		}
		oauthConfig["grant_type"] = grantType

		properties["oauth_config"] = oauthConfig
	}

	if state.Headers != nil {
		headers := make(map[string]string, len(state.Headers))
		for key, value := range state.Headers {
			if !value.IsNull() {
				headers[key] = value.ValueString()
			}
		}
		if len(headers) > 0 {
			properties["headers"] = headers
		}
	}

	e.Properties = properties

	return e, nil
}
