package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Option func(*PortClient)

type PortClient struct {
	Client                                *resty.Client
	ClientID                              string
	Token                                 string
	featureFlags                          []string
	JSONEscapeHTML                        bool
	BlueprintPropertyTypeChangeProtection bool

	// Rate limiting fields
	rateLimitInfo      *RateLimitInfo
	rateLimitMutex     sync.RWMutex
	rateLimitEnabled   bool
	rateLimitThreshold float64 // Threshold for when to start throttling (0.1 = 10% remaining)
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
		rateLimitEnabled:   true,
		rateLimitThreshold: 0.1, // Start throttling when 10% of requests remain
	}

	// Add rate limiting middleware
	c.Client.OnAfterResponse(c.rateLimitMiddleware)
	c.Client.OnBeforeRequest(c.preRequestRateLimitCheck)

	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// rateLimitMiddleware extracts rate limit information from response headers
func (c *PortClient) rateLimitMiddleware(client *resty.Client, resp *resty.Response) error {
	if !c.rateLimitEnabled {
		return nil
	}

	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	// Extract rate limit headers
	limitHeader := resp.Header().Get("x-ratelimit-limit")
	periodHeader := resp.Header().Get("x-ratelimit-period")
	remainingHeader := resp.Header().Get("x-ratelimit-remaining")
	resetHeader := resp.Header().Get("x-ratelimit-reset")

	// Only update if we have the headers
	if limitHeader != "" && remainingHeader != "" {
		rateLimitInfo := &RateLimitInfo{}

		if limit, err := strconv.Atoi(limitHeader); err == nil {
			rateLimitInfo.Limit = limit
		}
		if period, err := strconv.Atoi(periodHeader); err == nil {
			rateLimitInfo.Period = period
		}
		if remaining, err := strconv.Atoi(remainingHeader); err == nil {
			rateLimitInfo.Remaining = remaining
		}
		if reset, err := strconv.Atoi(resetHeader); err == nil {
			rateLimitInfo.Reset = reset
		}

		c.rateLimitInfo = rateLimitInfo
	}

	return nil
}

// preRequestRateLimitCheck checks if we should throttle before making a request
func (c *PortClient) preRequestRateLimitCheck(client *resty.Client, req *resty.Request) error {
	if !c.rateLimitEnabled {
		return nil
	}

	c.rateLimitMutex.RLock()
	rateLimitInfo := c.rateLimitInfo
	c.rateLimitMutex.RUnlock()

	if rateLimitInfo != nil && rateLimitInfo.ShouldThrottle(c.rateLimitThreshold) {
		// Calculate delay based on reset time and remaining requests
		var delay time.Duration

		if rateLimitInfo.Remaining > 0 && rateLimitInfo.Reset > 0 {
			// Spread remaining requests evenly over the reset period
			delay = time.Duration(rateLimitInfo.Reset) * time.Second / time.Duration(rateLimitInfo.Remaining)

			// Cap the delay to a reasonable maximum (30 seconds)
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
		} else if rateLimitInfo.Reset > 0 {
			// If no requests remaining, wait until reset
			delay = time.Duration(rateLimitInfo.Reset) * time.Second
		}

		if delay > 0 {
			time.Sleep(delay)
		}
	}

	return nil
}

// GetRateLimitInfo returns the current rate limit information
func (c *PortClient) GetRateLimitInfo() *RateLimitInfo {
	c.rateLimitMutex.RLock()
	defer c.rateLimitMutex.RUnlock()

	if c.rateLimitInfo == nil {
		return nil
	}

	// Return a copy to avoid race conditions
	return &RateLimitInfo{
		Limit:     c.rateLimitInfo.Limit,
		Period:    c.rateLimitInfo.Period,
		Remaining: c.rateLimitInfo.Remaining,
		Reset:     c.rateLimitInfo.Reset,
	}
}

// SetRateLimitEnabled enables or disables rate limiting
func (c *PortClient) SetRateLimitEnabled(enabled bool) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()
	c.rateLimitEnabled = enabled
}

// SetRateLimitThreshold sets the threshold for when to start throttling
func (c *PortClient) SetRateLimitThreshold(threshold float64) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()
	c.rateLimitThreshold = threshold
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
		pc.rateLimitEnabled = false
	}
}

// WithRateLimitThreshold sets the threshold for when to start throttling
// threshold should be between 0.0 and 1.0 (e.g., 0.1 means start throttling when 10% of requests remain)
func WithRateLimitThreshold(threshold float64) Option {
	return func(pc *PortClient) {
		if threshold >= 0.0 && threshold <= 1.0 {
			pc.rateLimitThreshold = threshold
		}
	}
}
