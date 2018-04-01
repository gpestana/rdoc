package operation

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/gpestana/crdt-json/types"
)

// A cursor identifies unambiguous a position in the JSON document by describing
// the path from the root until the leaf/node selected and the element ID
// Example:
// { "grosseries": { "food": ["pears", "onion"], "others": ["soap"]}, "aList": ["val1", 2, "val3"]}
// cursor: {path: ["grosseries", "food"], key: 1} ==> onion
// cursor: {path: ["aList"], key: 0} ==> val2
type Cursor struct {
	Path []interface{} `json:"path"`
	Id   interface{}   `json:"id"`
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
	for _, p := range c.Path {
		switch t := p.(type) {
		case float64:
		case string:
		default:
			return errors.New(
				fmt.Sprintf("Cursor path types can be Number or String, got a %T", t))
		}
	}

	switch t := c.Id.(type) {
	case float64:
	case string:
	default:
		return errors.New(
			fmt.Sprintf("Cursor Id types can be Number or String, got a %T", t))
	}

	return nil
}
