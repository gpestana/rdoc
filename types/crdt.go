// Implementation of CRDT types supported by the JSON representation: register,
// list and map.
package types

import (
	"log"
)

type CRDT interface {
	// Adds/remove presence ids
	AddOpPresence(string)
	RmOpPresence(string)
	// Deletes CRDT node
	Delete()
	// Data type needs to be unmarshalled into string
	String() string
}

// Implements the common operations for all CRDT types.
type Node struct {
	// Set of operation.Operation ids that have traversed or edited the node
	Presence []string
}

// Adds an operation id to the set pres(k). The set pres(k), is the set of all
//operations that have asserted the existence of the map.
func (n *Node) AddOpPresence(oid string) {
	log.Println("Add op presence")
}

// Removes an operation id from the set pres(k). The set pres(k), is the set of
//all operations that have asserted the existence of the map.
func (n *Node) RmOpPresence(oid string) {
	log.Println("Rm op presence")
}

func (n Node) Delete() {}
