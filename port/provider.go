package port

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"port_entity": newEntityResource(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

type AccessTokenResponse struct {
	Ok          bool   `json:"ok"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	secret := d.Get("secret").(string)
	clientID := d.Get("client_id").(string)
	token := d.Get("token").(string)
	baseURL := d.Get("base_url").(string)

	client := resty.New()
	client.SetBaseURL(baseURL)
	if token == "" {
		url := "v0.1/auth/access_token"
		resp, err := client.R().
			SetQueryParam("client_id", clientID).
			SetQueryParam("client_secret", secret).
			Get(url)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		var tokenResp AccessTokenResponse
		err = json.Unmarshal(resp.Body(), &tokenResp)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		token = tokenResp.AccessToken
		d.Set("token", tokenResp.AccessToken)
	}
	// else {
	// 	// TODO: verify token or regenerate
	// }
	client.SetAuthToken(token)
	return client, diags
}
