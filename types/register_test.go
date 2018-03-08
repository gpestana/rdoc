package types

import (
	"fmt"
	"testing"
)

func TestStringer(t *testing.T) {
	r := Register{}
	r.Set(1)
	r.Set("two")
	r.Set(3)
	str := r.String()
	expStr := "<int:1>, <string:two>, <int:3>"

	if str != expStr {
		t.Error(fmt.Sprintf("Register: string expected '%v', got '%v'", expStr, str))
	}
}
