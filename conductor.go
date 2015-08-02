package chorus

import (
	"sync"
	"reflect"
)

func Limit(f interface{}, l int) func(args ...interface{}) {
	limiter := make(chan bool, l)
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
