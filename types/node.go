package types

import (
	"log"
)

type CRDT interface {
	//InsertAfter(Node) // only for list
	//Delete()

	// Adds/remove presence ids
	AddOpPresence(string)
	RmOpPresence(string)
	// Data type needs to be unmarshalled into string
	String() string
}

type Node struct {
	// Set of operation.Operation ids that have traversed or edited the node
	Presence []string
}

// Adds operation.Operation id to set of presence
func (n *Node) AddOpPresence(oid string) {
	log.Println("Add op presence")
}

// Removes operation.Operation id to set of presence
func (n *Node) RmOpPresence(oid string) {
	log.Println("Rm op presence")
}
