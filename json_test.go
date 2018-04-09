package main

import (
	"encoding/json"
	"fmt"
	"github.com/gpestana/crdt-json/clock"
	"github.com/gpestana/crdt-json/operation"
	"github.com/gpestana/crdt-json/types"
	"testing"
)

func TestMarshalEmpty(t *testing.T) {
	jobj := New("clock1")
	expBytes, err := json.Marshal(map[string]types.Node{})
	if err != nil {
		t.Fatal(err)
	}

	enc, err := json.Marshal(jobj)
	if err != nil {
		t.Fatal(err)
	}

	for i := range enc {
		if enc[i] != expBytes[i] {
			t.Error(fmt.Sprintf("Marshal: expected %v, got %v", string(expBytes), string(enc)))
			break
		}
	}
}

func TestAddMapEmptyOperation(t *testing.T) {
	jobj := New("clock2")

	deps := []clock.Clock{}
	cursor := []byte(`[]`)
	mut := operation.NewMutation(operation.Insert, "")
	opId := jobj.Clock.Timestamp()

	op, err := operation.New(opId, deps, cursor, mut)
	if err != nil {
		t.Fatal(err)
	}
	err = jobj.NewOperation(op)
	if err != nil {
		t.Error(err)
	}

	head := jobj.Head.(types.Map)
	expHeadKeys := []string{}
	if len(head.Keys()) != len(expHeadKeys) {
		t.Error(fmt.Sprintf("Head keys expected to be %v, got %v", head.Keys(), expHeadKeys))
	}
}

func TestAddMapOperation(t *testing.T) {
	jobj := New("clock2")

	deps := []clock.Clock{}
	cursor := []byte(`[{"MapT": "map1"}]`)
	mut := operation.NewMutation(operation.Insert, "")
	opId := jobj.Clock.Timestamp()

	op, err := operation.New(opId, deps, cursor, mut)
	if err != nil {
		t.Fatal(err)
	}
	err = jobj.NewOperation(op)
	if err != nil {
		t.Error(err)
	}

	head := jobj.Head.(types.Map)
	expHeadKeysLen := 1
	if len(head.Keys()) != expHeadKeysLen {
		t.Error(fmt.Sprintf("Head keys length expected to be %v, got %v", expHeadKeysLen, len(head.Keys())))
	}

}
