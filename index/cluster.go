package index

import (
	t "root/b_plus_tree"
	"root/manager"
	"root/skip_list"
)

type Cluster struct {
	size int

	fileIndex     *t.Tree[int]
	memTableIndex *skip_list.SkipList
	byColumn      string
	*manager.Manager
}

func NewClusterIndex(manager *manager.Manager) *Cluster {
	return &Cluster{
		byColumn:      "id",
		memTableIndex: skip_list.New(32, 60_000, 0.25),
		fileIndex:     t.New[int](100),
		Manager:       manager,
	}
}

func (c *Cluster) Insert(key [][]byte, offset int) {
	c.size++
	c.memTableIndex.Insert(key, offset)
}

func (c *Cluster) GetByColumn() []string {
	return []string{c.byColumn}
}
