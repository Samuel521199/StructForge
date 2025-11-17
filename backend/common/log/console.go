package log

import (
	"io"
	"os"
	"sync"
)

// ConsoleWriter 控制台写入器
type ConsoleWriter struct {
	writer      io.Writer
	formatter   Formatter
	level       Level
	enableColor bool
	mu          sync.Mutex
}

// NewConsoleWriter 创建控制台写入器
func NewConsoleWriter(config ConsoleConfig, enableColor bool) *ConsoleWriter {
	var formatter Formatter
	if config.Format == JSONFormat {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{EnableColor: enableColor}
	}

	return &ConsoleWriter{
		writer:      os.Stdout,
		formatter:   formatter,
		level:       config.Level,
		enableColor: enableColor && config.Format == TextFormat,
	}
}

// Write 写入日志
func (w *ConsoleWriter) Write(entry *LogEntry) error {
	// 级别检查
	if entry.Level < w.level {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	// 错误和致命错误输出到stderr
	if entry.Level >= ErrorLevel {
		_, err = os.Stderr.Write(data)
	} else {
		_, err = w.writer.Write(data)
	}

	return err
}

// Sync 同步刷新
func (w *ConsoleWriter) Sync() error {
	return nil // 控制台不需要同步
}
