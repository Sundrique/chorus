package chorus

import (
	"sync"
	"reflect"
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

var limiter chan bool

func Limit(f interface{}, l int) func(args ...interface{}) {
	limiter = make(chan bool, l)
	fValue := reflect.ValueOf(f)

	wrapped := func(args ...interface{}) {
		limiter <- true

		argsValues := []reflect.Value{}
		for _, v := range args {
			argsValues = append(argsValues, reflect.ValueOf(v))
		}
		fValue.Call(argsValues)

		<-limiter
	}
	return wrapped
}

func Wait(f interface{}, wg sync.WaitGroup) func(args ...interface{}) {

	fValue := reflect.ValueOf(f)

	wrapped := func(args ...interface{}) {
		wg.Add(1)

		argsValues := []reflect.Value{}
		for _, v := range args {
			argsValues = append(argsValues, reflect.ValueOf(v))
		}
		fValue.Call(argsValues)

		wg.Done()
	}

	return wrapped
}
