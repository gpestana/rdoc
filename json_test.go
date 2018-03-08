package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalEmpty(t *testing.T) {
	jobj := New("clock1")
	expEncStr := "{}"

	enc, err := json.Marshal(jobj)
	if err != nil {
		t.Error(err)
	}

	if string(enc) != expEncStr {
		t.Error(fmt.Sprintf("Marshal: expected %v, got %v", expEncStr, enc))
	}
}
