package index

import (
	t "root/b_plus_tree"
	"root/filter"
	"root/manager"
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

func (c *NonCluster) Search(key [][]byte, filter filter.Filter) {

}
