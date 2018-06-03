package skiplist

import (
	"fmt"
	"strings"

	dot "github.com/asinglestep/godot"
)

// Key Key
type Key int

// less less
// k < key, 返回true
func (k Key) less(key Key) bool {
	if k < key {
		return true
	}

	return false
}

// equal equal
// k = key, 返回true
func (k Key) equal(key Key) bool {
	if k == key {
		return true
	}

	return false
}

// Value Value
type Value interface{}

// Entry Entry
type Entry struct {
	Key   Key
	Value Value
}

// NewEntry 创建Entry
func NewEntry(key Key, val Value) *Entry {
	return &Entry{
		Key:   key,
		Value: val,
	}
}

// SkipListNode SkipListNode
type SkipListNode struct {
	enrty   *Entry          // 节点数据
	forward []*SkipListNode // 指向后续节点
}

// NewSkipListNode 创建节点
func NewSkipListNode(level int, key Key, val Value) *SkipListNode {
	node := &SkipListNode{}
	node.enrty = NewEntry(key, val)
	node.forward = make([]*SkipListNode, level)

	return node
}

// addDotNode addDotNode
func (node *SkipListNode) addDotNode(curMaxLevel int) (dNode *dot.Node) {
	attrValues := make([]string, 0, len(node.forward))

	nodeLevel := len(node.forward)
	if nodeLevel > curMaxLevel {
		nodeLevel = curMaxLevel
	}

	for i := nodeLevel - 1; i >= 0; i-- {
		attrValues = append(attrValues, fmt.Sprintf("<td port=\"%d\"> %d </td>", i, node.enrty.Key))
	}

	attr := `<<table border="0" cellborder="1" cellspacing="0" align="left"><tr>` + strings.Join(attrValues, "</tr><tr>") + "</tr></table>>"

	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("node%d", node.enrty.Key)
	dNode.Attr = map[string]string{
		"label": attr,
	}

	return dNode
}

// addDotEdge addDotEdge
func (node *SkipListNode) addDotEdge(level int) (dEdge *dot.Edge) {
	dEdge = &dot.Edge{}
	dEdge.Src = fmt.Sprintf("node%d", node.enrty.Key)
	dEdge.SrcPort = fmt.Sprintf(":%d", level)
	dEdge.Dst = fmt.Sprintf("node%d", node.forward[level].enrty.Key)
	dEdge.DstPort = fmt.Sprintf(":%d", level)

	return dEdge
}
