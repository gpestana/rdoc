package types

import (
	"fmt"
	"testing"
)

func TestStringer(t *testing.T) {
	r := Register{}
	r.setValue(1)
	r.setValue("two")
	r.setValue(3)
	str := r.String()
	expStr := "<int:1>, <string:two>, <int:3>"

	if str != expStr {
		t.Error(fmt.Sprintf("Register: string expected '%v', got '%v'", expStr, str))
	}
}
