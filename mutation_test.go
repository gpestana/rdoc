package rdoc

import (
	"fmt"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

func TestMutateList(t *testing.T) {
	val := "list_element"
	mut, _ := op.NewMutation(op.Insert, 0, val)
	o, err := op.New("", []string{}, op.NewCursor("key"), mut)
	if err != nil {
		t.Error(err)
	}
	node := newNode("key")
	node.Mutate(*o)

	l := node.GetList()
	if l.Size() != 1 {
		t.Error(fmt.Sprintf("List should have lenght 1 after mutation, it has %v", l.Size()))
	}

	if !l.Contains(val) {
		t.Error(fmt.Sprintf("List should contain %v after mutation but it does not", val))
	}

}
