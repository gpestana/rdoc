package types

type Map struct {
	Node
	Key   string
	Value CRDT
}

func NewMap(k string, t CRDT) *Map {
	return &Map{
		Key:   k,
		Value: t,
	}
}

func NewEmptyMap() *Map {
	return &Map{}
}

func (m *Map) setValue(t CRDT) {
	m.Value = t
}
