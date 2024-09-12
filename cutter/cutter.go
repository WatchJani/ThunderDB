package cutter

import (
	"os"
	"root/linker"
)

type Cutter struct {
	reader *os.File
	linker.Linker
	//stack - buffer za citanje
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

	c := &Cutter{
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
		data, index := c.Receiver()

		//read

		//wreite

	}
}

func (c *Cutter) Write(file *os.File) {
	for {
		data := <-c.writeLink

		//
	}
}
