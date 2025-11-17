package log

// resolveLazyFields 解析延迟字段，将LazyField转换为普通Field
func resolveLazyFields(fields []Field) []Field {
	if len(fields) == 0 {
		return fields
	}

	resolved := make([]Field, 0, len(fields))
	for _, field := range fields {
		if lazyField, ok := field.(*LazyField); ok {
			// 延迟字段，计算值并转换为普通字段
			value := lazyField.Value()
			switch lazyField.Type() {
			case FieldTypeString:
				if str, ok := value.(string); ok {
					resolved = append(resolved, String(lazyField.Key(), str))
				} else {
					resolved = append(resolved, Any(lazyField.Key(), value))
				}
			case FieldTypeInt:
				if i, ok := value.(int); ok {
					resolved = append(resolved, Int(lazyField.Key(), i))
				} else if i64, ok := value.(int64); ok {
					resolved = append(resolved, Int64(lazyField.Key(), i64))
				} else {
					resolved = append(resolved, Any(lazyField.Key(), value))
				}
			case FieldTypeFloat:
				if f64, ok := value.(float64); ok {
					resolved = append(resolved, Float64(lazyField.Key(), f64))
				} else {
					resolved = append(resolved, Any(lazyField.Key(), value))
				}
			case FieldTypeBool:
				if b, ok := value.(bool); ok {
					resolved = append(resolved, Bool(lazyField.Key(), b))
				} else {
					resolved = append(resolved, Any(lazyField.Key(), value))
				}
			case FieldTypeError:
				if err, ok := value.(error); ok {
					resolved = append(resolved, ErrorField(err))
				} else if str, ok := value.(string); ok {
					resolved = append(resolved, String(lazyField.Key(), str))
				} else {
					resolved = append(resolved, Any(lazyField.Key(), value))
				}
			default:
				resolved = append(resolved, Any(lazyField.Key(), value))
			}
		} else {
			// 普通字段，直接添加
			resolved = append(resolved, field)
		}
	}
	return resolved
}
