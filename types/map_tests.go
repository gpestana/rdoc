package types

import (
	"testing"
)

func TestNewEmptyMap(t *testing.T) {
	m := NewEmptyMap()
	if len(m.KV) != 0 {
		t.Error("Map: New empty map should be empty")
	}
}
