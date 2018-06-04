package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/asinglestep/gods/utils"
)

const (
	MAX_LEVEL = 32 // 跳跃表最大层数
)

// List List
type List struct {
	level  int // 当前跳跃表的最大层数
	length int // 长度
	head   *Node

	comparator utils.Comparator
}

// NewList 创建跳跃表
func NewList(comparator utils.Comparator) *List {
	list := &List{}
	list.level = 1
	list.length = 0
	list.head = NewNode(MAX_LEVEL, nil, nil)
	list.comparator = comparator

	return list
}

// Search 查找
func (l *List) Search(key interface{}) *utils.Entry {
	x := l.head

	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil {
			res := l.comparator.Compare(x.level[i].forward.entry.GetKey(), key)
			if res == utils.Et {
				return x.level[i].forward.entry
			}

			if res == utils.Gt {
				break
			}

			x = x.level[i].forward
		}
	}

	return nil
}

// Insert 插入
func (l *List) Insert(key, val interface{}) {
	x := l.head
	level := l.randomLevel()
	update := make([]*Node, MAX_LEVEL)
	rank := make([]int, MAX_LEVEL)

	// 将第0层到第l.level层中最后一个小于key的节点保存到update中
	for i := l.level - 1; i >= 0; i-- {
		if i == l.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil && l.comparator.Compare(x.level[i].forward.entry.GetKey(), key) == utils.Lt {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}

		update[i] = x
	}

	// update的l.level到level层的前一个节点为head
	if level > l.level {
		for i := l.level; i < level; i++ {
			rank[i] = 0
			update[i] = l.head
			update[i].level[i].span = l.length
		}

		l.level = level
	}

	node := NewNode(level, key, val)
	for i := 0; i < level; i++ {
		// forward
		node.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = node

		// span
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = rank[0] - rank[i] + 1
	}

	for i := level; i < l.level; i++ {
		update[i].level[i].span++
	}

	node.backward = update[0]
	l.length++
}

// Delete 删除
func (l *List) Delete(key interface{}) {
	x := l.head
	update := make([]*Node, MAX_LEVEL)

	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && l.comparator.Compare(x.level[i].forward.entry.GetKey(), key) == utils.Lt {
			// 当前节点的key小于要删除的key，取下一个节点
			x = x.level[i].forward
		}

		// 最后一个小于key的节点保存到update中
		update[i] = x
	}

	dNode := x.level[0].forward
	if l.comparator.Compare(dNode.entry.GetKey(), key) == utils.Et {
		l.deleteNode(update, dNode)
	}
}

func (l *List) deleteNode(update []*Node, dNode *Node) {
	for i := 0; i < l.level; i++ {
		if update[i].level[i].forward == dNode {
			update[i].level[i].span += dNode.level[i].span - 1
			update[i].level[i].forward = dNode.level[i].forward
			dNode.level[i].forward = nil
		} else {
			update[i].level[i].span--
		}

		if l.level > 1 && l.head.level[i].forward == nil {
			l.level--
		}
	}

	if dNode.level[0].forward != nil {
		dNode.level[0].forward.backward = dNode.backward
	}

	l.length--
}

// String String
func (l *List) String() string {
	buffer := bytes.Buffer{}
	for i := l.level - 1; i >= 0; i-- {
		x := l.head
		buffer.WriteString(fmt.Sprintf("第%d层\t", i))

		for x != nil {
			n := 1
			tmp := x

			for j := 0; j < x.level[i].span-1; j++ {
				tmp = tmp.level[0].forward
				n += len(fmt.Sprintf("%v", tmp.entry.GetKey())) // entry
				n++                                             // space
			}

			buffer.WriteString(fmt.Sprintf("%v%v", x.entry.GetKey(), l.genSpace(n)))
			x = x.level[i].forward
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}

// randomLevel randomLevel
func (l *List) randomLevel() int {
	level := 1
	for rand.New(rand.NewSource(time.Now().UnixNano())).Float32() < 0.5 {
		level++
	}

	if level > MAX_LEVEL {
		level = MAX_LEVEL
	}

	return level
}

// genSpace genSpace
func (l *List) genSpace(n int) string {
	buffer := bytes.Buffer{}
	for i := 0; i < n; i++ {
		buffer.WriteString(" ")
	}

	return buffer.String()
}
