package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Option func(*PortClient)

type PortClient struct {
	Client         *resty.Client
	ClientID       string
	Token          string
	featureFlags   []string
	JSONEscapeHTML bool
}

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

// FeatureFlags Fetches the feature flags from the Organization API. It caches the feature flags locally to reduce call
// count.
func (c *PortClient) FeatureFlags(ctx context.Context) ([]string, error) {
	if c.featureFlags == nil {
		organization, _, err := c.ReadOrganization(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read organization data: %w", err)
		}
		c.featureFlags = organization.FeatureFlags
	}
	return slices.Clone(c.featureFlags), nil
}

func (c *PortClient) HasFeatureFlags(ctx context.Context, flags ...string) (bool, error) {
	orgFlags, err := c.FeatureFlags(ctx)
	if err != nil {
		return false, err
	}
	for _, flag := range flags {
		if !slices.Contains(orgFlags, flag) {
			return false, nil
		}
	}
	return true, nil
}

func (c *PortClient) Authenticate(ctx context.Context, clientID, clientSecret string) (string, error) {
	url := "v1/auth/access_token"
	resp, err := c.Client.R().
		SetBody(map[string]interface{}{
			"clientId":     clientID,
			"clientSecret": clientSecret,
		}).
		SetContext(ctx).
		Post(url)
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
