// Defines json-crdt document data structure and private methods
package main

import (
	"fmt"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
	"log"
)

type JSON struct {
	Clock clock.Clock
	Root  []types.CRDT
}

// Returns a new CRDT JSON data structure. It receives an ID which must be
// unique in the context of the network.
func New(uid string) *JSON {
	clk := clock.New([]byte(uid))

	return &JSON{
		Root:  []types.CRDT{},
		Clock: clk,
	}
}

// Checks whether operation to process is local or remote and calls the correct
// operation handler.
func (j *JSON) newOperation(op *operation.Operation) error {
	if op.NodeID() == j.Clock.ID() {
		j.handleLocalOperation(op)
	} else {
		j.handleRemoteOperation(op)
	}
	return nil
}

// Traverses the JSON document based on the operation cursor and applies the
// mutation on the correct leaf/node. During the traversal, it also handles the
// operation presence according to the mutation type
func (j *JSON) apply(op *operation.Operation) error {
	cpath := op.Cursor.Path

	// Selects correct node to apply the operation mutation. If necessary, create
	// the necesseary nodes so that the cursor can be traversable in the current
	// document
	err, node := traverse(j.Root, cpath)
	if err != nil {
		return err
	}

	// Apply operation mutation on the selected node
	mut := op.Mutation
	err = mutate(node, mut)
	if err != nil {
		return err
	}

	return nil
}

func traverse(root []types.CRDT, path []operation.PathItem) (error, types.CRDT) {
	return nil, nil
}

func mutate(node types.CRDT, mut operation.Mutation) error {
	return nil
}

func (j JSON) handleLocalOperation(op *operation.Operation) {
	err := j.apply(op)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error handling the operation, %v", err))
	}
}

func (j JSON) handleRemoteOperation(op *operation.Operation) {
	log.Println("AddRemoteOperation")
}
