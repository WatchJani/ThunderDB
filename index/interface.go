package index

import (
	"root/column"
	"root/filter"
)

type Index interface {
	Insert([][]byte, int)
	GetByColumn() []string
	GetIndexType() string
	Search([][]byte, []filter.FilterField, []column.Column) ([]byte, error)
}
