package safe

import "testing"

type value struct {
	test string
}

func (v *value) reset() {
	v = &value{}
}

func BenchmarkNewPool(b *testing.B) {
	p := NewPool(func() any { return &value{} })
	arr := []*value{}
	for i := 0; i < b.N; i++ {
		arr = append(arr, p.Get().(*value))
	}
	for i := 0; i < b.N; i++ {
		p.Put(arr[i])
	}
	if b.N != p.Len() {
		b.Errorf("want pool len:%d,have(%d)", p.Len(), b.N)
	}
}
func TestNewPool(t *testing.T) {
	p := NewPool(func() any { return &value{} })
	a := p.Get().(*value)
	b := p.Get().(*value)
	p.Put(a)
	p.Put(b)
	if p.Len() != 2 {
		t.Errorf("want pool len:%d,have(%d)", p.Len(), 2)
	}
}
