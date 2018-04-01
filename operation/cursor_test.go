package operation

import (
	"fmt"
	"testing"
)

func TestNewCursor(t *testing.T) {
	c := []byte(`{"path": ["some", "list", 2, "aMap"], "id": 1}`)
	cur, err := newCursor(c)
	if err != nil {
		t.Fatal(err)
	}

	expLenPath := 4
	if len(cur.Path) != expLenPath {
		t.Error(fmt.Sprintf("Lenght of cursor path should be %v, got %v", expLenPath, len(cur.Path)))
	}

	expId := float64(1)
	if cur.Id != expId {
		t.Error(fmt.Sprintf("Cursor ID should be %v(%T), got %v(%T)", expId, expId, cur.Id, cur.Id))
	}

	c = []byte(`{"path": ["ok", "list", 2, "aMap"], "id": "somthing"}`)
	_, err = newCursor(c)
	if err != nil {
		t.Error(fmt.Sprintf("Cursor representation `%v` should be accepted", string(c)))
	}

}

func TestErrNewCursor(t *testing.T) {
	c := []byte(`{"path": [{}, "list", 2, "aMap"], "id": 1}`)
	_, err := newCursor(c)
	if err == nil {
		t.Error(fmt.Sprintf("Cursor representation `%v` should not be accepted", string(c)))
	}

}
