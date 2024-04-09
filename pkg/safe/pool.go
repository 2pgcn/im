package safe

import (
	"container/list"
	"sync"
)

type Pool[T any] struct {
	look sync.Mutex
	*list.List
	New         func() T
	lastGropNum int
}

func NewPool[T any](n func() T) *Pool[T] {
	return &Pool[T]{
		look: sync.Mutex{},
		List: list.New(),
		New:  n,
	}
}

// Grow 初始化扩容,避免多次内存分配
func (p *Pool[T]) Grow(n int) {
	p.lastGropNum = n
	for i := 0; i < n; i++ {
		item := p.New()
		p.List.PushBack(item)
	}
}

func (p *Pool[T]) Get() T {
	p.look.Lock()
	defer p.look.Unlock()
	if p.List.Len() == 0 {
		p.Grow(p.lastGropNum)
	}
	return p.Remove(p.Front()).(T)
}

func (p *Pool[T]) Put(v T) {
	p.look.Lock()
	defer p.look.Unlock()
	p.List.PushBack(v)
}
