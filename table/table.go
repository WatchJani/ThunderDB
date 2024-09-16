package table

import (
	"os"
	"root/column"
	"root/index"
	"root/linker"
	"root/manager"
	"sync"
)

type Table struct {
	linker.Linker
	memTable []byte
	counter  int

	*manager.Manager
	cluster    *index.Cluster
	nonCluster []index.NonCluster
	columns    []column.Column
	wg         sync.WaitGroup
}

func New(columns []column.Column, linker linker.Linker, reader *os.File) (*Table, error) {
	memTable := make([]byte, 8*1024*1024)

	manager, err := manager.New(memTable, reader)
	if err != nil {
		return nil, err
	}

	return &Table{
		memTable: memTable,
		Linker:   linker,
		Manager:  manager,
		cluster:  index.NewClusterIndex(manager),
	}, nil
}

func (t *Table) Insert(data []byte) int {
	if !t.IsEnoughSpace(data) {
		//If my cutter is slower than filling up the entire buffer (memTable),
		//then the data prepared for sending to the cutter will be overwritten.
		//Since the cutter is slower, it won't know that new data has been
		//placed in the same location and will simply report the value as nil.
		//Because of this, I need to wait for my cutter to finish, confirm that
		//the buffer is nil, and only then can I assign a new value to that buffer.

		t.wg.Wait()
		t.wg.Add(1)
		t.SetOld(t.memTable) //set to manager to old
		t.Link <- linker.NewPayload(t.GetOld(), t.counter, &t.wg)

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
