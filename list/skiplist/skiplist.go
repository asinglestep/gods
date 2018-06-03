package skiplist

import (
	"fmt"
	"math/rand"
	"os/exec"
	"runtime"
	"time"

	dot "github.com/asinglestep/godot"
)

const (
	MAX_LEVEL = 32 // 跳跃表最大层数
)

// SkipList SkipList
type SkipList struct {
	head  *SkipListNode
	level int // 当前跳跃表的最大层数
}

// NewSkipList 创建跳跃表
func NewSkipList() *SkipList {
	list := &SkipList{}
	list.level = 1
	list.head = NewSkipListNode(MAX_LEVEL, 0, nil)

	return list
}

// Search 查找
func (sl *SkipList) Search(key Key) *Entry {
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil {
			if x.forward[i].enrty.Key.less(key) {
				x = x.forward[i]
			} else {
				if x.forward[i].enrty.Key.equal(key) {
					return x.forward[i].enrty
				}

				break
			}
		}
	}

	return nil
}

// Insert 插入
func (sl *SkipList) Insert(key Key, val Value) {
	x := sl.head
	level := randomLevel()
	update := make([]*SkipListNode, MAX_LEVEL)

	// 将0到sl.level层中最后一个小于key的节点保存到update中
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].enrty.Key.less(key) {
			x = x.forward[i]
		}

		update[i] = x
	}

	// update的sl.level到level层的前一个节点为head
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}

		sl.level = level
	}

	nNode := NewSkipListNode(level, key, val)
	for i := 0; i < level; i++ {
		nNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = nNode
	}
}

// Delete 删除
func (sl *SkipList) Delete(key Key) {
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil {
			if x.forward[i].enrty.Key.less(key) {
				// 当前节点的key小于要删除的key，取下一个节点
				x = x.forward[i]
			} else {
				// 找到要删除的key，删除这一层包含要删除key的节点
				if x.forward[i].enrty.Key.equal(key) {
					sl.deleteNextNode(x, i)
				}

				break
			}
		}
	}
}

// deleteNextNode 删除level层的node的下一个节点
func (sl *SkipList) deleteNextNode(node *SkipListNode, level int) {
	next := node.forward[level]
	node.forward[level] = node.forward[level].forward[level]
	next.forward[level] = nil
	if sl.head.forward[level] == nil {
		sl.level--
	}
}

// PrintSkipList 打印SkipList
func (sl *SkipList) PrintSkipList() {
	for i := sl.level - 1; i >= 0; i-- {
		x := sl.head
		fmt.Printf("第%d层\n", i)
		for x != nil {
			fmt.Printf("%v\t", x.enrty.Key)
			x = x.forward[i]
		}

		fmt.Println("")
	}
}

// Dot Dot
func (sl *SkipList) Dot() error {
	dGraph := dot.NewGraph()
	dGraph.SetNodeGlobalAttr(map[string]string{
		"height":  ".1",
		"shape":   "record",
		"width":   ".1",
		"rankdir": "LR",
		"rotate":  "90",
	})

	dGraph.SetNodeGlobalAttr(map[string]string{
		"shape": "plaintext",
	})

	x := sl.head
	for x != nil {
		dNode := x.addDotNode(sl.level)
		dGraph.AddNode(dNode)

		x = x.forward[0]
	}

	for i := sl.level - 1; i >= 0; i-- {
		x := sl.head
		for x.forward[i] != nil {
			dEdge := x.addDotEdge(i)
			dGraph.AddEdge(dEdge)

			x = x.forward[i]
		}
	}

	if err := dot.GenerateDotFile("skiplist.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "skiplist.dot", "-o", "skiplist.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "skiplist.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// randomLevel randomLevel
func randomLevel() int {
	level := 1
	for rand.New(rand.NewSource(time.Now().UnixNano())).Float32() < 0.5 {
		level++
	}

	if level > MAX_LEVEL {
		level = MAX_LEVEL
	}

	return level
}
