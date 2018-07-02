package main

import (
	"github.com/gpestana/rdoc/clock"
	"github.com/gpestana/rdoc/operation"
	"log"
)

func main() {
	uid := "a8fdc205a9f19cc1c7507a60c4f01b13d11d7fd0"
	obj := Init(uid)

	deps := []clock.Clock{}

	cursor := operation.NewCursor(
		operation.MapKey{"some"},
		operation.ListKey{0},
		operation.MapKey{"map"},
	)

	mut := operation.NewMutation(operation.Insert, "val")
	opId := obj.Clock.Timestamp()
	op, err := operation.New(opId, deps, cursor, mut)

	log.Printf("%+v\n", op)

	if err != nil {
		log.Fatal(err)
	}
}
