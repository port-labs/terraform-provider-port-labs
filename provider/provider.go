package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/port/action"
	"github.com/port-labs/terraform-provider-port-labs/port/action-permissions"
	"github.com/port-labs/terraform-provider-port-labs/port/blueprint"
	"github.com/port-labs/terraform-provider-port-labs/port/entity"
	"github.com/port-labs/terraform-provider-port-labs/port/scorecard"
	"github.com/port-labs/terraform-provider-port-labs/port/team"
	"github.com/port-labs/terraform-provider-port-labs/port/webhook"
	"github.com/port-labs/terraform-provider-port-labs/version"
)

var (
	_ provider.Provider = &PortLabsProvider{}
)

type PortLabsProvider struct{}

func New() provider.Provider {
	return &PortLabsProvider{}
}

func (p *PortLabsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = consts.ProviderName
}

func (p *PortLabsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Port-labs",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for Port-labs",
				Optional:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret for Port-labs",
				Sensitive:           true,
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token for Port-labs",
				Sensitive:           true,
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *PortLabsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data *cli.PortProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var baseUrl string

	if data.BaseUrl.IsNull() {
		baseUrl = os.Getenv("PORT_BASE_URL")
	} else {
		baseUrl = data.BaseUrl.ValueString()
	}

	if baseUrl == "" {
		baseUrl = consts.DefaultBaseUrl
	}

	c, err := cli.New(baseUrl, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Port-labs client", err.Error())
		return
	}

	if data.Token.ValueString() != "" {
		c.Client.SetAuthToken(data.Token.ValueString())
	} else {
		var clientID, secret string
		if data.ClientId.IsNull() {
			clientID = os.Getenv("PORT_CLIENT_ID")
		} else {
			clientID = data.ClientId.ValueString()
		}
		if clientID == "" {
			resp.Diagnostics.AddError("Unable to find client ID",
				"Client ID is required, either set in config or environment variable PORT_CLIENT_ID")
			return
		}

		if data.Secret.IsNull() {
			secret = os.Getenv("PORT_CLIENT_SECRET")
		} else {
			secret = data.Secret.ValueString()
		}
		if secret == "" {
			resp.Diagnostics.AddError("Unable to find client secret",
				"Client secret is required, either set in config or environment variable PORT_CLIENT_SECRET")
			return
		}

		_, err = c.Authenticate(ctx, clientID, secret)
		if err != nil {
			resp.Diagnostics.AddError("Failed to authenticate with Port-labs", err.Error())
			return
		}
	}

	resp.ResourceData = c

}

func (p *PortLabsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		blueprint.NewBlueprintResource,
		entity.NewEntityResource,
		action.NewActionResource,
		action_permissions.NewActionPermissionsResource,
		webhook.NewWebhookResource,
		scorecard.NewScorecardResource,
		team.NewTeamResource,
	}
}

func (p *PortLabsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
