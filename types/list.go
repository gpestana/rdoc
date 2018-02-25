package types

type List struct {
	Node
	Elements []*CRDT
}

func (l *List) setValue(t *CRDT) {
	l.Elements = append(l.Elements, t)
}
