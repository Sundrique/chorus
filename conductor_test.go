package chorus

import (
	. "gopkg.in/check.v1"
	"time"
	"sync"
)

type ConductorSuite struct{}

var _ = Suite(&ConductorSuite{})

func (s *ConductorSuite) TestLimit(c *C) {
	limit := 3

	concurrentCount := 0
	var wg sync.WaitGroup

	f := func() {
		concurrentCount++
		if concurrentCount > limit {
			c.Fail()
		}

		time.Sleep(100)

		concurrentCount--
		wg.Done()

	}

	wg.Add(10)

	limitedF := Limit(f, limit)
	for i := 0; i < 10; i++ {
		go limitedF()
	}
	wg.Wait()
}

func (s *ConductorSuite) TestWait(c *C) {
	var sum int

	f := func(args ...interface{}) {
		sum += args[0].(int)
	}

	var wg sync.WaitGroup
	fWait := Wait(f, wg)

	for i := 0; i < 10; i++ {
		fWait(i)
	}

	wg.Wait()

	c.Check(sum, Equals, 45)
}