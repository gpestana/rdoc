package operation

import (
	"fmt"
	"testing"
)

func TestNewCursor(t *testing.T) {
	c := []byte(`[{"MapT": "some"}, {"MapT":"list"}, {"ListT": 2}, {"MapT":"aMap"}]`)
	cur, err := newCursor(c)
	if err != nil {
		t.Fatal(err)
	}

	expLenPath := 4
	if len(cur.Path) != expLenPath {
		t.Error(fmt.Sprintf("Lenght of cursor path should be %v, got %v", expLenPath, len(cur.Path)))
	}

	expItemName := "list"
	if cur.Path[1].MapT != expItemName {
		t.Error(fmt.Sprintf("Cursor secod item should be %v, got %v", expItemName, cur.Path[1].MapT))
	}

	expItemValue := 2
	if cur.Path[2].ListT != expItemValue {
		t.Error(fmt.Sprintf("Cursor third item should be %v, got %v", expItemValue, cur.Path[2].MapT))
	}

	c = []byte(`[{"MapT":"ok"}, {"ListT": 2}]`)
	_, err = newCursor(c)
	if err != nil {
		t.Error(fmt.Sprintf("Cursor representation `%v` should be accepted", string(c)))
	}

}

//func TestErrNewCursor(t *testing.T) {
//	c := []byte(`{[{}, "list", 2, "aMap"]`)
//	_, err := newCursor(c)
//	if err == nil {
//		t.Error(fmt.Sprintf("Cursor representation `%v` should not be accepted", string(c)))
//	}
//}
