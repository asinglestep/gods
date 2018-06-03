package bptree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
)

// Key Key
type Key int32

// less less
// k < key, 返回true
func (k Key) less(key Key) bool {
	if k < key {
		return true
	}

	return false
}

// more more
// k > key, 返回true
func (k Key) more(key Key) bool {
	if k > key {
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
type Value int32

// Entry Entry
type Entry struct {
	Key   Key
	Value Value
}

// String String
func (e *Entry) String() string {
	return fmt.Sprintf("(%v, %v)  ", e.Key, e.Value)
}

// NewEntry NewEntry
func NewEntry(key Key, value Value) *Entry {
	e := &Entry{}
	e.Key = key
	e.Value = value

	return e
}

// Tree Tree
type Tree struct {
	maxKeyN int
	minKeyN int
	root    iNode // 指向根节点
}

// NewTree NewTree
//
// 参数
// t: 最小度数
func NewTree(t int) *Tree {
	tree := &Tree{}
	tree.maxKeyN = 2 * t
	tree.minKeyN = t
	tree.root = NewTreeLeaf()

	return tree
}

// iNode iNode
type iNode interface {
	// 分裂节点
	// 返回值: 节点分裂后的左节点、节点分裂后的右节点、分裂2个节点的key
	// split(*Tree) (iNode, iNode, Key)
	split(*Tree) iNode

	// 获取相邻节点
	adjacent(*Tree) iNode

	// 移动当前节点的key到相邻节点中或者将相邻节点的key移到当前节点中
	// 参数: 相邻节点
	// 返回值: 该节点的父节点
	moveKey(iNode) *TreeNode

	// 将相邻节点和当前节点合并
	// 参数: 相邻节点
	// 返回值: 该节点的父节点
	merge(iNode) iNode

	// 设置父节点
	// 参数: 父节点
	setParent(*TreeNode)

	// 获取父节点
	// 返回值: 节点的父节点
	getParent() *TreeNode

	// 找到key的位置
	findKeyPosition(Key) int

	// 获取pos位置的key
	getPosKey(pos int) Key

	// 获取节点key的数量
	getKeys() int

	// 节点是否是叶子节点
	isLeaf() bool

	// 节点是否是满节点
	isFull(int) bool

	// 验证节点
	verif(*Tree) bool

	// 打印节点
	print()

	// 生成dot
	// 参数: 当前节点dot name，父节点dot name
	dot(string, string) (*dot.Node, *dot.Edge)
}

// Insert 插入
func (t *Tree) Insert(entry *Entry) {
	keyPos := 0
	iNode := t.root

	for {
		// 满节点进行分裂
		if iNode.isFull(t.maxKeyN) {
			// left, right, midKey := iNode.split(t)
			// if !midKey.more(entry.Key) {
			// 	// 中间的key小于等于entry.Key，继续在右节点查找
			// 	iNode = right
			// } else {
			// 	iNode = left
			// }
			iNode = iNode.split(t)
		}

		keyPos = iNode.findKeyPosition(entry.Key)

		// 是叶子节点，退出
		if iNode.isLeaf() {
			break
		}

		node := iNode.(*TreeNode)
		if keyPos == 0 {
			// 比第一个关键字还小
			// 替换node.keys[0]
			node.keys[0] = entry.Key

			// 在第一个子节点中继续查找
			iNode = node.childrens[0]
		} else {
			// 在其子节点中继续查找
			iNode = node.childrens[keyPos-1]
		}
	}

	// 插入到叶子节点
	leaf := iNode.(*TreeLeaf)

	// 插入新的key
	newEntries := make([]*Entry, len(leaf.entries)+1)
	newEntries[keyPos] = entry

	for i, v := range leaf.entries {
		if i < keyPos {
			newEntries[i] = v
		} else {
			newEntries[i+1] = v
		}
	}

	leaf.entries = newEntries
}

// Delete 删除
func (t *Tree) Delete(key Key) {
	fixNode := t.deleteKey(key)
	if fixNode == nil {
		return
	}

	t.deleteFixUp(fixNode, key)
}

// deleteKey 删除key
func (t *Tree) deleteKey(key Key) (fixNode iNode) {
	iNode := t.root

	// 找到key在节点的位置
	pos := iNode.findKeyPosition(key)
	if iNode.getPosKey(0) > key {
		// 查找的key不存在，根节点的最小key是所有关键字的最小值
		return nil
	}

	for {
		if iNode.isLeaf() {
			// 是叶子节点，退出
			break
		}

		// 是根节点或者内节点继续在子节点中查找
		node := iNode.(*TreeNode)
		// 在当前节点中找到了要删除的key
		if pos < len(node.keys) && node.keys[pos] == key {
			// 找到pos对应的子节点
			iNode = node.childrens[pos]

			// 之后都选择第一个子节点
			// 当前节点包含要删除的key，则要删除的key一定是childrens[pos]的一个key
			// 不需要每次都去查key所在节点的位置
			for !iNode.isLeaf() {
				node := iNode.(*TreeNode)
				iNode = node.childrens[0]
			}

			pos = 0
			break
		}

		// 否则，继续在其子节点childrens[pos-1]中查找key的位置
		iNode = node.childrens[pos-1]
		pos = iNode.findKeyPosition(key)
	}

	if pos == iNode.getKeys() || iNode.getPosKey(pos) != key {
		// 叶子节点不存在这个key
		return nil
	}

	// 删除key
	leaf := iNode.(*TreeLeaf)
	leaf.entries = append(leaf.entries[:pos], leaf.entries[pos+1:]...)
	return leaf
}

// deleteFixUp 删除修复
//
// 参数
// node: 修复节点
// delKey: 删除的key
func (t *Tree) deleteFixUp(node iNode, delKey Key) {
	for {
		parent := node.getParent()
		if parent == nil {
			// 修复根节点
			if !node.isLeaf() && node.getKeys() == 1 {
				n := node.(*TreeNode)
				t.root = n.childrens[0]
				n.childrens[0].setParent(nil)
			}

			return
		}

		if node.getKeys() >= t.minKeyN {
			// 修复节点的key的数量大于等于t
			// 从该节点的父节点开始向上查找，替换所有包含delKey的节点中的key
			t.replaceDelKey(parent, delKey)
			return
		}

		// 修复节点的key的数量为t-1
		// 获取相邻节点
		adj := node.adjacent(t)
		if adj.getKeys() > t.minKeyN {
			// 相邻节点的key的数量大于t
			parent := node.moveKey(adj)
			// 从父节点开始向上查找，替换所有包含delKey的节点中的key
			t.replaceDelKey(parent, delKey)
			return
		}

		// 相邻节点的key的数量为t-1，将2个节点合并
		node = node.merge(adj)
		pos := node.findKeyPosition(delKey)
		if pos < node.getKeys() && node.getPosKey(pos) == delKey {
			// node节点中包含delKey
			// 用其子节点的第一个key替换
			n := node.(*TreeNode)
			n.keys[pos] = n.childrens[pos].getPosKey(0)
		}
	}
}

// replaceDelKey 如果node中包含delKey，则用子节点的第一个key替换
func (t *Tree) replaceDelKey(node *TreeNode, delKey Key) {
	for node != nil {
		// 修复节点的父节点是否包含要删除的key
		pos := node.findKeyPosition(delKey)
		if pos == node.getKeys() || node.getPosKey(pos) != delKey {
			// node中不包含要删除的key，退出
			return
		}

		// node节点中包含要删除的key
		// 用子节点的第一个key替换
		node.keys[pos] = node.childrens[pos].getPosKey(0)
		node = node.getParent()
	}
}

// Search 查找key对应的数据
func (t *Tree) Search(key Key) *Entry {
	iNode := t.root

	for {
		pos := iNode.findKeyPosition(key)
		switch {
		case pos == iNode.getKeys() || iNode.getPosKey(pos).more(key):
			if pos == 0 || iNode.isLeaf() {
				// 第一个key比要查找的key大 或者 node是叶子节点的key大于要查找的key
				return nil
			}

			node := iNode.(*TreeNode)
			iNode = node.childrens[pos-1]

		default:
			if iNode.isLeaf() {
				leaf := iNode.(*TreeLeaf)
				return leaf.entries[pos]
			}

			node := iNode.(*TreeNode)
			iNode = node.childrens[pos]

			for !iNode.isLeaf() {
				node := iNode.(*TreeNode)
				iNode = node.childrens[0]
			}

			leaf := iNode.(*TreeLeaf)
			return leaf.entries[0]
		}
	}
}

// SearchRange 查找[min, max]之间的数据
func (t *Tree) SearchRange(min Key, max Key) []*Entry {
	pos := 0
	iNode := t.root

	for {
		pos = iNode.findKeyPosition(min)
		if pos < iNode.getKeys() && iNode.getPosKey(pos) == min {
			// 找到key的位置
			if !iNode.isLeaf() {
				node := iNode.(*TreeNode)
				iNode = node.childrens[pos]
			}

			for {
				if iNode.isLeaf() {
					break
				}

				node := iNode.(*TreeNode)
				iNode = node.childrens[0]
				pos = 0
			}

			break
		}

		// pos == iNode.getKeys() 或者 iNode.getPosKey(pos) > min
		if pos != 0 {
			pos--
		}

		if iNode.isLeaf() {
			break
		}

		node := iNode.(*TreeNode)
		iNode = node.childrens[pos]
	}

	leaf := iNode.(*TreeLeaf)
	entries := []*Entry{}

	for _, v := range leaf.entries[pos:] {
		if v.Key > max {
			return entries
		}

		entries = append(entries, v)
	}

	for leaf.next != nil {
		leaf = leaf.next
		for _, v := range leaf.entries {
			if v.Key > max {
				return entries
			}

			entries = append(entries, v)
		}
	}

	return entries
}

// VerifBpTree VerifBpTree
func (t *Tree) VerifBpTree() bool {
	stack := list.New()
	stack.PushBack(t.root)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		iNode := e.(iNode)

		// 验证
		if !iNode.verif(t) {
			return false
		}

		// 将根节点和内节点的所有子节点加入到stack
		if !iNode.isLeaf() {
			node := iNode.(*TreeNode)
			for _, v := range node.childrens {
				stack.PushBack(v)
			}
		}
	}

	return true
}

// PrintBpTree PrintBpTree
func (t *Tree) PrintBpTree() {
	stack := list.New()
	stack.PushBack(t.root)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		iNode := e.(iNode)

		// 打印节点
		iNode.print()

		// 将根节点和内节点的所有子节点加入到stack
		if !iNode.isLeaf() {
			node := iNode.(*TreeNode)
			for _, v := range node.childrens {
				stack.PushBack(v)
			}
		}
	}
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

		dNode, dEdge := d.node.dot(d.nDotName, d.pDotName)
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
