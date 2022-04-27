package comet

import (
	"github.com/php403/gameim/api/comet"
	"sync"
)

type Area struct {
	Id     uint64
	Online uint32
	lock   sync.RWMutex
	head   *User
}

func NewArea(id uint64) (a *Area) {
	a = new(Area)
	a.Id = id
	a.Online = 0
	return
}
func (a *Area) Close() {

}

func (a *Area) JoinArea(user *User) {
	a.lock.Lock()
	if a.head == nil {
		a.head = user
		a.Online++
		a.lock.Unlock()
		return
	}
	a.lock.Unlock()

	user.Next = a.head
	a.head = user
	a.Online++
}

func (a *Area) Push(m *comet.Msg) error {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for u := a.head; u != nil; u = u.Next {
		err := u.Push(m)
		if err != nil {
			return err
		}
	}
	return nil
}
