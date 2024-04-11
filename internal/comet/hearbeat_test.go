package comet

import (
	"time"
)

var timeTest time.Time
var TestIndex int
var runNum int32

//
//func TestHeartbeat(t *testing.T) {
//	timeTest = time.UnixMicro(1)
//	hb := NewHeartbeat()
//	var runTmpNum = 100
//	tmpDuration := time.Microsecond
//	for i := 0; i < runTmpNum; i++ {
//		randDuration := rand.Intn(runTmpNum)
//		t.Run(fmt.Sprintf("user-%d", i), func(t *testing.T) {
//			tmpTime := timeTest.Add(tmpDuration*1 + time.Nanosecond*time.Duration(randDuration))
//			heap.Push(hb, &HeapItem{Id: userId(strconv.Itoa(i)), Time: tmpTime, fn: func() {
//				atomic.AddInt32(&runNum, 1)
//			}})
//		})
//	}
//	if hb.Len() != runTmpNum {
//		t.Errorf("push hearbeat error:len want(%d) have(%d)", runTmpNum, hb.Len())
//	}
//	hb.Start()
//	if atomic.LoadInt32(&runNum) != int32(runTmpNum) {
//		t.Errorf("heartbeat error:want(%d) num,have(%d)", runTmpNum, atomic.LoadInt32(&runNum))
//	}
//
//}
