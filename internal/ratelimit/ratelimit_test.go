package ratelimit

import (
	"bytes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerInit(t *testing.T) {
	manager := New(nil)
	t.Cleanup(manager.Close)

	assert.Nil(t, manager.GetInfo())
	assert.Nil(t, manager.lastRequestTime.Load())
}

func TestResponseMiddleware(t *testing.T) {
	manager := New(nil)
	t.Cleanup(manager.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-period", "300")
		w.Header().Set("x-ratelimit-remaining", "50")
		w.Header().Set("x-ratelimit-reset", "120")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL).
		SetRateLimiter(manager).
		OnAfterResponse(manager.ResponseMiddleware)

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
	manager := New(nil)
	t.Cleanup(manager.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL).
		SetRateLimiter(manager).
		OnAfterResponse(manager.ResponseMiddleware)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestMiddlewareDisabled(t *testing.T) {
	manager := New(&Options{Enabled: utils.PtrTo(false)})
	t.Cleanup(manager.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL).
		SetRateLimiter(manager).
		OnAfterResponse(manager.ResponseMiddleware)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	rateLimitInfo := manager.GetInfo()
	assert.Nil(t, rateLimitInfo)
}

func TestThrottling(t *testing.T) {
	manager := New(nil)
	t.Cleanup(manager.Close)
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

	client := resty.New().SetBaseURL(server.URL).
		SetRateLimiter(manager).
		OnAfterResponse(manager.ResponseMiddleware)

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
	manager := New(nil)
	t.Cleanup(manager.Close)
	var requestCount atomic.Int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-ratelimit-limit", "10")
		w.Header().Set("x-ratelimit-remaining", "2")
		w.Header().Set("x-ratelimit-reset", "5")

		requestCount.Add(1)
		time.Sleep(50 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL).
		SetRateLimiter(manager).
		OnAfterResponse(manager.ResponseMiddleware)

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

	assert.Equal(t, int64(numRequests), requestCount.Load())

	t.Logf("Concurrent requests took %v", elapsed)
}

func TestCalculateDelayNoRemainingRequests(t *testing.T) {
	manager := New(nil)
	t.Cleanup(manager.Close)

	// No remaining requests - should wait until reset
	rateLimitInfo := &Info{
		Limit:     100,
		Remaining: 0,
		Reset:     10,
	}
	delay := manager.calculateDelay(nil, rateLimitInfo)

	// Base delay: 10s, with jitter (1.1-1.2x): 11-12s
	assert.Greater(t, delay, 10*time.Second)
	assert.Less(t, delay, 13*time.Second)
}

func TestCalculateDelay(t *testing.T) {
	manager := New(nil)
	t.Cleanup(manager.Close)

	t.Run("Delay using MinRequestInterval", func(t *testing.T) {
		// Some remaining requests, no active requests
		rateLimitInfo := &Info{
			Limit:     100,
			Remaining: 5,
			Reset:     10,
		}
		delay := manager.calculateDelay(utils.PtrTo(time.Now()), rateLimitInfo)

		assert.Greater(t, delay, time.Duration(0))
		assert.Less(t, delay, time.Duration(1.5*float64(time.Second)))
	})

	t.Run("Delay using Info.Reset", func(t *testing.T) {
		// One remaining request, long reset time (should be capped at 30s)
		rateLimitInfo := &Info{
			Limit:     100,
			Remaining: 0,
			Reset:     60,
		}
		delay := manager.calculateDelay(utils.PtrTo(time.Now()), rateLimitInfo)

		// delay should be between [Info.Reset] and 1.1 times the [Info.Reset] (for jitter)
		assert.GreaterOrEqual(t, delay, time.Duration(rateLimitInfo.Reset)*time.Second)
		assert.LessOrEqual(t, delay, time.Duration(float64(rateLimitInfo.Reset)*1.1)*time.Second)
	})
}

func TestNewManagerWithOptions(t *testing.T) {
	t.Run("with nil options uses defaults", func(t *testing.T) {
		manager := New(nil)
		t.Cleanup(manager.Close)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
	})

	t.Run("with custom logger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		manager := New(&Options{
			Logger:  logger,
			Enabled: utils.PtrTo(false),
		})
		t.Cleanup(manager.Close)

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		// Other fields should use defaults
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)

		// Test that the logger works by triggering a log message
		// Create a simple test to trigger logging
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := resty.New().SetBaseURL(server.URL).SetRateLimiter(manager)

		_, err := client.R().Get("/test")
		require.NoError(t, err)

		// Should have logged the "disabled" message
		assert.Contains(t, buf.String(), "Rate limiting disabled")
	})

	t.Run("with custom min request interval", func(t *testing.T) {
		interval := 100 * time.Millisecond
		manager := New(&Options{
			MinRequestInterval: &interval,
		})
		t.Cleanup(manager.Close)

		assert.Equal(t, interval, manager.minRequestInterval)
		// Other fields should use defaults
	})

	t.Run("with all custom options", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		minInterval := 200 * time.Millisecond

		manager := New(&Options{
			Logger:             logger,
			MinRequestInterval: &minInterval,
		})
		t.Cleanup(manager.Close)

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		assert.Equal(t, minInterval, manager.minRequestInterval)
	})
}
