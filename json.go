package main

import (
	"encoding/json"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
	"log"
)

type JSON struct {
	// A valid JSON can by a map, list or register
	Head  types.CRDT
	Clock clock.Clock

	// Channel in which operation to be broadcast remotely are buffered
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
		handleLocalOperation(op)
	} else {
		handleRemoteOperation(op)
	}
	return nil
}

func handleLocalOperation(op *operation.Operation) {
	log.Println("AddLocalOperation")
}

func handleRemoteOperation(op *operation.Operation) {
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
