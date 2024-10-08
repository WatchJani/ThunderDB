package table

import (
	"fmt"
	"os"
	"root/column"
	"root/helper"
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
	nonCluster []*index.NonCluster
	columns    []column.Column
	wg         sync.WaitGroup
}

func (t *Table) GetClusterIndex() *index.Cluster {
	return t.cluster
}

func (t *Table) GetColumns() []column.Column {
	return t.columns
}

func containsAll(slice1 []column.Column, slice2 []string) (string, bool) {
	for _, str1 := range slice2 {
		found := false
		for _, str2 := range slice1 {
			if str1 == str2.GetName() {
				found = true
				break
			}
		}
		if !found {
			return str1, false
		}
	}
	return "", true
}

func (t *Table) NewIndex(columns ...string) error {
	if len(columns) < 1 {
		return fmt.Errorf("is not log to much")
	}

	if column, ok := containsAll(t.columns, columns); !ok {
		return fmt.Errorf("column %s is not founded", column)
	}

	t.nonCluster = append(t.nonCluster, index.NewNonCluster(t.Manager, columns...))

	return nil
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
		columns:  columns,
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
		t.SetOld(t.memTable, t.cluster.GetDataStructure()) //set to manager to old
		t.Link <- linker.NewPayload(t.GetOld(), t.counter, &t.wg, t.nonCluster, t.cluster, t.columns)

		t.memTable = make([]byte, 8*1024*1024)
		t.cluster.NewMemTableIndex()
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

func (t *Table) ReadSingleData(data []byte) ([][]byte, error) {
	return helper.ReadSingleData(data, t.columns)
}

func (t *Table) GetColumnNameIndex(name string) int {
	return helper.GetColumnNameIndex(name, t.columns)
}

func GenerateKey(index index.Index, columnData [][]byte, columns []column.Column) [][]byte {
	indexColumn := index.GetByColumn()
	key := make([][]byte, len(indexColumn))
	for index, column := range indexColumn {
		key[index] = columnData[helper.GetColumnNameIndex(column, columns)]
	}

	return key
}

func (t *Table) GetIndexes() []index.Index {
	index := make([]index.Index, len(t.nonCluster)+1)
	index[0] = t.cluster

	for i := 1; i < len(index); i++ {
		index[i] = t.nonCluster[i-1]
	}

	return index
}

func (t *Table) GetNonClusterIndex() []*index.NonCluster {
	return t.nonCluster
}
