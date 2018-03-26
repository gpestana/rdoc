package types

import (
	"fmt"
	"testing"
)

func TestNewEmptyList(t *testing.T) {
	l := NewListEmpty()
	expLength := 0
	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}
}

func TestAddElement(t *testing.T) {
	l := NewListEmpty()
	m := NewEmptyMap()
	m2 := NewEmptyMap()
	l.AddElement(0, *m)
	l.AddElement(1, *m2)

	expLength := 2
	length := l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}

	m3 := NewEmptyMap()
	l.AddElement(0, *m3)
	expLength = 3

	length = l.Length()
	if expLength != length {
		t.Error(fmt.Sprintf("New list expected to have %d lenght, had %d", expLength, length))
	}
}

func TestAddElementOutOfBound(t *testing.T) {
	l := NewListEmpty()
	m := NewEmptyMap()
	m2 := NewEmptyMap()
	l.AddElement(0, *m)
	err := l.AddElement(3, *m2)
	if err == nil {
		t.Error("Error out of bound expected")
	}

}
