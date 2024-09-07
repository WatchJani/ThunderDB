package index

import b "github.com/WatchJani/BPlustTree/btree"

type Index struct {
	// name string //not important to much
	size int
	//menage
	index    *b.Tree[int, int]
	byColumn []string // mene ne interesuje kolona vec key za search
	//dodati sve kombinacije indexa tako da samo mogu da ih izbucem poslije :D [id, idname, idnameage] :DD
}

func (i *Index) GetIndexSize() int {
	return i.size
}

func (i *Index) GetColumnNumber() int {
	return len(i.byColumn)
}

func New(byColumn ...string) Index {
	return Index{
		// name:  name,
		index: b.New[int, int](5),
		// byColumn: byColumn,
	}
}

