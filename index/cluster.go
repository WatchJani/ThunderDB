package index

import (
	"os"
	t "root/b_plus_tree"
	"root/filter"
	"root/manager"
	"root/skip_list"
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

func NewFileReader(file *os.File) FileReader {
	return FileReader{
		file: file,
	}
}

func (f *FileReader) Next() {

}

type NextData interface {
	Next()
}

type InMemory struct {
	tree   *skip_list.SkipList
	buffer []byte
}

func NewInMemory(tree *skip_list.SkipList, buffer []byte) InMemory {
	return InMemory{
		tree:   tree,
		buffer: buffer,
	}
}

func (f *InMemory) Next() {

}

func Generate(manager *manager.Manager, memTable *skip_list.SkipList) (InMemory, InMemory, FileReader) {
	store, memTableMemory, frozenMemory, skipListFrozen := manager.GetAllData()
	return NewInMemory(memTable, memTableMemory),
		NewInMemory(skipListFrozen, frozenMemory),
		NewFileReader(store)
}

//fix the index
//fix cutter cluster inset

func (c *Cluster) Search(key [][]byte, filter []filter.FilterField, index int) ([]byte, error) {
	memTable, frozenMem, store := Generate(c.Manager, c.memTableIndex)

	// node, _ := c.memTableIndex.Search(key, "==")
	// if node != nil {
	// 	fmt.Println(node.GetValue())

	// }
	// // fmt.Println(node.GetValue())
	// // // c.Manager.GetOldIndex().Search(key, filter[index].GetOperation()) //check if exist
	// c.fileIndex.Find(key, filter[0].GetOperation())

	// fmt.Println(key, filter[index:], index)

	return []byte{}, nil
}
