package index

import (
	"fmt"
	"log"
	"os"
	t "root/b_plus_tree"
	"root/column"
	"root/filter"
	"root/helper"
	"root/manager"
	"root/skip_list"
	"strconv"
)

type Cluster struct {
	size int

	indexType     string
	fileIndex     *t.Tree[int]
	memTableIndex *skip_list.SkipList
	byColumn      string
	*manager.Manager
}

func NewClusterIndex(manager *manager.Manager) *Cluster {
	return &Cluster{
		byColumn:      "id",
		indexType:     "cluster",
		memTableIndex: skip_list.New(32, 60_000, 0.25),
		fileIndex:     t.New[int](100),
		Manager:       manager,
	}
}

func (c *Cluster) NewMemTableIndex() {
	c.memTableIndex = skip_list.New(32, 60_000, 0.25)
}

func (c *Cluster) GetIndexType() string {
	return c.indexType
}

func (c *Cluster) Insert(key [][]byte, offset int) {
	c.size++
	c.memTableIndex.Insert(key, offset)
}

func (c *Cluster) GetDataStructure() *skip_list.SkipList {
	return c.memTableIndex
}

func (c *Cluster) GetByColumn() []string {
	return []string{c.byColumn}
}

//========================================================================

type FileReader struct {
	offset int
	file   *os.File
}

func NewFileReader(file *os.File) *FileReader {
	return &FileReader{
		file: file,
	}
}

func (f *FileReader) Next() {

}

type NextData interface {
	Next()
}

type InMemory struct {
	buffer []byte
	node   *skip_list.Node
}

func NewInMemory(
	tree *skip_list.SkipList,
	buffer []byte,
	key [][]byte,
	filter []filter.FilterField,
	tableFields []column.Column,
) *InMemory {
	if buffer == nil { //Check for frozen memory
		return nil
	}

	node, _ := tree.Search(key, "==") //search first node
	if node == nil {
		return nil
	}

	offset := node.GetValue() //read value from node
	if offset == -1 {
		return nil
	}

	size, err := strconv.Atoi(string(buffer[offset : offset+5])) //get data size
	if err != nil {
		log.Println(err)
		return nil
	}

	data := buffer[offset : offset+size+5]                   //our data
	col, err := helper.ReadSingleData(data[5:], tableFields) //read all column from data
	if err != nil {
		log.Println(err)
		return nil
	}

	var found bool = true
	for _, filterFn := range filter {
		index := helper.GetColumnNameIndex(filterFn.GetField(), tableFields)

		fmt.Println(string(col[index]))
		if !filterFn.GetFilter()(col[index]) {
			found = false
			break
		}
	}

	fmt.Println(col)
	fmt.Println(string(data))
	fmt.Println(found)

	return &InMemory{
		buffer: buffer,
		node:   node,
	}
}

func (f *InMemory) Next() {

}

func Generate(
	manager *manager.Manager,
	memTable *skip_list.SkipList,
	key [][]byte,
	filter []filter.FilterField,
	tableFields []column.Column,
) (*InMemory, *InMemory, *FileReader) {
	store, memTableMemory, frozenMemory, skipListFrozen := manager.GetAllData()
	return NewInMemory(memTable, memTableMemory, key, filter, tableFields),
		NewInMemory(skipListFrozen, frozenMemory, key, filter, tableFields),
		NewFileReader(store)
}

//fix the index
//fix cutter cluster inset

func (c *Cluster) Search(key [][]byte, filter []filter.FilterField, index int, tableFields []column.Column) ([]byte, error) {
	Generate(c.Manager, c.memTableIndex, key, filter, tableFields)

	return []byte{}, nil
}
