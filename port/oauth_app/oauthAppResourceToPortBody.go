package oauth_app

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func redirectURIsFromState(ctx context.Context, redirectURIs types.List) ([]string, error) {
	if redirectURIs.IsNull() || redirectURIs.IsUnknown() {
		return nil, fmt.Errorf("redirect_uris must be set")
	}

	items, err := utils.TerraformListToGoArray(ctx, redirectURIs, "string")
	if err != nil {
		return nil, err
	}

	redirectURIsSlice := make([]string, 0, len(items))
	for _, item := range items {
		redirectURIsSlice = append(redirectURIsSlice, item.(string))
	}

	if err := validateRedirectURIs(redirectURIsSlice); err != nil {
		return nil, err
	}

	return redirectURIsSlice, nil
}

func oauthAppResourceToPortBodyCreate(ctx context.Context, state *OAuthAppModel) (*cli.OAuthAppCreate, error) {
	redirectURIs, err := redirectURIsFromState(ctx, state.RedirectURIs)
	if err != nil {
		return nil, err
	}

	return &cli.OAuthAppCreate{
		Name:         state.Name.ValueString(),
		RedirectURIs: redirectURIs,
	}, nil
}

func oauthAppResourceToPortBodyUpdate(ctx context.Context, state *OAuthAppModel) (*cli.OAuthAppUpdate, error) {
	update := &cli.OAuthAppUpdate{}

	if !state.Name.IsNull() {
		name := state.Name.ValueString()
		update.Name = &name
	}

	if !state.RedirectURIs.IsNull() && !state.RedirectURIs.IsUnknown() {
		redirectURIs, err := redirectURIsFromState(ctx, state.RedirectURIs)
		if err != nil {
			return nil, err
		}
		update.RedirectURIs = redirectURIs
	}

	return update, nil
}

func refreshOAuthAppState(ctx context.Context, state *OAuthAppModel, app *cli.OAuthApp) error {
	state.ID = types.StringValue(app.ID)
	state.Name = types.StringValue(app.Name)
	state.RedirectURIs = flex.GoArrayStringToTerraformList(ctx, app.RedirectURIs)
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
