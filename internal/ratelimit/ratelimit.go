package ratelimit

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// RateLimitInfo holds rate limiting information from port api response headers
type RateLimitInfo struct {
	Limit     int // x-ratelimit-limit
	Period    int // x-ratelimit-period
	Remaining int // x-ratelimit-remaining
	Reset     int // x-ratelimit-reset (seconds until reset)
}

// IsNearLimit checks if we're close to hitting the rate limit
func (r *RateLimitInfo) IsNearLimit(threshold float64) bool {
	if r.Limit == 0 {
		return false
	}
	return float64(r.Remaining)/float64(r.Limit) < threshold
}

// ShouldThrottle determines if we should pause before the next request
func (r *RateLimitInfo) ShouldThrottle(threshold float64) bool {
	return r.IsNearLimit(threshold)
}

// Manager handles rate limiting logic and middleware
type Manager struct {
	mu                 sync.RWMutex
	rateLimitInfo      *RateLimitInfo
	enabled            bool
	threshold          float64
	activeRequests     int
	requestSemaphore   chan struct{} // Semaphore to limit concurrent requests
	lastRequestTime    time.Time
	minRequestInterval time.Duration
	debug              bool
}

// NewManager creates a new rate limit manager
func NewManager() *Manager {
	return &Manager{
		enabled:            os.Getenv("PORT_RATE_LIMIT_DISABLED") == "",
		threshold:          0.02,
		debug:              os.Getenv("PORT_DEBUG_RATE_LIMIT") != "",
		requestSemaphore:   make(chan struct{}, 50),
		minRequestInterval: 50 * time.Millisecond,
	}
}

// GetInfo returns a copy of the current rate limit information
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

// SetEnabled enables or disables rate limiting
func (m *Manager) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

// SetThreshold sets the threshold for when to start throttling
func (m *Manager) SetThreshold(threshold float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if threshold >= 0.0 && threshold <= 1.0 {
		m.threshold = threshold
	}
}

func (m *Manager) log(logString string) {
	if m.debug {
		fmt.Print(logString)
	}
}

// RequestMiddleware handles pre-request rate limiting
func (m *Manager) RequestMiddleware(client *resty.Client, req *resty.Request) error {
	m.log(fmt.Sprintf("DEBUG: RequestMiddleware called, enabled: %v\n", m.enabled))

	if !m.enabled {
		m.log("DEBUG: Rate limiting disabled - returning early\n")
		return nil
	}

	m.log(fmt.Sprintf("DEBUG: Attempting to acquire semaphore (capacity: %d)\n", cap(m.requestSemaphore)))

	// Acquire semaphore slot to limit concurrent requests with timeout
	select {
	case m.requestSemaphore <- struct{}{}:
		// Got a slot, continue
		m.log("DEBUG: Acquired semaphore slot\n")
	case <-time.After(30 * time.Second):
		// Timeout waiting for semaphore slot - this prevents indefinite blocking
		m.log("DEBUG: Semaphore timeout - proceeding anyway\n")
	}

	m.log("DEBUG: Getting rate limit state\n")

	m.mu.Lock()
	m.activeRequests++
	rateLimitInfo := m.rateLimitInfo
	lastRequestTime := m.lastRequestTime
	activeRequests := m.activeRequests
	m.mu.Unlock()

	m.log(fmt.Sprintf("DEBUG: Active requests: %d, Rate limit info: %+v\n", activeRequests, rateLimitInfo))

	// Ensure minimum interval between requests to prevent thundering herd
	timeSinceLastRequest := time.Since(lastRequestTime)
	if timeSinceLastRequest < m.minRequestInterval {
		m.log(fmt.Sprintf("DEBUG: Applying minimum interval delay: %v\n", m.minRequestInterval-timeSinceLastRequest))
		time.Sleep(m.minRequestInterval - timeSinceLastRequest)
	}

	// Apply rate limiting if we have rate limit information
	if rateLimitInfo != nil && rateLimitInfo.ShouldThrottle(m.threshold) {
		delay := m.calculateDelay(rateLimitInfo)
		if delay > 0 {
			m.log(fmt.Sprintf("DEBUG: Throttling request - delay: %v (remaining: %d, reset: %d, threshold: %.2f)\n",
				delay, rateLimitInfo.Remaining, rateLimitInfo.Reset, m.threshold))
			time.Sleep(delay)
		}
	} else {
		if rateLimitInfo == nil {
			m.log("DEBUG: No rate limit info available yet\n")
		} else {
			m.log(fmt.Sprintf("DEBUG: Not throttling - remaining: %d, limit: %d, ratio: %.2f, threshold: %.2f\n",
				rateLimitInfo.Remaining, rateLimitInfo.Limit,
				float64(rateLimitInfo.Remaining)/float64(rateLimitInfo.Limit), m.threshold))
		}
	}

	m.mu.Lock()
	m.lastRequestTime = time.Now()
	m.mu.Unlock()

	m.log("DEBUG: RequestMiddleware completed successfully\n")

	return nil
}

// ResponseMiddleware handles post-response rate limit information extraction
func (m *Manager) ResponseMiddleware(client *resty.Client, resp *resty.Response) error {
	m.log(fmt.Sprintf("DEBUG: ResponseMiddleware called, enabled: %v\n", m.enabled))

	if !m.enabled {
		m.log("DEBUG: Rate limiting disabled - response middleware returning early\n")
		return nil
	}

	m.log("DEBUG: Setting up defer function for semaphore release\n")

	// Release the semaphore slot and decrement active requests
	// Use a timeout to prevent hanging if semaphore is in bad state
	defer func() {
		m.log("DEBUG: Defer function executing - attempting to release semaphore\n")

		select {
		case <-m.requestSemaphore:
			m.log("DEBUG: Successfully released semaphore slot\n")
		case <-time.After(1 * time.Second):
			m.log("DEBUG: Semaphore release timeout - continuing anyway\n")
		}

		m.log("DEBUG: Updating active request count\n")

		m.mu.Lock()
		m.activeRequests--
		if m.activeRequests < 0 {
			m.activeRequests = 0 // Prevent negative count
		}
		m.log(fmt.Sprintf("DEBUG: Active requests now: %d\n", m.activeRequests))
		m.mu.Unlock()

		m.log("DEBUG: Defer function completed\n")
	}()

	m.log("DEBUG: Extracting rate limit headers\n")

	// Extract rate limit headers
	limitHeader := resp.Header().Get("x-ratelimit-limit")
	periodHeader := resp.Header().Get("x-ratelimit-period")
	remainingHeader := resp.Header().Get("x-ratelimit-remaining")
	resetHeader := resp.Header().Get("x-ratelimit-reset")

	m.log(fmt.Sprintf("DEBUG: Rate limit headers - limit: %q, period: %q, remaining: %q, reset: %q\n",
		limitHeader, periodHeader, remainingHeader, resetHeader))

	// Only update if we have the essential headers
	if limitHeader != "" && remainingHeader != "" {
		m.log("DEBUG: Parsing rate limit headers\n")

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

		m.log(fmt.Sprintf("DEBUG: Parsed rate limit info: %+v\n", rateLimitInfo))
	} else {
		m.log("DEBUG: No rate limit headers found or incomplete\n")
	}

	m.log("DEBUG: ResponseMiddleware completed successfully\n")

	return nil
}

// calculateDelay calculates the delay needed based on rate limit info and active requests
func (m *Manager) calculateDelay(rateLimitInfo *RateLimitInfo) time.Duration {
	if rateLimitInfo.Remaining <= 0 && rateLimitInfo.Reset > 0 {
		// No requests remaining, wait until reset but cap at 2 minutes for real API usage
		delay := time.Duration(rateLimitInfo.Reset) * time.Second
		if delay > 2*time.Minute {
			delay = 2 * time.Minute
		}
		// Add some jitter to prevent thundering herd when reset occurs
		jitter := time.Duration(float64(delay) * 0.1) // 10% jitter
		delay += jitter
		return delay
	}

	if rateLimitInfo.Remaining > 0 && rateLimitInfo.Reset > 0 {
		// Account for concurrent requests to avoid overshooting
		effectiveRemaining := rateLimitInfo.Remaining - m.activeRequests
		if effectiveRemaining <= 0 {
			effectiveRemaining = 1 // Always leave at least some room
		}

		// Be more conservative - spread requests over 80% of the reset period to leave buffer
		resetBuffer := float64(rateLimitInfo.Reset) * 0.8
		delay := time.Duration(resetBuffer) * time.Second / time.Duration(effectiveRemaining)

		// Cap the delay to a reasonable maximum (30 seconds for real API usage)
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}

		// Ensure minimum delay to prevent thundering herd
		if delay < m.minRequestInterval {
			delay = m.minRequestInterval
		}

		// Add small jitter to prevent synchronized requests
		jitter := time.Duration(float64(delay) * 0.1) // 10% jitter
		delay += jitter

		return delay
	}

	return 0
}
