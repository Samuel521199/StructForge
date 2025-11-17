package log

import (
	"os"
	"runtime"
)

// 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
	colorBold   = "\033[1m"
)

// 级别颜色映射（延迟初始化，在level.go之后）
var levelColors map[Level]string

// init 初始化颜色映射
func init() {
	levelColors = map[Level]string{
		DebugLevel: colorCyan,            // 青色
		InfoLevel:  colorGreen,           // 绿色
		WarnLevel:  colorYellow,          // 黄色
		ErrorLevel: colorRed,             // 红色
		FatalLevel: colorRed + colorBold, // 红色+粗体
	}
}

// isTerminal 检查是否为终端
func isTerminal() bool {
	if runtime.GOOS == "windows" {
		return false // Windows下简化处理，可通过环境变量控制
	}
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
