package rdoc

import (
	"encoding/json"
	"fmt"
	"strings"

	jpatch "github.com/evanphx/json-patch"
	"github.com/gpestana/rdoc/idset"
	"github.com/gpestana/rdoc/lclock"
)

// Doc represents a JSON CRDT document
type Doc struct {
	id                 string
	appliedIDs         *idset.Set
	clock              lclock.Clock
	operations         []Operation
	bufferedOperations map[string]Operation
}

// Init returns a new JSON CRDT document
func Init(id string) *Doc {

	return &Doc{
		id:                 id,
		appliedIDs:         idset.New(),
		clock:              lclock.New([]byte(id)),
		operations:         []Operation{},
		bufferedOperations: map[string]Operation{},
	}
}

// Apply applies a valid json patch on the document
func (doc *Doc) Apply(rawPatch []byte) error {
	patch, err := jpatch.DecodePatch(rawPatch)
	if err != nil {
		return err
	}

	appliedRemoteOperations := false

	for _, opRaw := range patch {
		op, err := operationFromPatch(opRaw)
		if err != nil {
			return err
		}

		isFromSameClock, err := doc.clock.CheckTick(op.id)
		if err != nil {
			return err
		}

		// apply local operations and continues
		if isFromSameClock {
			doc.applyOperation(*op)
			continue
		}

		// attempts to apply remote operations by checking if all operation
		// dependencies have been applied on the doc
		if len(doc.appliedIDs.Diff(op.deps)) != 0 {
			doc.bufferedOperations[op.id] = *op
		} else {
			appliedRemoteOperations = true
			delete(doc.bufferedOperations, op.id) // remove buffered operation in case it was buffered
			doc.applyOperation(*op)
		}
	}

	// if remote operation hasbeen applied, attemps to apply buffered operations
	if appliedRemoteOperations {
		doc.tryBufferedOperations()
	}

	return nil
}

func (doc *Doc) applyOperation(operation Operation) {
	doc.appliedIDs.Add(operation.id)
	doc.operations = append(doc.operations, operation)
}

func (doc *Doc) tryBufferedOperations() {
	buffer, err := doc.MarshalFullJSON()
	if err != nil {
		panic(fmt.Sprintf("Buffered operations are not valid -- this should never happen: %v\n", err))
	}

	err = doc.Apply(buffer)
	if err != nil {
		panic(fmt.Sprintf("Error applying buffered operations -- this should never happen: %v\n", err))
	}
}

// MarshalJSON marshals a Doc into a buffer, excluding the deps field on each
// operation
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

// MarshalFullJSON marshals a Doc into a buffer, including the dependencies
// field
func (doc Doc) MarshalFullJSON() ([]byte, error) {
	type operationNoDeps struct {
		ID    string      `json:"id"`
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Deps  []string    `json:"deps"`
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
			Deps:  operation.deps,
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
