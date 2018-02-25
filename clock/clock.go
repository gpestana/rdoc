package clock

import (
	"fmt"
	"hash/adler32"
)

type Clock struct {
	seed  uint32
	count int
}

// Initializes a clock. The `seed` is a string which uniquely identifies the
// clock in the network
func New(seed []byte) Clock {
	s := adler32.Checksum(seed)
	return Clock{
		seed:  s,
		count: 1,
	}
}

func (c *Clock) Tick() {
	c.count++
}

func (c *Clock) Update(rc Clock) {
	if rc.count > c.count {
		c.count = rc.count
	}
}

func (c *Clock) String() string {
	return fmt.Sprintf("%v.%v", c.count, c.seed)
}
