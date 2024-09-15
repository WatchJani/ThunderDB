package table

import (
	"os"
	"root/column"
	"root/index"
	"root/linker"
	"root/manager"
)

type Table struct {
	linker.Linker
	memTable []byte
	counter  int
	*manager.Manager
	cluster    index.Cluster
	nonCluster []index.NonCluster
	columns    []column.Column
}

func New(linker linker.Linker, reader *os.File) (*Table, error) {
	memTable := make([]byte, 8*1024*1024)

	manager, err := manager.New(memTable, reader)
	if err != nil {
		return nil, err
	}

	return &Table{
		memTable: memTable,
		Linker:   linker,
		Manager:  manager,
	}, nil
}

func (t *Table) Insert(data []byte) int {
	if !t.IsEnoughSpace(data) {
		t.Link <- t.memTable[:t.counter]

		t.memTable = make([]byte, 8*1024*1024)
		t.counter = 0
	}

	offset := t.counter
	copy(t.memTable[t.counter:], data)
	t.counter += len(data)

	return offset
}

func (t *Table) IsEnoughSpace(data []byte) bool {
	return cap(t.memTable) > t.counter+len(data)
}
