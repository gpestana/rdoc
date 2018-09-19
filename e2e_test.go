// e2e tests map to the examples of the original paper
package rdoc

import (
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

// Case A: different value assignment of a register in different replicas
func TestCaseA(t *testing.T) {
	id1 := "1"
	doc1 := Init(id1)

	id2 := "2"
	doc2 := Init(id2)
	emptyC := op.NewEmptyCursor()

	// contructs operation to initially populate the docs
	nmap1 := n.New("op-doc1")
	mut1, _ := op.NewMutation(op.Assign, "key", nmap1)
	op1, _ := op.New(id1+".0", []string{}, emptyC, mut1) // using id1 means that the operation was generated by id1

	nmap2 := n.New("op-doc2")
	mut2, _ := op.NewMutation(op.Assign, "key", nmap2)
	op2, _ := op.New(id1+".0", []string{}, emptyC, mut2)

	_, err := doc1.ApplyOperation(*op1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = doc2.ApplyRemoteOperation(*op2)
	if err != nil {
		t.Fatal(err)
	}

	// constructs and applies locally operation from replica 1
	mut3, _ := op.NewMutation(op.Assign, nil, "B")
	cur3 := op.NewCursor("key", op.MapKey{"key"})
	op3, _ := op.New(id1+".1", []string{id1 + ".0"}, cur3, mut3)
	doc1.ApplyOperation(*op3)

	// constructs and applies locally operation for replica 2
	mut4, _ := op.NewMutation(op.Assign, nil, "C")
	cur4 := op.NewCursor("key", op.MapKey{"key"})
	op4, _ := op.New(id2+".1", []string{id1 + ".0"}, cur4, mut4)
	doc2.ApplyOperation(*op4)

	// at this moment, we have:
	// doc1: {"key": MVR{1.1: "B"}}
	// doc2: {"key": MVR{2.1: "C"}}

	// network communication: cross-apply operations in replica 1 and 2
	doc1.ApplyRemoteOperation(*op4)
	doc2.ApplyRemoteOperation(*op3)

	// after network communication, we have:
	// doc1: {"key": MVR{1.1: "B", 2.1: "C"}}
	// doc2: {"key": MVR{2.1: "C", 1.1: "b"}}

	if len(doc1.Head.Map().Values()) != 1 {
		t.Error("In doc1, lenght of Head.Map should be 1")
	}

	if len(doc2.Head.Map().Values()) != 1 {
		t.Error("In doc2, lenght of Head.Map should be 2")
	}

	if len(nmap1.GetMVRegister()) != 2 {
		t.Error("In doc1, lenght of MVRegister should be 2")
	}

	if len(nmap2.GetMVRegister()) != 2 {
		t.Error("In doc2, lenght of MVRegister should be 2")
	}
}