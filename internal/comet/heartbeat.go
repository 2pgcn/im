package comet

import (
	"container/heap"
	"sort"
	"sync"
	"time"
)

const infiniteDuration = time.Duration(1<<63 - 1)

// todo 增加对象池
// HeartHeap An IntHeap is a min-heap of ints.
type HeartHeap struct {
	sort.Interface
	//堆具体的数据
	heapItems map[any]*HeapItem

	//堆index排序,数据对应heapItems里元素
	queues []any
	look   sync.RWMutex
}

// Heartbeat min heap
type HeapItem struct {
	Id userId //userId
	//过期时间
	Time time.Time
	fn   func()
}

func (h *HeartHeap) Len() int {
	return len(h.queues)
}
func (h *HeartHeap) Less(i, j int) bool {
	return h.heapItems[h.queues[i]].Time.Before(h.heapItems[j].Time)
}

func (h *HeartHeap) Swap(i, j int) {
	h.queues[i], h.queues[j] = h.queues[j], h.queues[i]
}

// Push 插入或者更新,根据id判定是否存在,更新可不传回掉fn
func (h *HeartHeap) Push(x any) {
	item := x.(*HeapItem)
	if h.heapItems[item.Id] == nil {
		h.queues = append(h.queues, item.Id)
	}
	if item.fn == nil {
		item.fn = h.heapItems[item.Id].fn
	}
	h.heapItems[item.Id] = item
}

//func (h *HeartHeap) Update(item *HeapItem) {
//
//}

func (h *HeartHeap) Pop() any {
	old := h.queues
	n := len(old)
	x := old[n-1]
	h.queues = old[0 : n-1]
	item := h.heapItems[x]
	delete(h.heapItems, x)
	return item
}

func NewHeartbeat() *HeartHeap {
	h := &HeartHeap{
		heapItems: map[any]*HeapItem{},
		queues:    []any{},
		look:      sync.RWMutex{},
	}
	return h
}

func (h *HeartHeap) Start() {
	h.look.Lock()
	defer h.look.Unlock()
	heap.Init(h)
	var (
		item    *HeapItem
		itemKey any
		fn      func()
		index   int
	)
	//每次取0位,有过期则移除,反复取堆顶元素
	for h.Len() > 0 {
		itemKey = h.queues[index]
		item = h.heapItems[itemKey]
		if item.Delay() > 0 {
			break
		}
		fn = item.fn
		fn()
		heap.Remove(h, index)
	}
}

func (h *HeapItem) Delay() time.Duration {
	return time.Until(h.Time)
}
