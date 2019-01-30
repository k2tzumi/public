package needyerror

import (
	"errors"
	"log"
	"runtime"
	"time"
)

func doSomething() error {
	return NewNeedyError(errors.New("Some error"))
}

func ExampleBasicUse() {
	const check = false

	if check {
		err := doSomething()
		if err != nil {
			log.Printf("Got error: %v", err)
		}
	} else {
		// ignoring the error!
		doSomething()
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		runtime.GC()
	}()
	time.Sleep(5 * time.Second)
}
