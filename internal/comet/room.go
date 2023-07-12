package comet

import (
	"github.com/2pgcn/gameim/api/comet"
	"sync"
)

type Room struct {
	Id     roomId
	Online uint32
	drop   bool
	lock   sync.RWMutex
	users  map[userId]*User //uid
}

func NewRoom(id roomId) (r *Room) {
	r = new(Room)
	r.Id = id
	r.drop = false
	r.users = make(map[userId]*User, 1024) //todo 从config读取
	r.Online = 0
	return
}

func (r *Room) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.drop = true
}

func (r *Room) JoinRoom(u *User) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.drop {
		return
	}
	r.users[u.Uid] = u
	r.Online++
}

// ExitRoom 一般退出工会,限定操作,一天仅一次
func (r *Room) ExitRoom(u *User) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.drop {
		return
	}
	delete(r.users, u.Uid)
	r.Online--
}

func (r *Room) Push(m *comet.Msg) error {
	r.lock.RLock()
	for _, u := range r.users {
		err := u.Push(m)
		if err != nil {
			r.lock.RUnlock()
			return err
		}
	}
	r.lock.RUnlock()
	return nil
}
