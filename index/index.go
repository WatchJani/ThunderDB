package index

import b "github.com/WatchJani/BPlustTree/btree"

type Index struct {
	name string //not important to much
	size int
	//menage
	//tree
	index    *b.Tree[int, int]
	byColumn []string
}

// func (i *Index) GetColumn() []string {
// 	return i.byColumn
// }

func (i *Index) GetIndexSize() int {
	return i.size
}

func (i *Index) GetColumnNumber() int {
	return len(i.byColumn)
}

func New(name string, byColumn ...string) Index {
	return Index{
		name:     name,
		index:    b.New[int, int](5),
		byColumn: byColumn,
	}
}
