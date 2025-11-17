package log

import (
	"fmt"
	"os"
	"sync"
)

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	HandleWriteError(err error, entry *LogEntry)
	HandleFormatError(err error, entry *LogEntry)
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	fallbackWriter Writer
	errorCounter   map[string]int64
	mu             sync.RWMutex
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(fallbackWriter Writer) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		fallbackWriter: fallbackWriter,
		errorCounter:   make(map[string]int64),
	}
}

// HandleWriteError 处理写入错误
func (h *DefaultErrorHandler) HandleWriteError(err error, entry *LogEntry) {
	h.mu.Lock()
	h.errorCounter["write_error"]++
	h.mu.Unlock()

	// 如果有备用写入器，尝试写入
	if h.fallbackWriter != nil {
		if writeErr := h.fallbackWriter.Write(entry); writeErr != nil {
			// 备用写入器也失败，输出到stderr
			fmt.Fprintf(os.Stderr, "日志写入失败（包括备用写入器）: %v, 原始错误: %v\n", writeErr, err)
		}
	} else {
		// 没有备用写入器，输出到stderr
		fmt.Fprintf(os.Stderr, "日志写入失败: %v\n", err)
	}
}

// HandleFormatError 处理格式化错误
func (h *DefaultErrorHandler) HandleFormatError(err error, entry *LogEntry) {
	h.mu.Lock()
	h.errorCounter["format_error"]++
	h.mu.Unlock()

	// 格式化错误，输出简化信息到stderr
	fmt.Fprintf(os.Stderr, "日志格式化失败: %v, 消息: %s\n", err, entry.Message)
}

// GetErrorCount 获取错误计数
func (h *DefaultErrorHandler) GetErrorCount(errorType string) int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.errorCounter[errorType]
}

// GetErrorStats 获取所有错误统计
func (h *DefaultErrorHandler) GetErrorStats() map[string]int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[string]int64)
	for k, v := range h.errorCounter {
		stats[k] = v
	}
	return stats
}

// Reset 重置错误统计
func (h *DefaultErrorHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errorCounter = make(map[string]int64)
}
