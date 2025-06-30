package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/ratelimit"
)

type Option func(*PortClient)

type PortClient struct {
	Client                                *resty.Client
	ClientID                              string
	Token                                 string
	featureFlags                          []string
	JSONEscapeHTML                        bool
	BlueprintPropertyTypeChangeProtection bool

	// Rate limiting
	rateLimitManager *ratelimit.Manager
}

func New(baseURL string, opts ...Option) (*PortClient, error) {
	rateLimitManager := ratelimit.NewManager()

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
		rateLimitManager: rateLimitManager,
	}

	c.Client.
		OnBeforeRequest(rateLimitManager.RequestMiddleware).
		OnAfterResponse(rateLimitManager.ResponseMiddleware)

	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// GetRateLimitInfo returns the current rate limit information
func (c *PortClient) GetRateLimitInfo() *ratelimit.RateLimitInfo {
	return c.rateLimitManager.GetInfo()
}

// SetRateLimitEnabled enables or disables rate limiting
func (c *PortClient) SetRateLimitEnabled(enabled bool) {
	c.rateLimitManager.SetEnabled(enabled)
}

// SetRateLimitThreshold sets the threshold for when to start throttling
func (c *PortClient) SetRateLimitThreshold(threshold float64) {
	c.rateLimitManager.SetThreshold(threshold)
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

// WithRateLimitDisabled disables rate limiting
func WithRateLimitDisabled() Option {
	return func(pc *PortClient) {
		pc.rateLimitManager.SetEnabled(false)
	}
}

// WithRateLimitThreshold sets the threshold for when to start throttling
// threshold should be between 0.0 and 1.0 (e.g., 0.1 means start throttling when 10% of requests remain)
func WithRateLimitThreshold(threshold float64) Option {
	return func(pc *PortClient) {
		if threshold >= 0.0 && threshold <= 1.0 {
			pc.rateLimitManager.SetThreshold(threshold)
		}
	}
}
