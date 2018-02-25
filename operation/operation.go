package operation

import (
	"github.com/gpestana/crdt-json/clock"
)

type Operation struct {
	id       clock.Timestamp
	deps     []clock.Timestamp
	cursor   []interface{}
	mutation Mutation
}

func New(id clock.Timestamp, deps []clock.Timestamp, c []interface{}, m Mutation) *Operation {
	return &Operation{
		id:       id,
		deps:     deps,
		cursor:   c,
		mutation: m,
	}
}

type Mutation struct {
	key   interface{}
	value interface{}
}

func NewMutation(k interface{}, v interface{}) Mutation {
	return Mutation{k, v}
}
