package cli

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-resty/resty/v2"
)

type (
	Option     func(*PortClient)
	PortClient struct {
		Client   *resty.Client
		ClientID string
	}
)

func New(baseURL string, opts ...Option) (*PortClient, error) {
	c := &PortClient{
		Client: resty.New().
			SetBaseURL(baseURL).
			SetRetryCount(5).
			SetRetryWaitTime(300).
			// retry when create permission fails because scopes are created async-ly and sometimes (mainly in tests) the scope doesn't exist yet.
			AddRetryCondition(func(r *resty.Response, err error) bool {
				if err != nil {
					return true
				}
				if !strings.Contains(r.Request.URL, "/permissions") {
					return false
				}
				b := make(map[string]interface{})
				err = json.Unmarshal(r.Body(), &b)
				return err != nil || b["ok"] != true
			}),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c *PortClient) Authenticate(ctx context.Context, clientID, clientSecret string) (string, error) {
	url := "v0.1/auth/access_token"
	resp, err := c.Client.R().
		SetContext(ctx).
		SetQueryParam("client_id", clientID).
		SetQueryParam("client_secret", clientSecret).
		Get(url)
	if err != nil {
		return "", err
	}
	var tokenResp AccessTokenResponse
	err = json.Unmarshal(resp.Body(), &tokenResp)
	if err != nil {
		return "", err
	}
	c.Client.SetAuthToken(tokenResp.AccessToken)
	return tokenResp.AccessToken, nil
}

func WithHeader(key, val string) Option {
	return func(pc *PortClient) {
		pc.Client.SetHeader(key, val)
	}
}

func WithClientID(clientID string) Option {
	return func(pc *PortClient) {
		pc.ClientID = clientID
	}
}

func WithToken(token string) Option {
	return func(pc *PortClient) {
		pc.Client.SetAuthToken(token)
	}
}
