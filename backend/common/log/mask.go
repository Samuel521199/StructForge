package log

import (
	"strings"
	"sync"
)

// MaskConfig 脱敏配置
type MaskConfig struct {
	Enabled  bool     // 是否启用脱敏
	Fields   []string // 需要脱敏的字段名（如：password, token, secret）
	KeepHead int      // 保留前N个字符（默认3）
	KeepTail int      // 保留后N个字符（默认0）
	MaskChar string   // 脱敏字符（默认*）
}

// masker 脱敏器
type masker struct {
	config MaskConfig
	fields map[string]bool // 快速查找
	mu     sync.RWMutex
}

// newMasker 创建脱敏器
func newMasker(config MaskConfig) *masker {
	m := &masker{
		config: config,
		fields: make(map[string]bool),
	}

	// 构建快速查找map
	for _, field := range config.Fields {
		m.fields[strings.ToLower(field)] = true
	}

	// 设置默认值
	if config.KeepHead == 0 {
		m.config.KeepHead = 3
	}
	if config.KeepTail == 0 {
		m.config.KeepTail = 0
	}
	if config.MaskChar == "" {
		m.config.MaskChar = "*"
	}

	return m
}

// shouldMask 检查字段是否需要脱敏
func (m *masker) shouldMask(key string) bool {
	if !m.config.Enabled {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// 检查字段名（不区分大小写）
	keyLower := strings.ToLower(key)
	if m.fields[keyLower] {
		return true
	}

	// 检查常见敏感字段模式
	sensitivePatterns := []string{
		"password", "pwd", "pass",
		"token", "secret", "key",
		"auth", "credential",
		"idcard", "id_card", "identity",
		"phone", "mobile", "tel",
		"email", "mail",
		"card", "credit", "bank",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(keyLower, pattern) {
			return true
		}
	}

	return false
}

// maskValue 脱敏值
func (m *masker) maskValue(value interface{}) interface{} {
	if !m.config.Enabled {
		return value
	}

	// 只对字符串类型脱敏
	str, ok := value.(string)
	if !ok {
		return value
	}

	return m.maskString(str)
}

// maskString 脱敏字符串
func (m *masker) maskString(s string) string {
	if !m.config.Enabled || s == "" {
		return s
	}

	length := len(s)

	// 如果长度太短，全部脱敏
	if length <= m.config.KeepHead+m.config.KeepTail {
		return strings.Repeat(m.config.MaskChar, length)
	}

	// 保留头部和尾部
	head := s[:m.config.KeepHead]
	tail := ""
	if m.config.KeepTail > 0 && length > m.config.KeepHead {
		tail = s[length-m.config.KeepTail:]
	}

	// 中间部分用脱敏字符替换
	maskLen := length - m.config.KeepHead - m.config.KeepTail
	if maskLen < 0 {
		maskLen = 0
	}

	return head + strings.Repeat(m.config.MaskChar, maskLen) + tail
}

// maskField 脱敏字段值
func (m *masker) maskField(field Field) Field {
	if !m.shouldMask(field.Key()) {
		return field
	}

	// 创建脱敏后的字段
	maskedValue := m.maskValue(field.Value())

	// 根据字段类型创建新的字段
	switch field.Type() {
	case FieldTypeString:
		if str, ok := maskedValue.(string); ok {
			return String(field.Key(), str)
		}
	case FieldTypeInt:
		if i, ok := maskedValue.(int); ok {
			return Int(field.Key(), i)
		}
	case FieldTypeError:
		if str, ok := maskedValue.(string); ok {
			return String("error", str)
		}
	}

	// 默认使用Any
	return Any(field.Key(), maskedValue)
}

// maskFields 脱敏字段列表
func (m *masker) maskFields(fields []Field) []Field {
	if !m.config.Enabled || len(fields) == 0 {
		return fields
	}

	masked := make([]Field, len(fields))
	for i, field := range fields {
		masked[i] = m.maskField(field)
	}
	return masked
}
