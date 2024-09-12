package linker

import "root/skip_list"

type Linker struct {
	link chan Payload
}

type Payload struct {
	data  []byte
	index *skip_list.SkipList
}

func New() Linker {
	return Linker{
		link: make(chan Payload),
	}
}

func (l *Linker) Send(data []byte, index *skip_list.SkipList) {
	l.link <- Payload{
		data:  data,
		index: index,
	}
}

func (l *Linker) Receiver() ([]byte, *skip_list.SkipList) {
	d := <-l.link
	return d.data, d.index
}
