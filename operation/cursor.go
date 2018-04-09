package operation

import (
	"encoding/json"
	_ "github.com/gpestana/crdt-json/types"
)

// A cursor identifies unambiguous a position in the JSON document by describing
// the path from the root until the leaf/node selected and the element ID
// Example:
// { "grosseries": { "food": ["pears", "onion"], "others": ["soap"]}, "aList": ["val1", 2, "val3"]}
// cursor: [{"MapT":"grosseries"}, {"ListT":"food"}, {"RegisterT": 1}] ==> onion
// cursor: [{"MaptT": "gorsseries"}, {"ListT": "aList"}] ==> ["val1", 2, "val3"]
type Cursor struct {
	Path []PathItem
}

type PathItem struct {
	MapT      string `json:"MapT"`
	ListT     int    `json:"ListT"`
	RegisterT string `json:"RegisterT"`
}

func newCursor(c []byte) (Cursor, error) {
	cur := Cursor{}
	path := []PathItem{}
	err := json.Unmarshal(c, &path)

	if err != nil {
		return Cursor{}, err
	}
	cur.Path = path
	err = cur.validate()
	if err != nil {
		return cur, err
	}
	return cur, nil
}

func (c *Cursor) validate() error {
	return nil
}
