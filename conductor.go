package chorus

import (
	"sync"
)

type Conductor struct {
	sync.WaitGroup
	limiter chan bool
	limited bool
}

func (s *Conductor) Go(fn func()) {
	s.Add(1)
	if s.limited {
		s.limiter <- true
	}

	wrapped := func() {
		fn()

		s.Done()
		if s.limited {
			<-s.limiter
		}
	}
	go wrapped()
}

func (s *Conductor) Limit(limit int) {
	if limit > 0 {
		s.limited = true
		s.limiter = make(chan bool, limit)
	}
}
