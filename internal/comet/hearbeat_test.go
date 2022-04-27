package comet

import (
	"container/heap"
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	hb := NewHeartbeat()
	heap.Init(hb)

	heap.Push(hb, &Heartbeat{
		Id:   3,
		Time: time.Now().Add(2 * time.Second).Unix(),
	})
	heap.Push(hb, &Heartbeat{
		Id:   1,
		Time: time.Now().Unix(),
	})
	heap.Push(hb, &Heartbeat{
		Id:   2,
		Time: time.Now().Add(1 * time.Second).Unix(),
	})
	t.Logf("%+v", heap.Pop(hb))
	t.Logf("%+v", heap.Pop(hb))
}
