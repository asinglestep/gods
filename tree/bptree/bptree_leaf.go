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

func (leaf *TreeLeaf) String() string {
	return fmt.Sprintf("叶子节点[%v]", leaf.entries)
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

	// 分裂成2个新节点，左右节点各占t个key，t个children
	// 新的右节点，修改父节点，next指针，prev指针，entries
	right := NewTreeLeaf()
	right.parent = parent
	right.next = leaf.next
	right.prev = leaf
	right.entries = make([]*utils.Entry, mid)
	for i := range right.entries {
		right.entries[i] = leaf.entries[mid+i]
	}

	// 新的左节点，修改父节点，next指针，entries
	leaf.parent = parent
	leaf.next = right
	leaf.entries = leaf.entries[:mid]

	// 找到key在父节点中的位置
	pos, _ := parent.findKeyPosition(t.comparator, midEntry.GetKey())

	// 修改父节点的keys，childrens
	// 将key插入到父节点中
	parentKeysLen := len(parent.keys)
	newKeys := make([]interface{}, parentKeysLen+1)
	copy(newKeys, parent.keys[:pos])
	newKeys[pos] = midEntry.GetKey()
	copy(newKeys[pos+1:], parent.keys[pos:])
	parent.keys = newKeys

	// 将right节点加入到父节点中，keys和childrens数量相同
	newCs := make([]iNode, parentKeysLen+1)
	copy(newCs, parent.childrens[:pos])
	newCs[pos] = right
	copy(newCs[pos+1:], parent.childrens[pos:])
	parent.childrens = newCs

	// return leaf, right, midEntry.Key
	return parent
}

// adjacent 获取相邻节点
func (leaf *TreeLeaf) adjacent(t *Tree) (adj iNode) {
	parent := leaf.parent

	// 在父节点中查找leaf.entries[0].Key
	pos, _ := parent.findKeyPosition(t.comparator, leaf.entries[0].GetKey())
	switch pos {
	case len(parent.keys):
		// parent的key都小于leaf.entries[0].Key，该节点的位置为pos-1
		// 相邻节点为parent.childrens[pos-2]
		adj = parent.childrens[pos-2]

	case 0:
		// 第一个子节点的相邻节点为parent.childrens[pos+1]
		adj = parent.childrens[1]

	default:
		var left, right iNode
		if parent.keys[pos] == leaf.entries[0].GetKey() {
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
func (leaf *TreeLeaf) moveKey(adj iNode) *TreeNode {
	adjNode := adj.(*TreeLeaf)
	parent := leaf.parent

	if leaf.entries[0].Key < adjNode.entries[0].Key {
		// 在相邻节点的左侧，找到相邻节点在父节点中的位置
		pos := parent.findKeyPosition(adjNode.entries[0].Key)

		// 将相邻节点的key加入到当前节点
		leaf.entries = append(leaf.entries, adjNode.entries[0])

		// 删除相邻节点的key
		adjNode.entries = adjNode.entries[1:]

		// 修改父节点的key
		parent.keys[pos] = adjNode.entries[0].Key
	} else {
		// 在相邻节点的右侧，找到当前节点在父节点中的位置
		pos := parent.findKeyPosition(leaf.entries[0].Key)
		if pos == len(parent.keys) || parent.keys[pos] > leaf.entries[0].Key {
			// 删除的是子节点的最小值，导致删除后在父节点中找不到这个key
			pos--
		}

		// 将相邻节点的key加入到当前节点
		entries := make([]*Entry, len(leaf.entries)+1)
		entries[0] = adjNode.entries[len(adjNode.entries)-1]
		copy(entries[1:], leaf.entries)
		leaf.entries = entries

		// 删除相邻节点的key
		adjNode.entries = adjNode.entries[:len(adjNode.entries)-1]

		// 修改父节点的key
		parent.keys[pos] = leaf.entries[0].Key
	}

	return parent
}

// merge 合并相邻节点
func (leaf *TreeLeaf) merge(adj iNode) iNode {
	var findKey Key

	pos := 0
	adjNode := adj.(*TreeLeaf)
	parent := adjNode.parent

	if leaf.entries[0].Key < adjNode.entries[0].Key {
		// 在相邻节点的左侧，找到相邻节点在父节点中的位置
		findKey = adjNode.entries[0].Key
		pos = parent.findKeyPosition(findKey)

		// 将相邻节点合并到当前节点
		leaf.entries = append(leaf.entries, adjNode.entries...)
		leaf.next = adjNode.next
		if adjNode.next != nil {
			adjNode.next.prev = leaf
		}
	} else {
		// 在相邻节点的右侧，找到当前节点在父节点中的位置
		findKey = leaf.entries[0].Key
		pos = parent.findKeyPosition(findKey)

		// 将当前节点合并到相邻节点
		adjNode.entries = append(adjNode.entries, leaf.entries...)
		adjNode.next = leaf.next
		if leaf.next != nil {
			leaf.next.prev = adjNode
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
func (leaf *TreeLeaf) setParent(parent *TreeNode) {
	leaf.parent = parent
}

// getParent 获取父节点
func (leaf *TreeLeaf) getParent() *TreeNode {
	return leaf.parent
}

// findKeyPosition 在节点中查找第一个大于等于key的位置，没有比key大的，则返回node.entries的长度
func (leaf *TreeLeaf) findKeyPosition(comparator utils.Comparator, key interface{}) (pos int, bFound bool) {
	i, j := 0, len(leaf.entries)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(leaf.entries[h].GetKey(), key) == utils.Lt {
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

// verif 验证
func (leaf *TreeLeaf) verif(t *Tree) bool {
	if leaf.parent != nil {
		// 非根节点至少有t个关键字
		if len(leaf.entries) < t.minKeyN {
			fmt.Printf("非根节点少于%v个关键字\n", t.minKeyN)
			return false
		}
	}

	// 最多有2t个关键字
	if len(leaf.entries) > t.maxKeyN {
		fmt.Printf("节点多于%v个关键字\n", t.maxKeyN)
		return false
	}

	return true
}

// print print
func (leaf *TreeLeaf) print() {
	if leaf.parent != nil {
		fmt.Printf("叶子节点key: %v, \t此节点为父节点key[%v]的子节点, \t父节点为[%v]\n", leaf.entries, leaf.entries[0].Key, leaf.parent)
	} else {
		if len(leaf.entries) == 0 {
			fmt.Printf("此b树为一个空树\n")
		} else {
			fmt.Printf("叶子节点key: %v, \t此节点为根节点\n", leaf.entries)
		}
	}
}

// dot dot
func (leaf *TreeLeaf) dot(nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加一个node
	attrValues := make([]string, 0, len(leaf.entries))

	for i, v := range leaf.entries {
		attrValues = append(attrValues, fmt.Sprintf("<f%d> %d ", i, v.Key))
	}

	attr := "\"" + strings.Join(attrValues, "|") + "\""

	dNode = &dot.Node{}
	dNode.Name = nName
	dNode.Attr = map[string]string{
		"label": attr,
	}

	// 添加一个edge
	if leaf.parent != nil {
		pos := leaf.parent.findKeyPosition(leaf.entries[0].Key)
		dEdge = &dot.Edge{}
		dEdge.Src = pName
		dEdge.SrcPort = ":f" + fmt.Sprintf("%d", pos)
		dEdge.Dst = nName
	}

	return dNode, dEdge
}

// insertEntry 将新的entry插入到pos位置上
func (leaf *TreeLeaf) insertEntry(entry *utils.Entry, pos int) {
	newEntries := make([]*Entry, len(leaf.entries)+1)
	newEntries[pos] = entry
	copy(newEntries[:pos], leaf.entries[:pos])
	copy(newEntries[pos+1:], leaf.entries[pos:])
	leaf.entries = newEntries
}
