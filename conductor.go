package chorus

import (
	"sync"
)

type Conductor struct {
	sync.WaitGroup
	limiter chan bool
	limited bool
}

func (c *Conductor) Go(fn func(args ...interface{}), args ...interface{}) {
	c.Add(1)
	if c.limited {
		c.limiter <- true
	}

	wrapped := func() {
		fn(args...)

		c.Done()
		if c.limited {
			<-c.limiter
		}
	}
	go wrapped()
}

func (c *Conductor) Limit(limit int) {
	if limit > 0 {
		c.limited = true
		c.limiter = make(chan bool, limit)
	}
}
