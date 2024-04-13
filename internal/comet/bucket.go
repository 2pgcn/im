package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"sync"
	"time"
)

type roomId string

type Bucket struct {
	ctx           context.Context
	log           gamelog.GameLog
	rooms         map[roomId]*Room
	users         map[userId]*User //所有用户
	lock          sync.RWMutex
	routines      []chan event.Event
	heartHeap     *HeartHeap
	onlineUserNum uint64
}

func (b *Bucket) GetLog() gamelog.GameLog {
	return b.log
}

func NewBucket(ctx context.Context, l gamelog.GameLog) *Bucket {
	return &Bucket{
		ctx:           ctx,
		log:           l,
		rooms:         make(map[roomId]*Room, 16),
		users:         make(map[userId]*User, 1024),
		routines:      make([]chan event.Event, 128),
		heartHeap:     NewHeartbeat(),
		onlineUserNum: 0,
	}
}

func (b *Bucket) StartBucket() {
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
		time.Sleep(time.Second * 5)
		//uheart := b.heartbeat.Pop()
		//if uheart. {
		//}
	}
}

func (b *Bucket) PutUser(user *User) {
	var room *Room
	var ok bool

	b.lock.Lock()
	b.users[user.Uid] = user
	room, ok = b.rooms[user.RoomId]
	b.onlineUserNum++
	b.lock.Unlock()

	user.lock.Lock()
	defer user.lock.Unlock()
	//roomid存在,但bucker里无room,自动创建
	if user.RoomId != "" && !ok {
		room = NewRoom(user.RoomId)
		room.JoinRoom(user)
		b.rooms[user.RoomId] = room
	}
	//heap.Push(b.heartHeap, &HeapItem{
	//	Id:   user.Uid,
	//	Time: time.Now().Add(time.Second * 30),
	//	fn: func() {
	//		user.Close()
	//	},
	//})
	return
}

func (b *Bucket) DeleteUser(uid userId) {
	b.lock.RLock()
	user, ok := b.users[uid]
	room, ok1 := b.rooms[user.RoomId]
	if !ok || !ok1 {
		//已经退出
		b.lock.RUnlock()
		return
	}
	room.ExitRoom(user.Uid)
	b.lock.RUnlock()

	b.lock.Lock()
	delete(b.users, user.Uid)
	b.lock.Unlock()
	user.Close()
}

// 更新心跳时间
func (b *Bucket) UpHeartTime(uid userId, t time.Time) {
	b.heartHeap.look.Lock()
	defer b.heartHeap.look.Unlock()
	//为nil会继承之前fn
	b.heartHeap.Push(&HeapItem{
		Id:   uid,
		Time: t,
		fn: func() {
			b.DeleteUser(uid)
		},
	})
}

// Room get a room by roomid.
func (b *Bucket) Room(rid roomId) (room *Room) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	room = b.rooms[rid]
	if room == nil {
		room = NewRoom(rid)
		b.rooms[rid] = room
	}
	return
}

// DelRoom delete a room by roomid.
func (b *Bucket) DelRoom(room *Room) {
	delete(b.rooms, room.Id)
	room.Close()
}

func (b *Bucket) Rooms() (res map[roomId]struct{}) {
	var (
		rid  roomId
		room *Room
	)
	res = make(map[roomId]struct{})
	b.lock.RLock()
	for rid, room = range b.rooms {
		if room.Online > 0 {
			res[rid] = struct{}{}
		}
	}
	b.lock.RUnlock()
	return
}

func (b *Bucket) broadcast(ev event.Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	c := ev.GetQueueMsg().Data
	switch c.Type {
	case protocol.Type_APP:
		for _, v := range b.users {
			err := v.Push(context.Background(), ev)
			if err != nil {
				b.GetLog().Errorf("room.Push error:%+v", ev)
			}
		}
	case protocol.Type_ROOM:
		if room := b.Room(roomId(c.ToId)); room != nil {
			err := room.Push(b.ctx, ev)
			if err != nil {
				b.log.Errorf("room.Push error:%+v", ev)
			}
		}
	case protocol.Type_CLOSE:
		if room := b.Room(roomId(c.ToId)); room != nil {
			err := room.Push(b.ctx, ev)
			if err != nil {
				b.log.Errorf("room.Push error:%+v", ev)
			}
		}
	case protocol.Type_PUSH:
		b.lock.RLock()
		user := b.users[userId(c.ToId)]
		b.lock.RUnlock()
		if user != nil {
			err := user.Push(user.ctx, ev)
			if err != nil {
				b.log.Errorf("room.Push error:%+v", ev)
			}
		}

	}
}

func (b *Bucket) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()

	for _, v := range b.rooms {
		b.DelRoom(v)
	}
}
