package rdoc

import (
	"fmt"
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
	n, trvN, crtN := doc.traverse(cursor1)

	ni1, _ := doc.Head.hmap.Get("root")
	ni2, _ := ni1.(*Node).hmap.Get("sub-node")
	ni3, _ := ni2.(*Node).list.Get(0)

	if n != ni3 {
		t.Error(fmt.Sprintf("Node created %v should be the same as in document %v", n, ni3))
	}

	if len(trvN) != 3 {
		t.Error(fmt.Sprintf("Set of traversed nodes should have length 3, got %v", len(trvN)))
	}

	for i, _ := range trvN {
		if trvN[i] != crtN[i] {
			t.Error(fmt.Sprintf("All traversed nodes were created. Got %v and %v", trvN, crtN))
		}
	}

	if ni1 != trvN[0] {
		t.Error(fmt.Sprintf("1st node traversed should be %v, got %v", ni1, trvN[0]))
	}
	if ni2 != trvN[1] {
		t.Error(fmt.Sprintf("1st node traversed should be %v, got %v", ni2, trvN[1]))
	}
	if ni3 != trvN[2] {
		t.Error(fmt.Sprintf("1st node traversed should be %v, got %v", ni3, trvN[2]))
	}
}

func TestAllChildren(t *testing.T) {
	docId := "doc1"
	doc := Init(docId)

	cursor1 := op.NewCursor(
		"key", // cursor's key
		op.MapKey{"root"},
		op.MapKey{"sub-node"},
		op.ListKey{0},
		op.MapKey{"some"},
	)
	doc.traverse(cursor1)

	ni1, _ := doc.Head.hmap.Get("root")
	ni2, _ := ni1.(*Node).hmap.Get("sub-node")
	ni3, _ := ni2.(*Node).list.Get(0)
	ni4, _ := ni3.(*Node).hmap.Get("some")

	expCh := 3
	ch := ni1.(*Node).allChildren()
	if len(ch) != expCh {
		t.Error(fmt.Sprintf("root of document must have %v children, got %v", expCh, len(ch)))
	}

	nPtrs := []*Node{ni3.(*Node), ni4.(*Node)}
	ch2 := ni2.(*Node).allChildren()
	for i, _ := range ch2 {
		if ch2[i] != nPtrs[i] {
			t.Error(fmt.Sprintf("fetching all childrens is not behaving as expected"))
		}
	}
}
