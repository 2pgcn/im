package comet

import (
	"github.com/php403/gameim/api/comet"
	"sync"
)

type Bucket struct {
	areas         map[uint64]*Area
	rooms         map[uint64]*Room
	users         map[uint64]*User
	lock          sync.RWMutex
	routines      []chan *comet.Msg
	onlineUserNum uint64
	heartbeat     *HeartHeap // 心跳
}

func NewBucket() *Bucket {
	return &Bucket{
		areas:         make(map[uint64]*Area, 1),
		rooms:         make(map[uint64]*Room, 16),
		users:         make(map[uint64]*User, 1025),
		routines:      make([]chan *comet.Msg, 128),
		heartbeat:     NewHeartbeat(),
		onlineUserNum: 0,
	}
}

// Room get a room by roomid.
func (b *Bucket) Room(rid uint64) (room *Room) {
	b.lock.RLock()
	room = b.rooms[rid]
	b.lock.RUnlock()
	return
}

func (b *Bucket) Area(rid uint64) (area *Area) {
	b.lock.RLock()
	area = b.areas[rid]
	b.lock.RUnlock()
	return
}

// DelRoom delete a room by roomid.
func (b *Bucket) DelRoom(room *Room) {
	b.lock.Lock()
	delete(b.rooms, room.Id)
	b.lock.Unlock()
	room.Close()
}

func (b *Bucket) Rooms() (res map[uint64]struct{}) {
	var (
		roomID uint64
		room   *Room
	)
	res = make(map[uint64]struct{})
	b.lock.RLock()
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = struct{}{}
		}
	}
	b.lock.RUnlock()
	return
}

func (b *Bucket) broadcast(c chan *comet.Msg) {
	for {
		arg := <-c
		switch arg.Type {
		case comet.Type_ROOM:
			if room := b.Room(arg.ToId); room != nil {
				room.Push(arg)
			}
		case comet.Type_AREA:
			if area := b.Area(arg.ToId); area != nil {
				area.Push(arg)
			}
		case comet.Type_PUSH:
			b.lock.RLock()
			user := b.users[arg.ToId]
			b.lock.RUnlock()
			if user != nil {
				user.Push(arg)
			}

		}
	}
}
