package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/ratelimit"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
)

type Option func(*PortClient)

type PortClient struct {
	Client                                *resty.Client
	ClientID                              string
	Token                                 string
	featureFlags                          []string
	JSONEscapeHTML                        bool
	BlueprintPropertyTypeChangeProtection bool
}

func isTooManyRequests(r *resty.Response, _ error) bool {
	return r.StatusCode() == 429
}

func New(baseURL string, opts ...Option) (*PortClient, error) {
	ratelimitOpts := &ratelimit.Options{
		Enabled: utils.PtrTo(os.Getenv("PORT_RATE_LIMIT_DISABLED") == ""),
	}
	if isDebug := os.Getenv("PORT_DEBUG_RATE_LIMIT") != ""; isDebug {
		ratelimitOpts.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	rateLimitManager := ratelimit.New(ratelimitOpts)

	c := &PortClient{
		Client: resty.New().
			SetRateLimiter(rateLimitManager).
			OnAfterResponse(rateLimitManager.ResponseMiddleware).
			SetBaseURL(baseURL).
			SetRetryCount(5).
			AddRetryCondition(isTooManyRequests).
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
		return "", fmt.Errorf("failed to connect to Port API for authentication: %w", err)
	}

	// Check for specific authentication errors first
	if resp.StatusCode() == http.StatusUnauthorized {
		return "", fmt.Errorf("authentication failed: invalid client ID or secret (status: 401). Please verify your PORT_CLIENT_ID and PORT_CLIENT_SECRET")
	}

	if resp.StatusCode() == http.StatusForbidden {
		return "", fmt.Errorf("authentication forbidden: access denied (status: 403). Please check your client permissions")
	}

	// General HTTP error check (matches organization.go pattern)
	if resp.IsError() {
		return "", fmt.Errorf("authentication request failed (status: %d), got: %s", resp.StatusCode(), string(resp.Body()))
	}

	var tokenResp AccessTokenResponse
	err = json.Unmarshal(resp.Body(), &tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse authentication response: %w", err)
	}

	if !tokenResp.Ok {
		return "", fmt.Errorf("authentication failed: Port API returned ok=false, got: %s", string(resp.Body()))
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("authentication response contained empty access token, got: %s", string(resp.Body()))
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
