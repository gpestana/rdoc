// Package rdoc is a native go implementation of a conflict-free replicated
// JSON data structure (JSON CRDT). A JSON CRDT is a data structure that
// automatically resolves concurrent modifications such that no updates are
// lost, and such that all replicas converge towards the same state.
package rdoc

import (
	"encoding/json"
	"fmt"
	"strings"

	jpatch "github.com/evanphx/json-patch"
	"github.com/gpestana/rdoc/idset"
	"github.com/gpestana/rdoc/lclock"
)

// Doc represents a JSON CRDT document; It maintains the metadata necessary to guarantee
// that the state of the document replicas converge over time without losing data.
type Doc struct {
	id                 string
	appliedIDs         *idset.Set
	clock              lclock.Clock
	operations         []Operation
	bufferedOperations map[string]Operation
}

// Operation the metadata of a valid operation to perform a mutation on the
// JSON CRDT's document state.
type Operation struct {
	id   string
	deps []string
	raw  jpatch.Operation
}

// Init initiates and returns a new JSON CRDT document. The input is a string encoding an
// unique ID of the replica. The application logic *must* ensure that the replica IDs are
// unique within its network
func Init(id string) *Doc {
	return &Doc{
		id:                 id,
		appliedIDs:         idset.New(),
		clock:              lclock.New([]byte(id)),
		operations:         []Operation{},
		bufferedOperations: map[string]Operation{},
	}
}

// Apply applies a valid operation represented as a JSON patch (https://tools.ietf.org/html/rfc6902)
// on the document. Apply handles both local and remote operations.
func (doc *Doc) Apply(rawPatch []byte) error {
	patch, err := jpatch.DecodePatch(rawPatch)
	if err != nil {
		return err
	}

	appliedRemoteOperations := false

	for _, opRaw := range patch {
		op, err := doc.getOperationFromPatch(opRaw)
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

// Operations returns the encoded operations applied to the document. The output
// of this function can be sent over the wire to other replicas, in order to
// achieve convergence.
func (doc Doc) Operations() ([]byte, error) {
	return doc.MarshalFullJSON()
}

// MarshalJSON marshals a Doc into a buffer, excluding the deps field of each
// operation. The returned buffer *contains only the applied operations* that
// mutated the document state. The buffered operations as not included in the
// document serialization.
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
// field of each operation.
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

func (doc *Doc) getOperationFromPatch(rawOp jpatch.Operation) (*Operation, error) {
	var id string
	deps := []string{}

	rawDeps := rawOp["deps"]
	rawID := rawOp["id"]

	// ID is set but dependency set is not (or vice-versa) means that the remote operation
	// is not valid
	if rawID == nil && rawDeps != nil || rawID != nil && rawDeps == nil {
		return nil,
			fmt.Errorf("Remote operation must have an associated id and set of dependencies")
	}

	// if id AND deps are not set, we assume it is a local operation. Thus, we should add
	// metadata based on local replica state (id and dependencies)
	if rawID == nil && rawDeps == nil {
		// local operation
		doc.clock.Tick()
		return &Operation{
			id:   doc.clock.String(),
			deps: doc.appliedIDs.GetIDs(),
			raw:  rawOp,
		}, nil
	}

	// remote operation
	id = string(*rawID)
	id = strings.TrimSuffix(id, "\"")
	id = strings.TrimPrefix(id, "\"")

	err := json.Unmarshal(*rawDeps, &deps)
	if err != nil {
		return nil, err
	}

	return &Operation{
		id:   id,
		deps: deps,
		raw:  rawOp,
	}, nil
}
