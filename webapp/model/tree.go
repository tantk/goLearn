package model

import "fmt"

// The keys in a binary search tree are always sorted.

// Tree struct
type Tree struct {
	root *Node // May contain multiple same value keys
}

// Node Struct
type Node struct {
	// for any node x
	// if y is left child of x then y.key < x.key
	// if y is right child of x then y.key > x.key
	key    int
	parent *Node
	left   *Node
	right  *Node
	bookID int
}

// newNode: create a new node from key, with satellite data
func newNodeA(key int, bookID int) *Node {
	n := Node{
		key: key,
	}
	return &n
}

// AddA key to tree, with satellite data
func (t *Tree) AddA(key int, bookID int) {
	n := newNodeA(key, bookID)
	t.insert(n)
}

func newNode(key int) *Node {
	n := Node{
		key: key,
	}
	return &n
}

// Add key to tree, no satellite data
func (t *Tree) Add(key int) {
	n := newNode(key)
	t.insert(n)
}

// insert node n into t
func (t *Tree) insert(new *Node) {
	if t.root == nil {
		t.root = new
		return
	}
	// find position
	var p *Node = nil
	n := t.root
	for n != nil {
		p = n
		if new.key == n.key {
			return
		} else if new.key < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}
	// insert node
	if new.key < p.key {
		p.left = new
	} else {
		p.right = new
	}
	new.parent = p
}

// Delete node with key from tree
// true if deleted
func (t *Tree) Delete(key int) bool {
	n := t.root.find(key)
	if n == nil {
		return false
	}
	t.delete(n)
	return true
}

func (t *Tree) delete(d *Node) {
	if d.left == nil {
		t.transplant(d, d.right)
	} else if d.right == nil {
		t.transplant(d, d.left)
	} else {
		n := d.right.min()
		if n.parent != d {
			t.transplant(n, n.right)
			n.right = d.right
			n.right.parent = n
		}
		t.transplant(d, n)
		n.left = d.left
		n.left.parent = n
	}
}

func (t *Tree) transplant(u, v *Node) {
	if u.parent == nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// Flatten : return inorder slice of keys
func (t *Tree) Flatten() []int {
	var keys []int
	fn := func(n *Node) {
		keys = append(keys, n.key)
	}
	t.root.inorder(fn)
	return keys
}

// FlattenA : return inorder slice of keys and bookID
func (t *Tree) FlattenA() ([]int, []int) {
	var keys []int
	var bookids []int
	fn := func(n *Node) {
		keys = append(keys, n.key)
		bookids = append(bookids, n.bookID)
	}
	t.root.inorder(fn)
	return keys, bookids
}

//recursive function
// Walk tree in order calling fn for each node
func (n *Node) inorder(fn func(n *Node)) {
	if n != nil {
		n.left.inorder(fn)
		fn(n)
		n.right.inorder(fn)
	}
}

// Max : return maximum key from tree
func (t *Tree) Max() (int, error) {
	if t.root == nil {
		return 0, fmt.Errorf("Max() called on empty tree")
	}
	n := t.root.max()
	return n.key, nil
}

// max: find node with maximum key from tree rooted at n
func (n *Node) max() *Node {
	for n.right != nil {
		n = n.right
	}
	return n
}

// Min : return minimum key from tree
func (t *Tree) Min() (int, error) {
	if t.root == nil {
		return 0, fmt.Errorf("Min() called on empty tree")
	}
	n := t.root.min()
	return n.key, nil
}

// min: find node with minimum key from tree rooted at n
func (n *Node) min() *Node {
	for n.left != nil {
		n = n.left
	}
	return n
}

// Successor : find smallest key value larger than key
// panic if key not present, ok if found
func (t *Tree) Successor(key int) (int, bool) {
	n := t.root.find(key)
	if n == nil {
		panic("Succesor() called with non-existant key")
	}
	next := n.successor()
	if next == nil {
		return 0, false
	}
	return next.key, true
}

// find node by key
func (n *Node) find(key int) *Node {
	for n != nil && key != n.key {
		if key < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}
	return n
}

// successor: return node with smallest key larger than n.key
func (n *Node) successor() *Node {
	if n.right != nil {
		return n.right.min()
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	return p
}

// Predecessor : find largest key value smaller than key
// panic if key not present, ok if found
func (t *Tree) Predecessor(key int) (int, bool) {
	n := t.root.find(key)
	if n == nil {
		panic("Predecessor() called with non-existant key")
	}
	prev := n.predecessor()
	if prev == nil {
		return 0, false
	}
	return prev.key, true
}

// predecessor: return node with largest key smaller than n.key
func (n *Node) predecessor() *Node {
	if n.left != nil {
		return n.left.max()
	}
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}
	return p
}

// Size of tree
func (t *Tree) Size() int {
	if t.root == nil {
		return 0
	}
	return t.root.size()
}

// size of tree rooted at node n
func (n *Node) size() int {
	total := 0
	fn := func(n *Node) {
		if n != nil {
			total++
		}
	}
	n.inorder(fn)
	return total
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
