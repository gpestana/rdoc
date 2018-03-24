package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalEmpty(t *testing.T) {
	jobj := New("clock1")
	expEncStr := "{}"
	expBytes, err := json.Marshal(expEncStr)
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
