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
	// Reset is extracted from x-ratelimit-reset (seconds until reset)
	Reset int
}

type Options struct {
	Logger             *slog.Logger
	MinRequestInterval *time.Duration
	Enabled            *bool
	Ctx                context.Context
}

func DefaultOptions() *Options {
	return &Options{
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		MinRequestInterval: utils.PtrTo(50 * time.Millisecond),
		Enabled:            utils.PtrTo(true),
		Ctx:                context.Background(),
	}
}

type Manager struct {
	enabled            bool
	lastRequestTime    atomic.Pointer[time.Time]
	minRequestInterval time.Duration
	logger             *slog.Logger

	info      atomic.Pointer[Info]
	remaining atomic.Int64

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
	info := m.info.Load()
	if info == nil {
		return nil
	}
	return &Info{
		Limit: info.Limit,
		Reset: info.Reset,
	}
}

// Close gracefully shuts down the rate limit manager
func (m *Manager) Close() {
	m.cancelCtxFunc()
}

func (m *Manager) Allow() bool {
	m.logger.Debug("ratelimit.Manager.Allow called")

	if !m.enabled {
		m.logger.Debug("Rate limiting disabled - returning early")
		return true
	}
	if utils.IsDone(m.ctx) {
		m.logger.Debug("Rate limiting context cancelled - returning early")
		return true
	}

	lastRequestTime := m.lastRequestTime.Load()
	defer m.lastRequestTime.Store(utils.PtrTo(time.Now()))

	remaining := m.remaining.Add(-1)
	info := m.GetInfo()

	if throttlingDelay := m.calculateDelay(lastRequestTime, remaining, info); throttlingDelay > 0 {
		m.logger.Debug("Throttling request", "delay", throttlingDelay, "remaining", remaining)
		select {
		case <-m.ctx.Done():
			m.logger.Debug("Rate limiting context cancelled - stopping delay", "remaining", remaining)
			return true
		case <-time.After(throttlingDelay):
		}
	} else {
		m.logger.Debug("Not throttling", "remaining", remaining)
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

	remainingHeader := resp.Header().Get("x-ratelimit-remaining")
	limitHeader := resp.Header().Get("x-ratelimit-limit")
	resetHeader := resp.Header().Get("x-ratelimit-reset")

	if limitHeader+remainingHeader+resetHeader == "" {
		m.logger.Debug("No rate limit headers found or incomplete", slog.Group("headers",
			"limit", limitHeader,
			"remaining", remainingHeader,
			"reset", resetHeader))
		return nil
	}
	m.logger.Debug("Parsing rate limit headers")

	remaining, err := strconv.ParseInt(remainingHeader, 10, 64)
	if err != nil {
		m.logger.Debug("Invalid RateLimit remaining header - ignoring all RateLimit headers", "remaining",
			remainingHeader, "error", err)
		return nil
	}

	limit, err := strconv.Atoi(limitHeader)
	if err != nil {
		m.logger.Debug("Invalid RateLimit limit header - ignoring all RateLimit headers", "limit",
			limitHeader, "error", err)
		return nil
	}

	reset, err := strconv.Atoi(resetHeader)
	if err != nil {
		m.logger.Debug("Invalid RateLimit reset header - ignoring all RateLimit headers", "reset", resetHeader,
			"error", err)
		return nil
	}

	rateLimitInfo := &Info{Limit: limit, Reset: reset}
	oldRateLimitInfo := m.info.Swap(rateLimitInfo)
	m.remaining.Store(remaining)

	m.logger.Debug("Parsed rate limit info", "new_rate_limit_info", rateLimitInfo,
		"old_rate_limit_info", oldRateLimitInfo, "remaining", remaining)

	return nil
}

func (m *Manager) calculateDelay(lastRequestTime *time.Time, remaining int64, info *Info) time.Duration {
	if lastRequestTime == nil {
		lastRequestTime = utils.PtrTo(time.Time{})
	}

	// Calculate minimum interval delay
	var minIntervalDelay time.Duration
	if timeSinceLastRequest := time.Since(*lastRequestTime); timeSinceLastRequest < m.minRequestInterval {
		minIntervalDelay = m.minRequestInterval - timeSinceLastRequest
	}

	if info == nil || info.Limit <= 0 || remaining > 0 {
		return minIntervalDelay
	}

	if info.Reset > 0 {
		delay := float64(info.Reset)
		jitterMultiplier := 1 + rand.Float64()*0.1
		return max(time.Duration(delay*jitterMultiplier*float64(time.Second)), minIntervalDelay)
	}

	return 0
}
