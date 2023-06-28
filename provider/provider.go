package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/port/action"
	"github.com/port-labs/terraform-provider-port-labs/port/blueprint"
	"github.com/port-labs/terraform-provider-port-labs/port/entity"
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
	resp.TypeName = "port-labs"
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
	var data *cli.PortProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var baseUrl string

	if data.BaseUrl.ValueString() == "" {
		baseUrl = consts.DefaultBaseUrl
	} else {
		baseUrl = data.BaseUrl.ValueString()
	}

	c, err := cli.New(baseUrl, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Port-labs client", err.Error())
		return
	}

	if data.Token.ValueString() != "" {
		c.Client.SetAuthToken(data.Token.ValueString())
	} else {
		_, err = c.Authenticate(ctx, data.ClientId.ValueString(), data.Secret.ValueString())
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
	}
}

func (p *PortLabsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
