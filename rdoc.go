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
func (doc Doc) Apply(rawPatch []byte) error {
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

		fmt.Println(op)
	}

	return nil
}

// MarshalJSON marshals a buffer into a crdt doc
func (doc Doc) MarshalJSON() ([]byte, error) {
	return nil, nil
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
