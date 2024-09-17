package cutter

import (
	"log"
	"os"
	"root/linker"
	"root/table"
	"strconv"
	"sync"
)

type Cutter struct {
	// reader *os.File
	link linker.Linker
	// stack.Stack[[]byte]
	writeLink chan WriteLink
	chunk     int
}

type WriteLink struct {
	data  []byte
	chunk int
	wg    *sync.WaitGroup
}

func New(linker linker.Linker, path string, numWorkers int) (*Cutter, error) {
	files := make([]*os.File, numWorkers)

	for index := range files {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		files[index] = file
	}

	c := &Cutter{
		// Stack:     stack,
		link:      linker,
		writeLink: make(chan WriteLink),
	}

	for index := range files {
		go c.Write(files[index])
	}

	return c, nil
}

func (c *Cutter) Cut() {
	singleData := make([]int, 0, 20)

	for {
		dataBlock, size, tableWg, nonCluster, columns := c.link.Receive()
		var wg sync.WaitGroup

		data, stack := (*dataBlock)[:size], make([]byte, 4096)

		offset, counter := 0, 0
		for offset < len(data) {
			end, err := SizeOf(data, offset)
			if err != nil {
				log.Println(err)
			}

			if end > len(data) {
				break
			}

			singleData = append(singleData, offset, end)

			if end > 4096 {
				wg.Add(1)
				c.writeLink <- WriteLink{
					data:  data[:offset],
					chunk: c.chunk,
					wg:    &wg,
				}

				for _, nonClusterIndex := range nonCluster { // update all nonCluster key
					for sdIndex := 0; sdIndex < len(singleData); sdIndex += 2 { // update all data index for one single chunk (block file)
						columnData, _ := table.ReadSingleData(data[singleData[sdIndex]+5:singleData[sdIndex+1]], columns)
						nonClusterIndex.UpdateIndex(table.GenerateKey(nonClusterIndex, columnData, columns), sdIndex*c.chunk)
					}
				}

				data = data[offset:]

				c.chunk++
				offset, counter = 0, 0
				singleData = singleData[:0] //reset all
				continue
			}

			copy(stack[counter:], data[offset:end])
			counter += end - offset
			offset = end
		}

		wg.Add(1)
		c.writeLink <- WriteLink{
			data:  data,
			chunk: c.chunk,
			wg:    &wg,
		}
		c.chunk++

		wg.Wait()

		*dataBlock = nil
		tableWg.Done()
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

		data.wg.Done()
	}
}

func SizeOf(data []byte, offset int) (int, error) {
	num, err := strconv.Atoi(string(data[offset : offset+5]))
	if err != nil {
		return -1, err
	}

	return offset + num + 5, nil
}
