package bptree

import (
	"fmt"
	"strings"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// TreeLeaf TreeLeaf
type TreeLeaf struct {
	parent  *TreeNode      // 指向父节点
	entries []*utils.Entry // 数据
	next    *TreeLeaf      // 指向下一个叶子节点的指针
	prev    *TreeLeaf      // 指向上一个叶子节点的指针
}

// NewTreeLeaf NewTreeLeaf
func NewTreeLeaf() *TreeLeaf {
	leaf := &TreeLeaf{}

	return leaf
}

// split 分裂节点
func (leaf *TreeLeaf) split(t *Tree) iNode {
	parent := leaf.parent
	if parent == nil {
		// 根节点为满的
		parent = t.splitRootNode(leaf, leaf.entries[0].GetKey())
	}

	// 满节点时有2t个key，所以中间key的位置为t
	mid := t.minKeys
	midEntry := leaf.entries[mid]

	// 新的右节点，修改父节点，next指针，prev指针，entries
	right := NewTreeLeaf()
	right.parent = parent
	right.next = leaf.next
	right.prev = leaf
	right.entries = leaf.entries[mid:]

	// 新的左节点，修改父节点，next指针，entries
	leaf.parent = parent
	leaf.next = right
	leaf.entries = leaf.entries[:mid]

	// 找到key在父节点中的位置
	pos, _ := parent.findKeyPosition(t.comparator, midEntry.GetKey())
	// 将key插入到父节点中
	parent.insertKey(midEntry.GetKey(), pos)
	// 将right节点加入到父节点中，keys和childrens数量相同
	parent.insertChildren(right, pos)

	return parent
}

// adjacent 获取相邻节点
func (leaf *TreeLeaf) adjacent(t *Tree) (adj iNode) {
	var left, right iNode
	parent := leaf.parent

	pos, bFound := parent.findKeyPosition(t.comparator, leaf.entries[0].GetKey())
	if pos == 0 {
		// 第一个子节点的相邻节点为parent.childrens[pos+1]
		return parent.childrens[1]
	}

	if bFound {
		if pos == len(parent.keys)-1 {
			// 当前节点是最后一个节点
			return parent.childrens[pos-1]
		}

		left = parent.childrens[pos-1]
		right = parent.childrens[pos+1]
	} else {
		if pos == 1 {
			// 当前节点是第一个节点
			return parent.childrens[1]
		}

		if pos == len(parent.keys) {
			// parent的key都小于leaf.entries[0].Key
			// 当前节点的位置为pos-1
			// 相邻节点为parent.childrens[pos-2]
			adj = parent.childrens[pos-2]
		}

		left = parent.childrens[pos-2]
		right = parent.childrens[pos]
	}

	if left.getKeys() > t.minKeys {
		return left
	}

	return right
}

// moveKey 移动相邻节点的key到leaf中
func (leaf *TreeLeaf) moveKey(t *Tree, adj iNode) *TreeNode {
	adjNode := adj.(*TreeLeaf)
	parent := leaf.parent

	if t.comparator.Compare(leaf.entries[0].GetKey(), adjNode.entries[0].GetKey()) == utils.Lt {
		leaf.moveCaseRightAdjacentNode(t, adjNode, parent)
	} else {
		leaf.moveCaseLeftAdjacentNode(t, adjNode, parent)
	}

	return parent
}

// moveCaseRightAdjacentNode 相邻节点在右侧
func (leaf *TreeLeaf) moveCaseRightAdjacentNode(t *Tree, adj *TreeLeaf, parent *TreeNode) {
	// 找到相邻节点在父节点中的位置
	pos, _ := parent.findKeyPosition(t.comparator, adj.entries[0].GetKey())

	// 将相邻节点的第一个key加入到当前节点
	leaf.entries = append(leaf.entries, adj.entries[0])

	// 删除相邻节点的key
	adj.entries = adj.entries[1:]

	// 修改父节点的key
	parent.keys[pos] = adj.entries[0].GetKey()
}

// moveCaseLeftAdjacentNode 相邻节点在左侧
func (leaf *TreeLeaf) moveCaseLeftAdjacentNode(t *Tree, adj *TreeLeaf, parent *TreeNode) {
	// 找到当前节点在父节点中的位置
	pos, bFound := parent.findKeyPosition(t.comparator, leaf.entries[0].GetKey())
	if !bFound {
		// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
		pos--
	}

	// 将相邻节点的最后一个key加入到当前节点
	leaf.insertEntry(adj.entries[len(adj.entries)-1], 0)

	// 删除相邻节点的key
	adj.entries = adj.entries[:len(adj.entries)-1]

	// 修改父节点的key
	parent.keys[pos] = leaf.entries[0].GetKey()
}

// merge 合并相邻节点
func (leaf *TreeLeaf) merge(t *Tree, adj iNode) iNode {
	var key interface{}
	adjNode := adj.(*TreeLeaf)
	parent := adjNode.parent

	if t.comparator.Compare(leaf.entries[0].GetKey(), adjNode.entries[0].GetKey()) == utils.Lt {
		key = adjNode.entries[0].GetKey()
		// 将相邻节点合并到当前节点
		appendLeaf(leaf, adjNode)
	} else {
		// 在相邻节点的右侧，找到当前节点在父节点中的位置
		key = leaf.entries[0].GetKey()
		// 将当前节点合并到相邻节点
		appendLeaf(adjNode, leaf)
	}

	pos, bFound := parent.findKeyPosition(t.comparator, key)
	if !bFound {
		// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
		pos--
	}

	// 删除父节点的key和子节点
	parent.keys = append(parent.keys[:pos], parent.keys[pos+1:]...)
	parent.childrens = append(parent.childrens[:pos], parent.childrens[pos+1:]...)
	return parent
}

// setParent 设置父节点
func (leaf *TreeLeaf) setParent(parent *TreeNode) {
	leaf.parent = parent
}

// getParent 获取父节点
func (leaf *TreeLeaf) getParent() *TreeNode {
	return leaf.parent
}

// findKeyPosition 在节点中查找第一个大于等于key的位置，没有比key大的，则返回node.entries的长度
func (leaf *TreeLeaf) findKeyPosition(comparator utils.Comparator, key interface{}) (pos int, bFound bool) {
	if len(leaf.entries) == 0 {
		return 0, false
	}

	i, j := 0, len(leaf.entries)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(leaf.entries[h].GetKey(), key) == utils.Lt {
			i = h + 1
		} else {
			j = h
		}
	}

	if i == len(leaf.entries) {
		return i, false
	}

	if comparator.Compare(leaf.entries[i].GetKey(), key) == utils.Et {
		bFound = true
	}

	return i, bFound
}

// getPosKey 获取pos位置的key
func (leaf *TreeLeaf) getPosKey(pos int) interface{} {
	return leaf.entries[pos].GetKey()
}

// getKeys 获取key的数量
func (leaf *TreeLeaf) getKeys() int {
	return len(leaf.entries)
}

// isLeaf 是否是叶子节点
func (leaf *TreeLeaf) isLeaf() bool {
	return true
}

// isFull 是否是满节点
func (leaf *TreeLeaf) isFull(max int) bool {
	return len(leaf.entries) == max
}

// free free
func (leaf *TreeLeaf) free() {
	leaf.parent = nil
	leaf.entries = nil
	leaf.prev = nil
	leaf.next = nil
}

// verif 验证
func (leaf *TreeLeaf) verify(t *Tree) bool {
	if leaf.parent != nil {
		// 非根节点至少有t个关键字
		if len(leaf.entries) < t.minKeys {
			fmt.Printf("非根节点少于%v个关键字\n", t.minKeys)
			return false
		}
	}

	// 最多有2t个关键字
	if len(leaf.entries) > t.maxKeys {
		fmt.Printf("节点多于%v个关键字\n", t.maxKeys)
		return false
	}

	return true
}

// String String
func (leaf *TreeLeaf) print() string {
	if leaf.parent != nil {
		return fmt.Sprintf("叶子节点key: %v, \t此节点为父节点key[%v]的子节点, \t父节点为[%v]\n", leaf.printKeys(), leaf.entries[0].GetKey(), leaf.parent.printKeys())
	}

	if len(leaf.entries) == 0 {
		return fmt.Sprintf("此b树为一个空树\n")
	}

	return fmt.Sprintf("叶子节点key: %v, \t此节点为根节点\n", leaf.printKeys())
}

func (leaf *TreeLeaf) printKeys() string {
	keys := make([]string, 0, len(leaf.entries))
	for _, entry := range leaf.entries {
		keys = append(keys, fmt.Sprintf("%v", entry.GetKey()))
	}

	return strings.Join(keys, ",")
}

// dot dot
func (leaf *TreeLeaf) dot(t *Tree, nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加一个node
	attrValues := make([]string, 0, len(leaf.entries))

	for i, v := range leaf.entries {
		attrValues = append(attrValues, fmt.Sprintf("<f%d> %d ", i, v.GetKey()))
	}

	attr := "\"" + strings.Join(attrValues, "|") + "\""

	dNode = &dot.Node{}
	dNode.Name = nName
	dNode.Attr = map[string]string{
		"label": attr,
	}

	// 添加一个edge
	if leaf.parent != nil {
		pos, _ := leaf.parent.findKeyPosition(t.comparator, leaf.entries[0].GetKey())
		dEdge = &dot.Edge{}
		dEdge.Src = pName
		dEdge.SrcPort = ":f" + fmt.Sprintf("%d", pos)
		dEdge.Dst = nName
	}

	return dNode, dEdge
}

// insertEntry 将新的entry插入到pos位置上
func (leaf *TreeLeaf) insertEntry(entry *utils.Entry, pos int) {
	newEntries := make([]*utils.Entry, len(leaf.entries)+1)
	newEntries[pos] = entry
	copy(newEntries[:pos], leaf.entries[:pos])
	copy(newEntries[pos+1:], leaf.entries[pos:])
	leaf.entries = newEntries
}

// appendLeaf src的entries追加到dst中
func appendLeaf(dst, src *TreeLeaf) {
	dst.entries = append(dst.entries, src.entries...)
	dst.next = src.next
	if src.next != nil {
		src.next.prev = dst
	}

	src.free()
}
