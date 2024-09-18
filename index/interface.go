package index

import "root/filter"

type Index interface {
	Insert([][]byte, int)
	GetByColumn() []string
	GetIndexType() string
	Search([][]byte, filter.Filter)
}
