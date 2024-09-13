package table

import (
	"fmt"
	"root/column"
	"root/index"
	"root/linker"
	"root/query"
)

const Cluster string = "cluster"

type Table struct {
	columns    []column.Column
	cluster    index.Cluster
	nonCluster []index.NonCluster
	linker.Linker
	// builder.Builder
	memTable []byte
	counter  int
}

func (t *Table) GetColumn() []column.Column {
	return t.columns
}

func (t *Table) GetColumnNum() int {
	return len(t.columns)
}

func (t *Table) GetNonClusterIndex() []index.NonCluster {
	return t.nonCluster
}

// add secondary index
func (t *Table) AddIndex(newIndex index.NonCluster) error {
	for _, index := range t.columns {
		if index.GetName() == newIndex.GetByColumn(0) {
			t.nonCluster = append(t.nonCluster, newIndex)
			return nil
		}
	}

	return fmt.Errorf("column %s not exist", newIndex.GetByColumn(0))
}

// !add default index if index not specified
func NewTable(columns []column.Column, clusterIndex index.Cluster, linker linker.Linker) (*Table, error) {
	return &Table{
		columns: columns,
		cluster: clusterIndex,
		// Builder:  builder.New(&stack),
		memTable: make([]byte, 8*1024*1024),
	}, nil
}

type Condition struct {
	Field    string
	Operator []byte
	Value    []byte
	Type     string
}

// new logic, just fix this code
func (ib *Table) Choice(userQuery []Condition) (index.Index, []func([]byte) bool, bool) {
	var i, j int

	//check cluster index
	clusterIndex := ib.cluster
	for i < clusterIndex.GetColumnNumber() && j < len(userQuery) {
		for j < len(userQuery) {
			if clusterIndex.GetByColumn(i) == userQuery[j].Field {
				userQuery[i], userQuery[j] = userQuery[j], userQuery[i]
				i++
				j = i
				break
			}
			j++
		}

		if i != j {
			break
		}
	}

	//non-cluster index

	if i != 0 {
		filter := make([]func([]byte) bool, len(userQuery)-i)
		CreateFilter(userQuery[i:], filter)
		return &clusterIndex, filter, false
	}

	//check non-cluster index
	for i := 1; i < len(ib.nonCluster); i++ {
		for j := 0; j < len(userQuery); j++ {
			if ib.nonCluster[i].GetByColumn(0) == userQuery[j].Field {

				filter := make([]func([]byte) bool, len(userQuery)-1)
				CreateFilter(userQuery[1:], filter)
				return &ib.nonCluster[i], filter, false
			}
		}
	}

	filter := make([]func([]byte) bool, len(userQuery))
	CreateFilter(userQuery, filter)
	return &clusterIndex, filter, true //full scan
}

func CreateFilter(userQuery []Condition, filter []func([]byte) bool) {
	for i := 0; i < len(userQuery); i++ {
		filter[i] = query.GenerateFilter(
			userQuery[i].Operator,
			userQuery[i].Value,
			userQuery[i].Type,
		)
	}
}

func (t *Table) Search(userQuery []Condition) error {
	// fmt.Println(userQuery)

	t.Choice(userQuery)
	// fmt.Println("index:", index)
	// fmt.Println("filter", filter)
	// fmt.Println("index", indexType)

	return nil
}

func (t *Table) Write(data []byte) int {
	if !t.IsEnoughSpace(data) {
		//send to cutter
		t.Linker.Send(t.memTable, t.cluster, t.nonCluster, t.columns)

		t.memTable = make([]byte, 8*1024*1024)
	}

	offset := t.counter
	copy(t.memTable[offset:], data)
	t.counter += len(data)

	return offset
}

func (t *Table) FindIndexColumn(name string) (int, error) {
	for index, value := range t.columns {
		if value.GetName() == name {
			return index, nil
		}
	}

	return -1, fmt.Errorf("cant find this key")
}

func (b *Table) IsEnoughSpace(data []byte) bool {
	return cap(b.memTable) > b.counter+len(data)
}
