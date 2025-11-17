package memory

import "time"

// Config 内存缓存配置
type Config struct {
	// MaxSize 最大内存使用（字节），0表示不限制
	MaxSize int64

	// MaxItems 最大条目数，0表示不限制
	MaxItems int

	// Strategy 淘汰策略："lru", "lfu", "fifo"
	Strategy string

	// CleanupInterval 清理间隔
	CleanupInterval time.Duration
}
