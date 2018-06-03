package rbtree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
)

// Tree Tree
type Tree struct {
	root *TreeNode
}

// NewTree NewTree
func NewTree() *Tree {
	t := &Tree{}
	t.root = NewSentinel()

	return t
}

// Insert 插入
func (t *Tree) Insert(key Key) {
	// 插入新节点
	fixNode := t.insertNode(NewTreeNode(key))
	if fixNode == nil {
		return
	}

	// 修复
	t.insertFixUp(fixNode)
	return
}

// insertNode 插入新节点
func (t *Tree) insertNode(newNode *TreeNode) (fixNode *TreeNode) {
	// 哨兵节点
	if t.root.isSentinel() {
		t.root = newNode
		return newNode
	}

	tmpNode := t.root

	for {
		if newNode.key.equal(tmpNode.key) {
			return nil
		}

		if newNode.key.less(tmpNode.key) {
			if tmpNode.left.isSentinel() {
				// 插入到左节点
				tmpNode.left = newNode
				break
			}

			tmpNode = tmpNode.left
		} else {
			if tmpNode.right.isSentinel() {
				// 插入到右节点
				tmpNode.right = newNode
				break
			}

			tmpNode = tmpNode.right
		}
	}

	newNode.parent = tmpNode
	return newNode
}

// insertFixUp 插入节点后进行修复
func (t *Tree) insertFixUp(fixNode *TreeNode) {
	for {
		// 修复节点为根节点，变成黑色节点
		parent := fixNode.parent
		if parent == nil {
			fixNode.color = BLACK
			t.root = fixNode
			return
		}

		// 修复节点的父节点为黑色，不处理
		if parent.isBlack() {
			return
		}

		// 父节点为红色
		// 获取叔叔节点
		uncle := parent.getBrother()
		grandfather := parent.parent
		fIsLeft := fixNode.isLeft()
		fIsRight := fixNode.isRight()
		pIsLeft := parent.isLeft()
		pIsRight := parent.isRight()

		if uncle.isRed() {
			// 叔叔节点也是红色
			parent.color = BLACK
			uncle.color = BLACK
			grandfather.color = RED
			fixNode = grandfather
		} else {
			// 叔叔节点是黑色
			switch {
			// 当前节点为左节点，父节点为左节点
			case fIsLeft && pIsLeft:
				parent.color = BLACK
				grandfatherParent := grandfather.parent
				isLeft := grandfather.isLeft()
				tmpNode := grandfather.rightRotate()

				if grandfatherParent == nil {
					t.root = tmpNode
				} else if isLeft {
					grandfatherParent.left = tmpNode
				} else {
					grandfatherParent.right = tmpNode
				}

				parent.right.color = RED
				return

			// 当前节点为右节点，父节点是右节点
			case fIsRight && pIsRight:
				parent.color = BLACK
				grandfatherParent := grandfather.parent
				isLeft := grandfather.isLeft()
				tmpNode := grandfather.leftRotate()

				if grandfatherParent == nil {
					t.root = tmpNode
				} else if isLeft {
					grandfatherParent.left = tmpNode
				} else {
					grandfatherParent.right = tmpNode
				}

				parent.left.color = RED
				return

			// 当前节点为左节点，父节点为右节点
			case fIsLeft && pIsRight:
				grandfather.right = parent.rightRotate()
				fixNode = parent

			// 当前节点为右节点，父节点为左节点
			case fIsRight && pIsLeft:
				grandfather.left = parent.leftRotate()
				fixNode = parent
			}
		}
	}
}

// Delete 删除
func (t *Tree) Delete(key Key) {
	// 删除节点
	fixNode, delNode := t.deleteNode(key)
	if delNode == nil || delNode.isRed() {
		// 如果没有找到删除的节点或者删除的节点是红色，不需要处理
		return
	}

	// 修复
	t.deleteFixUp(fixNode)
	return
}

// deleteNode 删除节点
// 返回值
// fixNode: 修复节点
// delNode: 删除节点
func (t *Tree) deleteNode(key Key) (fixNode *TreeNode, delNode *TreeNode) {
	delNode = t.root

	for {
		// 没找到
		if delNode.isSentinel() {
			return nil, nil
		}

		if key.equal(delNode.key) {
			break
		}

		if key.less(delNode.key) {
			delNode = delNode.left
		} else {
			delNode = delNode.right
		}
	}

	if !delNode.left.isSentinel() {
		// 有左子节点，找到前驱节点
		findNode := delNode
		delNode = delNode.findPrecursor()
		findNode.key = delNode.key
		fixNode = delNode.left
	} else if !delNode.right.isSentinel() {
		// 有右子节点，找到后继节点
		findNode := delNode
		delNode = delNode.findSuccessor()
		findNode.key = delNode.key
		fixNode = delNode.right
	} else {
		// 哨兵节点
		fixNode = delNode.left
	}

	// 删除节点的父节点
	parent := delNode.parent
	if delNode.isLeft() {
		// 删除节点的孩子节点变成父节点的左节点
		parent.left = fixNode
	} else if delNode.isRight() {
		// 删除节点的孩子节点变成父节点的右节点
		parent.right = fixNode
	} else {
		// 删除根节点
		t.root = fixNode
	}

	// 改变删除节点的孩子节点的父节点
	fixNode.parent = parent

	// 删除节点
	delNode.parent = nil
	delNode.left = nil
	delNode.right = nil
	return fixNode, delNode
}

// deleteFixUp 删除修复
// 参数
// fixNode: 修复节点
func (t *Tree) deleteFixUp(fixNode *TreeNode) {
	var tmpNode *TreeNode

	for {
		parent := fixNode.parent
		// 根节点
		if parent == nil {
			fixNode.color = BLACK
			t.root = fixNode
			return
		}

		brother := fixNode.getBrother()
		grandfather := parent.parent
		bIsRed := brother.isRed()
		blIsRed := brother.left.isRed()
		brIsRed := brother.right.isRed()

		if fixNode.isRed() {
			// 修复节点为红色
			fixNode.color = BLACK
			return
		}

		// 修复节点是黑色
		switch {
		// 修复节点的兄弟节点是红色
		case bIsRed:
			brother.color = BLACK
			parent.color = RED
			grandfather := parent.parent
			pIsLeft := parent.isLeft()

			if fixNode.isLeft() {
				// 修复节点是左节点
				tmpNode = parent.leftRotate()
			} else {
				// 修复节点是右节点
				tmpNode = parent.rightRotate()
			}

			if grandfather == nil {
				t.root = tmpNode
			} else if pIsLeft {
				grandfather.left = tmpNode
			} else {
				grandfather.right = tmpNode
			}

		// 修复节点的兄弟节点是黑色，修复节点的兄弟节点的左节点是黑色，修复节点的兄弟节点的右节点是黑色
		case !bIsRed && !blIsRed && !brIsRed:
			brother.color = RED
			// 继续修复父节点
			fixNode = parent

		// 修复节点的兄弟节点是黑色，修复节点的兄弟节点的左节点是红色，修复节点的兄弟节点的右节点是黑色
		case !bIsRed && blIsRed && !brIsRed:
			pIsLeft := parent.isLeft()

			if fixNode.isLeft() {
				// 修复节点是左节点
				brother.left.color = parent.color
				parent.color = BLACK
				parent.right = brother.rightRotate()
				tmpNode = parent.leftRotate()
			} else {
				// 修复节点是右节点
				brother.color = parent.color
				parent.color = BLACK
				brother.left.color = BLACK
				tmpNode = parent.rightRotate()
			}

			if grandfather == nil {
				t.root = tmpNode
			} else if pIsLeft {
				grandfather.left = tmpNode
			} else {
				grandfather.right = tmpNode
			}

			return

		// 修复节点的兄弟节点是黑色，修复节点的兄弟节点的右节点是红色
		case !bIsRed && brIsRed:
			pIsLeft := parent.isLeft()

			if fixNode.isLeft() {
				// 修复节点是左节点
				brother.color = parent.color
				parent.color = BLACK
				brother.right.color = BLACK
				tmpNode = parent.leftRotate()
			} else {
				// 修复节点是右节点
				brother.right.color = parent.color
				parent.color = BLACK
				parent.left = brother.leftRotate()
				tmpNode = parent.rightRotate()
			}

			if grandfather == nil {
				t.root = tmpNode
			} else if pIsLeft {
				grandfather.left = tmpNode
			} else {
				grandfather.right = tmpNode
			}

			return
		}
	}
}

// Search 查找key指定的节点
func (t *Tree) Search(key Key) *TreeNode {
	node := t.root

	for {
		if node.isSentinel() {
			return nil
		}

		if key.equal(node.key) {
			break
		}

		if key.less(node.key) {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node
}

// SearchRange 查找key在[min, max]之间的节点
func (t *Tree) SearchRange(min Key, max Key) []*TreeNode {
	stack := list.New()
	list := []*TreeNode{}
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key大于等于min的节点加入到stack中
		if !node.key.less(min) {
			stack.PushBack(node)
			node = node.left
		} else {
			node = node.right
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		// node.key大于max，退出
		if node.key.more(max) {
			break
		}

		list = append(list, node)
		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	return list
}

// SearchRangeLowerBoundKeyWithLimit 查找大于等于key的limit个节点
func (t *Tree) SearchRangeLowerBoundKeyWithLimit(key Key, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key大于等于key的节点加入到stack中
		if !node.key.less(key) {
			stack.PushBack(node)
			node = node.left
		} else {
			node = node.right
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		list = append(list, node)
		count++
		if count == limit {
			break
		}

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	return list
}

// SearchRangeUpperBoundKeyWithLimit 找到小于等于key的limit个节点
func (t *Tree) SearchRangeUpperBoundKeyWithLimit(key Key, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key小于等于key的节点加入到stack中
		if !node.key.more(key) {
			stack.PushBack(node)
			node = node.right
		} else {
			node = node.left
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		list = append(list, node)
		count++
		if count == limit {
			break
		}

		node = node.left
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.right
		}
	}

	return reverse(list)
}

// PrintRbTree PrintRbTree
func (t *Tree) PrintRbTree() {
	node := t.root
	stack := list.New()

	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		if node.parent != nil {
			fmt.Printf("节点为: %v, \t节点的颜色为: %v, \t父节点为: %v\n", node.key, node.color, node.parent.key)
		} else {
			fmt.Printf("节点为: %v, \t节点的颜色为: %v, \t此节点为根节点\n", node.key, node.color)
		}

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}
}

// VerifRbTree 验证是否是一个红黑树
func (t *Tree) VerifRbTree() bool {
	node := t.root
	stack := list.New()
	keys := make([]Key, 0)
	KeyMap := make(map[Key]*nodeColorStat)

	// 将树的左节点放到栈中
	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	// 从栈中弹出左节点
	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		if stat, ok := KeyMap[node.key]; !ok { // 如果key不在KeyMap中，将他的key放到keys数组中，如果他的右节点不是哨兵节点，将其放到栈中
			// 根节点为红色，返回false
			if node.parent == nil && node.isRed() {
				return false
			}

			// 节点的颜色和父节点的颜色都为红色，返回false
			if node.parent != nil && node.color == RED && node.color == node.parent.color {
				return false
			}

			KeyMap[node.key] = &nodeColorStat{}
			keys = append(keys, node.key)
			stack.PushBack(node)

			node = node.right

			for !node.isSentinel() {
				stack.PushBack(node)
				node = node.left
			}
		} else { // 如果key存在KeyMap中，更新LeftBlackCount、RightBlackCount、BlackCount
			if node.left.isSentinel() {
				// 左节点是哨兵节点，LeftBlackCount = 1
				stat.leftBlackCount = 1
			} else if node.left.isRed() {
				// 左节点是红色节点，LeftBlackCount = 左节点的BlackCount
				stat.leftBlackCount = KeyMap[node.left.key].blackCount
			} else {
				// LeftBlackCount = 左节点的BlackCount + 1
				stat.leftBlackCount = KeyMap[node.left.key].blackCount + 1
			}

			if node.right.isSentinel() {
				// 右节点是哨兵节点，RightBlackCount = 1
				stat.rightBlackCount = 1
			} else if node.right.isRed() {
				// 右节点是红色节点，RightBlackCount = 右节点的BlackCount
				stat.rightBlackCount = KeyMap[node.right.key].blackCount
			} else {
				// RightBlackCount = 右节点的BlackCount + 1
				stat.rightBlackCount = KeyMap[node.right.key].blackCount + 1
			}

			// 左右子节点的黑色数不等，返回false
			if stat.leftBlackCount != stat.rightBlackCount {
				return false
			}

			stat.blackCount = stat.leftBlackCount
		}
	}

	// 验证顺序
	for i := 0; i < len(keys)-1; i++ {
		if keys[i] > keys[i+1] {
			return false
		}
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

	if err := dot.GenerateDotFile("rbtree.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "rbtree.dot", "-o", "rbtree.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "rbtree.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
