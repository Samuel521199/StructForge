package server

import (
	"github.com/google/wire"
)

// ProviderSet 是 server 模块的依赖注入提供者集合
var ProviderSet = wire.NewSet(NewHTTPServer)
