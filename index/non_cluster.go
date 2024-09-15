package index

import "root/manager"

type NonCluster struct {
	size int
	// index    *t.Tree[Location]
	byColumn []string
	*manager.Manager
}
