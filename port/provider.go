package port

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
	"github.com/port-labs/terraform-provider-port-labs/version"
)

type PortLabsProvider struct {
	clientId string
	secret   string
	token    string
	baseUrl  string
}

func New() provider.Provider {
	return &PortLabsProvider{}
}

func (p *PortLabsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "examplecloud"
}

func (p *PortLabsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Port-labs",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for Port-labs",
				Required:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret for Port-labs",
				Sensitive:           true,
				Required:            true,
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
	p.clientId = os.Getenv("PORT_CLIENT_ID")
	p.secret = os.Getenv("PORT_SECRET")
	p.baseUrl = os.Getenv("PORT_BASE_URL")

	if p.clientId == "" {
		resp.Diagnostics.AddError("Missing PORT_CLIENT_ID", "PORT_CLIENT_ID is required")
		return
	}
	if p.secret == "" {
		resp.Diagnostics.AddError("Missing PORT_SECRET", "PORT_SECRET is required")
		return
	}

	if p.baseUrl == "" {
		p.baseUrl = "https://api.getport.io"
	}

	c, err := cli.New(p.baseUrl, cli.WithHeader("User-Agent", version.ProviderVersion), cli.WithClientID(p.clientId))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Port-labs client", err.Error())
		return
	}

	token, err := c.Authenticate(ctx, p.clientId, p.secret)
	if err != nil {
		resp.Diagnostics.AddError("Failed to authenticate with Port-labs", err.Error())
		return
	}

	p.token = token

	resp.ResourceData = p

}

func (p *PortLabsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *PortLabsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
