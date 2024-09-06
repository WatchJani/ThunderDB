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
func (b *Builder) reservations(dataSize int) int {
	b.Lock()
	defer b.Unlock()

	currentOffset := b.counter
	b.counter += dataSize

	return currentOffset
}

func (b *Builder) insert(data []byte, position int) {
	copy(b.buf[position:], data)
}

func (b *Builder) ParallelWrite(data []byte) {
	b.insert(data, b.reservations(len(data)))
}

func (b *Builder) Write(data []byte) {
	copy(b.buf[b.counter:], data)
	b.counter += len(data)
}

func (b *Builder) GetData() []byte {
	return b.buf[:b.counter]
}

func (b *Builder) IsEnoughSpace(data []byte) bool {
	return cap(b.buf) < b.counter+len(data)
}
