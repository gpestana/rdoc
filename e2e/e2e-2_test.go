// e2e tests map to the examples of the original paper
package rdoc

import (
	"fmt"
	"github.com/gpestana/rdoc"
	n "github.com/gpestana/rdoc/node"
	op "github.com/gpestana/rdoc/operation"
	"testing"
)

// Case B: Modifying the contents of a nested map while concurrently the entire
// map is overwritten.
func TestCaseB(t *testing.T) {
	id1, id2 := "1", "2"
	doc1, doc2 := rdoc.Init(id1), rdoc.Init(id2)

	// doc1: initial state: {"colors": { "blue": "#0000ff" }}
	curDoc1 := op.NewCursor("colors", op.MapKey{"colors"})
	mutDoc1, _ := op.NewMutation(op.Noop, nil, nil)
	opColors1, _ := op.New(id1+".1", []string{}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opColors1)

	curDoc1 = op.NewCursor("blue", op.MapKey{"colors"})
	mutDoc1, _ = op.NewMutation(op.Insert, "blue", "#0000ff")
	opBlue1, _ := op.New(id1+".2", []string{id1 + ".1"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opBlue1)

	// doc2: initial state: {"colors": { "blue": "#0000ff" }}
	curDoc2 := op.NewCursor("colors", op.MapKey{"colors"})
	mutDoc2, _ := op.NewMutation(op.Noop, nil, nil)
	opColors2, _ := op.New(id2+".1", []string{}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opColors2)

	curDoc2 = op.NewCursor("blue", op.MapKey{"colors"})
	mutDoc2, _ = op.NewMutation(op.Insert, "blue", "#0000ff")
	opBlue2, _ := op.New(id2+".2", []string{id2 + ".1"}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opBlue2)

	// sync
	doc1.ApplyRemoteOperation(*opColors2)
	doc1.ApplyRemoteOperation(*opBlue2)

	doc2.ApplyRemoteOperation(*opColors1)
	doc2.ApplyRemoteOperation(*opBlue1)

	// doc1: insert KV {"red": "#ff0000"} to map
	curDoc1 = op.NewCursor("red", op.MapKey{"colors"})
	mutDoc1, _ = op.NewMutation(op.Insert, "red", "#ff0000")
	opRed, _ := op.New(id1+".3", []string{id1 + ".1", id1 + ".2", id2 + ".1", id2 + ".2"}, curDoc1, mutDoc1)
	doc1.ApplyOperation(*opRed)

	// doc2: 1) clear map "colors"; 2) insert {"green": "#00ff00"} to "colors" map
	nEmpty := n.New("colors")
	curDoc2 = op.NewEmptyCursor()
	mutDoc2, _ = op.NewMutation(op.Assign, "color", nEmpty)
	opClear, _ := op.New(id2+".3", []string{id1 + ".1", id1 + ".2", id2 + ".1", id2 + ".2"}, curDoc2, mutDoc2)
	doc2.ApplyOperation(*opClear)

	curDoc2 = op.NewCursor("colors", op.MapKey{"colors"})
	mutDoc2, _ = op.NewMutation(op.Insert, "green", "#00ff00")
	opGreen, _ := op.New(id2+".4", []string{id1 + ".1", id1 + ".2", id2 + ".1", id2 + ".2", id2 + ".3"}, curDoc2, mutDoc2)

	doc2.ApplyOperation(*opGreen)

	// sync again
	doc1.ApplyRemoteOperation(*opClear)
	doc1.ApplyRemoteOperation(*opGreen)
	doc2.ApplyRemoteOperation(*opRed)

	// doc1: verifications
	doc1ColorsIf, _ := doc1.Head.Map().Get("colors")
	doc1Colors := doc1ColorsIf.(*n.Node)

	if doc1Colors.Map().Size() != 3 {
		t.Error(fmt.Printf("doc1.colors should have 3 elements, got %v", doc1Colors.Map().Size()))
	}

	doc1Blue, _ := doc1Colors.Map().Get("blue")
	doc1BlueDeps := doc1Blue.(*n.Node).Deps()
	if len(doc1BlueDeps) != 0 {
		t.Error("doc1: dependencies of 'blue' must be all cleared, got ", doc1BlueDeps)
	}

	doc1Red, _ := doc1Colors.Map().Get("red")
	doc1RedDeps := doc1Red.(*n.Node).Deps()
	if len(doc1RedDeps) == 0 {
		t.Error("doc1: dependencies of 'red' must NOT be all cleared, got ", doc1RedDeps)
	}

	doc1Green, _ := doc1Colors.Map().Get("green")
	doc1GreenDeps := doc1Green.(*n.Node).Deps()
	if len(doc1GreenDeps) == 0 {
		t.Error("doc1: dependencies of 'green' must NOT be all cleared, got ", doc1GreenDeps)
	}

	// doc2: verifications
	doc2ColorsIf, _ := doc2.Head.Map().Get("colors")
	doc2Colors := doc2ColorsIf.(*n.Node)

	if doc2Colors.Map().Size() != 3 {
		t.Error(fmt.Printf("doc2.colors should have 3 elements, got %v", doc2Colors.Map().Size()))
	}

	doc2Blue, _ := doc2Colors.Map().Get("blue")
	doc2BlueDeps := doc2Blue.(*n.Node).Deps()
	if len(doc2BlueDeps) != 0 {
		t.Error("doc2: dependencies of 'blue' must be all cleared, got ", doc2BlueDeps)
	}

	doc2Red, _ := doc2Colors.Map().Get("red")
	doc2RedDeps := doc2Red.(*n.Node).Deps()
	if len(doc2RedDeps) == 0 {
		t.Error("doc2: dependencies of 'red' must NOT be all cleared, got ", doc2RedDeps)
	}

	doc2Green, _ := doc2Colors.Map().Get("green")
	doc2GreenDeps := doc2Green.(*n.Node).Deps()
	if len(doc2GreenDeps) == 0 {
		t.Error("doc2: dependencies of 'green' must NOT be all cleared, got ", doc2GreenDeps)
	}
}
