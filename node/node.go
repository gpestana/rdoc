// low level package that manages rdoc nodes. the node primitives are used
// mainly by rdoc.Mutate() when applying mutations and rdoc.Traverse() when
// traversing the tree
package node

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/hashmap"
	"reflect"
)

type Node struct {
	deps []string
	// operation id that originated the node
	opId string
	// node may be a map
	hmap *hashmap.Map
	// node may be a list
	list *arraylist.List
	// node may be a register
	reg *hashmap.Map
}

func New(opId string) *Node {
	return &Node{
		deps: []string{},
		opId: opId,
		hmap: hashmap.New(),
		list: arraylist.New(),
		reg:  hashmap.New(),
	}
}

func (n *Node) AddDependency(dep string) {
	n.deps = append(n.deps, dep)
}

func (n *Node) ClearDependency(dep string) {
	n.deps = filter(n.deps, dep)
}

// returns a child node which is part of the list or map
func (n *Node) GetChild(k interface{}) (*Node, bool, error) {
	switch key := k.(type) {
	case string:
		ni, exists := n.hmap.Get(key)
		if exists {
			n := ni.(*Node)
			return n, exists, nil
		}
	case int:
		ni, exists := n.list.Get(key)
		if exists {
			n := ni.(*Node)
			return n, exists, nil
		}
	default:
		return nil, false, errors.New("Node child is stored in list or map, key must be int or string")
	}
	// child with key `k` does not exist
	return nil, false, nil
}

// returns map with all values associated to a given key. the map is indexed by
// operation ID - the operation that created the KV pair in first place
func (n *Node) GetMVRegister() map[string]interface{} {
	keys := n.reg.Keys()
	mvrMap := make(map[string]interface{})
	for _, k := range keys {
		v, _ := n.reg.Get(k)
		mvrMap[k.(string)] = v
	}
	return mvrMap
}

// adds a value to the node
func (n *Node) Add(k interface{}, v interface{}, opId string) error {
	switch key := k.(type) {
	case string:
		// adds to map
		node, ok := v.(*Node)
		if !ok {
			return errors.New(
				fmt.Sprintf("(map.Add) value must be of type Node. Got instead: (%v", reflect.TypeOf(v)))
		}
		n.hmap.Put(key, node)
	case int:
		// adds to list
		node, ok := v.(*Node)
		if !ok {
			return errors.New(
				fmt.Sprintf("(list.Add) value  must be of type Node. Got instead: (%v", reflect.TypeOf(v)))
		}
		n.list.Insert(key, node)
	case nil:
		// adds to mvregister
		n.reg.Put(opId, v)
	default:
		return errors.New("Key type must be of type string (map element), int (list element) or nil (register)")
	}
	return nil
}

// returns all direct non-leaf children (maps and lists) from node
func (n *Node) GetChildren() []*Node {
	var ich []interface{}
	var ch []*Node
	ich = append(ich, n.list.Values()...)
	ich = append(ich, n.hmap.Values()...)

	for i, c := range ich {
		ch[i] = c.(*Node)
	}
	return ch
}

func filter(deps []string, dep string) []string {
	ndeps := []string{}
	for _, d := range deps {
		if d != dep {
			ndeps = append(ndeps, d)
		}
	}
	return ndeps
}
