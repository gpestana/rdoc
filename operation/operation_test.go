package operation

import (
	"fmt"
	"github.com/gpestana/rdoc/clock"
	"testing"
)

func TestNodeID(t *testing.T) {
	nid := "123123123"
	cursor := NewCursor()
	op, err := New("10."+nid, []clock.Clock{}, cursor, Mutation{})
	if err != nil {
		t.Fatal(err)
	}

	actualNid := op.NodeID()
	if actualNid != nid {
		t.Error(fmt.Sprintf("Expected Node ID %v, had %v", nid, actualNid))
	}
}
