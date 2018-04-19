package rbtree

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"
)

type (
	rb uint32

	Direction  uint32
	InsertFlag uint32
	Compare    func(i interface{}, j interface{}) int
)

const (
	_BLACK rb = 0x00
	_RED   rb = 0x01
)

const (
	LEFT  Direction = 0x00
	RIGHT Direction = 0x01
)

const (
	ERROR_IF_EXIST   InsertFlag = 0x00
	SWAP_IF_EXIST    InsertFlag = 0x01
	COEXIST_IF_EXIST InsertFlag = 0x02
)

type Node struct {
	Item   interface{}
	parent *Node
	left   *Node
	right  *Node
	color  rb
}

func (n *Node) uncle() *Node {
	grandpa := n.parent.parent
	if n.parent == grandpa.right {
		return grandpa.left
	} else {
		return grandpa.right
	}
}

func rbtreePrint(w io.Writer, n *Node) {
	fmt.Fprintf(w, "{")
	if n != nil {
		if n.left != nil {
			rbtreePrint(w, n.left)
		}

		if n.color == _RED {
			fmt.Fprintf(w, "\033[1;31m %+v \033[0m", n.Item)
		} else {
			fmt.Fprintf(w, " %+v ", n.Item)
		}
		fmt.Fprintf(w, "\033[0;37m%x %x \033[0m",
			uintptr(unsafe.Pointer(n))&0xfff,
			uintptr(unsafe.Pointer(n.parent))&0xfff)

		if n.right != nil {
			rbtreePrint(w, n.right)
		}
	}
	fmt.Fprintf(w, "}")
}

type RBTree struct {
	root *Node
}

func (t *RBTree) Root() *Node {
	return t.root
}

func (t *RBTree) String() string {
	buf := new(bytes.Buffer)
	rbtreePrint(buf, t.root)
	return string(buf.Bytes())
}

/***********************************************************
 *         +---+                          +---+
 *         | q |                          | p |
 *         +---+                          +---+
 *        /     \     right rotation     /     \
 *     +---+   +---+  ------------->  +---+   +---+
 *     | p |   | z |                  | x |   | q |
 *     +---+   +---+                  +---+   +---+
 *    /     \                                /     \
 * +---+   +---+                          +---+   +---+
 * | x |   | y |                          | y |   | z |
 * +---+   +---+                          +---+   +---+
 **********************************************************/
func rotateRight(t *RBTree, q *Node) {
	p := q.left
	y := p.right
	if y != nil {
		y.parent = q
	}
	q.left = y

	fa := q.parent
	if fa == nil {
		t.root = p
	} else {
		if q == fa.left {
			fa.left = p
		} else {
			fa.right = p
		}
	}
	p.right = q
	p.parent = fa
	q.parent = p
}

/***********************************************************
 *         +---+                          +---+
 *         | q |                          | p |
 *         +---+                          +---+
 *        /     \                        /     \
 *     +---+   +---+                  +---+   +---+
 *     | p |   | z |                  | x |   | q |
 *     +---+   +---+  <-------------  +---+   +---+
 *    /     \          left rotation         /     \
 * +---+   +---+                          +---+   +---+
 * | x |   | y |                          | y |   | z |
 * +---+   +---+                          +---+   +---+
 **********************************************************/
func rotateLeft(t *RBTree, p *Node) {
	q := p.right
	if q == nil {
		return
	}

	y := q.left
	if y != nil {
		y.parent = p
	}
	p.right = y

	fa := p.parent
	if fa == nil {
		t.root = q
	} else {
		if p == fa.left {
			fa.left = q
		} else {
			fa.right = q
		}
	}
	q.left = p
	q.parent = fa
	p.parent = q
}

func deleteCase1(t *RBTree, n *Node, fa *Node, s **Node, d Direction) {
	(*s).color = _BLACK
	fa.color = _RED
	if d == LEFT {
		rotateLeft(t, fa)
		*s = fa.right
	} else {
		rotateRight(t, fa)
		*s = fa.left
	}
}

func deleteCase2(t *RBTree, n **Node, fa **Node, s *Node) {
	s.color = _RED
	*n = *fa
	*fa = (*n).parent
}

func deleteCase3(t *RBTree, n *Node, fa *Node, s **Node, d Direction) {
	if d == LEFT {
		if left := (*s).left; left != nil {
			left.color = _BLACK
		}
		(*s).color = _RED
		rotateRight(t, *s)
		*s = fa.right
	} else {
		if right := (*s).right; right != nil {
			right.color = _BLACK
		}
		(*s).color = _RED
		rotateLeft(t, *s)
		*s = fa.left
	}
}

func deleteCase4(t *RBTree, n **Node, fa *Node, s *Node, d Direction) {
	s.color = fa.color
	fa.color = _BLACK
	if d == LEFT {
		if right := s.right; right != nil {
			right.color = _BLACK
		}
		rotateLeft(t, fa)
	} else {
		if left := s.left; left != nil {
			left.color = _BLACK
		}
		rotateRight(t, fa)
	}
	*n = t.root
}

func deleteCase(t *RBTree, n *Node, fa *Node) {
	for (n == nil || n.color == _BLACK) && n != t.root {
		if n == fa.left {
			s := fa.right
			if s.color == _RED {
				deleteCase1(t, n, fa, &s, LEFT)
			}
			if (s.left == nil || s.left.color == _BLACK) &&
				(s.right == nil || s.right.color == _BLACK) {
				deleteCase2(t, &n, &fa, s)
			} else {
				if s.right == nil || s.right.color == _BLACK {
					deleteCase3(t, n, fa, &s, LEFT)
				}
				deleteCase4(t, &n, fa, s, LEFT)
				break
			}
		} else {
			s := fa.left
			if s.color == _RED {
				deleteCase1(t, n, fa, &s, RIGHT)
			}
			if (s.left == nil || s.left.color == _BLACK) &&
				(s.right == nil || s.right.color == _BLACK) {
				deleteCase2(t, &n, &fa, s)
			} else {
				if s.left == nil || s.left.color == _BLACK {
					deleteCase3(t, n, fa, &s, RIGHT)
				}
				deleteCase4(t, &n, fa, s, RIGHT)
				break
			}
		}
	}
	if n != nil {
		n.color = _BLACK
	}
}

func swapNodeWithLeaf(t *RBTree, n *Node, leaf *Node) {
	var (
		fa *Node
		rc *Node
	)
	color := leaf.color

	if fa = n.parent; fa == nil {
		t.root = leaf
	} else {
		if n == fa.left {
			fa.left = leaf
		} else {
			fa.right = leaf
		}
	}

	rc = leaf.right
	fa = leaf.parent
	if fa == n {
		fa = leaf
	} else {
		if rc != nil {
			rc.parent = fa
		}
		fa.left = rc
		leaf.right = n.right
		n.right.parent = leaf
	}

	leaf.parent = n.parent
	leaf.color = n.color
	leaf.left = n.left
	n.left.parent = leaf

	if color == _BLACK {
		deleteCase(t, rc, fa)
	}
}

func swapNodeWithSubtree(t *RBTree, n *Node, sub *Node) {
	if grandpa := n.parent; grandpa == nil {
		t.root = sub
		sub.color = _BLACK
		sub.parent = nil
	} else {
		if n == grandpa.left {
			grandpa.left = sub
		} else {
			grandpa.right = sub
		}
		sub.parent = grandpa

		if n.color == _RED || sub.color == _RED {
			sub.color = _BLACK
		} else {
			deleteCase(t, sub, grandpa)
		}
	}
}

func removeNodeHasOneChild(t *RBTree, ori *Node) {
	fa := ori.parent
	sub := ori.left
	if sub == nil {
		sub = ori.right
	}
	if sub != nil {
		sub.parent = fa
	}
	if fa != nil {
		if ori == fa.left {
			fa.left = sub
		} else {
			fa.right = sub
		}
	} else {
		t.root = sub
	}
	if ori.color == _BLACK {
		deleteCase(t, sub, fa)
	}
}

func DeleteNode(t *RBTree, n *Node) {
	if t.root == nil || n == nil {
		return
	}

	if n.left == nil || n.right == nil {
		removeNodeHasOneChild(t, n)
	} else {
		leaf := n.right
		for leaf.left != nil {
			leaf = leaf.left
		}
		swapNodeWithLeaf(t, n, leaf)
	}
}

func InsertCase(t *RBTree, n *Node) {
	if t.root == nil || n == nil {
		return
	}

	var (
		fa      *Node
		grandpa *Node
	)
	for {
		if fa = n.parent; fa == nil || fa.color != _RED {
			break
		}
		if uc := n.uncle(); uc != nil && uc.color == _RED {
			fa.color = _BLACK
			uc.color = _BLACK

			grandpa = fa.parent
			grandpa.color = _RED
			n = grandpa

			continue
		}

		grandpa = fa.parent
		if n == fa.right && fa == grandpa.left {
			rotateLeft(t, fa)
			n = n.left

			fa = n.parent
			grandpa = fa.parent
		} else if n == fa.left && fa == grandpa.right {
			rotateRight(t, fa)
			n = n.right

			fa = n.parent
			grandpa = fa.parent
		}

		fa.color = _BLACK
		grandpa.color = _RED
		if n == fa.left {
			rotateRight(t, grandpa)
		} else {
			rotateLeft(t, grandpa)
		}
	}
	t.root.color = _BLACK
}

func insertToLRMost(root **Node, n *Node, d Direction) {
	n.parent = nil
	n.left = nil
	n.right = nil
	n.color = _RED
	ptr := root
	if d == LEFT {
		for *ptr != nil {
			n.parent = *ptr
			ptr = &(*ptr).left
		}
	} else {
		for *ptr != nil {
			n.parent = *ptr
			ptr = &(*ptr).right
		}
	}
	*ptr = n
}

func RemoveSubtree(t *RBTree, sub *Node) {
	fa := sub.parent
	if fa == nil {
		t.root = nil
	} else if sub == fa.left {
		sil := fa.right
		swapNodeWithSubtree(t, fa, sil)
		insertToLRMost(&sil, fa, LEFT)
		InsertCase(t, fa)
	} else {
		sil := fa.left
		swapNodeWithSubtree(t, fa, sil)
		insertToLRMost(&sil, fa, RIGHT)
		InsertCase(t, fa)
	}
}

func InsertWithoutBalance(t *RBTree, n *Node, comp Compare, flag InsertFlag) (ori *Node) {
	n.parent = nil
	n.left = nil
	n.right = nil
	n.color = _RED

	ptr := &t.root
	for *ptr != nil {
		n.parent = *ptr
		if c := comp(n.Item, (*ptr).Item); c < 0 {
			ptr = &(*ptr).left
		} else if c > 0 {
			ptr = &(*ptr).right
		} else if flag != COEXIST_IF_EXIST {
			ori = *ptr
			if flag == ERROR_IF_EXIST {
				return
			} else if flag == SWAP_IF_EXIST {
				if left := (*ptr).left; left != nil {
					left.parent = n
					n.left = left
				}
				if right := (*ptr).right; right != nil {
					right.parent = n
					n.right = right
				}
				n.color = (*ptr).color
				*ptr = n
				return
			}
		}
	}
	*ptr = n
	return nil
}

func findAll(root *Node, item interface{}, comp Compare) (ns []*Node) {
	ns = []*Node{}
	ptr := root
	for ptr != nil {
		if c := comp(item, ptr.Item); c == 0 {
			left_ns := findAll(ptr.left, item, comp)
			right_ns := findAll(ptr.right, item, comp)
			ns = append(ns, ptr)
			ns = append(ns, left_ns...)
			ns = append(ns, right_ns...)
			break
		} else if c < 0 {
			ptr = ptr.left
		} else {
			ptr = ptr.right
		}
	}
	return ns
}

func FindAll(t *RBTree, item interface{}, comp Compare) (ns []*Node) {
	return findAll(t.root, item, comp)
}

func Find(t *RBTree, item interface{}, comp Compare) *Node {
	ptr := t.root
	for ptr != nil {
		if c := comp(item, ptr.Item); c == 0 {
			return ptr
		} else if c < 0 {
			ptr = ptr.left
		} else {
			ptr = ptr.right
		}
	}
	return nil
}
