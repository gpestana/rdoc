package rdoc

import (
	"encoding/json"
	"fmt"
	"strings"

	jpatch "github.com/evanphx/json-patch"
	"github.com/gpestana/rdoc/lclock"
)

// Doc represents a JSON CRDT document
type Doc struct {
	id                 string
	clock              lclock.Clock
	operations         []Operation
	bufferedOperations []Operation
}

// Init returns a new JSON CRDT document
func Init(id string) *Doc {

	return &Doc{
		id:                 id,
		clock:              lclock.New([]byte(id)),
		operations:         []Operation{},
		bufferedOperations: []Operation{},
	}
}

// Apply applies a valid json patch on the document
func (doc *Doc) Apply(rawPatch []byte) error {
	patch, err := jpatch.DecodePatch(rawPatch)
	if err != nil {
		return err
	}

	for _, opRaw := range patch {
		op, err := operationFromPatch(opRaw)
		if err != nil {
			return err
		}

		isFromSameClock, err := doc.clock.CheckTick(op.id)
		if err != nil {
			return err
		}
		if !isFromSameClock {
			fmt.Println("Remote operation")
		} else {
			fmt.Println("Local operation")
		}

		// when/where to append?
		doc.operations = append(doc.operations, *op)
	}

	return nil
}

// MarshalJSON marshals a buffer into a crdt doc
func (doc Doc) MarshalJSON() ([]byte, error) {
	type operationNoDeps struct {
		ID    string      `json:"id"`
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}

	buffer := []operationNoDeps{}

	for _, operation := range doc.operations {
		path, err := operation.raw.Path()
		if err != nil {
			return nil, err
		}
		value, err := operation.raw.ValueInterface()
		if err != nil {
			return nil, err
		}

		opNoDeps := operationNoDeps{
			ID:    operation.id,
			Op:    operation.raw.Kind(),
			Path:  path,
			Value: value,
		}

		buffer = append(buffer, opNoDeps)
	}
	return json.Marshal(buffer)
}

// Operation represents the CRDT operations
type Operation struct {
	id   string
	deps []string
	raw  jpatch.Operation
}

func operationFromPatch(rawOp jpatch.Operation) (*Operation, error) {
	rawID := rawOp["id"]
	if rawID == nil {
		return nil,
			fmt.Errorf("Operation must have an associated id, got: %v", rawID)
	}
	id := string(*rawID)
	id = strings.TrimSuffix(id, "\"")
	id = strings.TrimPrefix(id, "\"")

	rawDeps := rawOp["deps"]
	if rawDeps == nil {
		return nil,
			fmt.Errorf("Operation must have an associated dependency, got: %v", rawDeps)
	}

	deps := new([]string)
	err := json.Unmarshal(*rawDeps, deps)
	if err != nil {
		return nil, err
	}

	return &Operation{
		id:   id,
		deps: *deps,
		raw:  rawOp,
	}, nil
}
