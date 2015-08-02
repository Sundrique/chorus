package chorus

import (
	. "gopkg.in/check.v1"
	"time"
	"sync"
)

type ConductorSuite struct{}

var _ = Suite(&ConductorSuite{})

func (s *ConductorSuite) TestGoFinished(c *C) {
	var checker []bool

	f := func(args ...interface{}) {
		checker = append(checker, true)
	}

	var conductor Conductor

	for i := 0; i < 10; i++ {
		conductor.Go(f)
	}

	conductor.Wait()

	c.Check(len(checker), Equals, 10)
}

func (s *ConductorSuite) TestGoLimit(c *C) {
	limit := 3

	concurrentCount := 0

	f := func(args ...interface{}) {
		defer func() {
			concurrentCount--
		}()

		concurrentCount++
		if concurrentCount > limit {
			c.Fail()
		}

		time.Sleep(100)
	}

	conductor := new(Conductor)

	conductor.Limit(3)

	for i := 0; i < 10; i++ {
		conductor.Go(f)
	}
	conductor.Wait()
}

func (s *ConductorSuite) TestGoWithArguments(c *C) {
	var sum int

	f := func(args ...interface{}) {
		sum += args[0].(int)
	}

	var conductor Conductor

	for i := 0; i < 10; i++ {
		conductor.Go(f, i)
	}

	conductor.Wait()

	c.Check(sum, Equals, 45)
}

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