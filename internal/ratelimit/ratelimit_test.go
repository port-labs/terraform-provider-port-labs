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

	assert.True(t, rateLimitInfo.ShouldThrottle(0.2))   // 10/100 = 0.1 < 0.2
	assert.False(t, rateLimitInfo.ShouldThrottle(0.05)) // 10/100 = 0.1 > 0.05

	zeroLimitInfo := &RateLimitInfo{Limit: 0}
	assert.False(t, zeroLimitInfo.ShouldThrottle(0.5))
}

func TestManagerBasicFunctionality(t *testing.T) {
	manager := NewManager()

	assert.Nil(t, manager.GetInfo())

	manager.SetEnabled(false)
	assert.False(t, manager.enabled)

	manager.SetEnabled(true)
	assert.True(t, manager.enabled)

	manager.SetThreshold(0.25)
	assert.Equal(t, 0.25, manager.threshold)

	manager.SetThreshold(-0.1)
	assert.Equal(t, 0.25, manager.threshold)

	manager.SetThreshold(1.5)
	assert.Equal(t, 0.25, manager.threshold)
}

func TestResponseMiddleware(t *testing.T) {
	manager := NewManager()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-period", "300")
		w.Header().Set("x-ratelimit-remaining", "50")
		w.Header().Set("x-ratelimit-reset", "120")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := manager.GetInfo()
	require.NotNil(t, rateLimitInfo)
	assert.Equal(t, 1000, rateLimitInfo.Limit)
	assert.Equal(t, 300, rateLimitInfo.Period)
	assert.Equal(t, 50, rateLimitInfo.Remaining)
	assert.Equal(t, 120, rateLimitInfo.Reset)
}

func TestResponseMiddlewareNoHeaders(t *testing.T) {
	manager := NewManager()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestMiddlewareDisabled(t *testing.T) {
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "123")
	manager := NewManager()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestThrottling(t *testing.T) {
	manager := NewManager()
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "5") // Low remaining requests
		w.Header().Set("x-ratelimit-reset", "2")     // Reset in 2 seconds
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	start := time.Now()
	resp, err := client.R().Get("/test1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	resp, err = client.R().Get("/test2")
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 2, requestCount)

	assert.Greater(t, elapsed, 10*time.Millisecond, "Request should have been throttled")
}

func TestConcurrentRequests(t *testing.T) {
	manager := NewManager()
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "10")
		w.Header().Set("x-ratelimit-remaining", "2")
		w.Header().Set("x-ratelimit-reset", "5")

		requestCount++
		time.Sleep(50 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	start := time.Now()
	numRequests := 3
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			_, err := client.R().Get("/test")
			results <- err
		}(i)
	}

	for range numRequests {
		err := <-results
		assert.NoError(t, err)
	}

	elapsed := time.Since(start)

	assert.Equal(t, numRequests, requestCount)

	t.Logf("Concurrent requests took %v", elapsed)
}

func TestCalculateDelay(t *testing.T) {
	manager := NewManager()

	rateLimitInfo := &RateLimitInfo{
		Limit:     100,
		Remaining: 0,
		Reset:     10,
	}
	delay := manager.calculateDelay(rateLimitInfo)

	assert.Greater(t, delay, 10*time.Second)
	assert.Less(t, delay, 12*time.Second)

	rateLimitInfo = &RateLimitInfo{
		Limit:     100,
		Remaining: 5,
		Reset:     10,
	}
	manager.activeRequests = 0
	delay = manager.calculateDelay(rateLimitInfo)

	assert.Greater(t, delay, 1*time.Second)
	assert.Less(t, delay, 2*time.Second)

	manager.activeRequests = 2
	delay = manager.calculateDelay(rateLimitInfo)

	assert.Greater(t, delay, 2*time.Second)
	assert.Less(t, delay, 4*time.Second)

	rateLimitInfo = &RateLimitInfo{
		Limit:     100,
		Remaining: 1,
		Reset:     60,
	}
	manager.activeRequests = 0
	delay = manager.calculateDelay(rateLimitInfo)

	assert.Greater(t, delay, 30*time.Second)
	assert.Less(t, delay, 35*time.Second)
}
