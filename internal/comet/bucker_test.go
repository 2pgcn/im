package comet

//func newBucket(tb testing.TB, userNums int) *Bucket {
//	rId := 1
//	userS := make([]*User, userNums)
//	b := NewBucket(context.Background(), gamelog.GetGlobalog())
//	for i := 0; i < userNums; i++ {
//		u := newUserNotListen(tb)
//		u.Uid = userId(i)
//		u.RoomId = roomId(rId)
//
//		b.PutUser(u)
//		userS = append(userS, u)
//	}
//	return b
//}
//func TestNewBucket(t *testing.T) {
//	userNums := 10
//	b := newBucket(t, userNums)
//	if len(b.users) != userNums {
//		t.Fatalf("bucket user join nums is err,want:%d,have:%d", userNums, len(b.users))
//	}
//}

//
//func TestBucketBroadcast(t *testing.T) {
//	//测试bucket消息
//	userNums := 10
//	b := newBucket(t, userNums)
//	msgStr := "TestBucketBroadcast"
//	msg := &event.Msg{
//		Data: &comet.MsgData{
//			Type:   comet.Type_PUSH,
//			ToId:   0,
//			SendId: 0,
//			Msg:    []byte(msgStr),
//		},
//	}
//	b.broadcast(msg)
//	for _, v := range b.users {
//		uEMsg := v.Pop()
//		umsg := uEMsg.(*event.Msg)
//		if umsg.GetData().GetType() != msg.GetData().Type || string(umsg.GetData().GetMsg()) != msgStr {
//			t.Fatalf("TestBucketBroadcast broadcast error,want: %+v,have:%+v", msg, umsg)
//		}
//	}
//}
