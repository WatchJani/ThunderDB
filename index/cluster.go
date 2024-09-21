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

// ========================================================================
type NextData interface {
	Next()
}

type FileReader struct {
	offset int
	chunk  int
	file   *os.File
}

func NewFileReader(file *os.File,
	fileIndex *t.Tree[int],
	key [][]byte,
	tableFields []column.Column,
	filter []filter.FilterField,
) *FileReader {
	node, nodeIndex, _ := fileIndex.Find(key, filter[len(key)-1].GetOperation())
	if nodeIndex == -1 {
		return nil
	}

	//find chunk and go on that position
	chunk := node.GetValue(nodeIndex)
	file.Seek(int64(chunk), 0)

	//load data into buffer
	buffer := make([]byte, 4096)
	file.Read(buffer)

	return &FileReader{
		file:  file,
		chunk: chunk,
	}
}

func (f *FileReader) Next() {

}

//!====================================================================================

type InMemory struct {
	buffer      []byte
	node        *skip_list.Node
	tableFields []column.Column
	filter      []filter.FilterField
}

func NewInMemory(
	tree *skip_list.SkipList,
	buffer []byte,
	key [][]byte,
	filter []filter.FilterField,
	tableFields []column.Column,
) (*InMemory, []byte) {
	if buffer == nil { //Check for frozen memory
		return nil, []byte{}
	}

	node, _ := tree.Search(key, filter[len(key)-1].GetOperation()) //search first node

	for {
		found, data, err := checkValidity(node, buffer, tableFields, filter)
		if err != nil {
			return nil, []byte{}
		}

		if found {
			return &InMemory{
				buffer:      buffer,
				node:        node,
				tableFields: tableFields,
				filter:      filter,
			}, data
		}

		node = node.NextNode()
	}
}

func checkValidity(node *skip_list.Node, buffer []byte, tableFields []column.Column, filter []filter.FilterField) (bool, []byte, error) {
	if node == nil {
		return false, []byte{}, fmt.Errorf("node is not exist")
	}

	offset := node.GetValue() //read value from node
	if offset == -1 {
		return false, []byte{}, fmt.Errorf("wrong data for parsing")
	}

	size, err := strconv.Atoi(string(buffer[offset : offset+5])) //get data size
	if err != nil {
		return false, []byte{}, err
	}

	data := buffer[offset : offset+size+5]                   //our data
	col, err := helper.ReadSingleData(data[5:], tableFields) //read all column from data
	if err != nil {
		return false, []byte{}, err
	}

	found, counter := true, 0
	for _, filterFn := range filter {
		index := helper.GetColumnNameIndex(filterFn.GetField(), tableFields)

		if !filterFn.GetFilter()(col[index]) {
			found = false
			break
		}
		counter++
	}

	if counter == 0 {
		return found, []byte{}, fmt.Errorf("index is not usable anymore")
	}

	return found, buffer[offset : offset+5+size], nil
}

func (f *InMemory) Next() []byte {
	node := f.node

	for {
		node = node.NextNode()
		next, data, err := checkValidity(node, f.buffer, f.tableFields, f.filter)
		if err != nil {
			log.Println(err)
			return nil
		}

		if next {
			return data
		}
	}
}

func Generate(
	manager *manager.Manager,
	memTable *skip_list.SkipList,
	key [][]byte,
	filter []filter.FilterField,
	tableFields []column.Column,
	fileIndex *t.Tree[int],
) (*InMemory, []byte, *InMemory, []byte, *FileReader) {
	store, memTableMemory, frozenMemory, skipListFrozen := manager.GetAllData()

	memTableBuffer, dataMemTable := NewInMemory(memTable,
		memTableMemory,
		key,
		filter,
		tableFields,
	)
	frozenBuffer, dataFrozen := NewInMemory(skipListFrozen,
		frozenMemory,
		key,
		filter,
		tableFields,
	)

	return memTableBuffer, dataMemTable, frozenBuffer, dataFrozen, NewFileReader(store, fileIndex, key, tableFields, filter)
}

// fix cutter cluster inset
func (c *Cluster) Search(key [][]byte, filter []filter.FilterField, tableFields []column.Column) ([]byte, error) {
	_, memTableData, _, _, _ := Generate(c.Manager, c.memTableIndex, key, filter, tableFields, c.fileIndex)

	fmt.Println(string(memTableData))

	return []byte{}, nil
}
