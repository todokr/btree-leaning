package btree

import (
	"errors"
)

const (
	minChildren = 2 // 非リーフノードが指し示すことができる子ノードの最少数, degree
	maxChildren = minChildren * 2
	minItems    = minChildren - 1 // 子ノードを持つ非ルートノードが持つデータ項目の最少数
	maxItems    = maxChildren - 1
)

type Btree struct {
	root *node
}

func NewBtree() *Btree {
	return &Btree{}
}

func (t *Btree) Find(key []byte) ([]byte, error) {
	for next := t.root; next != nil; {
		pos, found := next.search(key)
		if found {
			return next.items[pos].val, nil
		}
		next = next.children[pos]
	}
	return nil, errors.New("key not found")
}

func (t *Btree) Insert(key, val []byte) {
	i := &item{key, val}

	if t.root == nil {
		t.root = &node{}
	}
	if t.root.numItems >= maxItems {
		t.splitRoot()
	}

	t.root.insert(i)
}

func (t *Btree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}
	deletedItem := t.root.delete(key, false)

	if t.root.numItems == 0 {
		if t.root.isLeaf() {
			t.root = nil
		} else {
			t.root = t.root.children[0]
		}
	}

	if deletedItem != nil {
		return true
	}
	return false
}

func (t *Btree) splitRoot() {
	newRoot := &node{}
	midItem, newNode := t.root.split()
	newRoot.insertItemAt(0, midItem)
	newRoot.insertChildAt(0, t.root)
	newRoot.insertChildAt(1, newNode)
	t.root = newRoot
}

func (t *Btree) String() string {
	v := &visualizer{t}
	return v.visualize()
}
