package types

import (
	"testing"
)

func TestTypesImplementCRDT(t *testing.T) {
	m := NewMap()
	r := NewRegister()
	l := NewList()

	implementsCRDT(*m)
	implementsCRDT(*r)
	implementsCRDT(*l)
}

func implementsCRDT(t CRDT) {}
