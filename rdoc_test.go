package rdoc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	jsonp "github.com/evanphx/json-patch"
	"github.com/gpestana/rdoc"
)

// 1. end-to-end tests mapping to the examples of the JSON CRDT paper

// 1.A: different value assignment of a register on multiple replicas
func Test_E2E_One(t *testing.T) {
	expectedFinalDoc1 := `[{"id":"2.64487784","op":"add","path":"/key","value":"A"},{"id":"3.64487784","op":"add","path":"/key","value":"B"},{"id":"2.64553321","op":"add","path":"/key","value":"A"},{"id":"3.64553321","op":"add","path":"/key","value":"C"}]`
	expectedFinalDoc2 := `[{"id":"2.64553321","op":"add","path":"/key","value":"A"},{"id":"3.64553321","op":"add","path":"/key","value":"C"},{"id":"2.64487784","op":"add","path":"/key","value":"A"},{"id":"3.64487784","op":"add","path":"/key","value":"B"}]`

	doc1 := rdoc.Init("doc1")
	doc2 := rdoc.Init("doc2")

	docInitPatch := `[{"op": "add", "path": "/key", "value": "A"}]`

	// apply initial state locally on both replicas
	err := doc1.Apply([]byte(docInitPatch))
	if err != nil {
		t.Error(err)
	}

	err = doc2.Apply([]byte(docInitPatch))
	if err != nil {
		t.Error(err)
	}

	// change value of `key` locally to different values
	doc1ChangePatch := `[{"op": "add", "path": "/key", "value": "B"}]`
	err = doc1.Apply([]byte(doc1ChangePatch))
	if err != nil {
		t.Error(err)
	}

	doc2ChangePatch := `[{"op": "add", "path": "/key", "value": "C"}]`
	err = doc2.Apply([]byte(doc2ChangePatch))
	if err != nil {
		t.Error(err)
	}

	// merge both by applying each other's remote operations:
	doc1PatchToApplyRemotely, err := doc1.Operations()
	if err != nil {
		t.Error(t)
	}

	doc2PatchToApplyRemotely, err := doc2.Operations()
	if err != nil {
		t.Error(t)
	}

	// merge doc2 state into doc1
	err = doc1.Apply(doc2PatchToApplyRemotely)
	if err != nil {
		t.Error(err)
	}

	// merge doc1 state into doc2
	err = doc2.Apply(doc1PatchToApplyRemotely)
	if err != nil {
		t.Error(err)
	}

	// check results
	bufferDoc1, err := json.Marshal(*doc1)
	if err != nil {
		t.Error(err)
	}
	if string(bufferDoc1) != expectedFinalDoc1 {
		t.Error("Doc1 not correct after Apply")
	}

	bufferDoc2, err := json.Marshal(*doc2)
	if err != nil {
		t.Error(err)
	}
	if string(bufferDoc2) != expectedFinalDoc2 {
		t.Error("Doc2 not correct after Apply")
	}
}

func TestApplyAndMarshall(t *testing.T) {

	patch := []byte(`[
{"op": "add", "path": "/", "value": "user", "id":"1.380503024", "deps": [] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.1", "deps": ["1.1"] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.380503024", "deps": ["1.380503024"] }
]`)

	expectedPatchAfterMarshaling := []byte(`[{"id":"1.380503024","op":"add","path":"/","value":"user"},{"id":"2.380503024","op":"add","path":"/name","value":"Jane"}]`)

	doc := rdoc.Init("document_1")

	err := doc.Apply(patch)
	if err != nil {
		t.Error("Applying valid patch should not err, got ", err)
	}

	buffer, err := json.Marshal(*doc)
	if err != nil {
		t.Error(err)
	}

	if len(buffer) < 10 {
		t.Error("Error encoding doc:", string(buffer))
	}

	if string(buffer) != string(expectedPatchAfterMarshaling) {
		t.Error("Expected marshaling mismatch:",
			string(buffer), string(expectedPatchAfterMarshaling))
	}

	_, err = jsonp.DecodePatch(buffer)
	if err != nil {
		t.Error(err)
	}
}

// 2. General tests

func TestApplyLocalOperations(t *testing.T) {
	// local operations do not require ID and dependencies
	patch := []byte(`[
{"op": "add", "path": "/", "value": "user" },
{"op": "add", "path": "/age", "value": 20 }
]`)

	expectedPatchAfterMarshaling := []byte(`[{"id":"2.380503024","op":"add","path":"/","value":"user"},{"id":"3.380503024","op":"add","path":"/age","value":20}]`)

	doc := rdoc.Init("document_1")

	err := doc.Apply(patch)
	if err != nil {
		t.Error("Applying valid patch should not err, got ", err)
	}

	buffer, err := json.Marshal(*doc)
	if err != nil {
		t.Error(err)
	}

	if string(buffer) != string(expectedPatchAfterMarshaling) {
		t.Error("Expected marshaling mismatch:",
			string(buffer), string(expectedPatchAfterMarshaling))
	}

	_, err = jsonp.DecodePatch(buffer)
	if err != nil {
		t.Error(err)
	}
}

func ExampleDoc_Apply() {
	// Initiates a new document with an unique ID (not enforced by underlying library)
	doc := rdoc.Init("document_id")

	remotePatch := []byte(`[
{"op": "add", "path": "/", "value": "user", "id":"1.380503024", "deps": [] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.1", "deps": ["1.1"] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.380503024", "deps": ["1.380503024"] }
]`)

	localPatch := []byte(`[
{"op": "add", "path": "/", "value": "user" },
{"op": "add", "path": "/age", "value": 20 }
]`)

	err := doc.Apply(remotePatch)
	if err != nil {
		// handle error, most likely reason is malformed patch
	}

	err = doc.Apply(localPatch)
	if err != nil {
		// handle error, most likely reason is malformed patch
	}

	// Marshals applied document into buffer
	buffer, err := json.Marshal(*doc)
	if err != nil {
		// handle error
	}

	fmt.Println(buffer)
}
