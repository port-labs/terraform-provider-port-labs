package oauth_app

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func oauthAppResourceToPortBodyCreate(state *OAuthAppModel) (*cli.OAuthAppCreate, error) {
	redirectURI := state.RedirectURI.ValueString()
	if err := validateRedirectURI(redirectURI); err != nil {
		return nil, err
	}

	return &cli.OAuthAppCreate{
		Name:        state.Name.ValueString(),
		RedirectURI: redirectURI,
	}, nil
}

func oauthAppResourceToPortBodyUpdate(state *OAuthAppModel) (*cli.OAuthAppUpdate, error) {
	update := &cli.OAuthAppUpdate{}

	if !state.Name.IsNull() {
		name := state.Name.ValueString()
		update.Name = &name
	}

	if !state.RedirectURI.IsNull() {
		redirectURI := state.RedirectURI.ValueString()
		if err := validateRedirectURI(redirectURI); err != nil {
			return nil, err
		}
		update.RedirectURI = &redirectURI
	}

	return update, nil
}

func refreshOAuthAppState(ctx context.Context, state *OAuthAppModel, app *cli.OAuthApp) error {
	state.ID = types.StringValue(app.ID)
	state.Name = types.StringValue(app.Name)
	state.RedirectURI = types.StringValue(app.RedirectURI)
	state.ClientID = types.StringValue(app.ClientID)

	if app.CreatedAt != nil {
		state.CreatedAt = types.StringValue(app.CreatedAt.String())
	}

	if app.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(app.UpdatedAt.String())
	}

	if app.ClientSecret != "" {
		state.ClientSecret = types.StringValue(app.ClientSecret)
	}

	return nil
}
