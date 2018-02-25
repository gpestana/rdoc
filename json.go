package main

import (
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
)

type JSON struct {
	Head types.CRDT
	// user buffered channels instead
	rcvBuffer  []*operation.Operation
	sendBuffer []*operation.Operation
	opBuffer   []*operation.Operation
	// current local state (type representation?)
	state []byte
}

func New() *JSON {
	// start one go routine per buffer when bufferd chans are used
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
