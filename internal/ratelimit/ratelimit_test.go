package ratelimit

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
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
	manager := NewManager(nil)

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
	manager := NewManager(nil)

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
	manager := NewManager(nil)

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
	manager := NewManager(nil)

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
	manager := NewManager(nil)
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
	manager := NewManager(nil)
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

	assert.Equal(t, int64(numRequests), requestCount.Load())

	t.Logf("Concurrent requests took %v", elapsed)
}

func TestCalculateDelayNoRemainingRequests(t *testing.T) {
	manager := NewManager(nil)

	// No remaining requests - should wait until reset
	rateLimitInfo := &RateLimitInfo{
		Limit:     100,
		Remaining: 0,
		Reset:     10,
	}
	delay := manager.calculateDelay(rateLimitInfo)

	// Base delay: 10s, with jitter (1.1-1.2x): 11-12s
	assert.Greater(t, delay, 10*time.Second)
	assert.Less(t, delay, 13*time.Second)
}

func TestCalculateDelaySomeRemainingRequests(t *testing.T) {
	manager := NewManager(nil)

	t.Run("normal reset time", func(t *testing.T) {
		// Some remaining requests, no active requests
		rateLimitInfo := &RateLimitInfo{
			Limit:     100,
			Remaining: 5,
			Reset:     10,
		}
		manager.activeRequestsMu.Lock()
		manager.activeRequests.Store(0)
		manager.activeRequestsMu.Unlock()
		delay := manager.calculateDelay(rateLimitInfo)

		// Base delay: (10*0.8)/5 = 1.6s, with jitter: 1.76-1.92s
		assert.Greater(t, delay, 1*time.Second)
		assert.Less(t, delay, time.Duration(2.5*float64(time.Second)))
	})

	t.Run("long reset time should be capped", func(t *testing.T) {
		// One remaining request, long reset time (should be capped at 30s)
		rateLimitInfo := &RateLimitInfo{
			Limit:     100,
			Remaining: 1,
			Reset:     60,
		}
		manager.activeRequestsMu.Lock()
		manager.activeRequests.Store(0)
		manager.activeRequestsMu.Unlock()
		delay := manager.calculateDelay(rateLimitInfo)

		// Base calculation: (60*0.8)/1 = 48s, but capped at 30s, with jitter: 33-36s
		assert.Greater(t, delay, 30*time.Second)
		assert.Less(t, delay, 37*time.Second)
	})
}

func TestCalculateDelayWithActiveRequestScaling(t *testing.T) {
	manager := NewManager(nil)

	// Some remaining requests, with active requests
	rateLimitInfo := &RateLimitInfo{
		Limit:     100,
		Remaining: 5,
		Reset:     10,
	}
	manager.activeRequestsMu.Lock()
	manager.activeRequests.Store(2)
	manager.activeRequestsMu.Unlock()
	delay := manager.calculateDelay(rateLimitInfo)

	// Base delay: 1.6s, scaling: 1.6*(1+2*0.2) = 2.24s, with jitter: 2.46-2.69s
	assert.Greater(t, delay, 2*time.Second)
	assert.Less(t, delay, time.Duration(3.5*float64(time.Second)))
}

func TestActiveRequestsCleanup(t *testing.T) {
	// Create manager with short cleanup interval for testing
	cleanupInterval := 100 * time.Millisecond
	manager := NewManager(&ManagerOptions{
		CleanupInterval: &cleanupInterval,
	})
	defer manager.Stop()

	// Artificially set activeRequests to simulate stuck state
	// Use write lock to ensure exclusive access during setup
	manager.activeRequestsMu.Lock()
	manager.activeRequests.Store(5)
	stuckCount := manager.activeRequests.Load()
	manager.activeRequestsMu.Unlock()

	assert.Equal(t, int64(5), stuckCount, "activeRequests should be set to 5")

	// Wait for cleanup to run (should happen within 100ms + some buffer)
	// Use a longer wait to ensure cleanup has definitely run
	time.Sleep(300 * time.Millisecond)

	// Check that activeRequests has been reset
	// Use read lock to safely check the value
	manager.activeRequestsMu.RLock()
	cleanedCount := manager.activeRequests.Load()
	manager.activeRequestsMu.RUnlock()

	assert.Equal(t, int64(0), cleanedCount, "activeRequests should be cleaned up to 0")
}

func TestNewManagerWithOptions(t *testing.T) {
	t.Run("with nil options uses defaults", func(t *testing.T) {
		manager := NewManager(nil)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})

	t.Run("with custom logger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		manager := NewManager(&ManagerOptions{
			Logger: logger,
		})

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		// Other fields should use defaults
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)

		// Test that the logger works by triggering a log message
		manager.SetEnabled(false)
		// Create a simple test to trigger logging
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := resty.New().SetBaseURL(server.URL)
		client.OnBeforeRequest(manager.RequestMiddleware)

		_, err := client.R().Get("/test")
		require.NoError(t, err)

		// Should have logged the "disabled" message
		assert.Contains(t, buf.String(), "Rate limiting disabled")
	})

	t.Run("with custom threshold", func(t *testing.T) {
		threshold := 0.15
		manager := NewManager(&ManagerOptions{
			Threshold: &threshold,
		})

		assert.Equal(t, threshold, manager.threshold)
		// Other fields should use defaults
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})

	t.Run("with custom semaphore weight", func(t *testing.T) {
		weight := int64(100)
		manager := NewManager(&ManagerOptions{
			SemaphoreWeight: &weight,
		})

		assert.NotNil(t, manager.requestSemaphore)
		// We can't directly test the semaphore weight, but we can verify it was set
		// Other fields should use defaults
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})

	t.Run("with custom min request interval", func(t *testing.T) {
		interval := 100 * time.Millisecond
		manager := NewManager(&ManagerOptions{
			MinRequestInterval: &interval,
		})

		assert.Equal(t, interval, manager.minRequestInterval)
		// Other fields should use defaults
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})

	t.Run("with custom cleanup interval", func(t *testing.T) {
		interval := 60 * time.Second
		manager := NewManager(&ManagerOptions{
			CleanupInterval: &interval,
		})

		assert.Equal(t, interval, manager.cleanupInterval)
		// Other fields should use defaults
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
	})

	t.Run("with all custom options", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		threshold := 0.25
		weight := int64(25)
		minInterval := 200 * time.Millisecond
		cleanupInterval := 45 * time.Second

		manager := NewManager(&ManagerOptions{
			Logger:             logger,
			Threshold:          &threshold,
			SemaphoreWeight:    &weight,
			MinRequestInterval: &minInterval,
			CleanupInterval:    &cleanupInterval,
		})

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		assert.Equal(t, threshold, manager.threshold)
		assert.Equal(t, minInterval, manager.minRequestInterval)
		assert.Equal(t, cleanupInterval, manager.cleanupInterval)
	})

	t.Run("NewManagerWithDebug respects environment", func(t *testing.T) {
		t.Setenv("PORT_DEBUG_RATE_LIMIT", "1")
		manager := NewManagerWithDebug()
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		// Should use defaults for other fields
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})

	t.Run("NewManagerWithDebug without debug env uses discard", func(t *testing.T) {
		// Ensure the env var is not set
		os.Unsetenv("PORT_DEBUG_RATE_LIMIT")
		manager := NewManagerWithDebug()
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.logger)
		// Should use defaults for other fields
		assert.Equal(t, 0.02, manager.threshold)
		assert.Equal(t, 50*time.Millisecond, manager.minRequestInterval)
		assert.Equal(t, 30*time.Second, manager.cleanupInterval)
	})
}

func TestDefaultManagerOptions(t *testing.T) {
	defaults := DefaultManagerOptions()

	assert.NotNil(t, defaults.Logger)
	assert.NotNil(t, defaults.Threshold)
	assert.Equal(t, 0.02, *defaults.Threshold)
	assert.NotNil(t, defaults.SemaphoreWeight)
	assert.Equal(t, int64(50), *defaults.SemaphoreWeight)
	assert.NotNil(t, defaults.MinRequestInterval)
	assert.Equal(t, 50*time.Millisecond, *defaults.MinRequestInterval)
	assert.NotNil(t, defaults.CleanupInterval)
	assert.Equal(t, 30*time.Second, *defaults.CleanupInterval)
}

func TestNoDoubleSleeping(t *testing.T) {
	// Create manager with longer minimum interval for testing
	minInterval := 200 * time.Millisecond
	manager := NewManager(&ManagerOptions{
		MinRequestInterval: &minInterval,
		Threshold:          ptrTo(0.1), // Low threshold to trigger throttling
	})

	// Set up rate limit info that will trigger throttling
	manager.mu.Lock()
	manager.rateLimitInfo = &RateLimitInfo{
		Limit:     100,
		Remaining: 5, // Low remaining to trigger throttling
		Reset:     1, // Short reset time
	}
	// Set lastRequestTime to trigger minimum interval delay
	manager.lastRequestTime = time.Now().Add(-100 * time.Millisecond) // 100ms ago
	manager.mu.Unlock()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	start := time.Now()
	_, err := client.R().Get("/test")
	elapsed := time.Since(start)

	require.NoError(t, err)

	// Should have slept for throttling delay (which is longer than min interval)
	// but NOT for both delays added together
	// Throttling delay calculation:
	// - resetBuffer = 1 * 0.8 = 0.8s
	// - baseDelay = 0.8 / 5 = 0.16s
	// - with 1 active request scaling: 0.16 * 1.2 = 0.192s
	// - with jitter: 0.192 * (1.1 to 1.2) = 0.211 to 0.230s
	// Min interval delay would be 100ms (200ms - 100ms already passed)
	// If we were double sleeping, it would be ~330ms+ (230ms + 100ms)
	// With the fix, it should be just the throttling delay (~210-230ms)

	assert.Greater(t, elapsed, 200*time.Millisecond, "Should have throttled")
	assert.Less(t, elapsed, 280*time.Millisecond, "Should not have double slept")

	t.Logf("Request took %v (should be ~210-230ms for throttling only, not ~330ms+ for double sleep)", elapsed)
}

func TestSemaphoreTimeoutNoReleasePanic(t *testing.T) {
	// Create manager with debug logging to capture log messages
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create manager with very small semaphore weight to force immediate "timeouts"
	semaphoreWeight := int64(1)
	manager := NewManager(&ManagerOptions{
		SemaphoreWeight: &semaphoreWeight,
		Logger:          logger,
	})
	defer manager.Stop()

	// First, exhaust the semaphore by acquiring it manually
	ctx := context.Background()
	err := manager.requestSemaphore.Acquire(ctx, 1)
	require.NoError(t, err, "Should be able to acquire the single semaphore slot")

	// Create a test that simulates the scenario without actually waiting
	// We'll test the logic directly by creating a request with the proper context
	req := &resty.Request{}
	req.SetContext(context.WithValue(context.Background(), semaphoreAcquiredKey, false))

	resp := &resty.Response{
		Request: req,
	}

	// Test ResponseMiddleware with semaphore_acquired = false
	// This should NOT try to release the semaphore and should NOT panic
	err = manager.ResponseMiddleware(nil, resp)
	require.NoError(t, err, "ResponseMiddleware should complete without error")

	// Check logs to verify semaphore release was skipped
	logs := buf.String()
	assert.Contains(t, logs, "Skipping semaphore release - was not acquired",
		"Should log that semaphore release was skipped")
	assert.NotContains(t, logs, "Successfully released semaphore slot",
		"Should not log successful semaphore release")

	// Clean up - release the semaphore we acquired manually
	manager.requestSemaphore.Release(1)
}

func TestSemaphoreNonBlockingLoadIndicator(t *testing.T) {
	// Create manager with debug logging to capture log messages
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create manager with small semaphore weight to test load indication
	semaphoreWeight := int64(2)
	manager := NewManager(&ManagerOptions{
		SemaphoreWeight: &semaphoreWeight,
		Logger:          logger,
	})
	defer manager.Stop()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := resty.New().SetBaseURL(server.URL)
	client.OnBeforeRequest(manager.RequestMiddleware)
	client.OnAfterResponse(manager.ResponseMiddleware)

	// First request should succeed and acquire semaphore
	_, err := client.R().Get("/test1")
	require.NoError(t, err, "First request should succeed")

	// Second request should succeed and acquire semaphore
	_, err = client.R().Get("/test2")
	require.NoError(t, err, "Second request should succeed")

	// Manually exhaust the semaphore to test high load behavior
	ctx := context.Background()
	err = manager.requestSemaphore.Acquire(ctx, 1)
	require.NoError(t, err, "Should be able to acquire remaining semaphore slot")
	err = manager.requestSemaphore.Acquire(ctx, 1)
	require.NoError(t, err, "Should be able to acquire final semaphore slot")

	// Clear the log buffer to focus on the high load case
	buf.Reset()

	// Third request should proceed anyway but log high load
	start := time.Now()
	_, err = client.R().Get("/test3")
	elapsed := time.Since(start)
	require.NoError(t, err, "Third request should still succeed despite high load")

	// Should complete quickly since it's non-blocking
	assert.Less(t, elapsed, 100*time.Millisecond, "Request should complete quickly (non-blocking)")

	// Check logs for high load warning
	logs := buf.String()
	assert.Contains(t, logs, "High concurrent load detected - proceeding anyway",
		"Should log high load warning")

	// Clean up - release the semaphores we acquired manually
	manager.requestSemaphore.Release(1)
	manager.requestSemaphore.Release(1)
}

func TestSemaphoreSuccessfulAcquireAndRelease(t *testing.T) {
	// Create manager with debug logging to capture log messages
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	manager := NewManager(&ManagerOptions{
		Logger: logger,
	})
	defer manager.Stop()

	// First acquire a semaphore slot so we can legitimately release it
	ctx := context.Background()
	err := manager.requestSemaphore.Acquire(ctx, 1)
	require.NoError(t, err, "Should be able to acquire semaphore slot")

	// Create a test that simulates successful semaphore acquisition
	req := &resty.Request{}
	req.SetContext(context.WithValue(context.Background(), semaphoreAcquiredKey, true))

	resp := &resty.Response{
		Request: req,
	}

	// Test ResponseMiddleware with semaphore_acquired = true
	// This should try to release the semaphore successfully
	err = manager.ResponseMiddleware(nil, resp)
	require.NoError(t, err, "ResponseMiddleware should complete without error")

	// Check logs to verify semaphore was released
	logs := buf.String()
	assert.Contains(t, logs, "Attempting to release semaphore",
		"Should log that semaphore release was attempted")
	assert.Contains(t, logs, "Successfully released semaphore slot",
		"Should log successful semaphore release")
	assert.NotContains(t, logs, "Skipping semaphore release",
		"Should not log that semaphore release was skipped")
}
