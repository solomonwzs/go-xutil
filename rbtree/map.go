package rbtree

type Map struct {
	tree *RBTree
	comp Compare
}

type kv struct {
	key   interface{}
	value interface{}
}

func NewMap(comp Compare) *Map {
	return &Map{
		tree: new(RBTree),
		comp: comp,
	}
}

func (m *Map) mapItemCompare(a, b interface{}) int {
	return m.comp(a.(kv).key, b.(kv).key)
}

func (m *Map) Set(key interface{}, value interface{}) {
	node := &Node{Item: kv{key, value}}
	InsertWithoutBalance(m.tree, node, m.mapItemCompare, SWAP_IF_EXIST)
	InsertCase(m.tree, node)
}

func (m *Map) Delete(key interface{}) {
	if n := Find(m.tree, kv{key, nil}, m.mapItemCompare); n != nil {
		DeleteNode(m.tree, n)
	}
}

func (m *Map) Get(key interface{}) (value interface{}, exist bool) {
	if n := Find(m.tree, kv{key, nil}, m.mapItemCompare); n != nil {
		return n.Item.(kv).value, true
	} else {
		return nil, false
	}
}
