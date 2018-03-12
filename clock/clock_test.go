package clock

import (
	"fmt"
	"testing"
)

func TestConstructor(t *testing.T) {
	seed := []byte("clock_1")
	clk := New(seed)
	initialCount := 1
	seedHash := int64(187368093)

	if clk.count != 1 {
		t.Error(fmt.Sprintf("Constructor: Clock should be initialized with count=%v", initialCount))
	}
	if clk.seed != seedHash {
		t.Error(fmt.Sprintf("Constructor: Clock should be initialized with seed=%v", seedHash))
	}
}

func TestClockTick(t *testing.T) {
	clk := New([]byte("clk"))
	clk.Tick()
	clk.Tick()
	clk.Tick()
	if clk.count != 4 {
		t.Error("Tick: Clock count should be 4")
	}
}

func TestTimestamp(t *testing.T) {
	clk := New([]byte("clock_1"))
	expectedTs := "1.187368093"
	actualTs := clk.Timestamp()

	if actualTs != expectedTs {
		t.Error(fmt.Sprintf("Timestamp: actual ts is %v, should be %v", actualTs, expectedTs))
	}

}

func TestString(t *testing.T) {
	clk := New([]byte("clock_1"))
	expectedStr := "1.187368093"
	actualStr := clk.String()

	if actualStr != expectedStr {
		t.Error(fmt.Sprintf("String: actual string is %v, should be %v", actualStr, expectedStr))
	}
}

func TestUpdateClock(t *testing.T) {
	clk1 := New([]byte("clk1"))

	clk1.Update("123.321")
	if clk1.count != 123 {
		t.Error(fmt.Sprintf("Clock count should be: 123, had  %v", clk1.count))
	}
}

func TestUpdateClockString(t *testing.T) {
	clk1 := New([]byte("clk1"))
	clk2 := New([]byte("clk2"))

	clk1.Tick()
	clk2.Update(clk1)
	if clk2.count != 2 {
		t.Error(fmt.Sprintf("Clock 2 should have same count as Clock 1: %v != %v", clk1.count, clk2.count))
	}

	clk2.Tick()
	clk2.Tick()
	clk1.Update(clk2)
	// test also idempotent update
	clk1.Update(clk2)
	clk1.Update(clk2)
	if clk1.count != 4 {
		t.Error(fmt.Sprintf("Clock 2 should have same count as Clock 1: %v != %v", clk1.count, clk2.count))
	}
}

func TestConvertString(t *testing.T) {
	c := "21"
	expC := int64(21)
	s := "31231233"
	expS := int64(31231233)

	clk, err := ConvertString(fmt.Sprintf("%v.%v", c, s))
	if err != nil {
		t.Fatal(err)
	}
	if clk.count != expC {
		t.Error(fmt.Sprintf("ConvertString: count is %v, should be %v", clk.count, expC))
	}
	if clk.seed != expS {
		t.Error(fmt.Sprintf("ConvertString: seed is %v, should be %v", clk.seed, expS))
	}
}
