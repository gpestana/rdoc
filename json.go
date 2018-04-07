package main

import (
	"encoding/json"
	"fmt"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
	"log"
)

type JSON struct {
	// A valid JSON can by a map, list or register
	Head  types.CRDT
	Clock clock.Clock
	// Buffer for operation to be broadcast remotely
	BroadcastBuffer chan *operation.Operation
}

// Returns a new CRDT JSON data structure. It receives an ID which must be
// unique in the context of the network.
func New(uid string) *JSON {
	m := types.NewMap()
	clk := clock.New([]byte(uid))

	return &JSON{
		Head:            *m,
		Clock:           clk,
		BroadcastBuffer: make(chan *operation.Operation),
	}
}

// Checks whether operation to process is local or remote and calls the correct
// operation handler.
func (j *JSON) NewOperation(op *operation.Operation) error {
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
	node := j.Head
	for _, p := range cpath {
		switch n := node.(type) {
		case types.Map:
			id := p.MapT
			node = node.(types.Map).Get(id)
			log.Println("MapT", id, n)
		case types.List:
			id := p.ListT
			node = node.(types.List).Get(id)
			log.Println("ListT", id, n)
		case types.Register:
			id := p.RegisterT
			//err, node = node.(types.Register).Get(id)
			log.Println("RegisterT", id, n)
		}
	}
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

// CRDT JSON can be marshaled into a valid JSON encoding for application
// consumption
func (j JSON) MarshalJSON() ([]byte, error) {
	b := []byte{}
	b, err := json.Marshal(j.Head)
	if err != nil {
		return b, err
	}
	return b, nil
}

// Handler for operations ready to be applied locally
func (j JSON) handleLocalOp(op *operation.Operation) {
	log.Println("Handle local operation")
}

// Handler for operations ready to be broacast to remote nodes
func (j JSON) handleSendOp(op *operation.Operation) {
	log.Println("Handle receive operation")
}

// Handler for operations received from remote nodes
func (j JSON) handleRcvOp(op *operation.Operation) {
	log.Println("Handle receive operation")
}
