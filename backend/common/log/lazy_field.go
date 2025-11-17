package log

import "fmt"

// LazyField 延迟计算字段
// 只有在日志真正需要写入时才计算字段值，避免不必要的计算
type LazyField struct {
	key string
	fn  func() interface{}
	typ FieldType
}

// NewLazyField 创建延迟字段
func NewLazyField(key string, fn func() interface{}) *LazyField {
	return &LazyField{
		key: key,
		fn:  fn,
		typ: FieldTypeObject, // 默认类型
	}
}

// NewLazyString 创建延迟字符串字段
func NewLazyString(key string, fn func() string) *LazyField {
	return &LazyField{
		key: key,
		fn: func() interface{} {
			return fn()
		},
		typ: FieldTypeString,
	}
}

// NewLazyInt 创建延迟整数字段
func NewLazyInt(key string, fn func() int) *LazyField {
	return &LazyField{
		key: key,
		fn: func() interface{} {
			return fn()
		},
		typ: FieldTypeInt,
	}
}

// Key 返回字段键
func (f *LazyField) Key() string {
	return f.key
}

// Value 计算并返回字段值
func (f *LazyField) Value() interface{} {
	if f.fn == nil {
		return nil
	}
	return f.fn()
}

// Type 返回字段类型
func (f *LazyField) Type() FieldType {
	return f.typ
}

// String 返回字段的字符串表示
func (f *LazyField) String() string {
	value := f.Value()
	return fmt.Sprintf("%s[%s]=%v", f.key, f.typ, value)
}
