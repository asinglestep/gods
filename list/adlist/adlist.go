package adlist

import (
	"github.com/asinglestep/gods/utils"
)

// List List
type List struct {
	head   *Node
	tail   *Node
	length int

	comparator utils.Comparator
}

// NewList NewList
func NewList(comparator utils.Comparator) *List {
	l := &List{}
	l.comparator = comparator

	return l
}

// NewListWithoutComparator NewListWithoutComparator
func NewListWithoutComparator() *List {
	l := &List{}
	l.comparator = nil

	return l
}

// Head 头节点
func (l *List) Head() *Node {
	return l.head
}

// Tail 尾节点
func (l *List) Tail() *Node {
	return l.tail
}

// Length 长度
func (l *List) Length() int {
	return l.length
}

// DeleteNode 从链表上移除节点
func (l *List) DeleteNode(node *Node) {
	if node.prev != nil {
		// 有上一个节点
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		// 有下一个节点
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	l.length--
	node.prev = nil
	node.next = nil
}

// InsertNode 插入节点
func (l *List) InsertNode(oldNode *Node, entry interface{}, after bool) *Node {
	newNode := NewNode(entry)

	if after {
		newNode.prev = oldNode
		newNode.next = oldNode.next
		if l.tail == oldNode {
			l.tail = newNode
		}

		if oldNode.next != nil {
			oldNode.next.prev = newNode
		}

		oldNode.next = newNode
	} else {
		newNode.prev = oldNode.prev
		newNode.next = oldNode
		if l.head == oldNode {
			l.head = newNode
		}

		if oldNode.prev != nil {
			oldNode.prev.next = newNode
		}

		oldNode.prev = newNode
	}

	l.length++
	return newNode
}

// AddNodeToTail 添加节点到链表尾部
func (l *List) AddNodeToTail(entry interface{}) *Node {
	node := NewNode(entry)

	if l.length == 0 {
		l.head = node
		l.tail = node
		node.prev = nil
		node.next = nil
	} else {
		node.prev = l.tail
		node.next = nil
		l.tail.next = node
		l.tail = node
	}

	l.length++
	return node
}

// AddNodeToHead 添加节点到链表头部
func (l *List) AddNodeToHead(entry interface{}) *Node {
	node := NewNode(entry)

	if l.length == 0 {
		l.head = node
		l.tail = node
		node.prev = nil
		node.next = nil
	} else {
		node.prev = nil
		node.next = l.head
		l.head.prev = node
		l.head = node
	}

	l.length++
	return node
}

// SearchNode 查找entry对应的节点
func (l *List) SearchNode(entry interface{}) *Node {
	node := l.head
	if l.comparator == nil {
		return nil
	}

	for node != nil {
		if l.comparator.Compare(node.entry, entry) == 0 {
			return node
		}

		node = node.next
	}

	return nil
}
