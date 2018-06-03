package btree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
)

// Tree Tree
type Tree struct {
	root *TreeNode // 指向根结点

	maxEntryN int // 每个节点最多有多少个entry
	minEntryN int // 非根节点最少有多少个entry
}

// NewTree 新建b树
//
// 参数
// t: 度数
func NewTree(t int) *Tree {
	tree := &Tree{}
	tree.maxEntryN = 2*t - 1
	tree.minEntryN = t - 1
	tree.root = NewLeafNode()

	return tree
}

// Insert 插入
func (t *Tree) Insert(entry *Entry) {
	keyPos := 0
	tmpNode := t.root

	// 找到key插入的叶子节点
	for {
		left, right, midEntry := t.splitNode(tmpNode)
		if right != nil && midEntry.Key.less(entry.Key) {
			// 分裂后，中间的key比key小，继续在右节点中查找
			tmpNode = right
		} else {
			tmpNode = left
		}

		// 查找key在节点中的位置
		keyPos = tmpNode.findKeyPosition(entry.Key)

		// 如果是叶子节点，退出
		if tmpNode.isLeaf() {
			break
		}

		// 指向下一个查找到节点
		tmpNode = tmpNode.childrens[keyPos]
	}

	// 插入新的key
	newEntries := make([]*Entry, len(tmpNode.entries)+1)
	newEntries[keyPos] = entry

	for i, v := range tmpNode.entries {
		if i < keyPos {
			newEntries[i] = v
		} else {
			newEntries[i+1] = v
		}
	}

	tmpNode.entries = newEntries
}

// Delete 删除
func (t *Tree) Delete(key Key) {
	node, entry := t.deleteNode(key)
	if node == nil {
		return
	}

	t.deleteFixUp(node, entry.Key)
}

// deleteNode deleteNode
func (t *Tree) deleteNode(key Key) (fixNode *TreeNode, entry *Entry) {
	tmpNode := t.root

	for {
		pos := tmpNode.findKeyPosition(key)
		switch {
		// 当前node的所有key都小于要查找的key 或者 查找的key在tmpNode.keys[pos]的左边
		case pos == len(tmpNode.entries) || tmpNode.entries[pos].Key.more(key):
			if tmpNode.isLeaf() {
				// 没有找到删除的key
				return nil, nil
			}

			tmpNode = tmpNode.childrens[pos]

		// 找到要删除的节点
		default:
			// 是叶子节点直接删除
			if tmpNode.isLeaf() {
				entry = tmpNode.entries[pos]
				tmpNode.entries = append(tmpNode.entries[:pos], tmpNode.entries[pos+1:]...)
				return tmpNode, entry
			}

			// 找到当前节点的前驱或者后继节点，删除后继或者前驱节点
			switch {
			// 右节点至少有t个关键字，找到前驱节点，找到后继节点，用后继节点的key替换key
			case len(tmpNode.childrens[pos+1].entries) > t.minEntryN:
				ssNode := tmpNode.findSuccessor(pos)
				key = ssNode.entries[0].Key
				tmpNode.entries[pos] = ssNode.entries[0]
				tmpNode = ssNode

			// 找到前驱节点，用前驱节点的key替换key
			default:
				preNode := tmpNode.findPrecursor(pos)
				key = preNode.entries[len(preNode.entries)-1].Key
				tmpNode.entries[pos] = preNode.entries[len(preNode.entries)-1]
				tmpNode = preNode
			}
		}
	}
}

// deleteFixUp 删除修复
func (t *Tree) deleteFixUp(node *TreeNode, key Key) {
	for {
		parent := node.parent
		if parent == nil {
			if len(node.entries) == 0 {
				if node.isLeaf() {
					// 删除之后，变成一个空树
					t.root = node
				} else {
					// 合并时，将根节点唯一的key合并到子节点中，需要修改t.root的指向
					t.root = node.childrens[0]
					node.childrens[0].parent = nil
				}
			}

			return
		}

		if len(node.entries) >= t.minEntryN {
			// 节点至少有t-1个关键字，不需要修复
			return
		}

		// 找到相邻节点
		bNode, pEntry, pos, bBig := t.getAdjacentNode(parent, key)
		if len(bNode.entries) > t.minEntryN {
			// 相邻节点至少有t个关键字
			entry := t.moveKey(parent, node, bNode, pEntry, pos, bBig)
			key = entry.Key
			return
		}

		// 节点和相邻节点都只有t-1个关键字
		entry := t.mergeNode(parent, node, bNode, pEntry, pos, bBig)
		key = entry.Key
		node = parent
	}
}

// Search 查找指定的key对应的Entry
func (t *Tree) Search(key Key) *Entry {
	node := t.root

	for {
		pos := node.findKeyPosition(key)
		switch {
		case pos == len(node.entries) || node.entries[pos].Key.more(key):
			// 如果是叶子节点，退出
			if node.isLeaf() {
				return nil
			}

			node = node.childrens[pos]

		default:
			return node.entries[pos]
		}
	}
}

// SearchRange 查找key在[min, max]之间的Entry
func (t *Tree) SearchRange(min Key, max Key) []*Entry {
	stack := list.New()
	node := t.root
	entries := []*Entry{}

	for {
		// 将所有包含min的节点加到stack中
		pos := node.findKeyPosition(min)
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
				if v.Key.more(max) {
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
				if node.entries[node.iOffset].Key.more(max) {
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

// getAdjacentNode 找到key所在节点的相邻节点
//
// 参数
// parent: key所在节点的父节点
// key: key
//
// 返回值
// bNode: 相邻节点
// pEntry: 父节点中和node、node相邻节点相关联的entry
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在右侧，false - 相邻节点在左侧
func (t *Tree) getAdjacentNode(parent *TreeNode, key Key) (bNode *TreeNode, pEntry *Entry, pos int, bBig bool) {
	pos = parent.findKeyPosition(key)
	switch pos {
	case 0:
		// 第一个子节点的相邻节点为pos+1
		bNode = parent.childrens[1]
		pEntry = parent.entries[0]
		bBig = true

	case len(parent.entries):
		// 最后一个节点的相邻节点为pos-1
		bNode = parent.childrens[pos-1]
		pEntry = parent.entries[pos-1]

	default:
		// node左侧的相邻节点
		lbNode := parent.childrens[pos-1]
		// node右侧的相邻节点
		rbNode := parent.childrens[pos+1]

		if len(lbNode.entries) > t.minEntryN {
			// node左侧的相邻节点至少有t个关键字，返回左侧的相邻节点
			bNode = lbNode
			pEntry = parent.entries[pos-1]
		} else {
			bNode = rbNode
			pEntry = parent.entries[pos]
			bBig = true
		}
	}

	return bNode, pEntry, pos, bBig
}

// moveKey 将相邻节点的key赋值给父节点，将父节点的key追加到node，将相邻节点的孩子节点移到node中
//
// 参数
// parent: 父节点
// node: node
// bNode: node的相邻节点
// pEntry: 父节点中和node、node相邻节点相关联的entry
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在node的右边，false - 相邻节点在node的左边
//
// 返回值
// entry: 移动之后，parent.entries[pos]对应的entry
func (t *Tree) moveKey(parent *TreeNode, node *TreeNode, bNode *TreeNode, pEntry *Entry, pos int, bBig bool) *Entry {
	if bBig {
		// 相邻节点在右侧，pKey的位置为pos
		parent.entries[pos] = bNode.entries[0]
		//删除相邻节点的key
		bNode.entries = bNode.entries[1:]
		// 将pKey追加到node
		node.entries = append(node.entries, pEntry)

		if !node.isLeaf() {
			// 相邻节点的孩子节点移到node中
			node.childrens = append(node.childrens, bNode.childrens[0])
			// 改变子节点的父节点
			bNode.childrens[0].parent = node
			// 删除相邻节点的孩子节点
			bNode.childrens = bNode.childrens[1:]
		}
	} else {
		// 删除相邻节点的key
		bNodeKeysLen := len(bNode.entries)
		// 相邻节点在左侧，pKey的位置为pos--
		pos--
		parent.entries[pos] = bNode.entries[bNodeKeysLen-1]
		bNode.entries = bNode.entries[:bNodeKeysLen-1]

		// 将key追加到node
		newKeys := make([]*Entry, len(node.entries)+1)
		newKeys[0] = pEntry
		copy(newKeys[1:], node.entries)
		node.entries = newKeys

		if !node.isLeaf() {
			bNodeCsLen := len(bNode.childrens)
			// 改变子节点的父节点
			bNode.childrens[bNodeCsLen-1].parent = node

			// 相邻节点的孩子节点移到node中
			newCs := make([]*TreeNode, len(node.childrens)+1)
			newCs[0] = bNode.childrens[bNodeCsLen-1]
			copy(newCs[1:], node.childrens)
			node.childrens = newCs

			// 删除相邻节点的孩子节点
			bNode.childrens = bNode.childrens[:bNodeCsLen-1]
		}
	}

	return parent.entries[pos]
}

// mergeNode 将node、bNode和key进行合并
//
// 参数
// parent: 父节点
// node: node
// bNode: node的相邻节点
// pEntry: 父节点中和node、node相邻节点相关联的entry
// pos: 父节点中第一个大于等于key的位置
// bBig: true - 相邻节点在node的右边，false - 相邻节点在node的左边
//
// 返回值
// entry: pEntry
func (t *Tree) mergeNode(parent *TreeNode, node *TreeNode, bNode *TreeNode, pEntry *Entry, pos int, bBig bool) *Entry {
	if bBig {
		// 合并node、bNode和key
		node.entries = append(node.entries, pEntry)
		node.entries = append(node.entries, bNode.entries...)

		if !node.isLeaf() {
			// 合并node、bNode的子节点
			node.childrens = append(node.childrens, bNode.childrens...)

			// 改变子节点的父节点
			for i := range bNode.childrens {
				bNode.childrens[i].parent = node
			}
		}

		// 删除bNode
		bNode.parent = nil
		bNode.childrens = nil
		bNode.entries = nil

		// 删除父节点key，相邻节点在右侧，pos为pKey所在的位置
		parent.entries = append(parent.entries[:pos], parent.entries[pos+1:]...)
		// 删除父节点childrens
		parent.childrens = append(parent.childrens[:pos+1], parent.childrens[pos+2:]...)
	} else {
		// 合并node、bNode和key
		bNode.entries = append(bNode.entries, pEntry)
		bNode.entries = append(bNode.entries, node.entries...)

		if !bNode.isLeaf() {
			// 合并node、bNode的子节点
			bNode.childrens = append(bNode.childrens, node.childrens...)

			// 改变子节点的父节点
			for i := range node.childrens {
				node.childrens[i].parent = bNode
			}
		}

		// 删除node
		node.parent = nil
		node.childrens = nil
		node.entries = nil

		// 删除父节点key，相邻节点在左侧，pos-1为pKey所在的位置
		parent.entries = append(parent.entries[:pos-1], parent.entries[pos:]...)
		// 删除父节点childrens
		parent.childrens = append(parent.childrens[:pos], parent.childrens[pos+1:]...)
	}

	return pEntry
}

// splitNode 分裂节点
//
// 参数
// node: 要分裂的节点
//
// 返回值
// left: 分裂之后的左节点
// right: 分裂之后的右节点
// midEntry: 中间的entry
func (t *Tree) splitNode(node *TreeNode) (left *TreeNode, right *TreeNode, midEntry *Entry) {
	// 是满节点，对该节点进行分裂
	if !node.isFull(t.maxEntryN) {
		return node, nil, nil
	}

	parent := node.parent
	if parent == nil {
		// 根节点为满的
		parent = &TreeNode{}
		parent.childrens = make([]*TreeNode, 1)
		parent.childrens[0] = node
		node.parent = parent
		t.root = parent
	}

	mid := int(uint(len(node.entries)) >> 1)
	midEntry = node.entries[mid]

	// 找到key在父节点中的位置
	pos := parent.findKeyPosition(midEntry.Key)

	// 将key插入到父节点中
	newKeys := make([]*Entry, len(parent.entries)+1)
	newKeys[pos] = midEntry

	for i, v := range parent.entries {
		if i < pos {
			newKeys[i] = v
		} else {
			newKeys[i+1] = v
		}
	}

	parent.entries = newKeys

	// 分裂成2个新节点
	// midKey的右节点
	right = &TreeNode{}
	right.parent = parent
	right.entries = make([]*Entry, mid)

	for i := range right.entries {
		right.entries[i] = node.entries[mid+i+1]
	}

	if !node.isLeaf() {
		right.childrens = make([]*TreeNode, mid+1)

		for i := range right.childrens {
			right.childrens[i] = node.childrens[mid+i+1]
			right.childrens[i].parent = right
		}
	}

	// midKey的左节点
	node.entries = node.entries[:mid]
	if !node.isLeaf() {
		node.childrens = node.childrens[:mid+1]

		for i := range node.childrens {
			node.childrens[i].parent = node
		}
	}

	// 将新分裂的2个节点加入到父节点中
	newCs := make([]*TreeNode, len(parent.childrens)+1)
	newCs[pos+1] = right

	for i, v := range parent.childrens {
		if i <= pos {
			newCs[i] = v
		} else {
			newCs[i+1] = v
		}
	}

	parent.childrens = newCs
	return node, right, midEntry
}

// VerifBTree 验证是否是一个b树
func (t *Tree) VerifBTree() bool {
	entires := make([]*Entry, 0)
	stack := list.New()
	stack.PushBack(t.root)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node := e.(*TreeNode)
		if node.parent != nil {
			// 每个非根节点至少有t-1个关键字
			if len(node.entries) < t.minEntryN {
				fmt.Printf("非根节点[%v]的关键字小于%v\n", node.entries, t.minEntryN)
				return false
			}
		}

		// 每个节点最多含有2t-1个关键字
		if len(node.entries) > t.maxEntryN {
			fmt.Printf("节点[%v]最多关键字大于%v\n", node.entries, t.maxEntryN)
			return false
		}

		if node.isLeaf() {
			entires = append(entires, node.entries...)
		} else {
			// 非根的内节点至少有t个子女
			if node.parent != nil && len(node.childrens) < t.minEntryN+1 {
				fmt.Printf("非根的内节点[%v]的子女小于%v\n", node.entries, t.minEntryN+1)
				return false
			}

			// 每个内节点最多有2t个子女
			if len(node.childrens) > t.maxEntryN+1 {
				fmt.Printf("内节点[%v]的子女大于%v\n", node.entries, t.maxEntryN+1)
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
		if entires[i].Key > entires[i+1].Key {
			fmt.Printf("b树顺序错误\n")
			return false
		}
	}

	// fmt.Printf("node count %v\n", len(entires))

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
			node.PrintBTreeNode()
		} else {
			// iOffset超过keys的长度跳过
			if node.iOffset == len(node.entries) {
				node.iOffset = 0
				node.bHandle = false
				continue
			}

			if !node.bHandle {
				// 第一次弹出这个node
				node.PrintBTreeNode()
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

		dNode, dEdge := d.node.dot(d.nDotName, d.pDotName)
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
