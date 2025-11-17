package cache

import "context"

// Adapter 缓存适配器接口
// 所有缓存适配器必须实现此接口
type Adapter interface {
	Cache

	// Name 返回适配器名称
	Name() string

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}
