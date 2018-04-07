package types

import (
	"fmt"
	"strconv"
	"strings"
)

// Single value register which can hold a string or int
type Register struct {
	Node
	Values []interface{}
}

func NewRegister() *Register {
	return &Register{}
}

// Sets new value to Register. Setting values in Registers is a one value
// operation, although reading from them may return more than one value
func (r *Register) Set(v interface{}) {
	r.Values = append(r.Values, v)
}

func (r *Register) Get() interface{} {
	return nil
}

//TODO
func (r Register) Delete()              {}
func (r Register) AddOpPresence(string) {}
func (r Register) RmOpPresence(string)  {}

// Returns string representation of type Register
func (r Register) String() string {
	out := []string{}
	for _, v := range r.Values {
		switch t := v.(type) {
		case int:
			out = append(out, fmt.Sprintf("<int:%v>", strconv.Itoa(t)))
		case string:
			out = append(out, fmt.Sprintf("<string:%v>", t))
		}
	}
	return strings.Join(out, ", ")
}
