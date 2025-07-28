package ratelimit

import (
	"context"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"io"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
)

// Info holds rate limit information extracted from Port HTTP Headers
type Info struct {
	// Limit is extracted from x-ratelimit-limit
	Limit int
	// Period is extracted from x-ratelimit-period
	Period int
	// Remaining is extracted from x-ratelimit-remaining
	Remaining int
	// Reset is extracted from x-ratelimit-reset (seconds until reset)
	Reset int
}

func (r *Info) ShouldThrottle() bool {
	if r == nil || r.Limit == 0 {
		return false
	}
	return r.Remaining <= 0
}

// Options holds configuration options for creating a new Manager
type Options struct {
	Logger             *slog.Logger
	MinRequestInterval *time.Duration
	Enabled            *bool
	Ctx                context.Context
}

// DefaultOptions returns a Options struct with sensible defaults
func DefaultOptions() *Options {
	return &Options{
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		MinRequestInterval: utils.PtrTo(50 * time.Millisecond),
		Enabled:            utils.PtrTo(true),
		Ctx:                context.Background(),
	}
}

type Manager struct {
	rateLimitInfo      atomic.Pointer[Info]
	enabled            bool
	lastRequestTime    atomic.Pointer[time.Time]
	minRequestInterval time.Duration
	logger             *slog.Logger

	ctx           context.Context
	cancelCtxFunc context.CancelFunc
}

func New(opts *Options) *Manager {
	if opts == nil {
		opts = &Options{}
	}

	// Apply defaults for nil fields
	defaults := DefaultOptions()
	logger := opts.Logger
	if logger == nil {
		logger = defaults.Logger
	}
	minRequestInterval := opts.MinRequestInterval
	if minRequestInterval == nil {
		minRequestInterval = defaults.MinRequestInterval
	}
	enabled := opts.Enabled
	if enabled == nil {
		enabled = defaults.Enabled
	}
	baseCtx := opts.Ctx
	if baseCtx == nil {
		baseCtx = defaults.Ctx
	}

	logger = logger.WithGroup("ratelimit").
		With("enabled", *enabled, "minRequestInterval", *minRequestInterval)

	ctx, cancel := context.WithCancel(baseCtx)

	manager := &Manager{
		enabled:            *enabled,
		minRequestInterval: *minRequestInterval,
		logger:             logger,
		ctx:                ctx,
		cancelCtxFunc:      cancel,
	}

	logger.Debug("ratelimit.Manager initialized")
	return manager
}

func (m *Manager) GetInfo() *Info {
	info := m.rateLimitInfo.Load()
	if info == nil {
		return nil
	}
	return &Info{
		Limit:     info.Limit,
		Period:    info.Period,
		Remaining: info.Remaining,
		Reset:     info.Reset,
	}
}

// Close gracefully shuts down the rate limit manager
func (m *Manager) Close() {
	m.cancelCtxFunc()
}

func (m *Manager) Allow() bool {
	rateLimitInfo := m.rateLimitInfo.Load()
	logger := m.logger.With("rate_limit_info", rateLimitInfo)
	logger.Debug("ratelimit.Manager.Allow called")

	if !m.enabled {
		logger.Debug("Rate limiting disabled - returning early")
		return true
	}
	if utils.IsDone(m.ctx) {
		logger.Debug("Rate limiting context cancelled - returning early")
		return true
	}

	lastRequestTime := m.lastRequestTime.Load()
	defer m.lastRequestTime.Store(utils.PtrTo(time.Now()))

	if throttlingDelay := m.calculateDelay(lastRequestTime, rateLimitInfo); throttlingDelay > 0 {
		logger.Debug("Throttling request", "delay", throttlingDelay)
		select {
		case <-m.ctx.Done():
			logger.Debug("Rate limiting context cancelled - stopping delay")
			return true
		case <-time.After(throttlingDelay):
		}
	} else {
		logger.Debug("Not throttling")
	}

	return true
}

func (m *Manager) ResponseMiddleware(_ *resty.Client, resp *resty.Response) error {
	m.logger.Debug("ResponseMiddleware called")

	if !m.enabled {
		m.logger.Debug("Rate limiting disabled - response middleware returning early")
		return nil
	}
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

	limit, err := strconv.Atoi(limitHeader)
	if err != nil {
		m.logger.Debug("Invalid RateLimit limit header - ignoring all RateLimit headers", "limit",
			limitHeader, "error", err)
		return nil
	}

	remaining, err := strconv.Atoi(remainingHeader)
	if err != nil {
		m.logger.Debug("Invalid RateLimit remaining header - ignoring all RateLimit headers", "remaining",
			remainingHeader, "error", err)
		return nil
	}

	period, err := strconv.Atoi(periodHeader)
	if err != nil {
		period = 0
		m.logger.Debug("Invalid RateLimit period header - ignoring this header", "period", periodHeader,
			"error", err)
	}

	reset, err := strconv.Atoi(resetHeader)
	if err != nil {
		reset = 0
		m.logger.Debug("Invalid RateLimit reset header - ignoring this header", "reset", resetHeader,
			"error", err)
	}

	rateLimitInfo := &Info{
		Limit:     limit,
		Period:    period,
		Remaining: remaining,
		Reset:     reset,
	}
	oldRateLimitInfo := m.rateLimitInfo.Swap(rateLimitInfo)

	m.logger.Debug("Parsed rate limit info", "new_rate_limit_info", rateLimitInfo,
		"old_rate_limit_info", oldRateLimitInfo)

	return nil
}

func (m *Manager) calculateDelay(lastRequestTime *time.Time, rateLimitInfo *Info) time.Duration {
	if lastRequestTime == nil {
		lastRequestTime = utils.PtrTo(time.Time{})
	}

	// Calculate minimum interval delay
	var minIntervalDelay time.Duration
	if timeSinceLastRequest := time.Since(*lastRequestTime); timeSinceLastRequest < m.minRequestInterval {
		minIntervalDelay = m.minRequestInterval - timeSinceLastRequest
	}

	if !rateLimitInfo.ShouldThrottle() {
		return minIntervalDelay
	}

	if rateLimitInfo.Reset > 0 {
		delay := float64(rateLimitInfo.Reset)
		jitterMultiplier := 1 + rand.Float64()*0.1
		return max(time.Duration(delay*jitterMultiplier*float64(time.Second)), minIntervalDelay)
	}

	return 0
}
