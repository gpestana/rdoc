package main

import (
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"log"
)

func main() {
	obj := New()

	clk := clock.New([]byte("a8fdc205a9f19cc1c7507a60c4f01b13d11d7fd0"))
	deps := []clock.Clock{}

	cursor := []interface{}{"root", "level_A", "level_AB", 1}
	id := clk.Timestamp()
	mut := operation.NewMutation(operation.Insert, "val")

	op := operation.New(id, deps, cursor, mut)
	obj.AddLocalOperation(op)

	log.Println(obj)
}
