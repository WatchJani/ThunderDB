package cutter

import (
	"os"
	"root/linker"
)

//!linker mora imati []index

type Cutter struct {
	reader *os.File
	linker.Linker
	//stack - buffer za citanje
	// stack.Stack[[]byte]
	writeLink chan WriteLink
}

type WriteLink struct {
	data  []byte
	chunk int
}

func New(linker linker.Linker, path string, numWorkers int) (*Cutter, error) {
	files := make([]*os.File, numWorkers+1)

	for index := range files {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		files[index] = file
	}

	// stack := stack.New[[]byte](200)
	// for range 200 {
	// 	stack.Push(make([]byte, 4096))
	// }

	c := &Cutter{
		// Stack:     stack,
		reader:    files[0],
		Linker:    linker,
		writeLink: make(chan WriteLink),
	}

	for index := range files[1:] {
		go c.Write(files[index])
	}

	return c, nil
}

func (c *Cutter) Cut() {
	for {
		// _, index := c.Receiver()

		// start := index[0].RootNode()
		// for node := start; node != nil; node = node.NextNode() {

		// }

		//read

		//wreite

	}
}

func (c *Cutter) Write(file *os.File) {
	for {
		// data := <-c.writeLink

		//
	}
}
