package index

import "root/manager"

type Cluster struct {
	size int

	// fileIndex     *t.Tree[int]
	// memTableIndex *skip_list.SkipList
	byColumn []string
	*manager.Manager
}
