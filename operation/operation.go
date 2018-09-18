// Implements CRDT operations.
// Everytime the state of the JSON document changes, an operations that
// describes the mutation is generated. The operation is kept for conflict
// resolution. The document state is the set of all operations ran against
// itself.
package operation

import (
	"strings"
)

const (
	Insert = iota
	Delete
	Assign
)

type Operation struct {
	// Lamport timestamp (implemented in clock.Clock) which uniquely identifies
	// the operation in the network
	ID string
	// Set of casual dependencies of the operation (all operations that
	// happened before the current operation)
	Deps []string
	// Ambiguously identifies the position in the JSON object to apply the
	// operation by describing a path from the root of the document tree to some
	// branch or leaf node
	Cursor Cursor
	// Mutation requested at the specific operation's position
	Mutation Mutation
}

// Returns new Operation object
func New(id string, deps []string, cursor Cursor, m Mutation) (*Operation, error) {
	return &Operation{
		ID:       id,
		Deps:     deps,
		Cursor:   cursor,
		Mutation: m,
	}, nil
}

// Returns ID of the node which generated the operation
func (op Operation) NodeID() string {
	splId := strings.Split(op.ID, ".")
	seed := splId[1]
	return seed
}

type Mutation struct {
	// Type of the mutation. Can be one of {insert(v), delete, assign(v)}
	Type  int
	Key   interface{}
	Value interface{}
}

// Returns new Mutation
func NewMutation(typ int, k interface{}, v interface{}) (Mutation, error) {
	// TODO: verify that mutation Value is of type of string or int
	return Mutation{Type: typ, Key: k, Value: v}, nil
}
