package ratelimit

import (
	"context"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/semaphore"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	semaphoreAcquiredKey contextKey = "semaphore_acquired"
)

// RateLimitInfo holds rate limit information extracted from Port HTTP Headers
type RateLimitInfo struct {
	// Limit is extracted from x-ratelimit-limit
	Limit int
	// Period is extracted from x-ratelimit-period
	Period int
	// Remaining is extracted from x-ratelimit-remaining
	Remaining int
	// Reset is extracted from x-ratelimit-reset (seconds until reset)
	Reset int
}

func (r *RateLimitInfo) ShouldThrottle(threshold float64) bool {
	if r.Limit == 0 {
		return false
	}
	return float64(r.Remaining)/float64(r.Limit) < threshold
}

// ManagerOptions holds configuration options for creating a new Manager
type ManagerOptions struct {
	Logger             *slog.Logger
	Threshold          *float64
	SemaphoreWeight    *int64
	MinRequestInterval *time.Duration
	CleanupInterval    *time.Duration
}

// DefaultManagerOptions returns a ManagerOptions struct with sensible defaults
func DefaultManagerOptions() *ManagerOptions {
	return &ManagerOptions{
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		Threshold:          ptrTo(0.02),
		SemaphoreWeight:    ptrTo(int64(50)),
		MinRequestInterval: ptrTo(50 * time.Millisecond),
		CleanupInterval:    ptrTo(30 * time.Second),
	}
}

// ptrTo returns a pointer to the given value
func ptrTo[T any](v T) *T {
	return &v
}

type Manager struct {
	mu                 sync.RWMutex
	rateLimitInfo      *RateLimitInfo
	enabled            bool
	threshold          float64
	activeRequests     atomic.Int64
	activeRequestsMu   sync.RWMutex
	requestSemaphore   *semaphore.Weighted
	lastRequestTime    time.Time
	minRequestInterval time.Duration
	logger             *slog.Logger

	cleanupCtx      context.Context
	cleanupCancel   context.CancelFunc
	cleanupInterval time.Duration
}

func NewManager(opts *ManagerOptions) *Manager {
	enabled := os.Getenv("PORT_RATE_LIMIT_DISABLED") == ""

	// Use defaults if opts is nil
	if opts == nil {
		opts = DefaultManagerOptions()
	}

	// Apply defaults for nil fields
	defaults := DefaultManagerOptions()
	if opts.Logger == nil {
		opts.Logger = defaults.Logger
	}
	if opts.Threshold == nil {
		opts.Threshold = defaults.Threshold
	}
	if opts.SemaphoreWeight == nil {
		opts.SemaphoreWeight = defaults.SemaphoreWeight
	}
	if opts.MinRequestInterval == nil {
		opts.MinRequestInterval = defaults.MinRequestInterval
	}
	if opts.CleanupInterval == nil {
		opts.CleanupInterval = defaults.CleanupInterval
	}

	// Add component context to logger
	logger := opts.Logger.With("component", "ratelimit", "enabled", enabled)

	// Create context for cleanup goroutine
	ctx, cancel := context.WithCancel(context.Background())

	manager := &Manager{
		enabled:            enabled,
		threshold:          *opts.Threshold,
		requestSemaphore:   semaphore.NewWeighted(*opts.SemaphoreWeight),
		minRequestInterval: *opts.MinRequestInterval,
		logger:             logger,
		cleanupCtx:         ctx,
		cleanupCancel:      cancel,
		cleanupInterval:    *opts.CleanupInterval,
	}

	// Start cleanup goroutine to handle activeRequests drift
	manager.startCleanup()

	return manager
}

// NewManagerWithDebug creates a new Manager with debug logging enabled
func NewManagerWithDebug() *Manager {
	debug := os.Getenv("PORT_DEBUG_RATE_LIMIT") != ""

	opts := DefaultManagerOptions()
	if debug {
		opts.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return NewManager(opts)
}

func (m *Manager) GetInfo() *RateLimitInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rateLimitInfo == nil {
		return nil
	}

	return &RateLimitInfo{
		Limit:     m.rateLimitInfo.Limit,
		Period:    m.rateLimitInfo.Period,
		Remaining: m.rateLimitInfo.Remaining,
		Reset:     m.rateLimitInfo.Reset,
	}
}

func (m *Manager) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

func (m *Manager) SetThreshold(threshold float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if threshold >= 0.0 && threshold <= 1.0 {
		m.threshold = threshold
	}
}

// Stop gracefully shuts down the rate limit manager
func (m *Manager) Stop() {
	if m.cleanupCancel != nil {
		m.cleanupCancel()
	}
}

// startCleanup starts a background goroutine that periodically resets activeRequests
// to handle cases where the count gets stuck due to timeouts, panics, etc.
func (m *Manager) startCleanup() {
	go func() {
		ticker := time.NewTicker(m.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check and reset activeRequests with write lock
				m.activeRequestsMu.Lock()
				currentActiveRequests := m.activeRequests.Load()
				if currentActiveRequests > 0 {
					m.logger.Debug("Cleaning up stuck activeRequests",
						"stuck_count", currentActiveRequests,
						"cleanup_interval", m.cleanupInterval)
					m.activeRequests.Store(0)
				}
				m.activeRequestsMu.Unlock()

			case <-m.cleanupCtx.Done():
				m.logger.Debug("Cleanup goroutine stopping")
				return
			}
		}
	}()
}

func (m *Manager) RequestMiddleware(client *resty.Client, req *resty.Request) error {
	m.logger.Debug("RequestMiddleware called")

	if !m.enabled {
		m.logger.Debug("Rate limiting disabled - returning early")
		return nil
	}

	m.logger.Debug("Checking load via semaphore")

	// Try to acquire semaphore immediately (non-blocking) as load indicator
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	semaphoreAcquired := false
	if err := m.requestSemaphore.Acquire(ctx, 1); err != nil {
		m.logger.Warn("High concurrent load detected - proceeding anyway")
	} else {
		m.logger.Debug("Normal load - acquired semaphore slot")
		semaphoreAcquired = true
	}

	// Store semaphore acquisition status in request context for ResponseMiddleware
	req.SetContext(context.WithValue(req.Context(), semaphoreAcquiredKey, semaphoreAcquired))

	m.logger.Debug("Getting rate limit state")

	m.mu.RLock()
	rateLimitInfo := m.rateLimitInfo
	lastRequestTime := m.lastRequestTime
	m.mu.RUnlock()

	// Increment activeRequests with read lock
	m.activeRequestsMu.RLock()
	activeRequests := m.activeRequests.Add(1)
	m.activeRequestsMu.RUnlock()

	m.logger.Debug("Active requests and rate limit info", "active_requests", activeRequests, "rate_limit_info", rateLimitInfo)

	// Calculate minimum interval delay
	var minIntervalDelay time.Duration
	timeSinceLastRequest := time.Since(lastRequestTime)
	if timeSinceLastRequest < m.minRequestInterval {
		minIntervalDelay = m.minRequestInterval - timeSinceLastRequest
	}

	// Calculate throttling delay
	var throttlingDelay time.Duration
	if rateLimitInfo != nil && rateLimitInfo.ShouldThrottle(m.threshold) {
		throttlingDelay = m.calculateDelay(rateLimitInfo)
	}

	// Use the maximum of the two delays to avoid double sleeping
	finalDelay := max(minIntervalDelay, throttlingDelay)

	if finalDelay > 0 {
		if throttlingDelay > minIntervalDelay {
			m.logger.Debug("Throttling request",
				"delay", finalDelay,
				"remaining", rateLimitInfo.Remaining,
				"reset", rateLimitInfo.Reset,
				"threshold", m.threshold)
		} else {
			m.logger.Debug("Applying minimum interval delay", "delay", finalDelay)
		}
		time.Sleep(finalDelay)
	} else {
		if rateLimitInfo == nil {
			m.logger.Debug("No rate limit info available yet")
		} else {
			ratio := float64(rateLimitInfo.Remaining) / float64(rateLimitInfo.Limit)
			m.logger.Debug("Not throttling",
				"remaining", rateLimitInfo.Remaining,
				"limit", rateLimitInfo.Limit,
				"ratio", ratio,
				"threshold", m.threshold)
		}
	}

	m.mu.Lock()
	m.lastRequestTime = time.Now()
	m.mu.Unlock()

	m.logger.Debug("RequestMiddleware completed successfully")

	return nil
}

func (m *Manager) ResponseMiddleware(client *resty.Client, resp *resty.Response) error {
	m.logger.Debug("ResponseMiddleware called")

	if !m.enabled {
		m.logger.Debug("Rate limiting disabled - response middleware returning early")
		return nil
	}

	m.logger.Debug("Setting up defer function for semaphore release")

	defer func() {
		// Use recover to catch any potential panics from semaphore operations
		if r := recover(); r != nil {
			m.logger.Debug("Recovered from panic in ResponseMiddleware defer", "panic", r)
		}

		m.logger.Debug("Defer function executing")

		// Only release semaphore if we actually acquired it
		if semaphoreAcquired, ok := resp.Request.Context().Value(semaphoreAcquiredKey).(bool); ok && semaphoreAcquired {
			m.logger.Debug("Attempting to release semaphore")
			m.requestSemaphore.Release(1)
			m.logger.Debug("Successfully released semaphore slot")
		} else {
			m.logger.Debug("Skipping semaphore release - was not acquired or acquisition status unknown")
		}

		m.logger.Debug("Updating active request count")

		// Decrement activeRequests with read lock
		m.activeRequestsMu.RLock()
		activeRequests := m.activeRequests.Add(-1)
		if activeRequests < 0 {
			// Reset to 0 if somehow went negative
			m.activeRequests.Store(0)
			activeRequests = 0
		}
		m.activeRequestsMu.RUnlock()
		m.logger.Debug("Active requests now", "active_requests", activeRequests)

		m.logger.Debug("Defer function completed")
	}()

	m.logger.Debug("Extracting rate limit headers")

	limitHeader := resp.Header().Get("x-ratelimit-limit")
	periodHeader := resp.Header().Get("x-ratelimit-period")
	remainingHeader := resp.Header().Get("x-ratelimit-remaining")
	resetHeader := resp.Header().Get("x-ratelimit-reset")

	m.logger.Debug("Rate limit headers",
		"limit", limitHeader,
		"period", periodHeader,
		"remaining", remainingHeader,
		"reset", resetHeader)

	if limitHeader == "" || remainingHeader == "" {
		m.logger.Debug("No rate limit headers found or incomplete")
		return nil
	}
	m.logger.Debug("Parsing rate limit headers")

	m.mu.Lock()
	defer m.mu.Unlock()

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

	m.rateLimitInfo = rateLimitInfo

	m.logger.Debug("Parsed rate limit info", "rate_limit_info", rateLimitInfo)

	return nil
}

func (m *Manager) calculateDelay(rateLimitInfo *RateLimitInfo) time.Duration {
	if rateLimitInfo.Remaining <= 0 && rateLimitInfo.Reset > 0 {
		delay := float64(rateLimitInfo.Reset)
		if delay > 120 {
			delay = 120
		}
		jitterMultiplier := 1.1 + rand.Float64()*0.1
		return time.Duration(delay * jitterMultiplier * float64(time.Second))
	}

	if rateLimitInfo.Remaining > 0 && rateLimitInfo.Reset > 0 {
		resetBuffer := float64(rateLimitInfo.Reset) * 0.8

		baseDelay := resetBuffer / float64(rateLimitInfo.Remaining)

		m.activeRequestsMu.RLock()
		activeRequests := m.activeRequests.Load()
		m.activeRequestsMu.RUnlock()
		activeRequestScaling := 1.0 + float64(activeRequests)*0.2
		delay := baseDelay * activeRequestScaling

		if delay > 30 {
			delay = 30
		}

		minDelaySeconds := m.minRequestInterval.Seconds()
		if delay < minDelaySeconds {
			delay = minDelaySeconds
		}

		jitterMultiplier := 1.1 + rand.Float64()*0.1
		return time.Duration(delay * jitterMultiplier * float64(time.Second))
	}

	return 0
}
