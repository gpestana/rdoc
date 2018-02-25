package main

import (
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"log"
)

func main() {
	obj := New()

	id := clock.New("ts_seed")
	deps := []clock.Timestamp{}
	cursor := []interface{}{"root", "level_A", "level_AB", 1}
	mut := operation.NewMutation("key_1", "value")

	op := operation.New(id, deps, cursor, mut)
	obj.AddLocalOperation(op)

	log.Println(obj)
}
