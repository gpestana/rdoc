// e2e tests map to the examples of the original paper
package rdoc

import (
	"github.com/gpestana/rdoc"
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

// Case D: 4. Concurrent editing of an ordered list of characters (i.e., a text
// document).
func TestCaseD(t *testing.T) {
	id1, id2 := "1", "2"
	doc1, doc2 := rdoc.Init(id1), rdoc.Init(id2)

	// doc1: populates head of doc with ["a", "b", "c"]
	curDoc1 := op.NewEmptyCursor()
	mutDoc1, _ := op.NewMutation(op.Insert, 0, "a")
	opInsert1, _ := op.New(id1+".1", []string{}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1)

	mutDoc1, _ = op.NewMutation(op.Insert, 1, "b")
	opInsert1, _ = op.New(id1+".2", []string{id1 + ".1"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1)

	mutDoc1, _ = op.NewMutation(op.Insert, 2, "c")
	opInsert1, _ = op.New(id1+".3", []string{id1 + ".1", id1 + ".2"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1)

	// doc1: delete element position 1 ("b")
	curDel := op.NewCursor(1, op.ListKey{1})
	mutDoc1, _ = op.NewMutation(op.Delete, nil, nil)
	opInsert1, _ = op.New(id1+".4", []string{id1 + ".1", id1 + ".2", id1 + ".3"}, curDel, mutDoc1)
	doc1.ApplyOperation(*opInsert1)

	// doc1: insert element "x" position 1
	mutDoc1, _ = op.NewMutation(op.Insert, 1, "x")
	opInsert1, _ = op.New(id1+".5", []string{id1 + ".1", id1 + ".2", id1 + ".3", id1 + ".4"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1)

	// doc1: initial verifications
	list1 := doc1.Head.List()

	for i := 0; i < list1.Size(); i++ {
		t.Error(list1.Get(i))
	}

	elA1, _ := list1.Get(0)
	t.Error(elA1.(*n.Node).Reg())
	elA1Val, _ := elA1.(*n.Node).Reg().Get("1.1")
	elA1DepsLen := len(elA1.(*n.Node).Deps())
	if elA1Val != "a" {
		t.Error("doc1: element 0 must be value 'a', got ", elA1Val)
	}
	if elA1DepsLen != 1 {
		t.Error("doc1: element '0:a'  must have 1 dependency, got", elA1DepsLen)
	}

	elB1, _ := list1.Get(1)
	t.Error(elB1.(*n.Node).Reg())
	elB1Val, _ := elB1.(*n.Node).Reg().Get("1.2")
	elB1DepsLen := len(elB1.(*n.Node).Deps())
	if elB1Val != "b" {
		t.Error("doc1: element 1 must be value 'b', got ", elB1Val)
	}
	if elB1DepsLen != 0 {
		t.Error("doc1: element '1:b'  must have 0 dependencies, got", elB1DepsLen)
	}

	elC1, _ := list1.Get(2)
	t.Error(elC1.(*n.Node).Reg())
	elC1Val, _ := elC1.(*n.Node).Reg().Get("1.3")
	elC1DepsLen := len(elB1.(*n.Node).Deps())
	if elC1Val != "c" {
		t.Error("doc1: element 2 must be value 'c', got ", elC1Val)
	}
	if elC1DepsLen != 3 {
		t.Error("doc1: element '2:c'  must have 3 dependencies, got", elC1DepsLen)
	}

	// doc2: populates head of doc with ["a", "b", "c"]
	curDoc2 := op.NewEmptyCursor()
	mutDoc2, _ := op.NewMutation(op.Insert, 0, "a")
	opInsert2, _ := op.New(id2+".1", []string{}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opInsert2)

	mutDoc2, _ = op.NewMutation(op.Insert, 1, "b")
	opInsert2, _ = op.New(id2+".2", []string{id1 + ".1"}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opInsert2)

	mutDoc2, _ = op.NewMutation(op.Insert, 3, "c")
	opInsert2, _ = op.New(id2+".2", []string{id1 + ".1", id1 + ".2"}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opInsert2)

	// doc2: insert element "y" position 0
	// doc2: insert element "x" position 3

	// sync

	// verifications

}
