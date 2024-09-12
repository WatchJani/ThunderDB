package linker

import (
	"root/index"
)

type Linker struct {
	link chan Payload
}

type Payload struct {
	data  []byte
	index []index.Index
}

func New() Linker {
	return Linker{
		link: make(chan Payload),
	}
}

func (l *Linker) Send(data []byte, index []index.Index) {
	l.link <- Payload{
		data:  data,
		index: index,
	}
}

func (l *Linker) Receiver() ([]byte, []index.Index) {
	d := <-l.link
	return d.data, d.index
}
