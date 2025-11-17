package handler

import (
	"github.com/google/wire"
)

// ProviderSet 是 handler 模块的依赖注入提供者集合
var ProviderSet = wire.NewSet(NewGatewayHandler)
