package index

import (
	"fmt"
	"log"
	t "root/b_plus_tree"
	"root/column"
	"root/filter"
	"root/manager"
	"strconv"
)

type NonCluster struct {
	size     int
	index    *t.Tree[Location]
	byColumn []string
	*manager.Manager
	indexType string
}

type Location struct {
	offset   int
	location byte
}

func NewNonCluster(manager *manager.Manager, byColumn ...string) *NonCluster {
	return &NonCluster{
		index:     t.New[Location](100),
		byColumn:  byColumn,
		Manager:   manager,
		indexType: "nonCluster",
	}
}

func (c *NonCluster) GetIndexType() string {
	return c.indexType
}

func (c *NonCluster) Insert(key [][]byte, offset int) {
	c.size++
	c.index.Insert(key, Location{
		offset:   offset,
		location: 'm',
	})
}

func (c *NonCluster) GetByColumn() []string {
	return c.byColumn
}

func (c *NonCluster) Update(key [][]byte, offset int) {
	c.index.Insert(key, Location{
		offset:   offset,
		location: 'f',
	})
}

func (c *NonCluster) UpdateIndex(key [][]byte, offset int) {
	c.Update(key, offset)
}

func (c *NonCluster) Search(key [][]byte, filter []filter.FilterField, tableFields []column.Column) ([]byte, error) {
	res := c.GetFreeByte()
	defer c.FlushFreeByte(res)

	dataBlock := c.GetFreeByte()
	defer c.FlushFreeByte(dataBlock)

	memTableMemory, frozenMemory, offset := c.Manager.GetMemTable(), c.Manager.GetFrozenMemory(), 0

	node, index, err := c.index.Find(key, filter[len(key)-1].GetOperation())
	if err != nil {
		log.Println(err)
	}

	for {
		location := node.GetValue(index)
		if location.location == 'm' { //search in memory
			offset := location.offset

			size, err := strconv.Atoi(string(memTableMemory[offset : offset+5])) //get data size
			if err == nil {
				offset += copy(res[offset:], memTableMemory[offset:offset+size+5])
				goto next
			}
			fmt.Println(err)

			size, err = strconv.Atoi(string(frozenMemory[offset : offset+5])) //get data size
			if err != nil {
				log.Println(err)
				break
			}

			offset += copy(res[offset:], memTableMemory[offset:offset+size+5])
		} else {
			file := c.Manager.GetStore()
			file.Seek(int64(offset), 0)

			file.Read(dataBlock)

			size, err := strconv.Atoi(string(dataBlock[offset : offset+5]))
			if err != nil {
				log.Println(err)
			}

			offset += copy(res[offset:], dataBlock[offset:offset+size+5])
		}
	next:
		node, index = node.GoForward(index)
	}

	return res, nil
}
