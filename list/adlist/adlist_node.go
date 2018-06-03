package adlist

// Node Node
type Node struct {
	prev  *Node
	next  *Node
	entry interface{}
}

// NewNode NewNode
func NewNode(entry interface{}) *Node {
	node := &Node{}
	node.entry = entry

	return node
}

// GetEntry 获取节点entry
func (n *Node) GetEntry() interface{} {
	return n.entry
}

// SetEntry 设置节点的entry
func (n *Node) SetEntry(entry interface{}) {
	n.entry = entry
}
