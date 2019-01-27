package needyerror // by Brad Fitzpatrick https://play.golang.org/p/JBQ3zeVMti

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

// NeedyError is an error which if unchecked will log an error at the
// next GC. It must be created with NewNeedyError.
type NeedyError struct {
	Err  error // the underlying error
	t    time.Time
	seen int32
}

func (e *NeedyError) Error() string {
	atomic.StoreInt32(&e.seen, 1)
	return e.Err.Error()
}

func NewNeedyError(err error) error {
	if err == nil {
		return nil
	}
	e := &NeedyError{
		Err: err,
		t:   time.Now(),
	}
	runtime.SetFinalizer(e, func(ne *NeedyError) {
		if atomic.LoadInt32(&ne.seen) != 0 {
			return
		}
		log.Printf("Unchecked error %v ago: %v", time.Now().Sub(ne.t), ne.Err)
	})
	return e
}
