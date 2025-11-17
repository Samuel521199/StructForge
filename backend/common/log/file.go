package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter 文件写入器
// 注意：lumberjack功能（按大小轮转、压缩）需要安装依赖: go get gopkg.in/natefinch/lumberjack.v2
// 当前实现仅支持按天轮转，如需完整功能请安装lumberjack
type FileWriter struct {
	formatter   Formatter
	level       Level
	config      FileConfig
	service     string
	writer      io.WriteCloser // 使用接口，支持lumberjack或标准文件
	file        *os.File       // 当前文件句柄（标准文件模式）
	mu          sync.Mutex
	currentDate string
}

// NewFileWriter 创建文件写入器
func NewFileWriter(config FileConfig, serviceName string) (*FileWriter, error) {
	var formatter Formatter
	if config.Format == TextFormat {
		formatter = &TextFormatter{EnableColor: false}
	} else {
		formatter = &JSONFormatter{}
	}

	// 确保日志目录存在
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 生成文件路径
	date := time.Now().Format("2006-01-02")
	filePath := fmt.Sprintf(config.Path, serviceName, date)

	// 打开文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	return &FileWriter{
		formatter:   formatter,
		level:       config.Level,
		config:      config,
		service:     serviceName,
		writer:      file,
		file:        file,
		currentDate: date,
	}, nil
}

// Write 写入日志
func (w *FileWriter) Write(entry *LogEntry) error {
	// 级别检查
	if entry.Level < w.level {
		return nil
	}

	// 检查是否需要轮转（按天，日期格式：2006-01-02）
	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		w.mu.Lock()
		if today != w.currentDate {
			// 关闭旧文件
			if w.file != nil {
				w.file.Close()
			}

			// 创建新文件
			filePath := fmt.Sprintf(w.config.Path, w.service, today)
			newFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				w.mu.Unlock()
				return fmt.Errorf("创建新日志文件失败: %w", err)
			}
			w.file = newFile
			w.writer = newFile
			w.currentDate = today
		}
		w.mu.Unlock()
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(data)
	if err == nil {
		_, err = w.writer.Write([]byte("\n"))
	}

	return err
}

// Sync 同步刷新
func (w *FileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}
