package ratelimit

import (
	"context"
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

// Extracted from Port HTTP Headers
type RateLimitInfo struct {
	// x-ratelimit-limit
	Limit int
	// x-ratelimit-period
	Period int
	// x-ratelimit-remaining
	Remaining int
	// x-ratelimit-reset (seconds until reset)
	Reset int
}

func (r *RateLimitInfo) ShouldThrottle(threshold float64) bool {
	if r.Limit == 0 {
		return false
	}
	return float64(r.Remaining)/float64(r.Limit) < threshold
}

type Manager struct {
	mu                 sync.RWMutex
	rateLimitInfo      *RateLimitInfo
	enabled            bool
	threshold          float64
	activeRequests     int64 // Changed to int64 for atomic operations
	requestSemaphore   *semaphore.Weighted
	lastRequestTime    time.Time
	minRequestInterval time.Duration
	logger             *slog.Logger

	// Cleanup mechanism for activeRequests drift
	cleanupStop     chan struct{}
	cleanupInterval time.Duration
}

func NewManager() *Manager {
	enabled := os.Getenv("PORT_RATE_LIMIT_DISABLED") == ""
	debug := os.Getenv("PORT_DEBUG_RATE_LIMIT") != ""

	// Create logger with appropriate level and output
	var logger *slog.Logger
	if debug {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})).With("component", "ratelimit", "enabled", enabled)
	} else {
		// Create a no-op logger that discards output
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelError, // Only log errors and above (effectively disabling debug)
		})).With("component", "ratelimit", "enabled", enabled)
	}

	manager := &Manager{
		enabled:            enabled,
		threshold:          0.02,
		requestSemaphore:   semaphore.NewWeighted(50),
		minRequestInterval: 50 * time.Millisecond,
		logger:             logger,
		cleanupStop:        make(chan struct{}),
		cleanupInterval:    30 * time.Second, // Reset stuck activeRequests every 30 seconds
	}

	// Start cleanup goroutine to handle activeRequests drift
	manager.startCleanup()

	return manager
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
	if m.cleanupStop != nil {
		close(m.cleanupStop)
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
				// Check and reset activeRequests atomically (no lock needed)
				currentActiveRequests := atomic.LoadInt64(&m.activeRequests)
				if currentActiveRequests > 0 {
					m.logger.Debug("Cleaning up stuck activeRequests",
						"stuck_count", currentActiveRequests,
						"cleanup_interval", m.cleanupInterval)
					atomic.StoreInt64(&m.activeRequests, 0)
				}

			case <-m.cleanupStop:
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

	m.logger.Debug("Attempting to acquire semaphore", "capacity", 50)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.requestSemaphore.Acquire(ctx, 1); err != nil {
		m.logger.Debug("Semaphore timeout - proceeding anyway")
	} else {
		m.logger.Debug("Acquired semaphore slot")
	}

	m.logger.Debug("Getting rate limit state")

	m.mu.Lock()
	rateLimitInfo := m.rateLimitInfo
	lastRequestTime := m.lastRequestTime
	m.mu.Unlock()

	// Increment activeRequests atomically (no lock needed)
	activeRequests := atomic.AddInt64(&m.activeRequests, 1)

	m.logger.Debug("Active requests and rate limit info", "active_requests", activeRequests, "rate_limit_info", rateLimitInfo)

	timeSinceLastRequest := time.Since(lastRequestTime)
	if timeSinceLastRequest < m.minRequestInterval {
		delay := m.minRequestInterval - timeSinceLastRequest
		m.logger.Debug("Applying minimum interval delay", "delay", delay)
		time.Sleep(delay)
	}

	if rateLimitInfo != nil && rateLimitInfo.ShouldThrottle(m.threshold) {
		delay := m.calculateDelay(rateLimitInfo)
		if delay > 0 {
			m.logger.Debug("Throttling request",
				"delay", delay,
				"remaining", rateLimitInfo.Remaining,
				"reset", rateLimitInfo.Reset,
				"threshold", m.threshold)
			time.Sleep(delay)
		}
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
		m.logger.Debug("Defer function executing - attempting to release semaphore")

		m.requestSemaphore.Release(1)
		m.logger.Debug("Successfully released semaphore slot")

		m.logger.Debug("Updating active request count")

		// Decrement activeRequests atomically (no lock needed)
		activeRequests := atomic.AddInt64(&m.activeRequests, -1)
		if activeRequests < 0 {
			// Reset to 0 if somehow went negative
			atomic.StoreInt64(&m.activeRequests, 0)
			activeRequests = 0
		}
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

	if limitHeader != "" && remainingHeader != "" {
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
	} else {
		m.logger.Debug("No rate limit headers found or incomplete")
	}

	m.logger.Debug("ResponseMiddleware completed successfully")

	return nil
}

func (m *Manager) calculateDelay(rateLimitInfo *RateLimitInfo) time.Duration {
	// Case 1: No remaining requests - wait until reset
	if rateLimitInfo.Remaining <= 0 && rateLimitInfo.Reset > 0 {
		delay := float64(rateLimitInfo.Reset)
		// Cap at 2 minutes
		if delay > 120 {
			delay = 120
		}
		// Apply random jitter (10-20% increase)
		jitterMultiplier := 1.1 + rand.Float64()*0.1
		return time.Duration(delay * jitterMultiplier * float64(time.Second))
	}

	// Case 2: Some requests remaining - calculate based on remaining and reset time
	if rateLimitInfo.Remaining > 0 && rateLimitInfo.Reset > 0 {
		// Use 80% of reset time as buffer
		resetBuffer := float64(rateLimitInfo.Reset) * 0.8

		// Calculate base delay: spread remaining time across remaining requests
		baseDelay := resetBuffer / float64(rateLimitInfo.Remaining)

		// Apply scaling factor based on active requests to be more conservative
		// When we have many active requests, increase the delay to avoid overwhelming
		activeRequests := atomic.LoadInt64(&m.activeRequests)
		activeRequestScaling := 1.0 + float64(activeRequests)*0.2
		delay := baseDelay * activeRequestScaling

		// Cap delay at 30 seconds
		if delay > 30 {
			delay = 30
		}

		// Ensure minimum delay
		minDelaySeconds := m.minRequestInterval.Seconds()
		if delay < minDelaySeconds {
			delay = minDelaySeconds
		}

		// Apply random jitter (10-20% increase)
		jitterMultiplier := 1.1 + rand.Float64()*0.1
		return time.Duration(delay * jitterMultiplier * float64(time.Second))
	}

	return 0
}
