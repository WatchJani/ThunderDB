package index

import (
	t "root/b_plus_tree"
)

type Location struct {
	offset   int
	location byte
}

type NonCluster struct {
	size     int
	index    *t.Tree[Location]
	byColumn []string
}

func (n *NonCluster) Search() {

}

func (n *NonCluster) Insert(key [][]byte, offset int) {
	n.size++
}

func (n *NonCluster) GetByColumn(index int) string {
	return n.byColumn[index]
}

func (n *NonCluster) GetColumnNumber() int {
	return len(n.byColumn)
}

func NewNonCluster(byColumn ...string) *NonCluster {
	return &NonCluster{
		index:    t.New[Location](100),
		byColumn: byColumn,
	}
}
