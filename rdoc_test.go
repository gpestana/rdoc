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
		op.MapKey{Key: "root"},
		op.MapKey{Key: "sub-node"},
		op.ListKey{Key: 0},
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
	opId := "opIdTest"
	v := "hello world"
	mut, _ := op.NewMutation(op.Insert, nil, v)
	op, _ := op.New(opId, []string{}, op.Cursor{}, mut)

	node := n.New("")

	err := Mutate(node, *op)
	if err != nil {
		t.Fatal(err)
	}

	mvr := node.GetMVRegister()
	if mvr[opId] != v {
		t.Error(fmt.Sprintf("MVR should be {opIdTest: 'hello world', got %v", mvr))
	}
}

func TestClearDeps(t *testing.T) {
	initDeps := []string{"1", "2", "3", "4"}
	removeDeps := []string{"2", "4"}
	node := n.New("test")
	node.SetDeps(initDeps)

	clearDeps([]*n.Node{node}, removeDeps)

	finalDeps := node.Deps()

	if len(finalDeps) != 2 {
		t.Error(fmt.Sprintf("Dependency set should have length 2 after clearing, got %v", len(finalDeps)))
	}

	if finalDeps[0] != initDeps[0] {
		t.Error(fmt.Sprintf("Dependency set element 0 should be %v after clearing, got %v", initDeps[0], finalDeps[0]))
	}

	if finalDeps[1] != initDeps[2] {
		t.Error(fmt.Sprintf("Dependency set element 1 should be %v after clearing, got %v", initDeps[1], finalDeps[2]))
	}

}
