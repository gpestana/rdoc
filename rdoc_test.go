package rdoc

import (
	"fmt"
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

func TestTraverseSimple(t *testing.T) {
	docId := "doc1"
	doc := Init(docId)

	cursor1 := op.NewCursor(
		0, // cursor's key
		op.MapKey{"root"},
		op.MapKey{"sub-node"},
		op.ListKey{0},
	)

	n1, trvN, crtN := doc.traverse(cursor1, "oId")

	if len(crtN) != 3 {
		t.Error(fmt.Sprintf("There should be 3 created nodes after first traversal, got %v", len(crtN)))
	}

	if len(trvN) != 0 {
		t.Error(fmt.Sprintf("There should be 0 traversed nodes after first traversal, got %v", len(trvN)))
	}

	n2, trvN, crtN := doc.traverse(cursor1, "opId")

	if len(crtN) != 0 {
		t.Error(fmt.Sprintf("There should be 0 created nodes after second traversal, got %v", len(crtN)))
	}

	if len(trvN) != 3 {
		t.Error(fmt.Sprintf("There should be 3 traversed nodes after second traversal, got %v", len(trvN)))
	}

	if n1 != n2 {
		t.Error(fmt.Sprintf("Both traverses should have finished in the same node, finished at: %v, %v instead", n1, n2))
	}
}

func TestMutateInsert(t *testing.T) {
	opId := "opBasic"
	k := "hello"
	v := "world"
	mut, _ := op.NewMutation(op.Insert, k, v)
	op, _ := op.New(opId, []string{}, op.Cursor{}, mut)

	node := n.New("")

	Mutate(node, *op)

	// TODO: finish
	mvr := node.GetMVRegister()
	fmt.Println(mvr)
}
