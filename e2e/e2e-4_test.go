// e2e tests map to the examples of the original paper
package rdoc

import (
	"fmt"
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
	opInsert1a, _ := op.New(id1+".1", []string{}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1a)

	mutDoc1, _ = op.NewMutation(op.Insert, 1, "b")
	opInsert1b, _ := op.New(id1+".2", []string{id1 + ".1"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1b)

	mutDoc1, _ = op.NewMutation(op.Insert, 2, "c")
	opInsert1c, _ := op.New(id1+".3", []string{id1 + ".1", id1 + ".2"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1c)

	// doc2: populates head of doc with ["a", "b", "c"] (through sync so that both
	// replicas have the same state)
	doc2.ApplyRemoteOperation(*opInsert1a)
	doc2.ApplyRemoteOperation(*opInsert1b)
	doc2.ApplyRemoteOperation(*opInsert1c)

	// doc1: delete element position 1 ("b")
	curDel := op.NewCursor(1, op.ListKey{1})
	mutDoc1, _ = op.NewMutation(op.Delete, nil, nil)
	opDelete1b, _ := op.New(id1+".4", []string{id1 + ".1", id1 + ".2", id1 + ".3"}, curDel, mutDoc1)
	doc1.ApplyOperation(*opDelete1b)

	// doc1: insert element "x" position 1
	mutDoc1, _ = op.NewMutation(op.Insert, 1, "x")
	opInsert1x, _ := op.New(id1+".5", []string{id1 + ".1", id1 + ".2", id1 + ".3", id1 + ".4"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opInsert1x)

	// doc1: initial verifications
	list1 := doc1.Head.List()

	elA1, _ := list1.Get(0)
	elA1Val, _ := elA1.(*n.Node).Reg().Get("1.1")
	elA1DepsLen := len(elA1.(*n.Node).Deps())
	if elA1Val != "a" {
		t.Error("doc1: element 0 must be value 'a', got ", elA1Val)
	}
	if elA1DepsLen != 1 {
		t.Error("doc1: element '0:a'  must have 1 dependency, got", elA1DepsLen)
	}

	elX1, _ := list1.Get(1)
	elX1Val, _ := elX1.(*n.Node).Reg().Get("1.5")
	elX1DepsLen := len(elX1.(*n.Node).Deps())
	if elX1Val != "x" {
		t.Error("doc1: element 1 must be value 'x', got ", elX1Val)
	}
	if elX1DepsLen != 5 {
		t.Error("doc1: element '1:b'  must have 0 dependencies, got", elX1DepsLen)
	}

	elB1, _ := list1.Get(2)
	elB1Val, _ := elB1.(*n.Node).Reg().Get("1.2")
	elB1DepsLen := len(elB1.(*n.Node).Deps())
	if elB1Val != "b" {
		t.Error("doc1: element 2 must be value 'b', got ", elB1Val)
	}
	if elB1DepsLen != 0 {
		t.Error("doc1: element '2:b'  must have 0 dependencies, got", elB1DepsLen)
	}

	elC1, _ := list1.Get(3)
	elC1Val, _ := elC1.(*n.Node).Reg().Get("1.3")
	elC1DepsLen := len(elC1.(*n.Node).Deps())
	if elC1Val != "c" {
		t.Error("doc1: element 3 must be value 'c', got ", elC1Val)
	}
	if elC1DepsLen != 3 {
		t.Error("doc1: element '3:c'  must have 3 dependencies, got", elC1DepsLen)
	}

	// doc2: insert element "y" position 0
	curDoc2 := op.NewEmptyCursor()
	mutDoc2, _ := op.NewMutation(op.Insert, 0, "y")
	opInsert2y, _ := op.New(id2+".1", []string{}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opInsert2y)

	// doc2: insert element "z" position 3
	mutDoc2, _ = op.NewMutation(op.Insert, 3, "z")
	opInsert2z, _ := op.New(id2+".2", []string{id2 + ".1"}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opInsert2z)

	// sync
	doc1.ApplyRemoteOperation(*opInsert2y)
	doc1.ApplyRemoteOperation(*opInsert2z)

	doc2.ApplyRemoteOperation(*opDelete1b)
	doc2.ApplyRemoteOperation(*opInsert1x)

	// verifications
	doc1List := doc1.Head.List()
	doc2List := doc2.Head.List()

	if doc1List.Size() != 6 {
		t.Error("doc1: size of list must be 6, got ", doc1List.Size())
	}

	if doc2List.Size() != 6 {
		t.Error("doc2: size of list must be 6, got ", doc2List.Size())
	}

	doc1Vals := []string{}
	for i := 0; i < doc1List.Size(); i++ {
		e, _ := list1.Get(i)
		el := e.(*n.Node)
		doc1Vals = append(doc1Vals, el.Reg().Values()[0].(string))
	}

	doc2Vals := []string{}
	for i := 0; i < doc2List.Size(); i++ {
		e, _ := doc2List.Get(i)
		el := e.(*n.Node)
		doc2Vals = append(doc2Vals, el.Reg().Values()[0].(string))
	}

	if len(doc1Vals) != len(doc2Vals) {
		t.Fatal(fmt.Sprintf("Lenght of doc1 and doc2 lists must be the same, got %v vs %v", doc1Vals, doc2Vals))
	}

	for i := 0; i < len(doc1Vals); i++ {
		if doc1Vals[i] != doc2Vals[i] {
			t.Error(fmt.Sprintf("Elements should be ordered equally, got: (%v:%v) vs (%v, %v)", i, doc1Vals[i], i, doc2Vals[i]))
		}
	}
}
