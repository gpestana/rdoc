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

// Mutate()

func TestMutateList(t *testing.T) {
	val := "list_element"
	oId := "op_x"
	mut, _ := op.NewMutation(op.Insert, 0, val)
	o, err := op.New(oId, []string{}, op.NewCursor("key"), mut)
	if err != nil {
		t.Error(err)
	}
	node := newNode("key")
	node.Mutate(*o)

	// get list from head node
	l := node.GetList()
	if l.Size() != 1 {
		t.Error(fmt.Sprintf("List should have lenght 1 after mutation, it has %v", l.Size()))
	}
	// get register from head node's list
	nn, _ := l.Get(0)
	// check if correct register kv pair exists
	nreg, _ := nn.(*Node).reg.Get(oId)
	if nreg != val {
		t.Error(fmt.Sprintf("List should contain a register with k:v (%v:%v)", oId, val))
	}
}

func TestMutateRegister(t *testing.T) {
	val := "value"
	mut, _ := op.NewMutation(op.Assign, nil, val)
	opId := "operation_id"
	o, err := op.New(opId, []string{}, op.NewEmptyCursor(), mut)
	if err != nil {
		t.Error(t)
	}

	node := newNode("test")
	node.Mutate(*o)

	r := node.GetReg()
	if r.Size() != 1 {
		t.Error(fmt.Sprintf("Register map should have 1 element, got %v instead", r.Size()))
	}
	rval, _ := r.Get(opId)
	if rval != val {
		t.Error(fmt.Sprintf("Register map should have an element associated with %v", opId))
	}

	// write diff value in same register
	val2 := "value_2"
	mut2, _ := op.NewMutation(op.Assign, nil, val2)
	opId2 := "operation_id_2"
	o2, err := op.New(opId2, []string{}, op.NewEmptyCursor(), mut2)
	if err != nil {
		t.Error(t)
	}

	node.Mutate(*o2)
	r = node.GetReg()

	if r.Size() != 2 {
		t.Error(fmt.Sprintf("Register map should have 2 element, got %v instead", r.Size()))
	}

	rval, _ = r.Get(opId2)
	if rval != val2 {
		t.Error(fmt.Sprintf("Register map should have an element associated with %v", opId2))
	}

}
