package index

import (
	t "root/b_plus_tree"
	"root/manager"
)

type NonCluster struct {
	size     int
	index    *t.Tree[Location]
	byColumn []string
	*manager.Manager
}

type Location struct {
	offset   int
	location byte
}

func NewNonCluster(manager *manager.Manager, byColumn ...string) *NonCluster {
	return &NonCluster{
		index:    t.New[Location](100),
		byColumn: byColumn,
		Manager:  manager,
	}
}
