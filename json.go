package main

import (
	"encoding/json"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
)

type JSON struct {
	// A valid JSON can by a map, list or register
	Head  types.CRDT
	Clock clock.Clock

	// Queue of operations to be applied locally
	opBuffer []*operation.Operation
	// Queue of received remote operations to be applied (if not yet)
	rcvBuffer []*operation.Operation
	// Queue of operations applied locally and ready to propagate over the network
	sendBuffer []*operation.Operation
}

// Returns a new CRDT JSON data structure. It receives an ID which must be
// unique in the context of the network.
func New(uid string) *JSON {
	m := types.NewMap()
	clk := clock.New([]byte(uid))
	rcvBuf := []*operation.Operation{}
	sendBuf := []*operation.Operation{}
	return &JSON{
		Head:       *m,
		Clock:      clk,
		rcvBuffer:  rcvBuf,
		sendBuffer: sendBuf,
	}
}

func (j *JSON) AddLocalOperation(op *operation.Operation) {
	j.opBuffer = append(j.opBuffer, op)
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
