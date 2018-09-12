// low level package that manages rdoc nodes. the node primitives are used
// mainly by rdoc.Mutate() when applying mutations and rdoc.Traverse() when
// traversing the tree
package node

import (
	"errors"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/hashmap"
)

type Node interface {
	// adds node dependency
	AddDependency(string)
	// clears dependency
	ClearDependency(string)
}

// MVRegister element is a multi-value register that holds concrete values in
// the rdoc structure.
type MVRegister struct {
	deps []string
	opId string
	// keeps the values of the multi-value register as a map, in which keys are
	// the operation IDs and values the value of the operation mutation
	values *hashmap.Map
}

func (r *MVRegister) AddDependency(dep string) {
	r.deps = append(r.deps, dep)
}

func (r *MVRegister) ClearDependency(dep string) {
	r.deps = filter(r.deps, dep)
}

// adds register value, indexed based on operation ID. current value type is
// string or int, otherwise it will return an error
func (r *MVRegister) AddValue(opId string, value interface{}) error {
	switch value.(type) {
	case int:
	case string:
	default:
		return errors.New("Value must be either string or int")
	}
	r.values.Put(opId, value)
	return nil
}

// ListElement is an element of a list which can point to another map, a list
// and/or a register
type ListElement struct {
	deps []string
	opId string
	// a list element may be a map
	hmap *hashmap.Map
	// a list element may be a list
	list *arraylist.List
	// a list element may be a register
	reg MVRegister
}

func (l *ListElement) AddDependency(dep string) {
	l.deps = append(l.deps, dep)
}

func (l *ListElement) ClearDependency(dep string) {
	l.deps = filter(l.deps, dep)
}

// adds a new node to the list
func (l *ListElement) Add(k interface{}, value interface{}, opId string) error {
	switch key := k.(type) {
	case string:
		// adds to map
		l.hmap.Put(key, node)
	case int:
		// adds to list
		l.list.Insert(key, node)
	case nil:
		// wrong. where is the key??
		l.reg.Insert(opId, value)
	default:
		return errors.New("Key type must be of type string (map element), int (list element) or nil (register)")
	}
	return nil
}

// MapElement is an element of a map which can point to another map, a list
// and/or a register
type MapElement struct {
	deps []string
	opId string
	// a map element may be a map
	hmap *hashmap.Map
	// a map element may be a list
	list *arraylist.List
	// a map element may be a register
	reg MVRegister
}

func (m *MapElement) AddDependency(dep string) {
	m.deps = append(m.deps, dep)
}

func (m *MapElement) ClearDependency(dep string) {
	m.deps = filter(m.deps, dep)
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

var _ Node = (*MVRegister)(nil)
var _ Node = (*ListElement)(nil)
var _ Node = (*MapElement)(nil)
