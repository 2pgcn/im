package comet

import (
	"testing"
)

func getUsers(tb testing.TB, nums int) (users []*User) {
	for i := 0; i < nums; i++ {
		u := newUserNotListen(tb)
		u.Uid = userId(i)
		users = append(users, u)
	}
	return
}

func getRoom(tb testing.TB, rid roomId, uNums int) (room *Room) {
	room = NewRoom(rid)
	users := getUsers(tb, uNums)
	for _, v := range users {
		room.JoinRoom(v)
	}
	return room
}

func roomUserExit(tb testing.TB, room *Room, uids []userId) *Room {
	for _, v := range uids {
		room.ExitRoom(v)
	}
	return room
}
func TestRoomJoinAndExit(t *testing.T) {
	uNums := 10
	room := getRoom(t, "1", uNums)
	if room.Online != uint64(uNums) {
		t.Errorf("JoinRoom after room.Online is error:want%d,have %d", uNums, room.Online)
	}
	for k, v := range room.users {
		if k%2 == 0 {
			room.ExitRoom(v.Uid)
		}
	}
	if room.Online != uint64(uNums/2) {
		t.Fatalf("TestRoomJoinAndExit error,online want:%d have:%d", room.Online, uNums/2)
	}

}

//func TestRoomPushMsg(t *testing.T) {
//	ctx := context.Background()
//	uNums := 10
//	room := getRoom(t, "1", uNums)
//	msgStr := "test"
//	msg := &event.Msg{Data: &comet.Msg{Type: comet.Type_ROOM, Msg: }}
//	err := room.Push(ctx, msg)
//	if err != nil {
//		t.Fatalf("room push msg error%s", err.Error())
//	}
//	for _, v := range room.users {
//		uMsgBytes, err := v.Pop(ctx)
//		if err != nil {
//			t.Error(err)
//		}
//		uMsg := uMsgBytes.(*comet.Msg)
//		if uMsg.GetType() != uMsg.GetType() && string(uMsg.GetMsg()) != msgStr {
//			t.Fatalf("room push msg error,want: type-%s,msg-%s have:type-%s,msg-%s", uMsg.GetType(), uMsg.GetMsg(),
//				uMsg.GetType(), uMsg.GetMsg())
//		}
//	}
//}
