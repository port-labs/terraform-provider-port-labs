package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
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

func TestManagerBasicFunctionality(t *testing.T) {
	manager := NewManager()

	// Test initial state
	assert.Nil(t, manager.GetInfo())

	// Test enabling/disabling
	manager.SetEnabled(false)
	assert.False(t, manager.enabled)

	manager.SetEnabled(true)
	assert.True(t, manager.enabled)

	// Test threshold setting
	manager.SetThreshold(0.25)
	assert.Equal(t, 0.25, manager.threshold)

	// Test invalid threshold (should not change)
	manager.SetThreshold(-0.1)
	assert.Equal(t, 0.25, manager.threshold)

	manager.SetThreshold(1.5)
	assert.Equal(t, 0.25, manager.threshold)
}

func TestResponseMiddleware(t *testing.T) {
	manager := NewManager()

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

	// Create resty client with our middleware
	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// Make a request
	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info was extracted
	rateLimitInfo := manager.GetInfo()
	require.NotNil(t, rateLimitInfo)
	assert.Equal(t, 1000, rateLimitInfo.Limit)
	assert.Equal(t, 300, rateLimitInfo.Period)
	assert.Equal(t, 50, rateLimitInfo.Remaining)
	assert.Equal(t, 120, rateLimitInfo.Reset)
}

func TestResponseMiddlewareNoHeaders(t *testing.T) {
	manager := NewManager()

	// Create a test server that doesn't return rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create resty client with our middleware
	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// Make a request
	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info is still nil
	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestMiddlewareDisabled(t *testing.T) {
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "123")
	manager := NewManager()

	// Create a test server that returns rate limit headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create resty client with our middleware
	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// Make a request
	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check that rate limit info was not extracted because middleware is disabled
	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestThrottling(t *testing.T) {
	manager := NewManager() // High threshold to trigger throttling
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

	// Create resty client with our middleware
	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// Make first request to establish rate limit info
	start := time.Now()
	resp, err := client.R().Get("/test1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Make second request - this should be throttled
	resp, err = client.R().Get("/test2")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 2, requestCount)

	// The second request should have been delayed due to throttling
	// With 5 remaining requests and 2 second reset, we expect some delay
	assert.Greater(t, elapsed, 10*time.Millisecond, "Request should have been throttled")
}

func TestConcurrentRequests(t *testing.T) {
	manager := NewManager()
	requestCount := 0

	// Create a test server that tracks concurrent requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate rate limit headers with very few remaining
		w.Header().Set("x-ratelimit-limit", "10")
		w.Header().Set("x-ratelimit-remaining", "2")
		w.Header().Set("x-ratelimit-reset", "5")

		requestCount++
		// Add a small delay to simulate server processing time
		time.Sleep(50 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create resty client with our middleware
	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// Make multiple concurrent requests
	start := time.Now()
	numRequests := 3
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			_, err := client.R().Get("/test")
			results <- err
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err)
	}

	elapsed := time.Since(start)

	// All requests should complete successfully
	assert.Equal(t, numRequests, requestCount)

	// Due to rate limiting, concurrent requests should be throttled
	// The semaphore and throttling should ensure we don't overwhelm the API
	t.Logf("Concurrent requests took %v", elapsed)
}

func TestCalculateDelay(t *testing.T) {
	manager := NewManager()

	// Test case 1: No remaining requests
	rateLimitInfo := &RateLimitInfo{
		Limit:     100,
		Remaining: 0,
		Reset:     10,
	}
	delay := manager.calculateDelay(rateLimitInfo)
	// Should be around 10s + 10% jitter
	assert.Greater(t, delay, 10*time.Second)
	assert.Less(t, delay, 12*time.Second)

	// Test case 2: Some remaining requests
	rateLimitInfo = &RateLimitInfo{
		Limit:     100,
		Remaining: 5,
		Reset:     10,
	}
	manager.activeRequests = 0
	delay = manager.calculateDelay(rateLimitInfo)
	// Should be around (10s * 0.8) / 5 = 1.6s + 10% jitter
	assert.Greater(t, delay, 1*time.Second)
	assert.Less(t, delay, 2*time.Second)

	// Test case 3: With active requests factored in
	manager.activeRequests = 2
	delay = manager.calculateDelay(rateLimitInfo)
	// With 2 active requests, effective remaining = 5 - 2 = 3
	// Should be around (10s * 0.8) / 3 = 2.67s + 10% jitter
	assert.Greater(t, delay, 2*time.Second)
	assert.Less(t, delay, 4*time.Second)

	// Test case 4: Maximum delay cap
	rateLimitInfo = &RateLimitInfo{
		Limit:     100,
		Remaining: 1,
		Reset:     60, // 60 seconds
	}
	manager.activeRequests = 0
	delay = manager.calculateDelay(rateLimitInfo)
	// Should be capped at 30 seconds, but with jitter it will be slightly more
	assert.Greater(t, delay, 30*time.Second) // Should be at least 30s
	assert.Less(t, delay, 35*time.Second)    // Should be less than 35s (30s + 10% jitter + buffer)
}
