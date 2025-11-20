package server

import "github.com/google/wire"

// ProviderSet 服务器依赖注入
var ProviderSet = wire.NewSet(
	NewGRPCServer,
	NewHTTPServer,
)
