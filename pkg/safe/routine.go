package safe

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"runtime/debug"
	"sync"

	"github.com/cenkalti/backoff/v4"
)

type routineCtx func(ctx context.Context)

// Pool is a pool of go routines.
type GoPool struct {
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a Pool.
func NewGoPool(parentCtx context.Context) *GoPool {
	ctx, cancel := context.WithCancel(parentCtx)
	return &GoPool{
		ctx:    ctx,
		cancel: cancel,
	}
}

// GoCtx starts a recoverable goroutine with a context.
func (p *GoPool) GoCtx(goroutine routineCtx) {
	p.waitGroup.Add(1)
	Go(func() {
		defer p.waitGroup.Done()
		goroutine(p.ctx)
	})
}

// Stop stops all started routines, waiting for their termination.
func (p *GoPool) Stop() {
	p.cancel()
	p.waitGroup.Wait()
}

// Go starts a recoverable goroutine.
func Go(goroutine func()) {
	GoWithRecover(goroutine, defaultRecoverGoroutine)
}

// GoWithRecover starts a recoverable goroutine using given customRecover() function.
func GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}

func defaultRecoverGoroutine(err interface{}) {
	gamelog.GetGlobalog().Errorf("Error in Go routine:%v", err)
	debug.PrintStack()
	//gamelog.GetGlobalog().Errorf("Stack: %+v", debug.Stack())
}

// OperationWithRecover wrap a backoff operation in a Recover.
func OperationWithRecover(operation backoff.Operation) backoff.Operation {
	return func() (err error) {
		defer func() {
			if res := recover(); res != nil {
				defaultRecoverGoroutine(res)
				err = fmt.Errorf("panic in operation: %w", err)
			}
		}()
		return operation()
	}
}
