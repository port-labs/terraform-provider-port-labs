package oauth_app

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OAuthAppModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	RedirectURI  types.String `tfsdk:"redirect_uri"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	CreatedAt    types.String `tfsdk:"created_at"`
}
