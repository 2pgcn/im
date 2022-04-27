package comet

import "testing"

func BenchmarkArea(b *testing.B) {
	room := NewRoom(1)
	for i := 0; i < b.N; i++ {
		user := NewUser()
		user.Uid = uint64(i)
		room.JoinRoom(user)
	}
}
