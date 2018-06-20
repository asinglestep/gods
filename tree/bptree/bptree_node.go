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
func NewTreeNode() *TreeNode {
	node := &TreeNode{}

	return node
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

	// 新的右节点，修改父节点，keys，childrens
	// rightKeys := make([]interface{}, mid)
	// copy(rKeys[:], node.keys[mid:])
	// rightChildrens := make([]iNode, mid)
	// copy(rightChildrens[:], node.childrens[mid:])
	right := NewTreeNode()
	right.parent = parent
	right.keys = node.keys[mid:]
	right.childrens = node.childrens[mid:]
	right.updateChildrensParent(right)

	// 新的左节点，修改父节点，keys，childrens
	node.parent = parent
	node.keys = node.keys[:mid]
	node.childrens = node.childrens[:mid]

	// 找到key在父节点中的位置
	pos, _ := parent.findKeyPosition(t.comparator, midKey)
	// 将key插入到父节点中
	parent.insertKey(midKey, pos)
	// 将right节点加入到父节点中
	parent.insertChildren(right, pos)
	return parent
}

// adjacent 获取相邻节点
func (node *TreeNode) adjacent(t *Tree) (adj iNode) {
	var left, right iNode
	parent := node.parent

	pos, bFound := parent.findKeyPosition(t.comparator, node.keys[0])
	if pos == 0 {
		// 第一个子节点的相邻节点为parent.childrens[1]
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
			// parent的key都小于node.keys[0]
			// 当前节点的位置为pos-1
			// 相邻节点为parent.childrens[pos-2]
			return parent.childrens[pos-2]
		}

		left = parent.childrens[pos-2]
		right = parent.childrens[pos]
	}

	if left.getKeys() > t.minKeys {
		return left
	}

	return right
}

// moveKey 移动相邻节点的key到node中
func (node *TreeNode) moveKey(t *Tree, adj iNode) *TreeNode {
	adjNode := adj.(*TreeNode)
	parent := node.parent

	if t.comparator.Compare(node.keys[0], adjNode.keys[0]) == utils.Lt {
		node.moveCaseRightAdjacentNode(t, adjNode, parent)
	} else {
		node.moveCaseLeftAdjacentNode(t, adjNode, parent)
	}

	return parent
}

// moveCaseRightAdjacentNode 相邻节点在右侧
func (node *TreeNode) moveCaseRightAdjacentNode(t *Tree, adj *TreeNode, parent *TreeNode) {
	// 找到相邻节点在父节点的位置
	pos, _ := parent.findKeyPosition(t.comparator, adj.keys[0])

	// 修改相邻节点的子节点的父节点
	adj.childrens[0].setParent(node)

	// 将相邻节点的第一个key和子节点移到当前节点
	node.keys = append(node.keys, adj.keys[0])
	node.childrens = append(node.childrens, adj.childrens[0])

	// 删除相邻节点的key和子节点
	adj.keys = adj.keys[1:]
	adj.childrens = adj.childrens[1:]

	// 修改父节点的key
	parent.keys[pos] = adj.keys[0]
}

// moveCaseLeftAdjacentNode 相邻节点在左侧
func (node *TreeNode) moveCaseLeftAdjacentNode(t *Tree, adj *TreeNode, parent *TreeNode) {
	// 找到当前节点在父节点的位置
	pos, bFound := parent.findKeyPosition(t.comparator, node.keys[0])
	if bFound {
		// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
		pos--
	}

	// 将相邻节点的最后一个key插入到当前节点
	node.insertKey(adj.keys[len(adj.keys)-1], 0)
	// 将相邻节点的最后一个子节点插入到当前节点
	node.insertChildren(adj.childrens[len(adj.childrens)-1], 0)
	// 修改当前节点的第一个子节点的父节点
	node.childrens[0].setParent(node)

	// 删除相邻节点的key和子节点
	adj.keys = adj.keys[:len(adj.keys)-1]
	adj.childrens = adj.childrens[:len(adj.childrens)-1]

	// 修改父节点的key
	parent.keys[pos] = node.keys[0]
}

// merge 合并相邻节点
func (node *TreeNode) merge(t *Tree, adj iNode) iNode {
	var key interface{}
	adjNode := adj.(*TreeNode)
	parent := node.parent

	if t.comparator.Compare(node.keys[0], adjNode.keys[0]) == utils.Lt {
		key = adjNode.keys[0]
		// 将相邻节点的keys和childrens加到当前节点中
		appendNode(node, adjNode)
	} else {
		key = node.keys[0]
		// 将当前节点的keys和childrens加入到相邻节点中
		appendNode(adjNode, node)
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
	if len(node.keys) == 0 {
		return 0, false
	}

	i, j := 0, len(node.keys)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(node.keys[h], key) == utils.Lt {
			i = h + 1
		} else {
			j = h
		}
	}

	if i == len(node.keys) {
		return i, false
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

// free free
func (node *TreeNode) free() {
	node.parent = nil
	node.keys = nil
	node.childrens = nil
}

// verif 验证
func (node *TreeNode) verify(t *Tree) bool {
	if node.parent != nil {
		// 非根节点至少有t个关键字
		if len(node.keys) < t.minKeys {
			fmt.Printf("非根节点少于%v个关键字\n", t.minKeys)
			return false
		}

		// 至少有t个子节点
		if len(node.childrens) < t.minKeys {
			fmt.Printf("非根节点少于%v个子节点\n", t.minKeys)
			return false
		}
	}

	// 最多有2t个关键字
	if len(node.keys) > t.maxKeys {
		fmt.Printf("节点多于%v个关键字\n", t.maxKeys)
		return false
	}

	// 最多有2t个子节点
	if len(node.childrens) > t.maxKeys {
		fmt.Printf("节点多于%v个子节点\n", t.maxKeys)
		return false
	}

	// 非叶子节点的key是其子节点key的最小值
	for i, v := range node.childrens {
		if v.isLeaf() {
			l := v.(*TreeLeaf)
			if t.comparator.Compare(l.entries[0].GetKey(), node.keys[i]) != utils.Et {
				fmt.Printf("父节点的第%v个key不是其叶子节点的第一个key, 父节点第%v个key: %v, 叶子节点的第一个key: %v\n", i, i, node.keys[i], l.entries[0].GetKey())
				return false
			}
		} else {
			n := v.(*TreeNode)
			if t.comparator.Compare(n.keys[0], node.keys[i]) != utils.Et {
				fmt.Printf("父节点的第%v个key不是其子节点的第一个key, 父节点第%v个key: %v, 子节点的第一个key: %v\n", i, i, node.keys[i], n.keys[0])
				return false
			}
		}
	}

	return true
}

// print print
func (node *TreeNode) print() string {
	if node.parent != nil {
		return fmt.Sprintf("内节点key: %v, \t此节点为父节点key[%v]的子节点, \t父节点为[%v]\n", node.printKeys(), node.keys[0], node.parent.printKeys())
	}

	return fmt.Sprintf("内节点key: %v, \t此节点为根节点\n", node.printKeys())
}

// String String
func (node *TreeNode) printKeys() string {
	keys := make([]string, 0, len(node.keys))
	for _, key := range node.keys {
		keys = append(keys, fmt.Sprintf("%v", key))
	}

	return strings.Join(keys, ",")
}

// dot dot
func (node *TreeNode) dot(t *Tree, nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
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
		pos, _ := node.parent.findKeyPosition(t.comparator, node.keys[0])
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
func (node *TreeNode) getChildrenAndUpdateFirstKeyIfNeed(pos int, key interface{}) iNode {
	if pos == 0 {
		// 比第一个关键字还小
		// 替换node.keys[0]
		node.keys[0] = key
		return node.childrens[0]
	}

	return node.childrens[pos-1]
}

// updateChildrensParent 更新的node的childrens的父节点
func (node *TreeNode) updateChildrensParent(parent *TreeNode) {
	for i := range node.childrens {
		node.childrens[i].setParent(parent)
	}
}

// insertKey 在pos位置插入key
func (node *TreeNode) insertKey(key interface{}, pos int) {
	keysLen := len(node.keys)
	newKeys := make([]interface{}, keysLen+1)
	copy(newKeys, node.keys[:pos])
	newKeys[pos] = key
	copy(newKeys[pos+1:], node.keys[pos:])
	node.keys = newKeys
}

// insertChildren 在pos位置插入children
func (node *TreeNode) insertChildren(children iNode, pos int) {
	csLen := len(node.childrens)
	newCs := make([]iNode, csLen+1)
	copy(newCs, node.childrens[:pos])
	newCs[pos] = children
	copy(newCs[pos+1:], node.childrens[pos:])
	node.childrens = newCs
}

// appendNode src的keys和childrens追加到dst中
func appendNode(dst, src *TreeNode) {
	dst.keys = append(dst.keys, src.keys...)
	dst.childrens = append(dst.childrens, src.childrens...)

	src.updateChildrensParent(dst)
	src.free()
}
