package comet

import (
	"container/heap"
	"sort"
	"sync"
)

// HeartHeap An IntHeap is a min-heap of ints.
type HeartHeap struct {
	sort.Interface
	hearts []*Heartbeat
	look   sync.RWMutex
}

type hearts []*Heartbeat

// Heartbeat min heap
type Heartbeat struct {
	Id   int //userId
	Time int64
}

func (h *HeartHeap) Len() int {
	return len(h.hearts)
}
func (h *HeartHeap) Less(i, j int) bool {
	h.look.RLock()
	defer h.look.RUnlock()
	return h.hearts[i].Time < h.hearts[j].Time
}
func (h *HeartHeap) Swap(i, j int) {
	h.look.Lock()
	defer h.look.Unlock()
	h.hearts[i], h.hearts[j] = h.hearts[j], h.hearts[i]
}

func (h *HeartHeap) Push(x interface{}) {
	h.look.Lock()
	defer h.look.Unlock()
	h.hearts = append(h.hearts, x.(*Heartbeat))
}

func (h *HeartHeap) Pop() interface{} {
	h.look.Lock()
	defer h.look.Unlock()
	old := h.hearts
	n := len(old)
	x := old[n-1]
	h.hearts = old[0 : n-1]
	return x
}

func NewHeartbeat() *HeartHeap {
	h := &HeartHeap{
		hearts: []*Heartbeat{},
	}
	heap.Init(h)
	return h
}
