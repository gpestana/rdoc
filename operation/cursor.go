package operation

const (
	MapT = iota
	ListT
	RegT
)

// A cursor identifies unambiguous a position in the JSON document by describing
// the path from the root until the leaf/node selected and the element ID
type Cursor struct {
	Path []CursorElement
	Key  interface{}
}

func NewCursor(key interface{}, path ...CursorElement) Cursor {
	c := Cursor{}
	for _, e := range path {
		c.Path = append(c.Path, e)
	}
	c.Key = key
	return c
}

func NewEmptyCursor() Cursor {
	return Cursor{}
}

type CursorElement interface {
	Get() interface{}
	Type() int
}

type MapKey struct {
	Key string
}

func (k MapKey) Get() interface{} {
	return k.Key
}

func (k MapKey) Type() int {
	return MapT
}

type ListKey struct {
	Key int
}

func (k ListKey) Get() interface{} {
	return k.Key
}

func (l ListKey) Type() int {
	return ListT
}
