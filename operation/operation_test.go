package operation

import (
	"fmt"
	"github.com/gpestana/crdt-json/clock"
	"testing"
)

func TestNodeID(t *testing.T) {
	nid := "123123123"
	op := New("10."+nid, []clock.Clock{}, []interface{}{}, Mutation{})

	actualNid := op.NodeID()
	if actualNid != nid {
		t.Error(fmt.Sprintf("Expected Node ID %v, had %v", nid, actualNid))
	}
}
