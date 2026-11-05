package oauth_app

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OAuthAppModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	RedirectURIs  types.List   `tfsdk:"redirect_uris"`
	ClientID      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}
