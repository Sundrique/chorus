package chorus

import (
	. "gopkg.in/check.v1"
	"time"
)

type ConductorSuite struct{}

var _ = Suite(&ConductorSuite{})

func (s *ConductorSuite) TestGoFinished(c *C) {
	var checker []bool

	f := func() {
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

	f := func() {
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