package btree

import (
	"bytes"
)

type item struct {
	key []byte
	val []byte
}

type node struct {
	children    [maxChildren]*node // 子ノード
	items       [maxItems]*item    // データ項目
	numItems    int
	numChildren int
}

func NewNode() *node {
	return &node{}
}

func (n *node) isLeaf() bool {
	return n.numChildren == 0
}

func (n *node) search(key []byte) (int, bool) {
	low, high := 0, n.numItems
	var mid int
	for low < high {
		mid = (low + high) / 2
		cmp := bytes.Compare(key, n.items[mid].key)
		switch {
		case cmp > 0:
			low = mid + 1
		case cmp < 0:
			high = mid
		case cmp == 0:
			return mid, true
		}
	}
	return low, false
}

// ノードを分割する
// *item: 分割元ノードの新たな終端
// *node: 新ノード
func (n *node) split() (*item, *node) {
	mid := minItems
	midItem := n.items[mid]
	newNode := NewNode()
	// データ項目の後半半分を新ノードへ
	copy(newNode.items[:], n.items[mid+1:])
	newNode.numItems = mid
	// 子ノードの後半部分を新ノードへ
	if !n.isLeaf() {
		copy(newNode.children[:], n.children[mid+1:])
		newNode.numChildren = mid + 1
	}
	// 移動済みのものを除去
	for i, l := mid, n.numItems; i < l; i++ {
		n.items[i] = nil
		n.numItems--

		if !n.isLeaf() {
			n.children[i+1] = nil
			n.numChildren--
		}
	}
	return midItem, newNode
}

func (n *node) insert(item *item) bool {
	pos, found := n.search(item.key)

	if found {
		n.items[pos] = item
		return false
	}

	if n.isLeaf() {
		n.insertItemAt(pos, item)
		return true
	}

	// いっぱいなので分割
	if n.children[pos].numItems >= maxItems {
		midItem, newNode := n.children[pos].split()
		n.insertItemAt(pos, midItem)
		n.insertChildAt(pos+1, newNode)

		// 目的となるリーフの走査方向のチェック
		switch cmp := bytes.Compare(item.key, n.items[pos].key); {
		case cmp < 0:
			// 探しているキーは走査方向の先にあるので、方向は同じまま
		case cmp > 0:
			// 子ノード中央のアイテムのキーは、探しているキーよりも小さいため、方向を変える必要がある
			pos++
		default:
			// 子ノード中央のキーが探しているキーなので、書き換える
			n.items[pos] = item
			return true
		}
	}
	return n.children[pos].insert(item)
}

func (n *node) insertItemAt(pos int, i *item) {
	if pos < n.numItems {
		// Make space for insertion if we are not appending to the very end of the items array.
		copy(n.items[pos+1:n.numItems+1], n.items[pos:n.numItems])
	}
	n.items[pos] = i
	n.numItems++
}

func (n *node) insertChildAt(pos int, c *node) {
	if pos < n.numChildren {
		// Make space for insertion if we are not appending to the very end of the children array.
		copy(n.children[pos+1:n.numChildren+1], n.children[pos:n.numChildren])
	}
	n.children[pos] = c
	n.numChildren++
}

func (n *node) removeItemAt(pos int) *item {
	removedItem := n.items[pos]
	n.items[pos] = nil
	// 歯抜けを埋める
	if lastPos := n.numItems - 1; pos < lastPos {
		copy(n.items[pos:lastPos], n.items[pos+1:lastPos+1])
		n.items[lastPos] = nil
	}
	n.numItems--

	return removedItem
}

func (n *node) removeChildAt(pos int) *node {
	removedChild := n.children[pos]
	n.children[pos] = nil
	// 歯抜けを埋める
	if lastPos := n.numChildren - 1; pos < lastPos {
		copy(n.children[pos:lastPos], n.children[pos+1:lastPos+1])
		n.children[lastPos] = nil
	}
	n.numChildren--

	return removedChild
}

// ある子ノードのノード数のアンダーフローを解消するために、隣接ノードとのマージ、あるいは兄弟ノードからの項目の借用を行う
func (n *node) fillChildAt(pos int) {
	switch {
	// 左ノードが存在し、かつ借用後に項目数の下限を下回らないなら、左ノードの終端項目を借用する
	//    ┌ M ┐  =>  ┌ C ┐
	//  A,B,C  Z      A,B  M,Z
	case pos > 0 && n.children[pos-1].numItems > minItems:
		left, right := n.children[pos-1], n.children[pos]
		copy(right.items[1:right.numItems+1], right.items[:right.numItems])
		right.items[0] = n.items[pos-1]
		right.numItems++

		if !right.isLeaf() {
			right.insertChildAt(0, left.removeChildAt(left.numChildren-1))
		}
		n.items[pos-1] = left.removeItemAt(left.numItems - 1)
	// 右ノードが存在し、かつ借用後に項目数の下限を下回らないなら、右ノードの始端項目を借用する
	//  ┌ M ┐    =>  ┌ X ┐
	//  A  X,Y,Z      A,M  Y,Z
	case pos+1 <= n.numChildren && n.children[pos].numItems > minItems:
		left, right := n.children[pos], n.children[pos+1]
		left.items[left.numItems] = n.items[pos]
		left.numItems++

		if !left.isLeaf() {
			left.insertChildAt(left.numChildren, right.removeChildAt(0))
		}

		n.items[pos] = right.removeItemAt(0)
	// 隣接する兄弟ノードをマージする
	//   ┌ C ┬ M ┬ X ┐   =>  ┌ C ┬ X ┐
	//   A    L    N    Z       A  L,M,N  Z
	default:
		if pos >= n.numItems {
			pos = n.numItems - 1
		}
		left, right := n.children[pos], n.children[pos+1]
		left.items[left.numItems] = n.removeItemAt(pos)
		left.numItems++
		copy(left.items[left.numItems:], right.items[:right.numItems])
		left.numItems += right.numItems

		if !left.isLeaf() {
			copy(left.children[left.numChildren:], right.children[:right.numChildren])
			left.numChildren += right.numChildren
		}

		n.removeChildAt(pos + 1)
		right = nil
	}
}

func (n *node) delete(key []byte, isSeekingSuccessor bool) *item {
	pos, found := n.search(key)
	var next *node

	if found {
		// リーフノードならそのまま削除するだけでOK
		if n.isLeaf() {
			return n.removeItemAt(pos)
		}
		// 右ノード以降にある削除対象の後継項目を探すフラグを立てる
		next, isSeekingSuccessor = n.children[pos+1], true
	} else {
		next = n.children[pos]
	}

	// 後継がリーフノードなら単純に削除でOK
	if n.isLeaf() && isSeekingSuccessor {
		return n.removeItemAt(0)
	}

	// キーに該当するノードがないなら何もしない
	if next == nil {
		return nil
	}

	deletedItem := next.delete(key, isSeekingSuccessor)
	if found && isSeekingSuccessor {
		n.items[pos] = deletedItem
	}

	// アンダーフローの修復
	if next.numItems < minItems {
		if found && isSeekingSuccessor {
			n.fillChildAt(pos + 1)
		} else {
			n.fillChildAt(pos)
		}
	}

	return deletedItem
}
