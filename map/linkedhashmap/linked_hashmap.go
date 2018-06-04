package linkedhashmap

import (
	"fmt"

	"github.com/asinglestep/gods/list/adlist"
	"github.com/asinglestep/gods/utils"
)

var (
	ErrEntryType = fmt.Errorf("Not a Entry type")
	ErrNotExist  = fmt.Errorf("No key exist")
)

// LinkedHashMap LinkedHashMap
type LinkedHashMap struct {
	capacity uint64                       // 容量
	list     *adlist.List                 // 双向链表
	m        map[interface{}]*adlist.Node // map
}

// NewLinkedHashMap NewLinkedHashMap
func NewLinkedHashMap(capacity uint64) *LinkedHashMap {
	l := &LinkedHashMap{}
	l.capacity = capacity
	l.list = adlist.NewListWithoutComparator()
	l.m = make(map[interface{}]*adlist.Node)

	return l
}

// Put 加入或者更新节点
func (l *LinkedHashMap) Put(entry *utils.Entry) error {
	node, ok := l.m[entry.GetKey()]
	if ok {
		// 节点存在，更新
		node.SetEntry(entry)
		l.list.DeleteNode(node)
		l.list.AddNodeToTail(node)
		return nil
	}

	if uint64(len(l.m)) >= l.capacity {
		// 超过LinkHashMap容量，删除第一个节点
		firstNode := l.list.Head()
		entry := firstNode.GetEntry()
		e, ok := entry.(*utils.Entry)
		if !ok {
			return ErrEntryType
		}

		l.list.DeleteNode(firstNode)
		delete(l.m, e.GetKey())
	}

	// 插入新节点
	node = l.list.AddNodeToTail(entry)
	l.m[entry.GetKey()] = node

	return nil
}

// Get 获取key指定的数据
func (l *LinkedHashMap) Get(key interface{}) (entry *utils.Entry, err error) {
	node, ok := l.m[key]
	if ok {
		// 节点存在，节点移动到链表尾部
		e, ok := node.GetEntry().(*utils.Entry)
		if !ok {
			return nil, ErrEntryType
		}

		l.list.DeleteNode(node)
		l.list.AddNodeToTail(node)
		return e, nil
	}

	return nil, ErrNotExist
}
