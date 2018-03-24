package main

import (
	"encoding/json"
	"fmt"
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
