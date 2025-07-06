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

func skipIfDisabled(t *testing.T) {
	if os.Getenv("PORT_RATE_LIMIT_DISABLED") != "" {
		t.Skip("Skipping rate limit test because PORT_RATE_LIMIT_DISABLED is set")
	}
}

func TestClientRateLimitIntegration(t *testing.T) {
	skipIfDisabled(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-period", "300")
		w.Header().Set("x-ratelimit-remaining", "50")
		w.Header().Set("x-ratelimit-reset", "120")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client, err := New(server.URL)
	require.NoError(t, err)

	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := client.GetRateLimitInfo()
	require.NotNil(t, rateLimitInfo)
	assert.Equal(t, 1000, rateLimitInfo.Limit)
	assert.Equal(t, 300, rateLimitInfo.Period)
	assert.Equal(t, 50, rateLimitInfo.Remaining)
	assert.Equal(t, 120, rateLimitInfo.Reset)
}

func TestClientRateLimitNoHeaders(t *testing.T) {
	skipIfDisabled(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client, err := New(server.URL)
	require.NoError(t, err)

	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestClientRateLimitDisabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client, err := New(server.URL, WithRateLimitDisabled())
	require.NoError(t, err)

	resp, err := client.Client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestClientRateLimitThrottling(t *testing.T) {
	skipIfDisabled(t)

	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "5")
		w.Header().Set("x-ratelimit-reset", "2")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client, err := New(server.URL, WithRateLimitThreshold(0.1))
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Client.R().Get("/test1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	resp, err = client.Client.R().Get("/test2")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 2, requestCount)

	assert.Greater(t, elapsed, 10*time.Millisecond, "Request should have been throttled")
}

func TestClientRateLimitSettings(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	// somewhat dummy tests since we really can't test, so we check that they don't panic
	client.SetRateLimitEnabled(false)
	client.SetRateLimitEnabled(true)
	client.SetRateLimitThreshold(0.25)
}

func TestClientRateLimitDisabledViaEnv(t *testing.T) {
	skipIfDisabled(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client, err := New(server.URL)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Client.R().Get("/test")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	assert.Less(t, elapsed, 100*time.Millisecond, "Request should not be throttled when rate limiting is disabled")

	rateLimitInfo := client.GetRateLimitInfo()
	assert.Nil(t, rateLimitInfo)
}
