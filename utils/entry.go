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

// GetKey GetKey
func (e *Entry) GetKey() interface{} {
	return e.key
}

// GetValue GetValue
func (e *Entry) GetValue() interface{} {
	return e.value
}
