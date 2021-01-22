package rdoc

import (
	"encoding/json"
	"testing"

	jsonp "github.com/evanphx/json-patch"
)

func TestApplyAndMarshall(t *testing.T) {

	patch := []byte(`[
{"op": "add", "path": "/", "value": "user", "id":"1.380503024", "deps": [] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.1", "deps": ["1.1"] },
{"op": "add", "path": "/name", "value": "Jane", "id":"2.380503024", "deps": ["1.380503024"] }
]`)

	expectedPatchAfterMarshaling := []byte(`[{"id":"1.380503024","op":"add","path":"/","value":"user"},{"id":"2.380503024","op":"add","path":"/name","value":"Jane"}]`)

	doc := Init("document_1")

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
