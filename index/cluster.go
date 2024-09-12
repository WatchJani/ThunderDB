package index

import (
	"root/skip_list"

	t "root/b_plus_tree"
)

type Cluster struct {
	size int

	fileIndex     *t.Tree[int]
	memTableIndex *skip_list.SkipList
	byColumn      []string
}

func (c *Cluster) Search() {

}

func (c *Cluster) Insert(key [][]byte, offset int) {
	c.size++
	c.memTableIndex.Insert(key, offset)
}

func (c *Cluster) GetIndexSize() int {
	return c.size
}

func (i *Cluster) GetColumnNumber() int {
	return len(i.byColumn)
}

func NewClusterIndex(byColumn ...string) *Cluster {
	return &Cluster{
		// name:  name,
		memTableIndex: skip_list.New(32, 60_000, 0.25),
		byColumn:      byColumn,
	}
}

func (i *Cluster) GetByColumn(index int) string {
	// fmt.Println(i.byColumn)
	return i.byColumn[index]
}

func (i *Cluster) LenByColumn() int {
	return len(i.byColumn)
}
