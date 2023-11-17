package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"strconv"
	"sync"
)

type roomId string
type userId uint64

type Bucket struct {
	ctx           context.Context
	log           gamelog.GameLog
	rooms         map[roomId]*Room
	users         map[userId]*User //所有用户
	lock          sync.RWMutex
	routines      []chan *event.Msg
	onlineUserNum uint64
	heartbeat     *HeartHeap // 心跳
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
		routines:      make([]chan *event.Msg, 128),
		heartbeat:     NewHeartbeat(),
		onlineUserNum: 0,
	}
}

func (b *Bucket) PutUser(user *User) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.users[user.Uid] = user
	b.onlineUserNum++
	var room *Room
	var ok bool
	//todo 原来没下线需要处理下线重新连接
	if user.RoomId != "" {
		if room, ok = b.rooms[user.RoomId]; !ok {
			room = NewRoom(user.RoomId)
			b.rooms[user.RoomId] = room
		}
		user.Room = room
	}
	room.JoinRoom(user)
	return

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
	b.lock.Lock()
	delete(b.rooms, room.Id)
	b.lock.Unlock()
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
	rawData := ev.RawValue()
	c := rawData.(*comet.MsgData)
	switch c.Type {
	case comet.Type_ROOM:
		if room := b.Room(roomId(c.ToId)); room != nil {
			room.Push(ev)
		}
	case comet.Type_PUSH:
		b.lock.RLock()
		uid, err := strconv.ParseUint(c.ToId, 10, 64)
		if err != nil {
			gamelog.Error(err)
		}
		user := b.users[userId(uid)]
		b.lock.RUnlock()
		if user != nil {
			user.Push(ev)
		}
	}
}

func (b *Bucket) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, v := range b.rooms {
		v.Close()
	}
}
