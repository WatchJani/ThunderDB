package index

import (
	"fmt"
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

func (c *Cluster) Search(key [][]byte, filter []filter.FilterField, index int) ([]byte, error) {
	node, _ := c.memTableIndex.Search(key, "==")
	fmt.Println(node.GetValue())
	// fmt.Println(node.GetValue())
	// // c.Manager.GetOldIndex().Search(key, filter[index].GetOperation()) //check if exist
	// c.fileIndex.Find(key, filter[index].GetOperation())

	// fmt.Println(key, filter[index:], index)

	return []byte{}, nil
}
