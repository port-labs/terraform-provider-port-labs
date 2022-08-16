package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PORT_CLIENT_ID", nil),
			},
			"secret": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PORT_CLIENT_SECRET", nil),
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
				Optional:  true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.getport.io",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"port-labs_entity":    newEntityResource(),
			"port-labs_blueprint": newBlueprintResource(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	secret := d.Get("secret").(string)
	clientID := d.Get("client_id").(string)
	// TODO: verify token or regenerate token
	// token := d.Get("token").(string)
	baseURL := d.Get("base_url").(string)

	c, err := cli.New(baseURL, cli.WithHeader("User-Agent", "terraform-provider-port-labs/0.1"), cli.WithClientID(clientID))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	token, err := c.Authenticate(ctx, clientID, secret)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	d.Set("token", token)
	return c, diags
}
