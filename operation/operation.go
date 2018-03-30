// Implements CRDT operations.
// Everytime the state of the JSON document changes, an operations that
// describes the mutation is generated. The operation is kept for conflict
// resolution. The document state is the set of all operations ran against
// itself.
package operation

import (
	"errors"
	"fmt"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/types"
	"reflect"
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
	cursor []interface{}
	// Mutation requested at the specific operation's position
	mutation Mutation
}

// Returns new Operation object
func New(id string, deps []clock.Clock, c []interface{}, m Mutation) *Operation {
	return &Operation{
		ID:       id,
		deps:     deps,
		cursor:   c,
		mutation: m,
	}
}

// Returns ID of the node which generated the operation
func (op Operation) NodeID() string {
	splId := strings.Split(op.ID, ".")
	seed := splId[1]
	return seed
}

// A cursor identifies unambiguous a position in the JSON document by describing
// the path from the root until the leaf/node selected and the element ID
type Cursor struct {
	Keys []cursorKey
	Id   string
}

type cursorKey struct {
	Type  string
	Value interface{}
}

func newCursorKey(t string, v interface{}) (cursorKey, error) {
	if t != types.MapT || t != types.ListT || t != types.RegisterT {
		return cursorKey{}, errors.New(
			fmt.Sprintf("Type of path not valid; It should be one of {%v, %v, %v}",
				types.MapT, types.ListT, types.RegisterT))
	}

	valT := reflect.TypeOf(v)
	switch {
	case valT.Kind() == reflect.Int:
		if t != types.ListT {
			return cursorKey{}, errors.New(
				fmt.Sprintf("Cursor value type int is valid to cursor type ListT, got %v", t))
		}

	case valT.Kind() == reflect.String:
		return cursorKey{}, errors.New(
			fmt.Sprintf("Cursor value type string is valid to cursor type MapT, got %v", t))

	default:
		return cursorKey{}, errors.New("Cursor value type can be one of string or int")
	}

	return cursorKey{t, v}, nil
}

func (op Operation) NewCursor(path []map[string]interface{}, id string) (Cursor, error) {
	c := Cursor{Keys: []cursorKey{}, Id: id}
	for _, p := range path {
		ck, err := newCursorKey(p["type"].(string), p["id"])
		if err != nil {
			return Cursor{}, err
		}
		c.Keys = append(c.Keys, ck)
	}

	return Cursor{}, nil
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
