package cutter

import (
	"bytes"
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

	stack := stack.New[[]byte](1)
	stack.Push(make([]byte, 4096*2))

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
		data, cluster, nonCluster, columnNames := c.Receiver()
		// for node := start; node != nil; node = node.NextNode() {

		memTable, fileIndex := cluster.GetMemTable(), cluster.GetFileIndex()

		buff, err := c.Stack.Pop()
		if err != nil {
			log.Println(err)
		}

		//I go through each node in the skip-list
		for node := memTable.RootNode(); node != nil; node = node.NextNode() {
			//get right chunk from file for insert node from skip list
			chunk, err := fileIndex.BetweenKey(node.Key())
			if err != nil {
				log.Println(err)
			}

			//jump on chunk of file
			_, err = c.reader.Seek(int64(chunk)*4096, 0)
			if err != nil {
				log.Fatal(err)
			}

			//read that chunk
			n, err := c.reader.Read(buff[:4096])
			if err != nil {
				log.Println(err)
			}

			var (
				counter          int
				writeDataCounter = 4096
			)

			for counter < n {
				maxSize, err := SizeOf(buff, counter) //read max size of single data
				if err != nil {
					log.Panic(err)
				}

				offset := make([]int, len(columnNames)*2)
				if err := CreateOffsetData(buff[5:maxSize], offset, columnNames); err != nil { // parse single data to column
					log.Println(err)
				}

				for index := 0; index < len(node.Key()); index++ {
					position, err := FindIndexColumn(columnNames, cluster.GetByColumn(index))
					if err != nil {
						log.Println(err)
					}

					num := bytes.Compare(node.Key()[index], buff[position:position+1])
					if num == -1 {
						end, err := SizeOf(data, node.GetValue())
						if err != nil {
							log.Println(err)
						}

						copy(buff[writeDataCounter:], data[node.GetValue():end])
						writeDataCounter += end - node.GetValue()

						node = node.NextNode()
					} else if num == 1 {
						copy(buff[writeDataCounter:], buff[counter:maxSize])
						writeDataCounter += maxSize - counter
					}
				}

				counter += maxSize
			}
		}

		c.Stack.Push(buff)
	}
}

func (c *Cutter) Write(file *os.File) {
	for {
		// data := <-c.writeLink

		//
	}
}

// [][]byte
//procitaj prvih 5bytova

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
