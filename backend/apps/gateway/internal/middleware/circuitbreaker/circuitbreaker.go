package circuitbreaker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"StructForge/backend/common/log"
)

// State 熔断器状态
type State int

const (
	// StateClosed 关闭状态：正常处理请求
	StateClosed State = iota
	// StateOpen 打开状态：拒绝所有请求
	StateOpen
	// StateHalfOpen 半开状态：允许少量请求通过，用于测试服务是否恢复
	StateHalfOpen
)

// String 返回状态字符串
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config 熔断器配置
type Config struct {
	// 失败率阈值（0-1），超过此值将打开熔断器
	FailureThreshold float64
	// 最小请求数，达到此数量后才开始计算失败率
	MinRequests int
	// 时间窗口（秒），在此时间窗口内统计失败率
	WindowSize int
	// 打开状态持续时间（秒），熔断器打开后等待此时间后进入半开状态
	OpenDuration int
	// 半开状态允许的请求数
	HalfOpenRequests int
	// 超时时间（秒），请求超过此时间视为失败
	Timeout int
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		FailureThreshold: 0.5, // 50% 失败率
		MinRequests:      10,  // 至少10个请求
		WindowSize:       60,  // 60秒时间窗口
		OpenDuration:     30,  // 打开状态持续30秒
		HalfOpenRequests: 3,   // 半开状态允许3个请求
		Timeout:          5,   // 5秒超时
	}
}

// RequestResult 请求结果
type RequestResult struct {
	Success  bool
	Error    error
	Duration time.Duration
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	config *Config
	state  State
	mu     sync.RWMutex

	// 统计信息
	failures      int       // 失败次数
	successes     int       // 成功次数
	lastFailure   time.Time // 最后失败时间
	lastStateTime time.Time // 状态切换时间
	halfOpenCount int       // 半开状态请求计数

	// 时间窗口内的请求记录
	requests []time.Time
	results  []RequestResult
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(config *Config) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig()
	}

	return &CircuitBreaker{
		config:        config,
		state:         StateClosed,
		lastStateTime: time.Now(),
		requests:      make([]time.Time, 0),
		results:       make([]RequestResult, 0),
	}
}

// Execute 执行请求（带熔断保护）
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// 检查熔断器状态
	if !cb.allowRequest() {
		return fmt.Errorf("circuit breaker is open")
	}

	// 记录请求开始时间
	startTime := time.Now()

	// 执行请求
	err := fn()

	// 计算请求耗时
	duration := time.Since(startTime)

	// 记录结果
	result := RequestResult{
		Success:  err == nil,
		Error:    err,
		Duration: duration,
	}

	// 检查是否超时
	if cb.config.Timeout > 0 && duration > time.Duration(cb.config.Timeout)*time.Second {
		result.Success = false
		result.Error = fmt.Errorf("request timeout: %v", duration)
	}

	// 更新状态
	cb.recordResult(result)

	return err
}

// allowRequest 检查是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否可以进入半开状态
		if time.Since(cb.lastStateTime) >= time.Duration(cb.config.OpenDuration)*time.Second {
			cb.mu.RUnlock()
			cb.mu.Lock()
			if cb.state == StateOpen && time.Since(cb.lastStateTime) >= time.Duration(cb.config.OpenDuration)*time.Second {
				cb.state = StateHalfOpen
				cb.lastStateTime = time.Now()
				cb.halfOpenCount = 0
			}
			cb.mu.Unlock()
			cb.mu.RLock()
			return cb.state == StateHalfOpen
		}
		return false
	case StateHalfOpen:
		// 半开状态只允许少量请求
		return cb.halfOpenCount < cb.config.HalfOpenRequests
	default:
		return false
	}
}

// recordResult 记录请求结果并更新状态
func (cb *CircuitBreaker) recordResult(result RequestResult) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	// 清理过期记录
	cb.cleanupOldRecords(now)

	// 记录请求
	cb.requests = append(cb.requests, now)
	cb.results = append(cb.results, result)

	// 更新统计
	if result.Success {
		cb.successes++
	} else {
		cb.failures++
		cb.lastFailure = now
	}

	// 根据状态更新逻辑
	switch cb.state {
	case StateClosed:
		// 检查是否需要打开熔断器
		if cb.shouldOpen() {
			cb.state = StateOpen
			cb.lastStateTime = now
			log.Warn(context.Background(), "熔断器已打开",
				log.Float64("failure_rate", cb.getFailureRate()),
				log.Int("failures", cb.failures),
				log.Int("total", len(cb.results)),
			)
		}
	case StateHalfOpen:
		cb.halfOpenCount++
		// 如果请求成功，进入关闭状态
		if result.Success {
			cb.state = StateClosed
			cb.lastStateTime = now
			cb.halfOpenCount = 0
			// 重置统计
			cb.resetStats()
			log.Info(context.Background(), "熔断器已关闭（服务恢复）",
				log.Int("half_open_requests", cb.halfOpenCount),
			)
		} else if cb.halfOpenCount >= cb.config.HalfOpenRequests {
			// 如果半开状态下的请求都失败，重新打开
			cb.state = StateOpen
			cb.lastStateTime = now
			cb.halfOpenCount = 0
			log.Warn(context.Background(), "熔断器重新打开（服务未恢复）",
				log.Int("half_open_requests", cb.halfOpenCount),
			)
		}
	case StateOpen:
		// 打开状态下不记录结果（请求已被拒绝）
	}
}

// shouldOpen 检查是否应该打开熔断器
func (cb *CircuitBreaker) shouldOpen() bool {
	// 如果请求数不足，不打开
	if len(cb.results) < cb.config.MinRequests {
		return false
	}

	// 计算失败率
	failureRate := cb.getFailureRate()
	return failureRate >= cb.config.FailureThreshold
}

// getFailureRate 获取失败率
func (cb *CircuitBreaker) getFailureRate() float64 {
	if len(cb.results) == 0 {
		return 0
	}

	failures := 0
	for _, result := range cb.results {
		if !result.Success {
			failures++
		}
	}

	return float64(failures) / float64(len(cb.results))
}

// cleanupOldRecords 清理过期的请求记录
func (cb *CircuitBreaker) cleanupOldRecords(now time.Time) {
	windowStart := now.Add(-time.Duration(cb.config.WindowSize) * time.Second)

	// 清理过期的请求时间记录
	validRequests := make([]time.Time, 0)
	for _, t := range cb.requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}
	cb.requests = validRequests

	// 清理过期的结果记录
	validResults := make([]RequestResult, 0)
	for i, t := range cb.requests {
		if t.After(windowStart) && i < len(cb.results) {
			validResults = append(validResults, cb.results[i])
		}
	}
	cb.results = validResults
}

// resetStats 重置统计信息
func (cb *CircuitBreaker) resetStats() {
	cb.failures = 0
	cb.successes = 0
	cb.requests = make([]time.Time, 0)
	cb.results = make([]RequestResult, 0)
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats 获取统计信息
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":          cb.state.String(),
		"failures":       cb.failures,
		"successes":      cb.successes,
		"total_requests": len(cb.results),
		"failure_rate":   cb.getFailureRate(),
		"last_failure":   cb.lastFailure.Format(time.RFC3339),
		"state_since":    cb.lastStateTime.Format(time.RFC3339),
	}
}

// CircuitBreakerManager 熔断器管理器
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager 创建熔断器管理器
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetBreaker 获取或创建熔断器
func (m *CircuitBreakerManager) GetBreaker(serviceName string, config *Config) *CircuitBreaker {
	m.mu.RLock()
	breaker, exists := m.breakers[serviceName]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if breaker, exists := m.breakers[serviceName]; exists {
		return breaker
	}

	// 创建新的熔断器
	breaker = NewCircuitBreaker(config)
	m.breakers[serviceName] = breaker

	return breaker
}

// GetBreakerStats 获取所有熔断器的统计信息
func (m *CircuitBreakerManager) GetBreakerStats() map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for serviceName, breaker := range m.breakers {
		stats[serviceName] = breaker.GetStats()
	}

	return stats
}

// IsOpen 检查服务是否处于打开状态
func (m *CircuitBreakerManager) IsOpen(serviceName string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	breaker, exists := m.breakers[serviceName]
	if !exists {
		return false
	}

	return breaker.GetState() == StateOpen
}

// Execute 执行请求（带熔断保护）
func (m *CircuitBreakerManager) Execute(ctx context.Context, serviceName string, config *Config, fn func() error) error {
	breaker := m.GetBreaker(serviceName, config)
	return breaker.Execute(ctx, fn)
}

// IsCircuitBreakerError 检查错误是否为熔断器错误
func IsCircuitBreakerError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "circuit breaker is open")
}
