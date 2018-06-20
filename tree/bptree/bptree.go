package bptree

import (
	"bytes"
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// iNode iNode
type iNode interface {
	// split 分裂节点
	// @return
	split(*Tree) iNode

	// adjacent 获取相邻节点
	adjacent(*Tree) iNode

	// moveKey 移动当前节点的key到相邻节点中或者将相邻节点的key移到当前节点中
	//
	// @param
	// adj: 相邻节点
	//
	// @return
	// 当前节点的父节点
	moveKey(t *Tree, adj iNode) *TreeNode

	// merge 将相邻节点和当前节点合并
	//
	// @param
	// adj: 相邻节点
	//
	// @return
	// 当前节点的父节点
	merge(t *Tree, adj iNode) iNode

	// setParent 设置父节点
	//
	// @param
	// 父节点
	setParent(*TreeNode)

	// getParent 获取父节点
	//
	// @return
	// 当前节点的父节点
	getParent() *TreeNode

	// findKeyPosition 找到key的位置
	findKeyPosition(comparator utils.Comparator, key interface{}) (pos int, bFound bool)

	// getPosKey 获取pos位置的key
	getPosKey(pos int) interface{}

	// getKeys 获取节点key的数量
	getKeys() int

	// isLeaf 节点是否是叶子节点
	isLeaf() bool

	// isFull 节点是否是满节点
	isFull(int) bool

	// verify 验证节点
	verify(*Tree) bool

	// print 打印节点
	print() string

	// dot 生成dot
	//
	// @param
	// dotName: 当前节点dot name
	// pDotName: 父节点dot name
	dot(t *Tree, dotName string, pDotName string) (*dot.Node, *dot.Edge)
}

// Tree Tree
type Tree struct {
	root       iNode // 指向根节点
	comparator utils.Comparator
	maxKeys    int
	minKeys    int
	size       int
}

// NewTree NewTree
//
// @param
// t: 最小度数
func NewTree(t int, comparator utils.Comparator) *Tree {
	tree := &Tree{}
	tree.root = NewTreeLeaf()
	tree.comparator = comparator
	tree.maxKeys = 2 * t
	tree.minKeys = t

	return tree
}

// Insert 插入
func (t *Tree) Insert(key, val interface{}) {
	keyPos := 0
	iNode := t.root

	for {
		// 满节点进行分裂
		if iNode.isFull(t.maxKeys) {
			iNode = iNode.split(t)
		}

		keyPos, _ = iNode.findKeyPosition(t.comparator, key)
		// 是叶子节点，退出
		if iNode.isLeaf() {
			break
		}

		node := iNode.(*TreeNode)
		iNode = node.getChildrenAndUpdateFirstKeyIfNeed(keyPos, key)
	}

	// 插入到叶子节点
	leaf := iNode.(*TreeLeaf)
	// 插入新的key
	leaf.insertEntry(utils.NewEntry(key, val), keyPos)
	t.size++
}

// Delete 删除
func (t *Tree) Delete(key interface{}) {
	t.deleteKey(key)
}

// deleteKey 删除key
func (t *Tree) deleteKey(key interface{}) {
	iNode := t.root

	// 找到key在节点的位置
	pos, bFound := iNode.findKeyPosition(t.comparator, key)
	if !bFound {
		// 查找的key不存在
		// 根节点的最小key是所有关键字的最小值
		return
	}

	for !iNode.isLeaf() {
		// 是根节点或者内节点继续在子节点中查找
		node := iNode.(*TreeNode)
		if bFound {
			// 在当前节点中找到了要删除的key
			// 找到pos对应的子节点
			iNode = node.childrens[pos]

			// 之后都选择第一个子节点
			// 当前节点包含要删除的key，则要删除的key一定是childrens[pos]的第一个key
			for !iNode.isLeaf() {
				node := iNode.(*TreeNode)
				iNode = node.childrens[0]
			}

			pos = 0
			break
		}

		// 继续在其子节点childrens[pos-1]中查找key的位置
		iNode = node.childrens[pos-1]
		pos, bFound = iNode.findKeyPosition(t.comparator, key)
	}

	if !bFound {
		// 要删除的key不存在
		return
	}

	// 删除key
	leaf := iNode.(*TreeLeaf)
	leaf.entries = append(leaf.entries[:pos], leaf.entries[pos+1:]...)
	t.deleteFixUp(leaf, key)
	t.size--
}

// deleteFixUp 删除修复
//
// @param
// node: 要修复的节点
// key: 删除的key
func (t *Tree) deleteFixUp(node iNode, key interface{}) {
	for {
		parent := node.getParent()
		if parent == nil {
			// 修复根节点
			t.dCaseRoot(node)
			return
		}

		if node.getKeys() >= t.minKeys {
			// 修复节点的key的数量大于等于t
			t.replaceKey(parent, key)
			return
		}

		// 获取相邻节点
		adj := node.adjacent(t)
		if adj.getKeys() > t.minKeys {
			// 相邻节点的key的数量大于t
			parent := node.moveKey(t, adj)
			t.replaceKey(parent, key)
			return
		}

		node = node.merge(t, adj)
		pos, bFound := node.findKeyPosition(t.comparator, key)
		if bFound {
			// node节点中包含delKey
			// 用其子节点的第一个key替换
			n := node.(*TreeNode)
			n.keys[pos] = n.childrens[pos].getPosKey(0)
		}
	}
}

// Search 查找key对应的数据
func (t *Tree) Search(key interface{}) *utils.Entry {
	iNode, pos, bFound := t.lookup(t.root, key)
	if !bFound {
		return nil
	}

	leaf := iNode.(*TreeLeaf)
	return leaf.entries[pos]
}

// SearchRange 查找[min, max]之间的数false
func (t *Tree) SearchRange(min, max interface{}) []*utils.Entry {
	iNode, pos, _ := t.lookup(t.root, min)
	if !iNode.isLeaf() {
		// 取pos位置的子节点
		iNode = t.getPosChildren(iNode, pos)
		pos = 0
	}

	for !iNode.isLeaf() {
		iNode = t.getPosChildren(iNode, 0)
	}

	leaf := iNode.(*TreeLeaf)
	entries := []*utils.Entry{}
	iter := NewIteratorWithLeaf(t, leaf, pos)
	for iter.Next() {
		if t.comparator.Compare(iter.entry.GetKey(), max) == utils.Gt {
			break
		}

		entries = append(entries, iter.entry)
	}

	return entries
}

// dCaseRoot 删除修复 - 修复节点为根节点
func (t *Tree) dCaseRoot(node iNode) {
	if !node.isLeaf() && node.getKeys() == 1 {
		n := node.(*TreeNode)
		t.root = n.childrens[0]
		n.childrens[0].setParent(nil)
	}
}

// replaceKey 从node节点开始向上查找，如果该节点中包含key，则用子节点的第一个key替换
func (t *Tree) replaceKey(node *TreeNode, key interface{}) {
	for node != nil {
		// 修复节点的父节点是否包含要删除的key
		pos, bFound := node.findKeyPosition(t.comparator, key)
		if !bFound {
			// node中不包含要删除的key，退出
			return
		}

		// node节点中包含要删除的key
		// 用子节点的第一个key替换
		node.keys[pos] = node.childrens[pos].getPosKey(0)
		node = node.getParent()
	}
}

// getPosChildren 获取pos位置的子节点
func (t *Tree) getPosChildren(iNode iNode, pos int) iNode {
	node := iNode.(*TreeNode)
	return node.childrens[pos]
}

// getPosEntry 获取pos位置的entry
func (t *Tree) getPosEntry(iNode iNode, pos int) *utils.Entry {
	leaf := iNode.(*TreeLeaf)
	return leaf.entries[pos]
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
// bFound: 是否找到key
func (t *Tree) lookup(sNode iNode, key interface{}) (node iNode, pos int, bFound bool) {
	for {
		pos, bFound = sNode.findKeyPosition(t.comparator, key)
		if bFound {
			if sNode.isLeaf() {
				return sNode, pos, bFound
			}

			sNode = t.getPosChildren(sNode, pos)
			for !sNode.isLeaf() {
				sNode = t.getPosChildren(sNode, 0)
			}

			return sNode, 0, bFound
		}

		if pos == 0 || sNode.isLeaf() {
			// 第一个key比要查找的key大 或者 node是叶子节点的key大于要查找的key
			return sNode, pos, bFound
		}

		sNode = t.getPosChildren(sNode, pos-1)
	}
}

// minimum 中序遍历后，树的最小节点
func (t *Tree) minimum() *TreeLeaf {
	iNode := t.root
	for !iNode.isLeaf() {
		iNode = t.getPosChildren(iNode, 0)
	}

	leaf := iNode.(*TreeLeaf)
	return leaf
}

// maximum 中序遍历后，树的最大节点
func (t *Tree) maximum() *TreeLeaf {
	iNode := t.root
	for !iNode.isLeaf() {
		iNode = t.getPosChildren(iNode, iNode.getKeys()-1)
	}

	leaf := iNode.(*TreeLeaf)
	return leaf
}

// splitRootNode 分裂根结点
//
// @param
// inode: 要分裂的节点
// key: 要分裂的节点的第一个key
func (t *Tree) splitRootNode(inode iNode, key interface{}) *TreeNode {
	parent := NewTreeNode()
	parent.childrens = make([]iNode, 1)
	parent.childrens[0] = inode
	parent.keys = make([]interface{}, 1)
	parent.keys[0] = key
	t.root = parent
	return parent
}

// Verify Verify
func (t *Tree) Verify() bool {
	queue := list.New()
	queue.PushBack(t.root)

	for queue.Len() != 0 {
		e := queue.Remove(queue.Front())
		iNode := e.(iNode)

		// 验证
		if !iNode.verify(t) {
			return false
		}

		// 将根节点和内节点的所有子节点加入到stack
		if !iNode.isLeaf() {
			node := iNode.(*TreeNode)
			for _, v := range node.childrens {
				queue.PushBack(v)
			}
		}
	}

	iter := NewIterator(t)
	keys := make([]interface{}, 0)
	for iter.Next() {
		keys = append(keys, iter.GetKey())
	}

	if len(keys) != t.size {
		fmt.Printf("len(keys) != t.size, len(keys): %v, t.size: %v\n", len(keys), t.size)
		return false
	}

	return true
}

// String String
func (t *Tree) String() string {
	buffer := bytes.Buffer{}
	queue := list.New()
	queue.PushBack(t.root)

	for queue.Len() != 0 {
		e := queue.Remove(queue.Front())
		iNode := e.(iNode)

		// 打印节点
		buffer.WriteString(iNode.print())

		// 将根节点和内节点的所有子节点加入到stack
		if !iNode.isLeaf() {
			node := iNode.(*TreeNode)
			for _, v := range node.childrens {
				queue.PushBack(v)
			}
		}
	}

	return buffer.String()
}

type dotNode struct {
	node     iNode
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

		dNode, dEdge := d.node.dot(t, d.nDotName, d.pDotName)
		dGraph.AddNode(dNode)
		if dEdge != nil {
			dGraph.AddEdge(dEdge)
		}

		// 将根节点和内节点的所有子节点加入到stack
		if !d.node.isLeaf() {
			node := d.node.(*TreeNode)
			for _, v := range node.childrens {
				stack.PushBack(&dotNode{v, fmt.Sprintf("node%d", nameIdx), d.nDotName})
				nameIdx++
			}
		}
	}

	if err := dot.GenerateDotFile("bptree.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "bptree.dot", "-o", "bptree.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "bptree.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
