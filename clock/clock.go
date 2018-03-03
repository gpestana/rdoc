package clock

import (
	"hash/adler32"
	"strconv"
	"strings"
)

const (
	BASE = 10
)

type Clock struct {
	seed  int64
	count int64
}

// Initializes a clock. The `seed` is a string which uniquely identifies the
// clock in the network
func New(seed []byte) Clock {
	s := adler32.Checksum(seed)
	return Clock{
		seed:  int64(s),
		count: 1,
	}
}

// Increments Clock count
func (c *Clock) Tick() {
	c.count++
}

// Returns Timestamp that uniquely identifies the state (clock and count) in the
// network
func (c Clock) Timestamp() string {
	return c.String()
}

// Updates a Clock based on another clock or string representation. If the
// current Clock count.seed value is higher, no changes are done. Othwerise, the
// clock updates to the upper count
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
	cnt := strconv.FormatInt(c.count, BASE)
	sd := strconv.FormatInt(c.seed, BASE)
	return cnt + "." + sd
}

// Convert string to Clock
func ConvertString(c string) (Clock, error) {
	return strToClock(c)
}
