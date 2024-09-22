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
type FileReader struct {
	offset      int
	chunk       int
	file        *os.File
	data        []byte
	buffer      []byte
	tableFields []column.Column
	filter      []filter.FilterField
	node        *t.Node[int]
	nodeIndex   int
}

func (f *FileReader) Read() []byte {
	return f.data
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

	buffer := make([]byte, 4096)

	for {
		//find chunk and go on that position
		chunk := node.GetValue(nodeIndex)
		file.Seek(int64(chunk), 0)

		//load data into buffer
		file.Read(buffer)

		data, offset, err := Offset(0, buffer, tableFields, filter)
		if err != nil {
			log.Println(err)
			return nil
		}

		if offset != -1 {
			return &FileReader{
				file:      file,
				chunk:     chunk,
				offset:    offset,
				data:      data,
				buffer:    buffer,
				filter:    filter,
				node:      node,
				nodeIndex: nodeIndex,
			}
		}

		node, nodeIndex = node.GoForward(nodeIndex)
	}
}

func (f *FileReader) Next() error {
	data, offset, err := Offset(f.offset, f.buffer, f.tableFields, f.filter)
	if err != nil {
		return err
	}

	if offset != -1 {
		f.data = data
		return nil
	}

	node, nodeIndex := f.node.GoForward(f.nodeIndex)

	for {
		chunk := node.GetValue(nodeIndex)
		f.file.Seek(int64(chunk), 0)

		f.file.Read(f.buffer)

		data, offset, err := Offset(0, f.buffer, f.tableFields, f.filter)
		if err != nil {
			return err
		}

		if offset != -1 {
			f.data = data
			return nil
		}

		node = node.NextRight()
	}
}

func Offset(offset int, buffer []byte, tableFields []column.Column, filter []filter.FilterField) ([]byte, int, error) {
	for offset < 4096 {
		size, err := strconv.Atoi(string(buffer[offset : offset+5])) //get data size
		if err != nil {
			return []byte{}, -1, err
		}

		data := buffer[offset : offset+size+5]                   //our data
		col, err := helper.ReadSingleData(data[5:], tableFields) //read all column from data
		if err != nil {
			return []byte{}, -1, err
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
			return []byte{}, -1, fmt.Errorf("index is not useful")
		}

		offset = offset + size + 5
		if found {
			return data, offset, nil
		}
	}

	return []byte{}, -1, nil
}

type InMemory struct {
	buffer      []byte
	node        *skip_list.Node
	tableFields []column.Column
	filter      []filter.FilterField
	data        []byte
}

func (f *InMemory) Read() []byte {
	return f.data
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

	node, _ := tree.Search(key, filter[len(key)-1].GetOperation()) //search first node

	for {
		found, data, err := checkValidity(node, buffer, tableFields, filter)
		if err != nil {
			return nil
		}

		if found {
			return &InMemory{
				buffer:      buffer,
				node:        node,
				tableFields: tableFields,
				filter:      filter,
				data:        data,
			}
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

func (f *InMemory) Next() error {
	node := f.node

	for {
		node = node.NextNode()
		next, data, err := checkValidity(node, f.buffer, f.tableFields, f.filter)
		if err != nil {
			return err
		}

		if next {
			f.data = data
			return nil
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
) []NextData {
	store, memTableMemory, frozenMemory, skipListFrozen := manager.GetAllData()
	res := make([]NextData, 0, 3)

	if memTableObj := NewInMemory(memTable,
		memTableMemory,
		key,
		filter,
		tableFields,
	); memTableObj != nil {
		res = append(res, memTableObj)
	}

	if frozenMemObj := NewInMemory(skipListFrozen,
		frozenMemory,
		key,
		filter,
		tableFields,
	); frozenMemObj != nil {
		res = append(res, frozenMemObj)
	}

	if storeMemObj := NewFileReader(store,
		fileIndex,
		key,
		tableFields,
		filter,
	); storeMemObj != nil {
		res = append(res, storeMemObj)
	}

	return res
}

// fix cutter cluster inset
func (c *Cluster) Search(key [][]byte, filter []filter.FilterField, tableFields []column.Column) ([]byte, error) {
	res := make([]byte, 4096)

	MergeSort(Generate(c.Manager, c.memTableIndex, key, filter, tableFields, c.fileIndex), res)
	fmt.Println(res)
	return res, nil
}

func MergeSort(tro []NextData, data []byte) {
	offset := 0

	for i := 0; i < len(tro); {
		if singleData := tro[i].Read(); singleData != nil {
			copy(data[offset:], singleData)
			offset += len(singleData)

			if err := tro[i].Next(); err != nil {
				tro = removeElement(tro, i)
				continue
			}
		} else {
			tro = removeElement(tro, i)
			continue
		}
		i++
	}
}

func removeElement(slice []NextData, i int) []NextData {
	return append(slice[:i], slice[i+1:]...)
}
