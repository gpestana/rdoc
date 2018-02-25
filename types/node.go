package types

import (
	"log"
)

type CRDT interface {
	AddOpPresence(string)
	RmOpPresence(string)
}

type Node struct {
	Presence []string
}

func (n *Node) AddOpPresence(oid string) {
	log.Println("Add op presence")
}
func (n *Node) RmOpPresence(oid string) {
	log.Println("Rm op presence")
}
