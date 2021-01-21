package rdoc

import (
	"testing"
)

func TestApply(t *testing.T) {

	patch := []byte(`[
{"op": "add", "path": "/", "value": "user", "id":"1.1", "deps": [] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.1", "deps": ["1.1"] },
{"op": "add", "path": "/name", "value": "Jane", "id":"1.380503024", "deps": [""] }
]`)

	doc := Init("document_1")

	err := doc.Apply(patch)
	if err != nil {
		t.Error("Applying valid patch should not err, got ", err)
	}

}
