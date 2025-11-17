package log

import "sync"

// entryPool LogEntry对象池
var entryPool = sync.Pool{
	New: func() interface{} {
		return &LogEntry{}
	},
}

// getEntry 从对象池获取LogEntry
func getEntry() *LogEntry {
	return entryPool.Get().(*LogEntry)
}

// putEntry 将LogEntry放回对象池
func putEntry(entry *LogEntry) {
	// 重置entry
	entry.Reset()
	entryPool.Put(entry)
}
