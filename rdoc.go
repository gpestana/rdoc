package rdoc

import (
	"errors"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/gpestana/rdoc/clock"
	op "github.com/gpestana/rdoc/operation"
	"log"
)

const (
	MapT = iota
	ListT
	RegT
)

type Doc struct {
	Id               string
	Clock            clock.Clock
	OperationsId     []string
	Head             *Node
	OperationsBuffer []op.Operation
}

// Returns a new rdoc data structure. It receives an ID which must be
// unique in the context of the network.
func Init(id string) *Doc {
	headNode := newNode(nil)
	c := clock.New([]byte(id))
	return &Doc{
		Id:               id,
		Clock:            c,
		OperationsId:     []string{},
		Head:             headNode,
		OperationsBuffer: []op.Operation{},
	}
}

func (d *Doc) ApplyRemoteOperation(o op.Operation) (*Doc, error) {
	// if operation has been applied already, skip
	if containsId(d.OperationsId, o.ID) {
		return d, nil
	}
	// if operation dependencies havent been all applied in the document, buffer
	// the operation
	missingOp := filter(o.Deps, d.OperationsId)
	if len(missingOp) != 0 {
		d.OperationsBuffer = append(d.OperationsBuffer, o)
		return d, nil
	}
	return d.ApplyOperation(o)
}

func (d *Doc) ApplyOperation(o op.Operation) (*Doc, error) {
	nPtr, travNodes, createdNodes := d.traverse(o.Cursor)

	// updates dependencies of traversed and created nodes
	var deps []*Node
	deps = append(deps, travNodes...)
	deps = append(deps, createdNodes...)
	for _, n := range deps {
		n.AddDependency(o.ID)
	}

	//TODO: let's assume the Mutate() call never fails for now.
	//TODO: how to rollback side effects of traverse if Mutate() fails?
	err := nPtr.Mutate(o)
	if err != nil {
		return d, err
	}

	d.OperationsId = append(d.OperationsId, o.ID)
	return d, nil
}

// Traverses the document form root element to the node indicated by the cursor
// input. When a path does not exist in the current document, create the node
// and link it to the document.
// The traverse function returns a pointer to the last node, a list of pointers
// od nodes traversed and a list of pointers of nodes created
func (d *Doc) traverse(cursor op.Cursor) (*Node, []*Node, []*Node) {
	var nPtr *Node
	var travNodes []*Node
	var createdNodes []*Node

	nPtr = d.Head

	for _, c := range cursor.Path {
		k := c.Get()
		switch k.(type) {
		// MapT
		case string:
			nodeType := MapT
			nif, exists := nPtr.hmap.Get(k.(string))
			if !exists {
				nn := newNode(k)
				_ = nPtr.link(nodeType, nn)
				nPtr = nn
				travNodes = append(travNodes, nPtr)
				createdNodes = append(createdNodes, nPtr)
				continue
			}
			nPtr = nif.(*Node)
			travNodes = append(travNodes, nPtr)

		// ListT
		case int:
			nodeType := ListT
			nif, exists := nPtr.list.Get(k.(int))
			if !exists {
				nn := newNode(k)
				_ = nPtr.link(nodeType, nn)
				nPtr = nn
				travNodes = append(travNodes, nPtr)
				createdNodes = append(createdNodes, nPtr)
				continue
			}
			nPtr = nif.(*Node)
			travNodes = append(travNodes, nPtr)
		}
	}
	return nPtr, travNodes, createdNodes
}

type Node struct {
	key  interface{}
	deps *arraylist.List
	hmap *hashmap.Map
	list *arraylist.List
	reg  *hashmap.Map
}

func newNode(key interface{}) *Node {
	return &Node{
		key:  key,
		deps: arraylist.New(),
		hmap: hashmap.New(),
		list: arraylist.New(),
		reg:  hashmap.New(),
	}
}

func (n *Node) GetList() *arraylist.List {
	return n.list
}

// applies operation mutation to the node
// note: assumes that mutation never fails for now
func (n *Node) Mutate(o op.Operation) error {
	mut := o.Mutation
	var err error

	switch mut.Key.(type) {
	case int:
		err = mut.List(n.list)
	case string:
		err = mut.Map(n.hmap)
	case nil:
		err = mut.Reg(n.hmap)
	default:
		err = errors.New("Invalid mutation")
	}
	return err
}

// appends new dependency to Node
func (n *Node) AddDependency(d string) {
	// TODO: should check if dep is valid with clock.Clock primitves?
	n.deps.Add(d)
}

// Links a node to the current node. The new node is linked depending on the
// type of linking required. It can be of type MapT, ListT or RegT.
func (n *Node) link(linkType int, node *Node) error {
	switch linkType {
	case MapT:
		key, ok := node.key.(string)
		if !ok {
			return errors.New("Map key must be string")
		}
		n.hmap.Put(key, node)

	case ListT:
		key, ok := node.key.(int)
		if !ok {
			return errors.New("List key must be an int")
		}
		n.list.Insert(key, node)

	case RegT:
		log.Println("linking RegT")
	default:
		return errors.New("linking type not correct")
	}

	return nil
}

// Returns all subsequent nodes from a particular Node
func (n *Node) allChildren() []*Node {
	var children []*Node
	var tmp []*Node
	tmp = append(tmp, directChildren(n)...)

	for {
		if len(tmp) == 0 {
			break
		}
		nextTmp := tmp[:1]
		tmp = tmp[1:]

		c := nextTmp[0]
		tmp = append(tmp, directChildren(c)...)
		children = append(children, c)
	}

	return children
}

func directChildren(n *Node) []*Node {
	var ch []*Node
	var in []interface{}
	in = append(in, n.hmap.Values()...)
	in = append(in, n.list.Values()...)
	in = append(in, n.reg.Values()...)

	// type cast to *Node
	for i, _ := range in {
		ch = append(ch, in[i].(*Node))
	}
	return ch
}

// checks if `sl` stice contains `id` string
func containsId(sl []string, id string) bool {
	for i, _ := range sl {
		if sl[i] == id {
			return true
		}
	}
	return false
}

// returns all strings in `deps` slice which do not exist in `ops`
func filter(deps []string, ops []string) []string {
	var diff []string
	for i, _ := range deps {
		contains := containsId(ops, deps[i])
		if !contains {
			diff = append(diff, deps[i])
		}
	}
	return diff
}
