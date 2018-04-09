// Implements CRDT operations.
// Everytime the state of the JSON document changes, an operations that
// describes the mutation is generated. The operation is kept for conflict
// resolution. The document state is the set of all operations ran against
// itself.
package operation

import (
	"github.com/gpestana/crdt-json/clock"
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
	deps []clock.Clock
	// Ambiguously identifies the position in the JSON object to apply the
	// operation by describing a path from the root of the document tree to some
	// branch or leaf node
	Cursor Cursor
	// Mutation requested at the specific operation's position
	Mutation Mutation
}

// Returns new Operation object
func New(id string, deps []clock.Clock, c []byte, m Mutation) (*Operation, error) {
	cursor, err := newCursor(c)
	if err != nil {
		return &Operation{}, err
	}
	return &Operation{
		ID:       id,
		deps:     deps,
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
	// Type of the mutation typ := {insert(v), delete, assign(v)}
	typ int
	// Value of the mutation; Value can be {string, int, list, obj}
	value interface{}
}

// Returns new Mutation
func NewMutation(t int, v interface{}) Mutation {
	return Mutation{t, v}
}
