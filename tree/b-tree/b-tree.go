package btree

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
	root       *TreeNode // 指向根结点
	comparator utils.Comparator
	size       int // 节点数
	maxEntry   int // 每个节点最多有多少个entry
	minEntry   int // 非根节点最少有多少个entry
}

// NewTree 新建b树
//
// @param
// t: 度数
// comparator: 比较器
func NewTree(t int, comparator utils.Comparator) *Tree {
	tree := &Tree{}
	tree.maxEntry = 2*t - 1
	tree.minEntry = t - 1
	tree.root = NewLeafNode()
	tree.comparator = comparator

	return tree
}

// Insert 插入
func (t *Tree) Insert(key, val interface{}) {
	keyPos := 0
	node := t.root

	// 找到key插入的叶子节点
	for {
		left, right, midEntry := t.splitNode(node)
		if right != nil && t.comparator.Compare(midEntry.GetKey(), key) == utils.Lt {
			// 分裂后，中间的key比key小，继续在右节点中查找
			node = right
		} else {
			node = left
		}

		// 查找key在节点中的位置
		keyPos = node.findKeyPosition(t.comparator, key)

		// 如果是叶子节点，退出
		if node.isLeaf() {
			break
		}

		// 指向下一个查找到节点
		node = node.childrens[keyPos]
	}

	// 插入新的entry
	node.insertEntry(utils.NewEntry(key, val), keyPos)
	t.size++
}

// Delete 删除
func (t *Tree) Delete(key interface{}) {
	t.deleteKey(key)
}

// deleteKey deleteKey
func (t *Tree) deleteKey(key interface{}) {
	pos := 0
	var node = t.root

	for {
		node, pos = t.lookup(node, key)
		if node == nil {
			return
		}

		if node.isLeaf() {
			// 是叶子节点直接删除
			entry := node.entries[pos]
			node.entries = append(node.entries[:pos], node.entries[pos+1:]...)
			t.deleteFixUp(node, entry.GetKey())
			t.size--
			return
		}

		// 找到当前节点的前驱或者后继节点，删除后继或者前驱节点
		if len(node.childrens[pos+1].entries) > t.minEntry {
			// 右节点至少有t个关键字，找到后继节点，用后继节点的key替换key
			ssNode := node.findSuccessor(pos)
			node.entries[pos] = ssNode.entries[0]
			key = ssNode.entries[0].GetKey()
			node = ssNode
		} else {
			// 找到前驱节点，用前驱节点的key替换key
			preNode := node.findPrecursor(pos)
			node.entries[pos] = preNode.entries[len(preNode.entries)-1]
			key = preNode.entries[len(preNode.entries)-1].GetKey()
			node = preNode
		}
	}
}

// deleteFixUp 删除修复
func (t *Tree) deleteFixUp(node *TreeNode, key interface{}) {
	for {
		parent := node.parent
		if parent == nil {
			t.dCaseRoot(node)
			return
		}

		if len(node.entries) >= t.minEntry {
			// 节点至少有t-1个关键字，不需要修复
			return
		}

		// 找到相邻节点
		adjNode, pos, bBig := t.getAdjacentNode(parent, key)
		if len(adjNode.entries) > t.minEntry {
			// 相邻节点至少有t个关键字
			t.moveEntry(parent, node, adjNode, pos, bBig)
			return
		}

		// 节点和相邻节点都只有t-1个关键字
		key = t.mergeNode(parent, node, adjNode, pos, bBig).GetKey()
		node = parent
	}
}

// Search 查找指定的key对应的Entry
func (t *Tree) Search(key interface{}) *utils.Entry {
	node, pos := t.lookup(t.root, key)
	if node == nil {
		return nil
	}

	return node.entries[pos]
}

// SearchRange 查找key在[min, max]之间的Entry
func (t *Tree) SearchRange(min, max interface{}) []*utils.Entry {
	stack := list.New()
	node := t.root
	entries := []*utils.Entry{}

	for {
		// 将所有包含min的节点加到stack中
		pos := node.findKeyPosition(t.comparator, min)
		if pos < len(node.entries) {
			node.iOffset = pos
			node.bHandle = true
			stack.PushBack(node)
		}

		if node.isLeaf() {
			break
		}

		node = node.childrens[pos]
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node := e.(*TreeNode)

		if node.isLeaf() {
			for _, v := range node.entries[node.iOffset:] {
				if t.comparator.Compare(v.GetKey(), max) == utils.Gt {
					break
				}

				entries = append(entries, v)
			}
		} else {
			if node.iOffset == len(node.entries) {
				node.iOffset = 0
				continue
			}

			if !node.bHandle {
				node.bHandle = true
			} else {
				// iOffset加一
				if t.comparator.Compare(node.entries[node.iOffset].GetKey(), max) == utils.Gt {
					break
				}

				entries = append(entries, node.entries[node.iOffset])
				node.iOffset++
			}

			stack.PushBack(node)
			stack.PushBack(node.childrens[node.iOffset])
		}
	}

	return entries
}

// getAdjacentNode 获取key所在节点的相邻节点
//
// @param
// parent: key所在节点的父节点
// key: key
//
// @return
// adjNode: 相邻节点
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在右侧，false - 相邻节点在左侧
func (t *Tree) getAdjacentNode(parent *TreeNode, key interface{}) (adjNode *TreeNode, pos int, bBig bool) {
	pos = parent.findKeyPosition(t.comparator, key)
	if pos == 0 || (pos != len(parent.entries) && len(parent.childrens[pos-1].entries) <= t.minEntry) {
		// 第一个key 或者 key在parent节点上且左侧相邻节点的关键字小于t个，返回右侧相邻节点
		adjNode = parent.childrens[pos+1]
		// pEntry = parent.entries[pos]
		bBig = true
	} else {
		// key不在parent节点上 或者 左侧相邻节点的关键字至少有t个，返回左侧相邻节点
		adjNode = parent.childrens[pos-1]
		// pEntry = parent.entries[pos-1]
	}

	return adjNode, pos, bBig
}

// moveEntry 将父节点和node、adjNode相关的entry放到node中，将adjNode的entry赋值给父节点，将adjNode的孩子节点移到node中
//
// @param
// parent: node的父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在node的右边，false - 相邻节点在node的左边
func (t *Tree) moveEntry(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int, bBig bool) {
	if bBig {
		t.moveCaseRightAdjacentNode(parent, node, adjNode, pos)
	} else {
		pos-- // 相邻节点在左侧，pos-1 为 父节点和node、adjNode相关的entry 所在的位置
		t.moveCaseLeftAdjacentNode(parent, node, adjNode, pos)
	}
}

// moveCaseRightAdjacent 相邻节点在右侧
//
// @param
// parent: node的父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
func (t *Tree) moveCaseRightAdjacentNode(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int) {
	pEntry := parent.entries[pos]
	// 将adjNode的第一个entry放到parent的pos上
	parent.entries[pos] = adjNode.entries[0]
	//删除相邻节点的key
	adjNode.entries = adjNode.entries[1:]
	// 将pEntry追加到node
	node.entries = append(node.entries, pEntry)

	if !node.isLeaf() {
		// 改变子节点的父节点
		adjNode.childrens[0].parent = node
		// 相邻节点的孩子节点移到node中
		node.childrens = append(node.childrens, adjNode.childrens[0])
		// 删除相邻节点的孩子节点
		adjNode.childrens = adjNode.childrens[1:]
	}
}

// moveCaseLeftAdjacentNode 相邻节点在左侧
//
// @param
// parent: node的父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
func (t *Tree) moveCaseLeftAdjacentNode(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int) {
	pEntry := parent.entries[pos]
	adjNodeEntryLen := len(adjNode.entries)
	// 将adjNode最后一个entry放到parent的pos上
	parent.entries[pos] = adjNode.entries[adjNodeEntryLen-1]
	// 删adjNode的entry
	adjNode.entries = adjNode.entries[:adjNodeEntryLen-1]
	// 将pEntry追加到node
	node.insertEntry(pEntry, 0)

	if !node.isLeaf() {
		adjNodeChildLen := len(adjNode.childrens)
		// 改变子节点的父节点
		adjNode.childrens[adjNodeChildLen-1].parent = node
		// 相邻节点的孩子节点移到node中
		node.insertChildren(adjNode.childrens[adjNodeChildLen-1], 0)
		// 删除相邻节点的孩子节点
		adjNode.childrens = adjNode.childrens[:adjNodeChildLen-1]
	}
}

// mergeNode 将父节点和node、adjNode相关的entry、node 和 adjNode 进行合并
//
// @param
// parent: 父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在node的右边，false - 相邻节点在node的左边
//
// @return
// pEntry: 父节点中和node、adjNode相关联的entry
func (t *Tree) mergeNode(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int, bBig bool) (pEntry *utils.Entry) {
	if bBig {
		pEntry = t.mergeCaseRightAdjacentNode(parent, node, adjNode, pos)
	} else {
		pos-- // 相邻节点在左侧，pos-1 为 父节点和node、adjNode相关的entry 所在的位置
		pEntry = t.mergeCaseLeftAdjacentNode(parent, node, adjNode, pos)
	}

	// 删除父节点key
	parent.entries = append(parent.entries[:pos], parent.entries[pos+1:]...)
	// 删除父节点childrens
	parent.childrens = append(parent.childrens[:pos+1], parent.childrens[pos+2:]...)
	return pEntry
}

// mergeCaseRightAdjacentNode 相邻节点在右侧
//
// @param
// parent: 父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
//
// @retrun
// pEntry: 父节点中和node、adjNode相关联的entry
func (t *Tree) mergeCaseRightAdjacentNode(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int) *utils.Entry {
	// 合并node、adjNode和pEntry
	pEntry := parent.entries[pos]
	node.entries = append(node.entries, pEntry)
	node.entries = append(node.entries, adjNode.entries...)

	if !node.isLeaf() {
		// 合并node、adjNode的子节点
		node.childrens = append(node.childrens, adjNode.childrens...)
		// 改变adjNode的子节点的父节点
		adjNode.updateChildrensParent(node)
	}

	// 删除adjNode
	adjNode.free()
	return pEntry
}

// mergeCaseLeftAdjacentNode 相邻节点在左侧
//
// @param
// parent: 父节点
// node: node
// adjNode: node的相邻节点
// pos: 父节点中第一个大于等于key的位置
//
// @retrun
// pEntry: 父节点中和node、adjNode相关联的entry
func (t *Tree) mergeCaseLeftAdjacentNode(parent *TreeNode, node *TreeNode, adjNode *TreeNode, pos int) *utils.Entry {
	// 合并node、adjNode和pEntry
	pEntry := parent.entries[pos]
	adjNode.entries = append(adjNode.entries, pEntry)
	adjNode.entries = append(adjNode.entries, node.entries...)

	if !adjNode.isLeaf() {
		// 合并node、adjNode的子节点
		adjNode.childrens = append(adjNode.childrens, node.childrens...)
		// 改变node的子节点的父节点
		node.updateChildrensParent(adjNode)
	}

	// 删除node
	node.free()
	return pEntry
}

// splitNode 分裂节点
//
// @param
// node: 要分裂的节点
//
// @return
// left: 分裂之后，midEntry的左节点
// right: 分裂之后，midEntry的右节点
// midEntry: node中间的entry
func (t *Tree) splitNode(node *TreeNode) (left *TreeNode, right *TreeNode, midEntry *utils.Entry) {
	// 是满节点，对该节点进行分裂
	if !node.isFull(t.maxEntry) {
		return node, nil, nil
	}

	parent := node.parent
	if parent == nil {
		// 根节点为满的
		parent = t.splitRootNode(node)
	}

	mid := int(uint(len(node.entries)) >> 1)
	midEntry = node.entries[mid]

	// 找到midEntry.key在父节点中的位置
	pos := parent.findKeyPosition(t.comparator, midEntry.GetKey())
	// 将entry插入到父节点中
	parent.insertEntry(midEntry, pos)

	// midEntry的右节点
	right = &TreeNode{}
	right.parent = parent
	right.entries = make([]*utils.Entry, mid)
	copy(right.entries[:], node.entries[mid+1:])
	if !node.isLeaf() {
		right.childrens = make([]*TreeNode, mid+1)
		copy(right.childrens[:], node.childrens[mid+1:])
		right.updateChildrensParent(right)
	}

	// midEntry的左节点
	node.entries = node.entries[:mid]
	if !node.isLeaf() {
		node.childrens = node.childrens[:mid+1]
		node.updateChildrensParent(node)
	}

	// 将right节点加入到父节点的childrens中
	parent.insertChildren(right, pos+1)
	return node, right, midEntry
}

// splitRootNode 分裂根节点
func (t *Tree) splitRootNode(node *TreeNode) *TreeNode {
	parent := &TreeNode{}
	parent.childrens = make([]*TreeNode, 1)
	parent.childrens[0] = node
	node.parent = parent
	t.root = parent
	return parent
}

// lookup 查找key所在的节点
//
// @param
// sNode: 从sNode开始查找
// key: 要查找的key
//
// @return
// node: key所在的节点
// pos: key在节点中的位置
func (t *Tree) lookup(sNode *TreeNode, key interface{}) (node *TreeNode, pos int) {
	node = sNode

	for {
		pos = node.findKeyPosition(t.comparator, key)
		if pos != len(node.entries) && t.comparator.Compare(node.entries[pos].GetKey(), key) == utils.Et {
			// 找到了
			return node, pos
		}

		if node.isLeaf() {
			// 如果是叶子节点，退出
			return nil, -1
		}

		node = node.childrens[pos]
	}
}

// dCaseRoot 删除修复 - 修复根节点
func (t *Tree) dCaseRoot(node *TreeNode) {
	if len(node.entries) != 0 {
		return
	}

	if node.isLeaf() {
		// 删除之后，变成一个空树
		t.root = node
	} else {
		// 合并时，将根节点唯一的key合并到子节点中
		// 需要修改t.root的指向
		t.root = node.childrens[0]
		node.childrens[0].parent = nil
	}
}

// VerifBTree 验证是否是一个b树
func (t *Tree) VerifBTree() bool {
	entires := make([]*utils.Entry, 0)
	stack := list.New()
	stack.PushBack(t.root)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node := e.(*TreeNode)
		if node.parent != nil {
			// 每个非根节点至少有t-1个关键字
			if len(node.entries) < t.minEntry {
				fmt.Printf("非根节点[%v]的关键字小于%v\n", node.entries, t.minEntry)
				return false
			}
		}

		// 每个节点最多含有2t-1个关键字
		if len(node.entries) > t.maxEntry {
			fmt.Printf("节点[%v]最多关键字大于%v\n", node.entries, t.maxEntry)
			return false
		}

		if node.isLeaf() {
			entires = append(entires, node.entries...)
		} else {
			// 非根的内节点至少有t个子女
			if node.parent != nil && len(node.childrens) < t.minEntry+1 {
				fmt.Printf("非根的内节点[%v]的子女小于%v\n", node.entries, t.minEntry+1)
				return false
			}

			// 每个内节点最多有2t个子女
			if len(node.childrens) > t.maxEntry+1 {
				fmt.Printf("内节点[%v]的子女大于%v\n", node.entries, t.maxEntry+1)
				return false
			}

			// iOffset超过entires的长度跳过
			if node.iOffset == len(node.entries) {
				node.iOffset = 0
				node.bHandle = false
				continue
			}

			if !node.bHandle {
				// 第一次弹出这个node
				node.bHandle = true
			} else {
				// 将entry加入entires，iOffset加一
				entires = append(entires, node.entries[node.iOffset])
				node.iOffset++
			}

			stack.PushBack(node)
			stack.PushBack(node.childrens[node.iOffset])
		}
	}

	// 验证顺序
	for i := 0; i < len(entires)-1; i++ {
		if t.comparator.Compare(entires[i].GetKey(), entires[i+1].GetKey()) == utils.Gt {
			fmt.Printf("Key顺序错误\n")
			return false
		}
	}

	if len(entires) != t.size {
		fmt.Printf("len(entires) != t.size, len(entires) %v, t.size %v", len(entires), t.size)
		return false
	}

	return true
}

// PrintBTree PrintBTree
func (t *Tree) PrintBTree() {
	stack := list.New()
	stack.PushBack(t.root)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node := e.(*TreeNode)

		if node.isLeaf() {
			node.printBTreeNode()
		} else {
			// iOffset超过keys的长度跳过
			if node.iOffset == len(node.entries) {
				node.iOffset = 0
				node.bHandle = false
				continue
			}

			if !node.bHandle {
				// 第一次弹出这个node
				node.printBTreeNode()
				node.bHandle = true
			} else {
				// iOffset加一
				node.iOffset++
			}

			stack.PushBack(node)
			stack.PushBack(node.childrens[node.iOffset])
		}
	}
}

type dotNode struct {
	node     *TreeNode
	nDotName string // 当前节点的dot name
	pDotName string // 父节点的dot name
}

// Dot Dot
func (t *Tree) Dot() error {
	nameIdx := 0
	stack := list.New()
	stack.PushBack(&dotNode{t.root, fmt.Sprintf("node%d", nameIdx), ""})
	nameIdx++

	dGraph := dot.NewGraph()
	dGraph.SetNodeGlobalAttr(map[string]string{
		"height": ".1",
		"shape":  "record",
		"width":  ".1",
	})

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		d := e.(*dotNode)

		dNode, dEdge := d.node.dot(t.comparator, d.nDotName, d.pDotName)
		dGraph.AddNode(dNode)
		if dEdge != nil {
			dGraph.AddEdge(dEdge)
		}

		// 将根节点和内节点的所有子节点加入到stack
		if !d.node.isLeaf() {
			for _, v := range d.node.childrens {
				stack.PushBack(&dotNode{v, fmt.Sprintf("node%d", nameIdx), d.nDotName})
				nameIdx++
			}
		}
	}

	if err := dot.GenerateDotFile("b-tree.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "b-tree.dot", "-o", "b-tree.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "b-tree.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
