// Package lclock implements Lamport logical clocks, with helps to be used in
// the context of rdocEvery operation has an unique ID in the
// network. Lamport timestamps ensure that if two operations in different
// network nodes have occurred concurrently, their order is arbitrary but
// deterministic
package lclock

import (
	"fmt"
	"hash/adler32"
	"strconv"
	"strings"
)

const (
	base = 10
)

// Clock holds a Lamport logical clock
type Clock struct {
	seed  int64
	count int64
}

// New initializes and returns a new clock. The `seed` is a string which
// uniquely identifies the clock in the network
func New(seed []byte) Clock {
	s := adler32.Checksum(seed)
	return Clock{
		seed:  int64(s),
		count: 1,
	}
}

// ID returns the id of the clock
func (c *Clock) ID() string {
	return strconv.FormatInt(c.seed, 10)
}

// Tick increments the clock counter
func (c *Clock) Tick() {
	c.count++
}

// Timestamp returns a timestamp that  uniquely identifies the state (id and
// counter) in the network
func (c Clock) Timestamp() string {
	return c.String()
}

// Update performs a clock update based on another clock or string
// representation. If the current Clock count.seed value is higher, no
// changes are done. Othwerise, the clock updates to the upper count
func (c *Clock) Update(rcv interface{}) error {
	var err error
	rcvC := Clock{}
	switch t := rcv.(type) {
	case Clock:
		rcvC = t

	case string:
		rcvC, err = strToClock(t)
		if err != nil {
			return err
		}
	}

	rcvCan, err := rcvC.canonical()
	if err != nil {
		return err
	}

	currCan, err := c.canonical()
	if err != nil {
		return err
	}

	if rcvCan > currCan {
		c.count = rcvC.count
	}
	return nil
}

// CheckTick checks if tick belongs to the clock, or if tick representation is
// invalid
func (c Clock) CheckTick(tick string) (bool, error) {
	tickClock, err := strToClock(tick)
	if err != nil {
		return false,
			fmt.Errorf("Operation ID invalid. Expected <counter>.<seed>, got %v", tick)
	}

	return c.ID() == tickClock.ID(), nil
}

// Returns the canonical value of clock. The canonical value of the logical
// clock is a float64 type in the form of <Clock.count>.<Clock.seed>. The
// Clock.seed value must be unique per Clock in the network.
func (c Clock) canonical() (float64, error) {
	fc, err := strconv.ParseFloat(c.String(), 10)
	return fc, err
}

// Converts string to Clock. The input string is expected to have format
// "counter>.<seed>"
func strToClock(s string) (Clock, error) {
	c := Clock{}
	str := strings.Split(s, ".")
	count, err := strconv.Atoi(str[0])
	if err != nil {
		return c, err
	}
	seed, err := strconv.Atoi(str[1])
	if err != nil {
		return c, err
	}

	c.count = int64(count)
	c.seed = int64(seed)
	return c, nil
}

func (c *Clock) String() string {
	cnt := strconv.FormatInt(c.count, base)
	sd := strconv.FormatInt(c.seed, base)
	return cnt + "." + sd
}

// ConvertString converts a string to a clock representation, or returns an
// error if string representation is invalid
func ConvertString(c string) (Clock, error) {
	return strToClock(c)
}
