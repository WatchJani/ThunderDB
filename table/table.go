package table

import (
	"fmt"
	"root/builder"
	"root/column"
	"root/index"
	"root/query"
)

const Cluster string = "cluster"

type Table struct {
	columns []column.Column
	index   []index.Index
	builder.Builder
}

func (t *Table) GetColumn() []column.Column {
	return t.columns
}

func (t *Table) GetColumnNum() int {
	return len(t.columns)
}

func (t *Table) GetIndex() []index.Index {
	return t.index
}

// add secondary index
func (t *Table) AddIndex(newIndex index.Index) error {
	for _, index := range t.columns {
		if index.GetName() == newIndex.GetByColumn(0) {
			t.index = append(t.index, newIndex)
			return nil
		}
	}

	return fmt.Errorf("column %s not exist", newIndex.GetByColumn(0))
}

// !add default index if index not specified
func NewTable(columns []column.Column, clusterIndex index.Index) (*Table, error) {
	// counter, i := 0, 0
	// for i < clusterIndex.LenByColumn() {
	// 	for _, column := range columns {
	// 		if clusterIndex.GetByColumn(i) == column.GetName() {
	// 			counter++
	// 			break
	// 		}
	// 	}
	// 	i++

	// 	if i != counter {
	// 		return nil, fmt.Errorf("index column field [%s] not exit", clusterIndex.GetByColumn(i))
	// 	}
	// }

	return &Table{
		columns: columns,
		index:   []index.Index{clusterIndex},
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
	clusterIndex := ib.index[0]
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
		return clusterIndex, filter, false
	}

	//check non-cluster index
	for i := 1; i < len(ib.index); i++ {
		for j := 0; j < len(userQuery); j++ {
			if ib.index[i].GetByColumn(0) == userQuery[j].Field {

				filter := make([]func([]byte) bool, len(userQuery)-1)
				CreateFilter(userQuery[1:], filter)
				return ib.index[i], filter, false
			}
		}
	}

	filter := make([]func([]byte) bool, len(userQuery))
	CreateFilter(userQuery, filter)
	return clusterIndex, filter, true //full scan
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

func (ib *Table) Search(userQuery []Condition) error {
	// fmt.Println(userQuery)

	ib.Choice(userQuery)
	// fmt.Println("index:", index)
	// fmt.Println("filter", filter)
	// fmt.Println("index", indexType)

	return nil
}
