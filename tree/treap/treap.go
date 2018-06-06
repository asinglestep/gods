package treap

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// Tree Tree
type Tree struct {
	root       *TreeNode // 根节点
	seed       uint32
	size       int // 节点数
	comparator utils.Comparator
}

// NewTree NewTree
func NewTree(comparator utils.Comparator) *Tree {
	t := &Tree{}
	t.root = SentinelNode
	t.seed = 1
	t.comparator = comparator

	return t
}

// Insert 插入
func (t *Tree) Insert(key, val interface{}) {
	node := NewTreeNode(utils.NewEntry(key, val), t.rand())
	t.insertNode(node)
}

// insertNode 插入节点
func (t *Tree) insertNode(node *TreeNode) {
	next := &t.root
	var parent *TreeNode

	for cur := *next; !cur.isSentinel(); cur = *next {
		res := t.comparator.Compare(cur.GetKey(), node.GetKey())
		if res == utils.Et {
			cur.entry.SetValue(node.GetValue())
			return
		}

		parent = cur
		if res == utils.Lt {
			// 在右子树查找
			next = &cur.right
		} else {
			// 在左子树查找
			next = &cur.left
		}
	}

	*next = node
	node.parent = parent
	t.insertFixUp(node)
	t.size++
}

// insertFixUp 插入修复
func (t *Tree) insertFixUp(node *TreeNode) {
	for node.parent != nil && node.parent.priority > node.priority {
		grandfather := node.parent.parent
		isLeft := node.parent.isLeft()

		if node.isLeft() {
			// 父节点的左节点
			node.parent.rightRotate()
		} else {
			// 父节点的右节点
			node.parent.leftRotate()
		}

		if grandfather == nil {
			t.root = node
		} else if isLeft {
			grandfather.left = node
		} else {
			grandfather.right = node
		}
	}
}

// Delete 删除
func (t *Tree) Delete(key interface{}) {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(node.GetKey(), key)
		if res == utils.Et {
			t.deleteNode(node)
			return
		}

		if res == utils.Lt {
			node = node.right
		} else {
			node = node.left
		}
	}
}

// deleteNode 删除节点
func (t *Tree) deleteNode(node *TreeNode) {
	var tmpNode *TreeNode

	for !node.left.isSentinel() || !node.right.isSentinel() {
		parent := node.parent
		isLeft := node.isLeft()

		if node.left.isSentinel() || (!node.right.isSentinel() && node.left.priority > node.right.priority) {
			// 左节点为空 或者 左右节点都不为空，且左节点的优先级大于右节点的优先级
			tmpNode = node.leftRotate()
		} else {
			// 右节点为空 或者 左右节点都不为空，且右节点的优先级大于左节点的优先级
			tmpNode = node.rightRotate()
		}

		if parent == nil {
			t.root = tmpNode
		} else if isLeft {
			parent.left = tmpNode
		} else {
			parent.right = tmpNode
		}
	}

	// 左右节点都为空，删除节点
	if node.parent == nil {
		// 删除节点为根节点
		t.root = SentinelNode
	} else if node.isLeft() {
		// 删除节点为左节点
		node.parent.left = SentinelNode
	} else {
		// 删除节点为右节点
		node.parent.right = SentinelNode
	}

	node.free()
	t.size--
}

// Search 查找
func (t *Tree) Search(key interface{}) *TreeNode {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(node.GetKey(), key)
		if res == utils.Et {
			return node
		}

		if res == utils.Lt {
			node = node.right
		} else {
			node = node.left
		}
	}

	return nil
}

// minimum 中序遍历后，树的最小节点
func (t *Tree) minimum() *TreeNode {
	return t.root.minimum()
}

// maximum 中序遍历后，树的最大节点
func (t *Tree) maximum() *TreeNode {
	return t.root.maximum()
}

// rand 生成随机数 xorshift
func (t *Tree) rand() uint32 {
	x := t.seed
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	t.seed = x
	return x
}

// VerifTreap 验证是否是treap
func (t *Tree) VerifTreap() bool {
	node := t.root
	entries := make([]*utils.Entry, 0)
	stack := list.New()

	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		// 验证左节点的优先级
		if !node.left.isSentinel() {
			if node.left.priority < node.priority {
				fmt.Printf("节点[%v]的优先级错误, 左节点的优先级[%v]小于当前节点的优先级[%v]\n", node.GetKey(), node.left.priority, node.priority)
				return false
			}
		}

		// 验证右节点的优先级
		if !node.right.isSentinel() {
			if node.right.priority < node.priority {
				fmt.Printf("节点[%v]的优先级错误, 右节点的优先级[%v]小于当前节点的优先级[%v]\n", node.GetKey(), node.right.priority, node.priority)
				return false
			}
		}

		entries = append(entries, node.entry)

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	// 验证顺序
	for i := 0; i < len(entries)-1; i++ {
		if t.comparator.Compare(entries[i].GetKey(), entries[i+1].GetKey()) == utils.Gt {
			fmt.Printf("Key顺序错误\n")
			return false
		}
	}

	if t.size != len(entries) {
		fmt.Printf("t.size != len(entries)\n")
		return false
	}

	return true
}

// Dot Dot
func (t *Tree) Dot() error {
	node := t.root
	if node.isSentinel() {
		return nil
	}

	stack := list.New()

	dGraph := dot.NewGraph()
	dGraph.SetNodeGlobalAttr(map[string]string{
		"height": ".1",
		"shape":  "record",
		"width":  ".1",
	})

	stack.PushBack(node)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		dNode, dEdge := node.dot()
		dGraph.AddNode(dNode)
		if dEdge != nil {
			dGraph.AddEdge(dEdge)
		}

		// 将左右子节点压入stack
		if !node.left.isSentinel() {
			stack.PushBack(node.left)
		}

		if !node.right.isSentinel() {
			stack.PushBack(node.right)
		}
	}

	if err := dot.GenerateDotFile("treap.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "treap.dot", "-o", "treap.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "treap.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
