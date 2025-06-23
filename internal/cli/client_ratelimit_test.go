package cli

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimitInfo(t *testing.T) {
	rateLimitInfo := &RateLimitInfo{
		Limit:     100,
		Period:    300,
		Remaining: 10,
		Reset:     60,
	}

	// Test IsNearLimit
	assert.True(t, rateLimitInfo.IsNearLimit(0.2))   // 10/100 = 0.1 < 0.2
	assert.False(t, rateLimitInfo.IsNearLimit(0.05)) // 10/100 = 0.1 > 0.05

	// Test ShouldThrottle
	assert.True(t, rateLimitInfo.ShouldThrottle(0.2))
	assert.False(t, rateLimitInfo.ShouldThrottle(0.05))

	// Test with zero limit
	zeroLimitInfo := &RateLimitInfo{Limit: 0}
	assert.False(t, zeroLimitInfo.IsNearLimit(0.5))
	assert.False(t, zeroLimitInfo.ShouldThrottle(0.5))
}

func TestRateLimitMiddleware(t *testing.T) {
	// Create a test server that returns rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-period", "300")
		w.Header().Set("x-ratelimit-remaining", "50")
		w.Header().Set("x-ratelimit-reset", "120")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client with rate limiting
	client, err := New(server.URL)
	require.NoError(t, err)

	// Make a request
	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info was extracted
	rateLimitInfo := client.GetRateLimitInfo()
	require.NotNil(t, rateLimitInfo)
	assert.Equal(t, 1000, rateLimitInfo.Limit)
	assert.Equal(t, 300, rateLimitInfo.Period)
	assert.Equal(t, 50, rateLimitInfo.Remaining)
	assert.Equal(t, 120, rateLimitInfo.Reset)
}

func TestRateLimitMiddlewareNoHeaders(t *testing.T) {
	// Create a test server that doesn't return rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client with rate limiting
	client, err := New(server.URL)
	require.NoError(t, err)

	// Make a request
	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info is still nil
	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestRateLimitDisabled(t *testing.T) {
	// Create a test server that returns rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client with rate limiting disabled
	client, err := New(server.URL, WithRateLimitDisabled())
	require.NoError(t, err)

	// Make a request
	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info was not extracted
	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestRateLimitThrottling(t *testing.T) {
	requestCount := 0

	// Create a test server that simulates approaching rate limit
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "5") // Low remaining requests
		w.Header().Set("x-ratelimit-reset", "2")     // Reset in 2 seconds
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client with high threshold (should throttle)
	client, err := New(server.URL, WithRateLimitThreshold(0.1)) // Throttle when < 10% remaining
	require.NoError(t, err)

	// Make first request to establish rate limit info
	start := time.Now()
	resp, err := client.Client.R().Get("/test1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Make second request - this should be throttled
	resp, err = client.Client.R().Get("/test2")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 2, requestCount)

	// The second request should have been delayed
	// (5 remaining / 2 seconds reset = 0.4 seconds between requests)
	assert.Greater(t, elapsed, 300*time.Millisecond, "Request should have been throttled")
}

func TestSetRateLimitSettings(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	// Test enabling/disabling
	client.SetRateLimitEnabled(false)
	assert.False(t, client.rateLimitEnabled)

	client.SetRateLimitEnabled(true)
	assert.True(t, client.rateLimitEnabled)

	// Test threshold setting
	client.SetRateLimitThreshold(0.25)
	assert.Equal(t, 0.25, client.rateLimitThreshold)
}
