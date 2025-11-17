package log

// Hook 钩子接口
type Hook interface {
	// BeforeWrite 在写入日志前调用
	// 可以修改entry或返回错误阻止写入
	BeforeWrite(entry *LogEntry) error

	// AfterWrite 在写入日志后调用
	// 可以用于统计、告警等
	AfterWrite(entry *LogEntry) error
}

// hookManager 钩子管理器
type hookManager struct {
	hooks []Hook
}

// newHookManager 创建钩子管理器
func newHookManager(hooks []Hook) *hookManager {
	return &hookManager{
		hooks: hooks,
	}
}

// beforeWrite 执行所有BeforeWrite钩子
func (hm *hookManager) beforeWrite(entry *LogEntry) error {
	for _, hook := range hm.hooks {
		if err := hook.BeforeWrite(entry); err != nil {
			return err
		}
	}
	return nil
}

// afterWrite 执行所有AfterWrite钩子
func (hm *hookManager) afterWrite(entry *LogEntry) {
	for _, hook := range hm.hooks {
		// AfterWrite的错误不应该阻止日志写入
		if err := hook.AfterWrite(entry); err != nil {
			// 可以记录错误，但不影响日志写入
			// 这里简化处理，实际可以记录到错误统计
		}
	}
}
