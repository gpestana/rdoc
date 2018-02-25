package types

type Register struct {
	Node
	Values []string // not quite
}

func (r *Register) setValue(v string) {
	r.Values = append(r.Values, v)
}
