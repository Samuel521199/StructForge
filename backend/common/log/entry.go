package log

import "time"

// LogEntry 日志条目
type LogEntry struct {
	Timestamp  time.Time
	Level      Level
	Service    string
	ServiceID  string
	InstanceID string
	Message    string
	Fields     []Field
	Caller     string
	Container  []Field
	Context    []Field
}

// Reset 重置LogEntry（用于对象池复用）
func (e *LogEntry) Reset() {
	e.Timestamp = time.Time{}
	e.Level = 0
	e.Service = ""
	e.ServiceID = ""
	e.InstanceID = ""
	e.Message = ""
	e.Fields = nil
	e.Caller = ""
	e.Container = nil
	e.Context = nil
}

