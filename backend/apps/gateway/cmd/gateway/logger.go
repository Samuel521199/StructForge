package main

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

// newLogger 创建Kratos Logger
// 注意：这里使用Kratos的标准logger，后续可以替换为我们的日志系统适配器
func newLogger() log.Logger {
	return log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
}
