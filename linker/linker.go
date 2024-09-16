package linker

import "sync"

type Linker struct {
	Link chan Payload
}

func New() Linker {
	return Linker{
		Link: make(chan Payload),
	}
}

type Payload struct {
	memTable *[]byte
	end      int
	*sync.WaitGroup
}

func (l *Linker) Receive() (*[]byte, int, *sync.WaitGroup) {
	payload := <-l.Link
	return payload.memTable, payload.end, payload.WaitGroup
}

func NewPayload(memTable *[]byte, end int, wg *sync.WaitGroup) Payload {
	return Payload{
		memTable:  memTable,
		end:       end,
		WaitGroup: wg,
	}
}
