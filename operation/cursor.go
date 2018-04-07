package operation

import (
	"encoding/json"
	_ "github.com/gpestana/crdt-json/types"
)

// A cursor identifies unambiguous a position in the JSON document by describing
// the path from the root until the leaf/node selected and the element ID
// Example:
// { "grosseries": { "food": ["pears", "onion"], "others": ["soap"]}, "aList": ["val1", 2, "val3"]}
// cursor: {path: ["grosseries", "food"], key: 1} ==> onion
// cursor: {path: ["aList"], key: 0} ==> val2
type Cursor struct {
	Path []PathItem  `json:"path"`
	Id   interface{} `json:"id"`
}

type PathItem struct {
	MapT      string `json:"MapT"`
	ListT     int    `json:"ListT"`
	RegisterT string `json:"RegisterT"`
}

func newCursor(c []byte) (Cursor, error) {
	cur := Cursor{}
	err := json.Unmarshal(c, &cur)
	if err != nil {
		return cur, err
	}
	err = cur.validate()
	if err != nil {
		return cur, err
	}
	return cur, nil
}

func (c *Cursor) validate() error {
	return nil
}
