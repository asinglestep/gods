package skiplist

import (
	"github.com/asinglestep/gods/utils"
)

// Level Level
type Level struct {
	forward *Node // 指向后续节点
	span    int   // 跨度
}

// Node Node
type Node struct {
	entry    *utils.Entry // 节点数据
	level    []Level
	backward *Node
}

// NewNode 创建节点
func NewNode(level int, key, val interface{}) *Node {
	node := &Node{}
	node.entry = utils.NewEntry(key, val)
	node.level = make([]Level, level)

	return node
}
