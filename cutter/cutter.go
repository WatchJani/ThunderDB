package cutter

import (
	"fmt"
	"log"
	"os"
	"root/column"
	"root/linker"
	"strconv"

	"github.com/WatchJani/stack"
)

type Cutter struct {
	reader *os.File
	linker.Linker
	//stack - buffer za citanje
	stack.Stack[[]byte]
	writeLink chan WriteLink
	chunk     int
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

	stack := stack.New[[]byte](200)
	for range 200 {
		stack.Push(make([]byte, 4096))
	}

	c := &Cutter{
		Stack:     stack,
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
		data, _, _, _ := c.Receiver()

		// memTable, fileIndex := cluster.GetMemTable(), cluster.GetFileIndex()

		offset, counter := 0, 0

		stack, err := c.Stack.Pop()
		if err != nil {
			log.Println(err)
		}

		for offset < len(data) {
			end, err := SizeOf(data, offset)
			if err != nil {
				log.Println(err)
			}

			if end > 4096 {
				c.writeLink <- WriteLink{
					data:  data[:offset],
					chunk: counter,
				}
				data = data[offset:]
				c.chunk++
				counter = 0
			}

			copy(stack[counter:], data[offset:end])
			counter += end - offset
			offset = end
		}
	}
}

func (c *Cutter) Write(file *os.File) {
	for {
		data := <-c.writeLink

		_, err := file.Seek(int64(data.chunk)*4096, 0)
		if err != nil {
			log.Println(err)
		}

		if _, err := file.Write(data.data); err != nil {
			log.Println(err)
		}
	}
}

func SizeOf(data []byte, offset int) (int, error) {
	num, err := strconv.Atoi(string(data[offset : offset+5]))
	if err != nil {
		return -1, err
	}

	return offset + num + 5, nil
}

func CreateOffsetData(data []byte, offset []int, queryTable []column.Column) error {
	index := 0
	for range queryTable {
		index += 5
		size := data[index-5 : index]

		num, err := strconv.Atoi(string(size))
		if err != nil {
			return err
		}

		offset = append(offset, index, index+num)
		// fmt.Println(column, string(conditionsPart[index:index+num]))
		index += num
	}

	return nil
}

func FindIndexColumn(column []column.Column, name string) (int, error) {
	for index, value := range column {
		if value.GetName() == name {
			return index, nil
		}
	}

	return -1, fmt.Errorf("cant find this key")
}
