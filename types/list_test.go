package types

import (
	"fmt"
	"testing"
)

func TestNewEmptyList(t *testing.T) {
	l := NewList()
	expLength := 0
	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}
}

func TestAddElement(t *testing.T) {
	l := NewList()
	m := NewMap()
	m2 := NewMap()
	l.AddElement(0, *m)
	l.AddElement(1, *m2)

	expLength := 2
	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}

	m3 := NewMap()
	l.AddElement(0, *m3)
	expLength = 3

	length = l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}
}

func TestAddElementOutOfBound(t *testing.T) {
	l := NewList()
	m := NewMap()
	m2 := NewMap()
	l.AddElement(0, *m)
	err := l.AddElement(3, *m2)
	if err == nil {
		t.Error("Error out of bound expected")
	}
}

func TestDeleteHead(t *testing.T) {
	l := NewList()
	m := NewMap()
	m2 := NewMap()
	l.AddElement(0, *m)
	l.AddElement(1, *m2)

	l.DeleteElement(0)
	expLength := 1

	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("List expected to have %d lenght, had %d", expLength, length))
	}
}

func TestDeleteElement(t *testing.T) {
	l := NewList()
	m := NewMap()
	m2 := NewMap()
	l.AddElement(0, *m)
	l.AddElement(1, *m2)

	l.DeleteElement(1)
	expLength := 1

	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("List expected to have %d lenght, had %d", expLength, length))
	}
}

func TestDeleteElementOutOfBound(t *testing.T) {
	l := NewList()
	m := NewMap()
	m2 := NewMap()
	l.AddElement(0, *m)
	l.AddElement(1, *m2)
	err := l.DeleteElement(4)

	if err == nil {
		t.Error("Error out of bound expected")
	}
}
