package main

import (
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
)

type JSON struct {
	Head types.CRDT

	// Queue of operations to be applied locally
	opBuffer []*operation.Operation
	// Queue of received remote operations to be applied (if not yet)
	rcvBuffer []*operation.Operation
	// Queue of operations applied locally and ready to propagate over the network
	sendBuffer []*operation.Operation
}

func New() *JSON {
	// TODO:
	// Start one go routine per buffer when bufferd chans are used
	// Or expect this to happen in the upper layer?
	rcvBuf := []*operation.Operation{}
	sendBuf := []*operation.Operation{}
	return &JSON{
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
	return []byte{}, nil
}
