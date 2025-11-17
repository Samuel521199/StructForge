package log

import (
	"context"
	"sync"
	"time"
)

// AsyncWriter 异步写入器
type AsyncWriter struct {
	writer        Writer
	queue         chan *LogEntry
	batchSize     int
	flushInterval time.Duration
	dropOnFull    bool
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.Mutex
	closed        bool
}

// NewAsyncWriter 创建异步写入器
func NewAsyncWriter(writer Writer, config AsyncConfig) *AsyncWriter {
	ctx, cancel := context.WithCancel(context.Background())

	aw := &AsyncWriter{
		writer:        writer,
		queue:         make(chan *LogEntry, config.QueueSize),
		batchSize:     config.BatchSize,
		flushInterval: config.FlushInterval,
		dropOnFull:    config.DropOnFull,
		ctx:           ctx,
		cancel:        cancel,
	}

	// 启动后台goroutine处理队列
	aw.wg.Add(1)
	go aw.processQueue()

	return aw
}

// Write 写入日志（异步）
func (aw *AsyncWriter) Write(entry *LogEntry) error {
	aw.mu.Lock()
	if aw.closed {
		aw.mu.Unlock()
		// 已关闭，直接同步写入
		return aw.writer.Write(entry)
	}
	aw.mu.Unlock()

	select {
	case aw.queue <- entry:
		return nil
	default:
		// 队列满了
		if aw.dropOnFull {
			// 丢弃日志
			return nil
		}
		// 阻塞等待（可能影响性能，但保证不丢失）
		select {
		case aw.queue <- entry:
			return nil
		case <-aw.ctx.Done():
			// 已关闭，直接同步写入
			return aw.writer.Write(entry)
		}
	}
}

// processQueue 处理队列中的日志
func (aw *AsyncWriter) processQueue() {
	defer aw.wg.Done()

	batch := make([]*LogEntry, 0, aw.batchSize)
	ticker := time.NewTicker(aw.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-aw.ctx.Done():
			// 关闭时，处理剩余日志
			aw.flushBatch(batch)
			// 处理队列中剩余的日志
			for {
				select {
				case entry := <-aw.queue:
					batch = append(batch, entry)
					if len(batch) >= aw.batchSize {
						aw.flushBatch(batch)
						batch = batch[:0]
					}
				default:
					// 队列已空，处理最后一批
					if len(batch) > 0 {
						aw.flushBatch(batch)
					}
					return
				}
			}

		case entry := <-aw.queue:
			batch = append(batch, entry)
			// 达到批量大小时立即写入
			if len(batch) >= aw.batchSize {
				aw.flushBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// 定时刷新
			if len(batch) > 0 {
				aw.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// flushBatch 批量写入日志
func (aw *AsyncWriter) flushBatch(batch []*LogEntry) {
	for _, entry := range batch {
		if err := aw.writer.Write(entry); err != nil {
			// 写入失败，记录到stderr（避免循环）
			// 这里可以添加错误统计
		}
		// 写入完成后，将entry放回对象池
		putEntry(entry)
	}
}

// Sync 同步刷新（等待所有日志写入完成）
func (aw *AsyncWriter) Sync() error {
	aw.mu.Lock()
	aw.closed = true
	aw.mu.Unlock()

	// 取消context，停止处理
	aw.cancel()

	// 等待处理完成
	aw.wg.Wait()

	// 同步底层写入器
	return aw.writer.Sync()
}

// QueueSize 返回当前队列大小
func (aw *AsyncWriter) QueueSize() int {
	return len(aw.queue)
}
