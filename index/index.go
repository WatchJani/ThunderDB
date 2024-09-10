package index

import (
	"root/skip_list"
)

type Index struct {
	// name string //not important to much
	size int
	//menage
	index    *skip_list.SkipList
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
		index:    skip_list.New(32, 60_000, 0.25),
		byColumn: byColumn,
	}
}

func (i *Index) GetByColumn(index int) string {
	return i.byColumn[index]
}

func (i *Index) LenByColumn() int {
	return len(i.byColumn)
}
