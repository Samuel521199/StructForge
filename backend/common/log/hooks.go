package log

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// AlertHook 告警钩子
type AlertHook struct {
	webhookURL string
	client     *http.Client
	levels     []Level // 需要告警的级别
}

// NewAlertHook 创建告警钩子
func NewAlertHook(webhookURL string, levels ...Level) *AlertHook {
	if len(levels) == 0 {
		levels = []Level{ErrorLevel, FatalLevel}
	}

	return &AlertHook{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		levels: levels,
	}
}

// BeforeWrite 写入前不做处理
func (h *AlertHook) BeforeWrite(entry *LogEntry) error {
	return nil
}

// AfterWrite 写入后检查是否需要告警
func (h *AlertHook) AfterWrite(entry *LogEntry) error {
	// 检查级别是否需要告警
	needAlert := false
	for _, level := range h.levels {
		if entry.Level == level {
			needAlert = true
			break
		}
	}

	if !needAlert {
		return nil
	}

	// 发送告警（简化实现，实际可以发送到钉钉、企业微信等）
	return h.sendAlert(entry)
}

// sendAlert 发送告警
func (h *AlertHook) sendAlert(entry *LogEntry) error {
	// 这里简化实现，实际应该发送HTTP请求到webhook
	// 示例：发送到钉钉、企业微信、Slack等
	_ = entry
	_ = h.webhookURL
	// TODO: 实现实际的告警发送逻辑
	return nil
}

// StatsHook 统计钩子
type StatsHook struct {
	counters map[Level]int64
	mu       sync.RWMutex
}

// NewStatsHook 创建统计钩子
func NewStatsHook() *StatsHook {
	return &StatsHook{
		counters: make(map[Level]int64),
	}
}

// BeforeWrite 写入前不做处理
func (h *StatsHook) BeforeWrite(entry *LogEntry) error {
	return nil
}

// AfterWrite 写入后更新统计
func (h *StatsHook) AfterWrite(entry *LogEntry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.counters[entry.Level]++
	return nil
}

// GetCount 获取指定级别的日志数量
func (h *StatsHook) GetCount(level Level) int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.counters[level]
}

// GetStats 获取所有统计
func (h *StatsHook) GetStats() map[Level]int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[Level]int64)
	for level, count := range h.counters {
		stats[level] = count
	}
	return stats
}

// Reset 重置统计
func (h *StatsHook) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.counters = make(map[Level]int64)
}

// FilterHook 过滤钩子
type FilterHook struct {
	filter func(*LogEntry) bool
}

// NewFilterHook 创建过滤钩子
func NewFilterHook(filter func(*LogEntry) bool) *FilterHook {
	return &FilterHook{
		filter: filter,
	}
}

// BeforeWrite 写入前过滤
func (h *FilterHook) BeforeWrite(entry *LogEntry) error {
	if h.filter != nil && !h.filter(entry) {
		// 返回错误阻止写入
		return fmt.Errorf("日志被过滤")
	}
	return nil
}

// AfterWrite 写入后不做处理
func (h *FilterHook) AfterWrite(entry *LogEntry) error {
	return nil
}
