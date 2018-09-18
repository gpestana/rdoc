package rdoc

import (
	"fmt"
	"github.com/gpestana/rdoc/clock"
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
)

type Doc struct {
	Id               string
	Clock            clock.Clock
	OperationsId     []string
	Head             *n.Node
	OperationsBuffer []op.Operation
}

// Returns a new rdoc data structure. It receives an ID which must be
// unique in the context of the network.
func Init(id string) *Doc {
	headNode := n.New("")
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
	missingOp := diff(o.Deps, d.OperationsId)
	if len(missingOp) != 0 {
		d.OperationsBuffer = append(d.OperationsBuffer, o)
		return d, nil
	}
	return d.ApplyOperation(o)
}

func (d *Doc) ApplyOperation(o op.Operation) (*Doc, error) {
	nPtr, travNodes, createdNodes := d.traverse(o.Cursor, o.ID)

	// updates dependencies of traversed and created nodes
	var deps []*n.Node
	deps = append(deps, travNodes...)
	deps = append(deps, createdNodes...)
	for _, n := range deps {
		n.AddDependency(o.ID)
	}

	//TODO: let's assume the Mutate() call never fails for now.
	//TODO: how to rollback side effects of traverse if Mutate() fails?
	err := Mutate(nPtr, o)
	if err != nil {
		return d, err
	}

	d.OperationsId = append(d.OperationsId, o.ID)
	return d, nil
}

// Traverses the document from root element to the node indicated by the cursor
// input. When a path does not exist in the current document, create the node
// and link it to the document.
// The traverse function returns a pointer to the last node, a list of pointers
// of nodes traversed and a list of pointers of nodes created
func (d *Doc) traverse(cursor op.Cursor, opId string) (*n.Node, []*n.Node, []*n.Node) {
	var nPtr *n.Node
	var travNodes []*n.Node
	var createdNodes []*n.Node

	// traverse starts from headNode
	nPtr = d.Head

	for _, c := range cursor.Path {
		switch c.Type() {
		case op.MapT:
			k := c.Get().(string)
			nn, exists, _ := nPtr.GetChild(k)
			if !exists {
				nn = n.New(opId)
				nPtr.Add(k, nn, opId)
				createdNodes = append(createdNodes, nn)
			} else {
				travNodes = append(travNodes, nPtr)
			}
			nPtr = nn
		case op.ListT:
			k := c.Get().(int)
			nn, exists, _ := nPtr.GetChild(k)
			if !exists {
				nn = n.New(opId)
				nPtr.Add(k, nn, opId)
				createdNodes = append(createdNodes, nn)
			} else {
				travNodes = append(travNodes, nPtr)
			}
			nPtr = nn
		}
	}

	return nPtr, travNodes, createdNodes
}

func Mutate(node *n.Node, o op.Operation) error {
	mut := o.Mutation
	_ = o.Deps

	switch mut.Type {
	case op.Delete:
		// clear ops
		// TODO: clear all children of opDeps
		_ = allChildren(node)
		return nil
	case op.Assign:
		// clear ops
		// TODO: clear all children of opDeps
		_ = allChildren(node)
		// continue to the insertion
	}

	// Insert
	err := node.Add(mut.Key, mut.Value, o.ID)

	return err
}

func (d Doc) String() string {
	ids := fmt.Sprintf("ID: %v; ClockId: %v", d.Id, d.Clock)
	ops := fmt.Sprintf("Operations: applied: %v, buffered: %v", d.OperationsId, d.OperationsBuffer)
	node := fmt.Sprintf("Head: %v", d.Head)
	return fmt.Sprintf("%v\n%v\n%v\n", ids, ops, node)
}
