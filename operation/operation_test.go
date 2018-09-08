package operation

import (
	"fmt"
	"testing"
)

func TestNodeID(t *testing.T) {
	nid := "123123123"
	cursor := NewCursor(0)
	op, err := New("10."+nid, []string{}, cursor, Mutation{})
	if err != nil {
		t.Fatal(err)
	}

	actualNid := op.NodeID()
	if actualNid != nid {
		t.Error(fmt.Sprintf("Expected Node ID %v, had %v", nid, actualNid))
	}
}
