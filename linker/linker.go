package linker

type Linker struct {
	Link chan []byte
}

func New() Linker {
	return Linker{
		Link: make(chan []byte),
	}
}
