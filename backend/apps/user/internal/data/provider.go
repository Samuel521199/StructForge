package data

import "github.com/google/wire"

// ProviderSet 数据访问层依赖注入
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewUserProfileRepo,
	NewEmailVerificationRepo,
)

