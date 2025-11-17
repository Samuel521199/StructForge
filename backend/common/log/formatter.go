package log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// colorize 为文本添加颜色（在formatter中使用）
func colorize(text string, level Level, enableColor bool) string {
	if !enableColor || !isTerminal() {
		return text
	}
	color, ok := levelColors[level]
	if !ok {
		return text
	}
	return color + text + colorReset
}

// colorizeLevel 为级别添加颜色
func colorizeLevel(level Level, enableColor bool) string {
	levelStr := fmt.Sprintf("[%s]", level.String())
	return colorize(levelStr, level, enableColor)
}

// Formatter 格式化器接口
type Formatter interface {
	Format(entry *LogEntry) ([]byte, error)
}

// TextFormatter 文本格式化器
type TextFormatter struct {
	EnableColor bool
}

// Format 格式化日志为文本
func (f *TextFormatter) Format(entry *LogEntry) ([]byte, error) {
	var buf strings.Builder

	// 时间戳 [2025-01-15 14:30:45.123]
	timestamp := entry.Timestamp.Format("2006-01-15 15:04:05.000")
	buf.WriteString(colorize(fmt.Sprintf("[%s]", timestamp), DebugLevel, f.EnableColor))
	buf.WriteString(" ")

	// 级别 [INFO]
	buf.WriteString(colorizeLevel(entry.Level, f.EnableColor))
	buf.WriteString(" ")

	// 服务名 [gateway]
	if entry.Service != "" {
		serviceStr := fmt.Sprintf("[%s]", entry.Service)
		buf.WriteString(colorize(serviceStr, InfoLevel, f.EnableColor))
		buf.WriteString(" ")
	}

	// 消息
	buf.WriteString(entry.Message)
	buf.WriteString("\n")

	// 字段
	if len(entry.Fields) > 0 || len(entry.Context) > 0 || len(entry.Container) > 0 {
		allFields := append(entry.Container, entry.Context...)
		allFields = append(allFields, entry.Fields...)

		for _, field := range allFields {
			buf.WriteString("  ")
			buf.WriteString(field.String())
			buf.WriteString("\n")
		}
	}

	// 调用者信息
	if entry.Caller != "" {
		buf.WriteString("  ")
		buf.WriteString(colorize(fmt.Sprintf("caller: %s", entry.Caller), DebugLevel, f.EnableColor))
		buf.WriteString("\n")
	}

	return []byte(buf.String()), nil
}

// JSONFormatter JSON格式化器
type JSONFormatter struct{}

// Format 格式化日志为JSON
func (f *JSONFormatter) Format(entry *LogEntry) ([]byte, error) {
	data := make(map[string]interface{})

	// 基础字段
	data["timestamp"] = entry.Timestamp.Format(time.RFC3339Nano)
	data["level"] = entry.Level.String()
	data["message"] = entry.Message

	// 服务信息
	if entry.Service != "" {
		data["service"] = entry.Service
	}
	if entry.ServiceID != "" {
		data["service_id"] = entry.ServiceID
	}
	if entry.InstanceID != "" {
		data["instance_id"] = entry.InstanceID
	}

	// 调用者信息
	if entry.Caller != "" {
		data["caller"] = entry.Caller
	}

	// 合并所有字段
	fields := make(map[string]interface{})
	for _, field := range entry.Container {
		fields[field.Key()] = field.Value()
	}
	for _, field := range entry.Context {
		fields[field.Key()] = field.Value()
	}
	for _, field := range entry.Fields {
		fields[field.Key()] = field.Value()
	}

	if len(fields) > 0 {
		data["fields"] = fields
	}

	return json.Marshal(data)
}

// getCaller 获取调用者信息
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, line)
	}

	// 提取文件名（不含路径）
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]

	// 提取函数名
	funcName := fn.Name()
	parts = strings.Split(funcName, ".")
	if len(parts) > 0 {
		funcName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%s:%d", fileName, funcName, line)
}
