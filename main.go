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

	cursor := []byte(`{"path": ["some", "list", 2, "aMap"], "id": 1}`)
	mut := operation.NewMutation(operation.Insert, "val")
	opId := obj.Clock.Timestamp()
	op, err := operation.New(opId, deps, cursor, mut)
	if err != nil {
		log.Fatal(err)
	}

	err = obj.NewOperation(op)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(obj)
}
