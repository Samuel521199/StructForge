package cache

import "errors"

// 标准缓存错误
var (
	// ErrNotFound 键不存在
	ErrNotFound = errors.New("cache: key not found")

	// ErrExpired 键已过期
	ErrExpired = errors.New("cache: key expired")

	// ErrSerialization 序列化错误
	ErrSerialization = errors.New("cache: serialization error")

	// ErrDeserialization 反序列化错误
	ErrDeserialization = errors.New("cache: deserialization error")

	// ErrConnection 连接错误
	ErrConnection = errors.New("cache: connection error")

	// ErrTimeout 操作超时
	ErrTimeout = errors.New("cache: operation timeout")

	// ErrInvalidKey 无效的键
	ErrInvalidKey = errors.New("cache: invalid key")

	// ErrInvalidValue 无效的值
	ErrInvalidValue = errors.New("cache: invalid value")

	// ErrNotSupported 不支持的操作
	ErrNotSupported = errors.New("cache: operation not supported")

	// ErrAdapterNotInitialized 适配器未初始化
	ErrAdapterNotInitialized = errors.New("cache: adapter not initialized")
)
