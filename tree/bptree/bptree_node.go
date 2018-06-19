package bptree

import (
	"fmt"
	"strings"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// TreeNode TreeNode
type TreeNode struct {
	parent    *TreeNode     // 指向父节点
	childrens []iNode       // 子节点
	keys      []interface{} // 关键字
}

// NewTreeNode NewTreeNode
func NewTreeNode(parent *TreeNode, childrens []iNode, keys []interface{}) *TreeNode {
	node := &TreeNode{}
	node.parent = parent
	node.childrens = childrens
	node.keys = keys

	return node
}

// String String
func (node *TreeNode) String() string {
	return fmt.Sprintf("内节点[%v]", node.keys)
}

// split 分裂节点
func (node *TreeNode) split(t *Tree) iNode {
	parent := node.parent
	if parent == nil {
		// 根节点为满的
		parent = t.splitRootNode(node, node.keys[0])
	}

	// 满节点时有2t个key，所以中间key的位置为t
	mid := t.minKeys
	midKey := node.keys[mid]

	// 分裂成2个新节点，左右节点各包含t个key，t个children
	// 新的右节点，修改父节点，keys，childrens
	// rightKeys := make([]interface{}, mid)
	// copy(rKeys[:], node.keys[mid:])
	// rightChildrens := make([]iNode, mid)
	// copy(rightChildrens[:], node.childrens[mid:])
	right := NewTreeNode(parent, node.keys[mid:], node.childrens[mid:])
	right.updateChildrensParent()

	// 新的左节点，修改父节点，keys，childrens
	// node.parent = parent
	// node.keys = node.keys[:mid]
	// node.childrens = node.childrens[:mid]
	// node.updateChildrensParent()
	left := NewTreeNode(parent, node.keys[:mid], node.childrens[:mid])
	left.updateChildrensParent()

	// 找到key在父节点中的位置
	pos, _ := parent.findKeyPosition(t.comparator, midKey)

	// 修改父节点的keys，childrens
	// 将key插入到父节点中
	parentKeysLen := len(parent.keys)
	newKeys := make([]Key, parentKeysLen+1)
	copy(newKeys, parent.keys[:pos])
	newKeys[pos] = midKey
	copy(newKeys[pos+1:], parent.keys[pos:])
	parent.keys = newKeys

	// 将right节点加入到父节点中，keys和childrens数量相同
	newCs := make([]iNode, parentKeysLen+1)
	copy(newCs, parent.childrens[:pos])
	newCs[pos] = right
	copy(newCs[pos+1:], parent.childrens[pos:])
	parent.childrens = newCs

	// return node, right, midKey
	return parent
}

// adjacent 获取相邻节点
func (node *TreeNode) adjacent(t *Tree) (adj iNode) {
	parent := node.parent

	// 在父节点中查找lnode.keys[0]
	pos := parent.findKeyPosition(node.keys[0])
	switch pos {
	case len(parent.keys):
		// parent的key都小于node.keys[0]，该节点的位置为pos-1
		// 相邻节点为parent.childrens[pos-2]
		adj = parent.childrens[pos-2]

	case 0:
		// 第一个子节点的相邻节点为parent.childrens[1]
		adj = parent.childrens[1]

	default:
		var left, right iNode
		if parent.keys[pos] == node.keys[0] {
			if pos == len(parent.keys)-1 {
				// 最后一个节点
				adj = parent.childrens[pos-1]
				return adj
			}

			left = parent.childrens[pos-1]
			right = parent.childrens[pos+1]
		} else {
			if pos == 1 {
				// 第一个节点
				adj = parent.childrens[1]
				return adj
			}

			left = parent.childrens[pos-2]
			right = parent.childrens[pos]
		}

		if left.getKeys() > t.minKeyN {
			adj = left
		} else {
			adj = right
		}
	}

	return adj
}

// moveKey 移动相邻节点的key到leaf中
func (node *TreeNode) moveKey(adj iNode) *TreeNode {
	adjNode := adj.(*TreeNode)
	parent := node.parent

	if node.keys[0] < adjNode.keys[0] {
		// 在相邻节点的左侧，找到相邻节点在父节点的位置
		pos := parent.findKeyPosition(adjNode.keys[0])

		// 将相邻节点的key和子节点移到当前节点
		node.keys = append(node.keys, adjNode.keys[0])
		node.childrens = append(node.childrens, adjNode.childrens[0])

		// 修改相邻节点的子节点的父节点
		adjNode.childrens[0].setParent(node)

		// 删除相邻节点的key和子节点
		adjNode.keys = adjNode.keys[1:]
		adjNode.childrens = adjNode.childrens[1:]

		// 修改父节点的key
		parent.keys[pos] = adjNode.keys[0]
	} else {
		// 在相邻节点的右侧，找到当前节点在父节点的位置
		pos := parent.findKeyPosition(node.keys[0])
		if pos == len(parent.keys) || parent.keys[pos] > node.keys[0] {
			// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
			pos--
		}

		// 将相邻节点的key和子节点移到当前节点
		lastKeyIdx := len(adjNode.keys) - 1
		keys := make([]Key, len(node.keys)+1)
		keys[0] = adjNode.keys[lastKeyIdx]
		copy(keys[1:], node.keys)
		node.keys = keys

		lastChildrenIdx := len(adjNode.childrens) - 1
		childrens := make([]iNode, len(node.childrens)+1)
		childrens[0] = adjNode.childrens[lastChildrenIdx]
		copy(childrens[1:], node.childrens)
		node.childrens = childrens

		// 修改当前节点的第一个子节点的父节点
		node.childrens[0].setParent(node)

		// 删除相邻节点的key和子节点
		adjNode.keys = adjNode.keys[:lastKeyIdx]
		adjNode.childrens = adjNode.childrens[:lastChildrenIdx]

		// 修改父节点的key
		parent.keys[pos] = node.keys[0]
	}

	return parent
}

// merge 合并相邻节点
func (node *TreeNode) merge(adj iNode) iNode {
	var findKey Key

	pos := 0
	adjNode := adj.(*TreeNode)
	parent := node.parent

	if node.keys[0] < adjNode.keys[0] {
		// 在相邻节点的左侧，找到相邻节点在父节点的位置
		findKey = adjNode.keys[0]
		pos = parent.findKeyPosition(findKey)

		// 将相邻节点的key和子节点加到当前节点中
		node.keys = append(node.keys, adjNode.keys...)
		node.childrens = append(node.childrens, adjNode.childrens...)

		// 修改相邻节点的子节点的父节点
		for i := range adjNode.childrens {
			adjNode.childrens[i].setParent(node)
		}
	} else {
		// 在相邻节点的右侧，找到当前节点在父节点的位置
		findKey = node.keys[0]
		pos = parent.findKeyPosition(findKey)

		// 将当前节点的key和子节点加入到相邻节点中
		adjNode.keys = append(adjNode.keys, node.keys...)
		adjNode.childrens = append(adjNode.childrens, node.childrens...)

		// 修改当前节点的子节点的父节点
		for i := range node.childrens {
			node.childrens[i].setParent(adjNode)
		}
	}

	if pos == len(parent.keys) || parent.keys[pos] > findKey {
		// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
		pos--
	}

	// 删除父节点的key和子节点
	parent.keys = append(parent.keys[:pos], parent.keys[pos+1:]...)
	parent.childrens = append(parent.childrens[:pos], parent.childrens[pos+1:]...)
	return parent
}

// setParent 设置父节点
func (node *TreeNode) setParent(parent *TreeNode) {
	node.parent = parent
}

// getParent 获取父节点
func (node *TreeNode) getParent() *TreeNode {
	return node.parent
}

// findKeyPosition
// 在节点中查找第一个大于等于key的位置
// 没有比key大的，则返回node.keys的长度
func (node *TreeNode) findKeyPosition(comparator utils.Comparator, key interface{}) (pos int, bFound bool) {
	i, j := 0, len(node.keys)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(node.keys[h], key) == utils.Lt {
			i = h + 1
		} else {
			j = h
		}
	}

	if comparator.Compare(node.keys[i], key) == utils.Et {
		bFound = true
	}

	return i, bFound
}

// getPosKey 获取pos位置的key
func (node *TreeNode) getPosKey(pos int) interface{} {
	return node.keys[pos]
}

// getKeys 获取key的数量
func (node *TreeNode) getKeys() int {
	return len(node.keys)
}

// isLeaf 是否是叶子节点
func (node *TreeNode) isLeaf() bool {
	return false
}

// isFull 是否是满节点
func (node *TreeNode) isFull(max int) bool {
	return len(node.keys) == max
}

// verif 验证
func (node *TreeNode) verif(t *Tree) bool {
	if node.parent != nil {
		// 非根节点至少有t个关键字
		if len(node.keys) < t.minKeyN {
			fmt.Printf("非根节点少于%v个关键字\n", t.minKeyN)
			return false
		}

		// 至少有t个子节点
		if len(node.childrens) < t.minKeyN {
			fmt.Printf("非根节点少于%v个子节点\n", t.minKeyN)
			return false
		}
	}

	// 最多有2t个关键字
	if len(node.keys) > t.maxKeyN {
		fmt.Printf("节点多于%v个关键字\n", t.maxKeyN)
		return false
	}

	// 最多有2t个子节点
	if len(node.childrens) > t.maxKeyN {
		fmt.Printf("节点多于%v个子节点\n", t.maxKeyN)
		return false
	}

	// 非叶子节点的key是其子节点key的最小值
	for i, v := range node.childrens {
		if v.isLeaf() {
			l := v.(*TreeLeaf)
			if l.entries[0].Key != node.keys[i] {
				fmt.Printf("l.entries[0].Key != node.keys[i], l.entries[0].Key %v, node.keys[i] %v\n", l.entries[0].Key, node.keys[i])
				return false
			}
		} else {
			n := v.(*TreeNode)
			if n.keys[0] != node.keys[i] {
				fmt.Printf("n.keys[0] != node.keys[i], n.keys[0] %v, node.keys[i] %v\n", n.keys[0], node.keys[i])
				return false
			}
		}
	}

	return true
}

// print print
func (node *TreeNode) print() {
	if node.parent != nil {
		fmt.Printf("内节点key: %v, \t此节点为父节点key[%v]的子节点, \t父节点为[%v]\n", node.keys, node.keys[0], node.parent)
	} else {
		fmt.Printf("内节点key: %v, \t此节点为根节点\n", node.keys)
	}
}

// dot dot
func (node *TreeNode) dot(nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加一个node
	attrValues := make([]string, 0, len(node.keys))

	for i, key := range node.keys {
		attrValues = append(attrValues, fmt.Sprintf("<f%d> %d ", i, key))
	}

	attr := "\"" + strings.Join(attrValues, "|") + "\""

	dNode = &dot.Node{}
	dNode.Name = nName
	dNode.Attr = map[string]string{
		"label": attr,
	}

	// 添加一个edge
	if node.parent != nil {
		pos := node.parent.findKeyPosition(node.keys[0])
		dEdge = &dot.Edge{}
		dEdge.Src = pName
		dEdge.SrcPort = ":f" + fmt.Sprintf("%d", pos)
		dEdge.Dst = nName
	}

	return dNode, dEdge
}

// getChildrenAndUpdateFirstKeyIfNeed
// 如果pos等于0，则更新第一个key，返回pos位置的子节点
// 如果pos大于0，则返回pos-1位置的子节点
func (node *TreeNode) getChildrenAndUpdateFirstKeyIfNeed(pos int) iNode {
	if pos == 0 {
		// 比第一个关键字还小
		// 替换node.keys[0]
		node.keys[0] = key
		return node.childrens[0]
	}

	return node.childrens[pos-1]
}

// updateChildrensParent 更新的node的childrens的父节点
func (node *TreeNode) updateChildrensParent() {
	for i := range node.childrens {
		node.childrens[i].setParent(node)
	}
}

// updateKeys 更新node.keys
func (node *TreeNode) updateKeys() {
	parentKeysLen := len(parent.keys)
	newKeys := make([]Key, parentKeysLen+1)
	copy(newKeys, parent.keys[:pos])
	newKeys[pos] = midKey
	copy(newKeys[pos+1:], parent.keys[pos:])
	parent.keys = newKeys
}

// updateChildrens 更新node.childrens
func (node *TreeNode) updateChildrens() {

}
