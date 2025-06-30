package cli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRateLimitIntegration(t *testing.T) {
	// Skip if rate limiting is disabled
	if os.Getenv("PORT_RATE_LIMIT_DISABLED") != "" {
		t.Skip("Skipping rate limit test because PORT_RATE_LIMIT_DISABLED is set")
	}

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

	// Create client with rate limiting enabled
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

func TestClientRateLimitNoHeaders(t *testing.T) {
	// Skip if rate limiting is disabled
	if os.Getenv("PORT_RATE_LIMIT_DISABLED") != "" {
		t.Skip("Skipping rate limit test because PORT_RATE_LIMIT_DISABLED is set")
	}

	// Create a test server that doesn't return rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client with rate limiting enabled
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

func TestClientRateLimitDisabled(t *testing.T) {
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

func TestClientRateLimitThrottling(t *testing.T) {
	// Skip if rate limiting is disabled
	if os.Getenv("PORT_RATE_LIMIT_DISABLED") != "" {
		t.Skip("Skipping rate limit test because PORT_RATE_LIMIT_DISABLED is set")
	}

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

	// The second request should have been delayed due to throttling
	assert.Greater(t, elapsed, 10*time.Millisecond, "Request should have been throttled")
}

func TestClientRateLimitSettings(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	// Test enabling/disabling
	client.SetRateLimitEnabled(false)
	// We can't directly test the internal state, but we can verify the methods don't panic

	client.SetRateLimitEnabled(true)
	// We can't directly test the internal state, but we can verify the methods don't panic

	// Test threshold setting
	client.SetRateLimitThreshold(0.25)
	// We can't directly test the internal state, but we can verify the methods don't panic
}

func TestClientRateLimitDisabledViaEnv(t *testing.T) {
	// This test verifies rate limiting is disabled when PORT_RATE_LIMIT_DISABLED is set
	if os.Getenv("PORT_RATE_LIMIT_DISABLED") == "" {
		t.Skip("Skipping test because PORT_RATE_LIMIT_DISABLED is not set")
	}

	// Create a test server that returns rate limit headers that would normally trigger throttling
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1") // Very low remaining
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create client - should have rate limiting disabled due to environment variable
	client, err := New(server.URL)
	require.NoError(t, err)

	// Make a request - should complete quickly without throttling
	start := time.Now()
	resp, err := client.Client.R().Get("/test")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Should complete very quickly since rate limiting is disabled
	assert.Less(t, elapsed, 100*time.Millisecond, "Request should not be throttled when rate limiting is disabled")

	// Rate limit info should be nil since rate limiting is disabled
	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}
