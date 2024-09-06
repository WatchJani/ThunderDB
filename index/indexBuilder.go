package index

import (
	"fmt"
	"root/builder"
	"root/column"
)

const Cluster string = "cluster"

type IndexBuilder struct {
	columns []column.Column
	index   []Index
	builder.Builder
}

// add secondary index
func (ib *IndexBuilder) AddIndex(newIndex Index) error {
	for _, index := range ib.columns {
		if index.GetName() == newIndex.byColumn[0] {
			ib.index = append(ib.index, newIndex)
			return nil
		}
	}

	return fmt.Errorf("column %s not exist", newIndex.name)
}

func NewIndexBuilder(columns []column.Column, index Index) (*IndexBuilder, error) {
	counter, i := 0, 0
	for i < len(index.byColumn) {
		for _, column := range columns {
			if index.byColumn[i] == column.GetName() {
				counter++
				break
			}
		}
		i++

		if i != counter {
			return nil, fmt.Errorf("index column field [%s] not exit", index.byColumn[i])
		}
	}

	return &IndexBuilder{
		columns: columns,
		index:   []Index{index},
	}, nil
}

// new logic, just fix this code
func (ib *IndexBuilder) Choice(userField []string) (Index, bool) {
	var i, j int

	//check cluster index
	clusterIndex := ib.index[0]
	for i < clusterIndex.GetColumnNumber() && j < len(userField) {
		for j < len(userField) {
			if clusterIndex.byColumn[i] == userField[j] {
				userField[i], userField[j] = userField[j], userField[i]
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

	if i != 0 {
		return clusterIndex, false
	}

	//check non-cluster index
	for i := 1; i < len(ib.index); i++ {
		for j := 0; j < len(userField); j++ {
			if ib.index[i].name == userField[j] {
				return ib.index[i], false
			}
		}
	}

	return clusterIndex, true
}
