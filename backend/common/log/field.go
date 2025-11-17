package log

import (
	"fmt"
	"time"
)

// FieldType 字段类型代号
type FieldType string

const (
	FieldTypeString   FieldType = "S" // 字符串
	FieldTypeInt      FieldType = "I" // 整数
	FieldTypeBool     FieldType = "B" // 布尔
	FieldTypeFloat    FieldType = "F" // 浮点数
	FieldTypeError    FieldType = "E" // 错误
	FieldTypeObject   FieldType = "O" // 对象
	FieldTypeArray    FieldType = "A" // 数组
	FieldTypeDuration FieldType = "D" // 时长
	FieldTypeTime     FieldType = "T" // 时间
	FieldTypeBytes    FieldType = "X" // 字节数组（十六进制）
)

// Field 日志字段接口
type Field interface {
	Key() string
	Value() interface{}
	Type() FieldType
	String() string
}

// field 字段实现
type field struct {
	key   string
	value interface{}
	typ   FieldType
}

func (f *field) Key() string {
	return f.key
}

func (f *field) Value() interface{} {
	return f.value
}

func (f *field) Type() FieldType {
	return f.typ
}

func (f *field) String() string {
	return fmt.Sprintf("%s[%s]=%v", f.key, f.typ, f.value)
}

// String 创建字符串字段
func String(key, value string) Field {
	return &field{key: key, value: value, typ: FieldTypeString}
}

// Int 创建整数字段
func Int(key string, value int) Field {
	return &field{key: key, value: value, typ: FieldTypeInt}
}

// Int64 创建64位整数字段
func Int64(key string, value int64) Field {
	return &field{key: key, value: value, typ: FieldTypeInt}
}

// Int32 创建32位整数字段
func Int32(key string, value int32) Field {
	return &field{key: key, value: value, typ: FieldTypeInt}
}

// Bool 创建布尔字段
func Bool(key string, value bool) Field {
	return &field{key: key, value: value, typ: FieldTypeBool}
}

// Float64 创建64位浮点数字段
func Float64(key string, value float64) Field {
	return &field{key: key, value: value, typ: FieldTypeFloat}
}

// Float32 创建32位浮点数字段
func Float32(key string, value float32) Field {
	return &field{key: key, value: value, typ: FieldTypeFloat}
}

// ErrorField 创建错误字段
// 注意：由于与日志级别函数Error(ctx, msg, fields...)冲突，使用ErrorField创建错误字段
func ErrorField(err error) Field {
	if err == nil {
		return &field{key: "error", value: nil, typ: FieldTypeError}
	}
	return &field{key: "error", value: err.Error(), typ: FieldTypeError}
}

// Object 创建对象字段
func Object(key string, value interface{}) Field {
	return &field{key: key, value: value, typ: FieldTypeObject}
}

// Duration 创建时长字段
func Duration(key string, d time.Duration) Field {
	return &field{key: key, value: d.String(), typ: FieldTypeDuration}
}

// Time 创建时间字段
func Time(key string, t time.Time) Field {
	return &field{key: key, value: t.Format(time.RFC3339Nano), typ: FieldTypeTime}
}

// Stringer 创建Stringer字段
func Stringer(key string, s fmt.Stringer) Field {
	if s == nil {
		return &field{key: key, value: nil, typ: FieldTypeString}
	}
	return &field{key: key, value: s.String(), typ: FieldTypeString}
}

// Bytes 创建字节数组字段
func Bytes(key string, b []byte) Field {
	return &field{key: key, value: fmt.Sprintf("%x", b), typ: FieldTypeBytes}
}

// StringSlice 创建字符串切片字段
func StringSlice(key string, ss []string) Field {
	return &field{key: key, value: ss, typ: FieldTypeArray}
}

// IntSlice 创建整数切片字段
func IntSlice(key string, is []int) Field {
	return &field{key: key, value: is, typ: FieldTypeArray}
}

// Int64Slice 创建64位整数切片字段
func Int64Slice(key string, is []int64) Field {
	return &field{key: key, value: is, typ: FieldTypeArray}
}

// Map 创建Map字段
func Map(key string, m map[string]interface{}) Field {
	return &field{key: key, value: m, typ: FieldTypeObject}
}

// Any 创建任意类型字段（自动判断类型）
func Any(key string, value interface{}) Field {
	typ := FieldTypeObject
	switch value.(type) {
	case string:
		typ = FieldTypeString
	case int, int8, int16, int32, int64:
		typ = FieldTypeInt
	case uint, uint8, uint16, uint32, uint64:
		typ = FieldTypeInt
	case bool:
		typ = FieldTypeBool
	case float32, float64:
		typ = FieldTypeFloat
	case error:
		typ = FieldTypeError
	case []interface{}:
		typ = FieldTypeArray
	case []string:
		typ = FieldTypeArray
	case []int, []int8, []int16, []int32, []int64:
		typ = FieldTypeArray
	case []uint, []uint16, []uint32, []uint64:
		typ = FieldTypeArray
	case []byte: // []byte 是 []uint8 的别名，需要单独处理
		typ = FieldTypeBytes
	case time.Duration:
		typ = FieldTypeDuration
	case time.Time:
		typ = FieldTypeTime
	case map[string]interface{}:
		typ = FieldTypeObject
	}
	return &field{key: key, value: value, typ: typ}
}
