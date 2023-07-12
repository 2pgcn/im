package comet

import (
	"github.com/2pgcn/gameim/api/comet"
	"sync"
)

type areaId uint64
type roomId uint64
type userId uint64

type Bucket struct {
	areas         map[areaId]*Area
	rooms         map[roomId]*Room
	users         map[userId]*User //所有用户
	lock          sync.RWMutex
	routines      []chan *comet.Msg
	onlineUserNum uint64
	heartbeat     *HeartHeap // 心跳
}

func NewBucket() *Bucket {
	return &Bucket{
		areas:         make(map[areaId]*Area, 1),
		rooms:         make(map[roomId]*Room, 16),
		users:         make(map[userId]*User, 1025),
		routines:      make([]chan *comet.Msg, 128),
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
	var area *Area
	//todo 原来没下线需要处理下线重新连接
	if user.RoomId > 0 {
		if room, ok = b.rooms[user.RoomId]; !ok {
			room = NewRoom(user.RoomId)
			b.rooms[user.RoomId] = room
		}
		user.Room = room
	}
	if user.AreaId > 0 {
		if area, ok = b.areas[user.AreaId]; !ok {
			area = NewArea(user.AreaId)
			b.areas[user.AreaId] = area
		}
		user.Area = area
	}
	room.JoinRoom(user)
	area.JoinArea(user)
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

func (b *Bucket) Area(aid areaId) (area *Area) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	area = b.areas[aid]
	if area == nil {
		area = NewArea(aid)
		b.areas[aid] = area
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

func (b *Bucket) broadcast(c *comet.Msg) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	switch c.Type {
	case comet.Type_ROOM:
		if room := b.Room(roomId(c.ToId)); room != nil {
			room.Push(c)
		}
	case comet.Type_AREA:
		if area := b.Area(areaId(c.ToId)); area != nil {
			area.Push(c)
		}
	case comet.Type_PUSH:
		b.lock.RLock()
		user := b.users[userId(c.ToId)]
		b.lock.RUnlock()
		if user != nil {
			user.Push(c)
		}

	}
}
