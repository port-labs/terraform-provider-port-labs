package mcp_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshMcpServerState(ctx context.Context, state *McpServerModel, entity *cli.Entity) error {
	state.ID = types.StringValue(entity.Identifier)
	state.Identifier = types.StringValue(entity.Identifier)
	state.Title = types.StringValue(entity.Title)

	if entity.Icon != "" {
		state.Icon = types.StringValue(entity.Icon)
	} else {
		state.Icon = types.StringNull()
	}

	if entity.CreatedAt != nil {
		state.CreatedAt = types.StringValue(entity.CreatedAt.String())
	}

	if entity.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(entity.UpdatedAt.String())
	}

	state.CreatedBy = types.StringValue(entity.CreatedBy)
	state.UpdatedBy = types.StringValue(entity.UpdatedBy)

	if entity.Team != nil {
		teams := make([]types.String, len(entity.Team))
		for i, team := range entity.Team {
			teams[i] = types.StringValue(team)
		}
		state.Teams = teams
	} else {
		state.Teams = nil
	}

	if url, ok := entity.Properties["url"].(string); ok {
		state.URL = types.StringValue(url)
	} else {
		state.URL = types.StringNull()
	}

	if authMethod, ok := entity.Properties["auth_method"].(string); ok {
		state.AuthMethod = types.StringValue(authMethod)
	} else {
		state.AuthMethod = types.StringNull()
	}

	existingClientSecret := types.StringNull()
	if state.OAuthConfig != nil && !state.OAuthConfig.ClientSecret.IsNull() {
		existingClientSecret = state.OAuthConfig.ClientSecret
	}

	if oauthConfigRaw, ok := entity.Properties["oauth_config"].(map[string]interface{}); ok {
		oauthConfig := &OAuthConfigModel{}

		if clientID, ok := oauthConfigRaw["client_id"].(string); ok {
			oauthConfig.ClientID = types.StringValue(clientID)
		} else {
			oauthConfig.ClientID = types.StringNull()
		}

		if clientSecret, ok := oauthConfigRaw["client_secret"].(string); ok && clientSecret != "" {
			oauthConfig.ClientSecret = types.StringValue(clientSecret)
		} else if !existingClientSecret.IsNull() && existingClientSecret.ValueString() != "" {
			oauthConfig.ClientSecret = existingClientSecret
		} else {
			oauthConfig.ClientSecret = types.StringNull()
		}

		if scope, ok := oauthConfigRaw["scope"].(string); ok {
			oauthConfig.Scope = types.StringValue(scope)
		} else {
			oauthConfig.Scope = types.StringNull()
		}

		if grantType, ok := oauthConfigRaw["grant_type"].(string); ok && grantType != "" {
			oauthConfig.GrantType = types.StringValue(grantType)
		} else {
			oauthConfig.GrantType = types.StringValue(GrantTypeAuthorizationCode)
		}

		state.OAuthConfig = oauthConfig
	} else {
		state.OAuthConfig = nil
	}

	if headersRaw, ok := entity.Properties["headers"].(map[string]interface{}); ok {
		headers := make(map[string]types.String, len(headersRaw))
		for key, value := range headersRaw {
			if strValue, ok := value.(string); ok {
				headers[key] = types.StringValue(strValue)
			}
		}
		state.Headers = headers
	} else {
		state.Headers = nil
	}

	return nil
}
