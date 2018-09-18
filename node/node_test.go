package node

import (
	"fmt"
	"testing"
)

func TestNodeConstruction(t *testing.T) {
	n1, n2, n3 := New("n1"), New("n2"), New("n3")

	// { "hello": ["world"] }
	err := n1.Add("hello", n2, "op1")
	checkError(err, t)
	err = n2.Add(0, n3, "op2")
	checkError(err, t)
	err = n3.Add(nil, "world", "op3")
	checkError(err, t)

	n2r, exists, err := n1.GetChild("hello")
	checkError(err, t)
	if !exists {
		t.Error("Key 'hello' should exist in n1 map")
	}

	n3r, exists, err := n2r.GetChild(0)
	checkError(err, t)
	if !exists {
		t.Error("Key 0 should exist in n2 list")
	}

	mvr := n3r.GetMVRegister()
	if mvr["op3"] != "world" {
		t.Error(fmt.Sprintf("Register should be map['op3']:world, got %v", mvr))
	}
}

func TestDependenciesMgmt(t *testing.T) {
	n := New("node")
	d1, d2, d3 := "d1", "d2", "d3"
	n.AddDependency(d1)
	n.AddDependency(d2)
	n.AddDependency(d3)

	if len(n.deps) != 3 {
		t.Error(fmt.Sprintf("Node should have 3 dependencies, got %v", len(n.deps)))
	}

	n.ClearDependency(d1)
	n.ClearDependency(d2)

	if len(n.deps) != 1 {
		t.Error(fmt.Sprintf("Node should have 1 dependency after clearing, got %v", len(n.deps)))
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error("Error should not happen in this case: ", err)
	}
}
