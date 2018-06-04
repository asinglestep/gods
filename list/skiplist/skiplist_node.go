package skiplist

// Entry Entry
type Entry struct {
	key   interface{}
	value interface{}
}

// NewEntry 创建Entry
func NewEntry(key, val interface{}) *Entry {
	return &Entry{
		key:   key,
		value: val,
	}
}

// GetKey GetKey
func (e *Entry) GetKey() interface{} {
	return e.key
}

// GetValue GetValue
func (e *Entry) GetValue() interface{} {
	return e.value
}

// Level Level
type Level struct {
	forward *Node // 指向后续节点
	span    int   // 跨度
}

// Node Node
type Node struct {
	entry    *Entry // 节点数据
	level    []Level
	backward *Node
}

// NewNode 创建节点
func NewNode(level int, key, val interface{}) *Node {
	node := &Node{}
	node.entry = NewEntry(key, val)
	node.level = make([]Level, level)

	return node
}
