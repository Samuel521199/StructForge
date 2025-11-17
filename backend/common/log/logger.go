package log

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger 日志接口
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
	Log(ctx context.Context, level Level, msg string, fields ...Field)
	Sync() error
	Shutdown(ctx context.Context) error // 优雅关闭
	With(fields ...Field) Logger
	Enabled(level Level) bool // 检查指定级别是否启用
}

// Writer 写入器接口
type Writer interface {
	Write(entry *LogEntry) error
	Sync() error
}

// logger 日志实现
type logger struct {
	config       Config
	writers      []Writer
	baseFields   []Field
	sampler      *sampler
	masker       *masker
	hookManager  *hookManager
	errorHandler ErrorHandler
	closed       bool
	mu           sync.RWMutex
}

// NewLogger 创建新的Logger
func NewLogger(config Config) (Logger, error) {
	// 设置错误处理器
	errorHandler := config.ErrorHandler
	if errorHandler == nil {
		// 默认使用stderr作为备用输出
		errorHandler = NewDefaultErrorHandler(nil)
	}

	l := &logger{
		config:       config,
		writers:      make([]Writer, 0),
		baseFields:   make([]Field, 0),
		sampler:      &sampler{config: config.Sampling},
		masker:       newMasker(config.Mask),
		hookManager:  newHookManager(config.Hooks),
		errorHandler: errorHandler,
	}

	// 添加服务信息字段
	if config.ServiceName != "" {
		l.baseFields = append(l.baseFields, String("service", config.ServiceName))
	}
	if config.ServiceID != "" {
		l.baseFields = append(l.baseFields, String("service_id", config.ServiceID))
	}
	if config.InstanceID != "" {
		l.baseFields = append(l.baseFields, String("instance_id", config.InstanceID))
	}

	// 创建控制台写入器
	if config.Console.Enabled {
		consoleWriter := NewConsoleWriter(config.Console, config.EnableColor)
		l.writers = append(l.writers, consoleWriter)
	}

	// 创建文件写入器
	if config.File.Enabled {
		fileWriter, err := NewFileWriter(config.File, config.ServiceName)
		if err != nil {
			return nil, fmt.Errorf("创建文件写入器失败: %w", err)
		}

		// 如果启用异步写入，包装为异步写入器
		if config.File.AsyncEnabled && config.Async.Enabled {
			asyncWriter := NewAsyncWriter(fileWriter, config.Async)
			l.writers = append(l.writers, asyncWriter)
		} else {
			l.writers = append(l.writers, fileWriter)
		}
	}

	// 如果没有配置任何输出，默认使用控制台
	if len(l.writers) == 0 {
		consoleWriter := NewConsoleWriter(ConsoleConfig{
			Enabled: true,
			Format:  TextFormat,
			Level:   config.Level,
		}, config.EnableColor)
		l.writers = append(l.writers, consoleWriter)
	}

	return l, nil
}

// log 内部日志方法
func (l *logger) log(ctx context.Context, level Level, msg string, fields []Field) {
	// 检查是否已关闭
	l.mu.RLock()
	if l.closed {
		l.mu.RUnlock()
		// 已关闭，直接同步写入到stderr（确保关键日志不丢失）
		if level >= ErrorLevel {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", level.String(), msg)
		}
		return
	}
	l.mu.RUnlock()

	// 级别检查
	if level < l.config.Level {
		return
	}

	// 采样检查（使用trace_id保证采样一致性）
	traceID := ""
	if ctx != nil {
		if tid := ctx.Value(CtxTraceID); tid != nil {
			if str, ok := tid.(string); ok {
				traceID = str
			}
		}
	}

	if traceID != "" {
		// 使用trace_id进行采样，保证同一请求的日志采样一致
		if !l.sampler.shouldSampleWithKey(level, traceID) {
			return
		}
	} else {
		// 没有trace_id，使用级别采样
		if !l.sampler.shouldSample(level) {
			return
		}
	}

	// 处理延迟字段（计算值）
	resolvedFields := resolveLazyFields(fields)

	// 脱敏处理字段
	maskedFields := l.masker.maskFields(resolvedFields)
	maskedContext := l.masker.maskFields(extractContextFields(ctx))
	maskedContainer := l.masker.maskFields(extractContainerInfo())

	// 从对象池获取LogEntry
	entry := getEntry()
	// 注意：不能立即defer putEntry，因为异步写入器可能会持有entry的引用
	// 只有在同步写入时才立即放回，异步写入由写入器负责

	// 设置日志条目
	entry.Timestamp = time.Now()
	entry.Level = level
	entry.Service = l.config.ServiceName
	entry.ServiceID = l.config.ServiceID
	entry.InstanceID = l.config.InstanceID
	entry.Message = msg
	entry.Fields = maskedFields
	entry.Caller = getCaller(4) // 跳过log方法调用栈
	entry.Container = maskedContainer
	entry.Context = maskedContext

	// 执行BeforeWrite钩子
	if err := l.hookManager.beforeWrite(entry); err != nil {
		// 钩子返回错误，阻止写入，立即放回对象池
		putEntry(entry)
		return
	}

	// 写入所有输出
	// 注意：如果有异步写入器，entry会被放入队列，不能立即放回对象池
	// 异步写入器会在处理完成后负责放回对象池（需要修改AsyncWriter）
	hasAsyncWriter := false
	for _, writer := range l.writers {
		if _, ok := writer.(*AsyncWriter); ok {
			hasAsyncWriter = true
		}
		if err := writer.Write(entry); err != nil {
			// 使用错误处理器处理写入错误
			l.errorHandler.HandleWriteError(err, entry)
		}
	}

	// 如果没有异步写入器，所有写入都是同步的，可以立即放回对象池
	if !hasAsyncWriter {
		// 执行AfterWrite钩子
		l.hookManager.afterWrite(entry)
		// 所有写入完成，放回对象池
		putEntry(entry)
	}
	// 如果有异步写入器，entry会被异步写入器持有，由AsyncWriter负责放回
	// AsyncWriter会在flushBatch中调用putEntry将entry放回对象池
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, DebugLevel, msg, fields)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, InfoLevel, msg, fields)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, WarnLevel, msg, fields)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, ErrorLevel, msg, fields)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, FatalLevel, msg, fields)
	os.Exit(1)
}

func (l *logger) Log(ctx context.Context, level Level, msg string, fields ...Field) {
	l.log(ctx, level, msg, fields)
}

func (l *logger) Sync() error {
	var errs []error
	for _, writer := range l.writers {
		if err := writer.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("同步日志失败: %v", errs)
	}
	return nil
}

// Shutdown 优雅关闭Logger
func (l *logger) Shutdown(ctx context.Context) error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return nil
	}
	l.closed = true
	l.mu.Unlock()

	// 停止接收新日志（通过closed标志）
	// 等待所有写入完成
	done := make(chan error, 1)
	go func() {
		// 同步所有写入器
		done <- l.Sync()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *logger) With(fields ...Field) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make([]Field, len(l.baseFields))
	copy(newFields, l.baseFields)
	newFields = append(newFields, fields...)

	return &logger{
		config:       l.config,
		writers:      l.writers,
		baseFields:   newFields,
		sampler:      l.sampler,
		masker:       l.masker,
		hookManager:  l.hookManager,
		errorHandler: l.errorHandler,
	}
}

// Enabled 检查指定级别是否启用
func (l *logger) Enabled(level Level) bool {
	// 检查全局级别
	return level >= l.config.Level
}

// 全局Logger
var (
	globalLogger Logger
	globalMu     sync.RWMutex
)

// SetGlobalLogger 设置全局Logger
func SetGlobalLogger(l Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalLogger = l
}

// getGlobalLogger 获取全局Logger
func getGlobalLogger() Logger {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalLogger
}

// 全局便捷函数（带级别检查优化）
func Debug(ctx context.Context, msg string, fields ...Field) {
	if l := getGlobalLogger(); l != nil {
		// 提前检查级别，避免不必要的字段构建
		if !l.Enabled(DebugLevel) {
			return
		}
		l.Debug(ctx, msg, fields...)
	}
}

func Info(ctx context.Context, msg string, fields ...Field) {
	if l := getGlobalLogger(); l != nil {
		// 提前检查级别，避免不必要的字段构建
		if !l.Enabled(InfoLevel) {
			return
		}
		l.Info(ctx, msg, fields...)
	}
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	if l := getGlobalLogger(); l != nil {
		// 提前检查级别，避免不必要的字段构建
		if !l.Enabled(WarnLevel) {
			return
		}
		l.Warn(ctx, msg, fields...)
	}
}

func Error(ctx context.Context, msg string, fields ...Field) {
	if l := getGlobalLogger(); l != nil {
		// Error级别通常总是启用，但为了统一性也检查
		if !l.Enabled(ErrorLevel) {
			return
		}
		l.Error(ctx, msg, fields...)
	}
}

func Fatal(ctx context.Context, msg string, fields ...Field) {
	if l := getGlobalLogger(); l != nil {
		// Fatal级别总是启用，但为了统一性也检查
		if !l.Enabled(FatalLevel) {
			return
		}
		l.Fatal(ctx, msg, fields...)
	}
}

func Sync() error {
	if l := getGlobalLogger(); l != nil {
		return l.Sync()
	}
	return nil
}

// Shutdown 优雅关闭全局Logger
func Shutdown(ctx context.Context) error {
	if l := getGlobalLogger(); l != nil {
		return l.Shutdown(ctx)
	}
	return nil
}
