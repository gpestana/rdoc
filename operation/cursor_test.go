package operation

import (
	"testing"
)

func TestCursor(t *testing.T) {
	k := "some_key"
	// c points at: map.anothermap[0]
	c := NewCursor(
		k,
		MapKey{"map"},
		MapKey{"another_map"},
		ListKey{0},
	)

	if c.Path[1].Type() != MapT {
		t.Error("Cursor should be of type MapT")
	}

	if c.Path[2].Type() != ListT {
		t.Error("Cursor should be of type ListT")
	}
}
