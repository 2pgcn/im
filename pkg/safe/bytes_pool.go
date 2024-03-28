package safe

import (
	"sync"
)

// Buffer buffer.
type Buffer struct {
	buf  []byte
	next *Buffer // next free buffer
}

// Bytes.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// BytePool is a buffer pool.
type BytePool struct {
	lock sync.Mutex
	free *Buffer
	max  int
	num  int
	size int
}

// NewBytePool new a memory buffer pool struct.
func NewBytePool(num, size int) (p *BytePool) {
	p = new(BytePool)
	p.init(num, size)
	return
}

// Init the memory buffer.
func (p *BytePool) Init(num, size int) {
	p.init(num, size)
}

// init the memory buffer.
func (p *BytePool) init(num, size int) {
	p.num = num
	p.size = size
	p.max = num * size
	p.grow()
}

// grow the memory buffer size, and update free pointer.
func (p *BytePool) grow() {
	var (
		i   int
		b   *Buffer
		bs  []Buffer
		buf []byte
	)
	buf = make([]byte, p.max)
	bs = make([]Buffer, p.num)
	p.free = &bs[0]
	b = p.free
	for i = 1; i < p.num; i++ {
		b.buf = buf[(i-1)*p.size : i*p.size]
		b.next = &bs[i]
		b = b.next
	}
	b.buf = buf[(i-1)*p.size : i*p.size]
	b.next = nil
}

// Get a free memory buffer.
func (p *BytePool) Get() (b *Buffer) {
	p.lock.Lock()
	if b = p.free; b == nil {
		p.grow()
		b = p.free
	}
	p.free = b.next
	p.lock.Unlock()
	return
}

// Put back a memory buffer to free.
func (p *BytePool) Put(b *Buffer) {
	p.lock.Lock()
	b.next = p.free
	p.free = b
	p.lock.Unlock()
}
