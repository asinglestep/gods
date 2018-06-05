package utils

// Entry Entry
type Entry struct {
	key   interface{}
	value interface{}
}

// NewEntry NewEntry
func NewEntry(key, val interface{}) *Entry {
	entry := &Entry{}
	entry.key = key
	entry.value = val

	return entry
}

// GetKey 获取key
func (e *Entry) GetKey() interface{} {
	return e.key
}

// GetValue 获取value
func (e *Entry) GetValue() interface{} {
	return e.value
}
