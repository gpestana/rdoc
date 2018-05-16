// Implements API for crdt-json
// The API exponses methods for changing the local state of the document and
// merging of remote replicas locally.
// The API provides immutability by creating and returning a new document
// everytime a mutation is applied
package main

import (
	"encoding/json"
)

func (j *JSON) Init()   {}
func (j *JSON) Change() {}
func (j *JSON) Merge()  {}

// CRDT JSON can be marshaled into a valid JSON encoding for application
// consumption
func (j JSON) MarshalJSON() ([]byte, error) {
	b := []byte{}
	b, err := json.Marshal(j.Root)
	if err != nil {
		return b, err
	}
	return b, nil
}

func copyDoc(j JSON) *JSON {
	newDoc := j
	return &newDoc
}
