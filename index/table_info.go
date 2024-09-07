package index

import (
	"fmt"
	"root/column"
)

type TableInfo struct {
	Table   string
	Columns []column.Column
	Indexes []string
}

func NewTableInfo() *TableInfo {
	return &TableInfo{
		Columns: make([]column.Column, 0),
		Indexes: make([]string, 0),
	}
}

func (t *TableInfo) String() string {
	return fmt.Sprintf("table: %s | columns: %v | index: %v", t.Table, t.Columns, t.Indexes)
}
