package builder

import "sync"

type Builder struct {
	buf     []byte
	counter int
	sync.RWMutex
}

func New(buf []byte) Builder {
	return Builder{
		buf: buf,
	}
}

func (b *Builder) Reset() {
	b.counter = 0
}

// use for parallel writing in memory
func (b *Builder) Reservations(data []byte) int {
	b.Lock()
	defer b.Unlock()

	currentOffset := b.counter
	b.counter += len(data)

	return currentOffset
}

func (b *Builder) Insert(data []byte, position int) {
	copy(b.buf[position:], data)
}
