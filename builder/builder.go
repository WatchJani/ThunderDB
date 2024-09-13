package builder

import (
	"root/linker"
	"sync"

	"github.com/WatchJani/stack"
)

type Builder struct {
	buf []byte
	*stack.Stack[[]byte]
	counter int
	sync.RWMutex
	linker.Linker
}

func New(stack *stack.Stack[[]byte]) Builder {
	buf, _ := stack.Pop()

	return Builder{
		Stack: stack,
		buf:   buf,
	}
}

func (b *Builder) Reset() {
	b.counter = 0
}

func (b *Builder) Insert(data []byte) int {
	if !b.IsEnoughSpace(data) {
		//send to cutter
		// b.Linker.Send(data)

		newBuffer, _ := b.Pop()
		b.buf = newBuffer

		b.Reset()
	}

	offset := b.counter
	copy(b.buf[offset:], data)
	b.counter += len(data)

	return offset
}

func (b *Builder) Write(data []byte) {
	copy(b.buf[b.counter:], data)
	b.counter += len(data)
}

func (b *Builder) GetData() []byte {
	return b.buf[:b.counter]
}

func (b *Builder) IsEnoughSpace(data []byte) bool {
	return cap(b.buf) > b.counter+len(data)
}

func (b *Builder) GetDataOnSpecificPosition(start, end int) []byte {
	return b.buf[start:end]
}
