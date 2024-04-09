package safe

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/cenkalti/backoff/v4"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type routineCtx func(ctx context.Context)

// GoPool todo add name,for debug
// Pool is a pool of go routines.
type GoPool struct {
	waitGroup    sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	poolnames    sync.Map
	isDebugStack bool
	//for debug
	curNum int64
	name   string
}

// NewPool creates a Pool.
func NewGoPool(parentCtx context.Context, name string) *GoPool {
	ctx, cancel := context.WithCancel(parentCtx)
	return &GoPool{
		ctx:    ctx,
		cancel: cancel,
		name:   name,
		//pro close
		isDebugStack: true,
	}
}

// GoCtx starts a recoverable goroutine with a context.
func (p *GoPool) GoCtx(goroutine routineCtx) {
	p.waitGroup.Add(1)
	key := ""
	if p.isDebugStack {
		_, file, line, _ := runtime.Caller(1)
		key = fmt.Sprintf("%s:%d", file, line)
		p.poolnames.Store(key, true)
	}
	Go(func() {
		defer p.waitGroup.Done()
		defer p.poolnames.Delete(key)
		atomic.AddInt64(&p.curNum, 1)
		goroutine(p.ctx)
		atomic.AddInt64(&p.curNum, -1)
	})
}

// Stop stops all started routines, waiting for their termination.
func (p *GoPool) Stop() {
	p.cancel()
	Go(func() {
		for {
			time.Sleep(time.Second * 30)
			gamelog.GetGlobalog().Infof("pool(%v) Waiting for all goroutines to finish,cur num(%d)", p.name, atomic.LoadInt64(&p.curNum))
			p.poolnames.Range(func(key, value any) bool {
				gamelog.GetGlobalog().Info(key)
				return true
			})

		}
	})
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
