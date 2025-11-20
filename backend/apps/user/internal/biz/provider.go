package biz

import "github.com/google/wire"

// ProviderSet 业务逻辑层依赖注入
var ProviderSet = wire.NewSet(
	NewUserUseCase,
	NewJWTManager,
)
