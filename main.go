package main

import (
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"log"
)

func main() {
	uid := "a8fdc205a9f19cc1c7507a60c4f01b13d11d7fd0"
	obj := New(uid)

	deps := []clock.Clock{}

	cursor := []byte(`[{"MapT": "some"}, {"MapT" :"aMap"}, {"ListT": 2}]`)
	mut := operation.NewMutation(operation.Insert, "val")
	opId := obj.Clock.Timestamp()
	op, err := operation.New(opId, deps, cursor, mut)

	log.Printf("%+v\n", op)

	if err != nil {
		log.Fatal(err)
	}

	err = obj.newOperation(op)
	if err != nil {
		log.Fatal(err)
	}
}
